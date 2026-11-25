import type { Anchor, AnnotationThread } from './types'

export const DEFAULT_EXCLUDE_SELECTOR =
  'pre, code, svg, .mermaid, .katex, .atk-annotation-host, script, style'

function closestExcluded(node: Node | null, selector: string): Element | null {
  if (!node) return null
  const el = node.nodeType === Node.ELEMENT_NODE ? (node as Element) : node.parentElement
  return el ? el.closest(selector) : null
}

export function selectionAllowed(
  range: Range,
  scope: HTMLElement,
  excludeSelector = DEFAULT_EXCLUDE_SELECTOR,
): boolean {
  if (!scope.contains(range.commonAncestorContainer)) return false
  if (
    closestExcluded(range.startContainer, excludeSelector) ||
    closestExcluded(range.endContainer, excludeSelector)
  ) {
    return false
  }
  try {
    const holder = document.createElement('div')
    holder.appendChild(range.cloneContents())
    if (holder.querySelector(excludeSelector)) return false
  } catch {
    return false
  }
  return true
}

function collectTextNodes(scope: HTMLElement): Text[] {
  const nodes: Text[] = []
  const walker = document.createTreeWalker(scope, NodeFilter.SHOW_TEXT)
  let current = walker.nextNode()
  while (current) {
    const parent = (current as Text).parentElement
    if (!parent || !parent.closest('script, style')) nodes.push(current as Text)
    current = walker.nextNode()
  }
  return nodes
}

function positionOfRange(
  range: Range,
  scope: HTMLElement,
): { start: number; end: number } | undefined {
  if (
    range.startContainer.nodeType !== Node.TEXT_NODE ||
    range.endContainer.nodeType !== Node.TEXT_NODE
  ) {
    return undefined
  }
  let start = -1
  let end = -1
  let offset = 0
  for (const node of collectTextNodes(scope)) {
    if (node === range.startContainer) start = offset + range.startOffset
    if (node === range.endContainer) end = offset + range.endOffset
    offset += node.data.length
    if (start >= 0 && end >= 0) break
  }
  if (start < 0 || end < 0) return undefined
  return { start, end }
}

function findSection(scope: HTMLElement, range: Range): string {
  const headings = Array.from(scope.querySelectorAll<HTMLElement>('h1, h2, h3, h4'))
  let label = ''
  for (const heading of headings) {
    const relation = heading.compareDocumentPosition(range.startContainer)
    if (relation & Node.DOCUMENT_POSITION_FOLLOWING) label = (heading.textContent || '').trim()
    else break
  }
  return label.slice(0, 200)
}

export async function describeSelection(range: Range, scope: HTMLElement): Promise<Anchor> {
  const dq: any = await import('dom-anchor-text-quote')
  const quote = dq.fromRange(scope, range)
  return {
    quote: {
      type: 'TextQuoteSelector',
      exact: quote.exact,
      prefix: quote.prefix || '',
      suffix: quote.suffix || '',
    },
    position: positionOfRange(range, scope),
    section: findSection(scope, range),
  }
}

export async function highlightThreads(
  scope: HTMLElement,
  list: AnnotationThread[],
  onOrphan: (id: number) => void,
): Promise<() => void> {
  const cleanups: Array<() => void> = []
  if (list.length === 0) return () => {}

  const [{ highlightText }, dq]: any[] = await Promise.all([
    import('@apache-annotator/dom'),
    import('dom-anchor-text-quote'),
  ])

  for (const thread of list) {
    const quote = thread.root.anchor?.quote
    if (!quote?.exact) {
      onOrphan(thread.root.id)
      continue
    }
    const selector: Record<string, string> = { exact: quote.exact }
    if (quote.prefix) selector.prefix = quote.prefix
    if (quote.suffix) selector.suffix = quote.suffix

    let range: Range | null = null
    try {
      range = dq.toRange(scope, selector)
    } catch (err) {
      console.warn('[artalk-annotation] anchor match failed', thread.root.id, err)
    }
    if (!range) {
      onOrphan(thread.root.id)
      continue
    }
    try {
      cleanups.push(
        highlightText(range, 'mark', {
          class: 'atk-annotation-mark',
          'data-annotation-id': String(thread.root.id),
        }),
      )
    } catch (err) {
      console.warn('[artalk-annotation] highlight failed', thread.root.id, err)
      onOrphan(thread.root.id)
    }
  }

  return () => {
    for (const cleanup of cleanups) {
      try {
        cleanup()
      } catch {
        /* ignore */
      }
    }
    try {
      scope.normalize()
    } catch {
      /* ignore */
    }
  }
}

export function sortByDocumentOrder(
  scope: HTMLElement,
  list: AnnotationThread[],
): AnnotationThread[] {
  const order = new Map<number, number>()
  scope.querySelectorAll<HTMLElement>('mark.atk-annotation-mark').forEach((mark, index) => {
    const id = Number(mark.dataset.annotationId)
    if (id && !order.has(id)) order.set(id, index)
  })
  return [...list].sort((a, b) => {
    const ai = order.get(a.root.id) ?? Number.MAX_SAFE_INTEGER
    const bi = order.get(b.root.id) ?? Number.MAX_SAFE_INTEGER
    return ai - bi
  })
}

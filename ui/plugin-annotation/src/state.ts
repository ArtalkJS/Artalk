import type { AnnotationCommentData, AnnotationThread, Filter, SortMode } from './types'

export interface AnnotationState {
  threads: AnnotationThread[]
  activeId: number | null
  filter: Filter
  sort: SortMode
  search: string
  orphans: Set<number>
  loading: boolean
  loadingMore: boolean
  total: number
  error: string
}

export function initialState(): AnnotationState {
  return {
    threads: [],
    activeId: null,
    filter: 'all',
    sort: 'position',
    search: '',
    orphans: new Set(),
    loading: false,
    loadingMore: false,
    total: 0,
    error: '',
  }
}

export function buildThreads(comments: AnnotationCommentData[]): AnnotationThread[] {
  const roots = comments.filter((c) => c.rid === 0)
  const children = new Map<number, AnnotationCommentData[]>()
  for (const c of comments) {
    if (c.rid === 0) continue
    const list = children.get(c.rid) || []
    list.push(c)
    children.set(c.rid, list)
  }

  const collect = (root: AnnotationCommentData): AnnotationCommentData[] => {
    const out: AnnotationCommentData[] = []
    const walk = (id: number) => {
      for (const child of children.get(id) || []) {
        out.push(child)
        walk(child.id)
      }
    }
    walk(root.id)
    return out
  }

  return roots.map((root) => ({ root, replies: collect(root), orphan: false }))
}

export function upsertThread(
  threads: AnnotationThread[],
  thread: AnnotationThread,
): AnnotationThread[] {
  let found = false
  const out = threads.map((item) => {
    if (item.root.id !== thread.root.id) return item
    found = true
    return thread
  })
  if (!found) out.push(thread)
  return out
}

export function appendReply(
  threads: AnnotationThread[],
  rootId: number,
  reply: AnnotationCommentData,
): AnnotationThread[] {
  return threads.map((thread) => {
    if (thread.root.id !== rootId) return thread
    if (thread.replies.some((item) => item.id === reply.id)) return thread
    return { ...thread, replies: [...thread.replies, reply] }
  })
}

export function removeThread(threads: AnnotationThread[], id: number): AnnotationThread[] {
  return threads.filter((thread) => thread.root.id !== id)
}

export function removeReply(
  threads: AnnotationThread[],
  rootId: number,
  replyId: number,
): AnnotationThread[] {
  return threads.map((thread) =>
    thread.root.id === rootId
      ? { ...thread, replies: thread.replies.filter((reply) => reply.id !== replyId) }
      : thread,
  )
}

export function sortForDisplay(threads: AnnotationThread[], mode: SortMode): AnnotationThread[] {
  if (mode !== 'newest') return threads
  return [...threads].sort((a, b) => Date.parse(b.root.date) - Date.parse(a.root.date))
}

export function searchThreads(threads: AnnotationThread[], query: string): AnnotationThread[] {
  const keyword = query.trim().toLowerCase()
  if (!keyword) return threads
  return threads.filter((thread) => {
    const haystack = [
      thread.root.content,
      thread.root.nick,
      thread.root.anchor?.quote.exact || '',
      ...thread.replies.map((reply) => `${reply.nick} ${reply.content}`),
    ]
      .join(' ')
      .toLowerCase()
    return haystack.includes(keyword)
  })
}

export function filterThreads(
  threads: AnnotationThread[],
  filter: Filter,
  userId?: number,
): AnnotationThread[] {
  switch (filter) {
    case 'mine':
      return threads.filter((thread) => userId !== undefined && thread.root.user_id === userId)
    case 'resolved':
      return threads.filter((thread) => thread.root.anchor_resolved === true)
    case 'pending':
      return threads.filter((thread) => thread.root.is_pending === true)
    default:
      return threads
  }
}

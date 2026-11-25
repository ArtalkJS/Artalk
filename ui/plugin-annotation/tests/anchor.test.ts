import { describe, expect, it } from 'vitest'
import { describeSelection, highlightThreads } from '../src/anchor'
import type { AnnotationThread } from '../src/types'

describe('anchor', () => {
  it('describes a selection and restores the highlight', async () => {
    document.body.innerHTML = '<main><p>这是第一段可划线的正文，用于锚点测试。</p></main>'
    const root = document.querySelector('main') as HTMLElement
    const paragraph = root.querySelector('p') as HTMLElement
    const textNode = paragraph.firstChild as Text

    const range = document.createRange()
    range.setStart(textNode, 0)
    range.setEnd(textNode, 10)

    const anchor = await describeSelection(range, root)
    expect(anchor.quote.exact).toBe(textNode.data.slice(0, 10))
    expect(anchor.section).toBe('')

    const thread: AnnotationThread = {
      root: {
        id: 1,
        content: 'note',
        nick: 'alice',
        email_encrypted: '',
        link: '',
        rid: 0,
        ua: '',
        date: new Date().toISOString(),
        is_collapsed: false,
        is_pending: false,
        is_pinned: false,
        is_verified: false,
        is_allow_reply: true,
        page_key: '/doc',
        site_name: 'Test',
        visible: true,
        vote_up: 0,
        vote_down: 0,
        anchor,
      },
      replies: [],
      orphan: false,
    }

    const orphans: number[] = []
    const cleanup = await highlightThreads(root, [thread], (id) => orphans.push(id))
    expect(orphans).toEqual([])
    expect(root.querySelector('mark.atk-annotation-mark')?.textContent).toBe(anchor.quote.exact)

    cleanup()
    expect(root.querySelector('mark.atk-annotation-mark')).toBeNull()
  })
})

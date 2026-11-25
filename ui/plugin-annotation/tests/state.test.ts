import { describe, expect, it } from 'vitest'
import {
  appendReply,
  buildThreads,
  filterThreads,
  removeReply,
  removeThread,
  searchThreads,
  sortForDisplay,
  upsertThread,
} from '../src/state'
import type { AnnotationCommentData, AnnotationThread } from '../src/types'

function comment(id: number, rid = 0, content = `text-${id}`): AnnotationCommentData {
  return {
    id,
    rid,
    content,
    nick: 'alice',
    email_encrypted: '',
    link: '',
    ua: '',
    date: '2026-01-01 00:00:00',
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
    user_id: 1,
    anchor: { quote: { type: 'TextQuoteSelector', exact: content } },
    is_annotation: rid === 0,
    anchor_resolved: false,
  }
}

describe('state', () => {
  it('builds threads with nested replies', () => {
    const threads = buildThreads([comment(1), comment(2, 1), comment(3, 2)])
    expect(threads).toHaveLength(1)
    expect(threads[0].replies.map((reply) => reply.id)).toEqual([2, 3])
  })

  it('upserts and removes threads', () => {
    const thread: AnnotationThread = { root: comment(1), replies: [], orphan: false }
    const next = upsertThread([], thread)
    expect(next).toHaveLength(1)
    expect(removeThread(next, 1)).toHaveLength(0)
  })

  it('appends replies without duplicates', () => {
    const thread: AnnotationThread = { root: comment(1), replies: [], orphan: false }
    const once = appendReply([thread], 1, comment(2, 1))
    const twice = appendReply(once, 1, comment(2, 1))
    expect(twice[0].replies).toHaveLength(1)
  })

  it('removes a reply from a thread', () => {
    const thread: AnnotationThread = {
      root: comment(1),
      replies: [comment(2, 1), comment(3, 1)],
      orphan: false,
    }
    const next = removeReply([thread], 1, 2)
    expect(next[0].replies.map((reply) => reply.id)).toEqual([3])
  })

  it('filters and searches', () => {
    const threads: AnnotationThread[] = [
      { root: { ...comment(1), anchor_resolved: true }, replies: [], orphan: false },
      { root: comment(2, 0, 'other'), replies: [], orphan: false },
    ]
    expect(filterThreads(threads, 'resolved')).toHaveLength(1)
    expect(filterThreads(threads, 'mine', 1)).toHaveLength(2)
    expect(searchThreads(threads, 'other')).toHaveLength(1)
  })

  it('sorts newest first', () => {
    const threads: AnnotationThread[] = [
      { root: { ...comment(1), date: '2026-01-01 00:00:00' }, replies: [], orphan: false },
      { root: { ...comment(2), date: '2026-02-01 00:00:00' }, replies: [], orphan: false },
    ]
    expect(sortForDisplay(threads, 'newest')[0].root.id).toBe(2)
  })
})

import { expect, it, vi } from 'vitest'
import CommentNode from '@/comment/comment-node'
import { Defaults } from '@/defaults'
import { parseDate } from '@/lib/utils'
import { makeNestCommentNodeList } from '@/list/nest'
import type { CommentData } from '@/types'

it('keeps timestamp parsing, comment display and date sorting compatible', () => {
  const instant = Date.parse('2025-01-02T03:04:05Z')
  const legacy = '2025-01-01 03:04:05'
  const base: CommentData = {
    id: 1,
    rid: 0,
    content: '',
    nick: 'Synthetic user',
    email_encrypted: '',
    link: '',
    ua: '',
    date: '2025-01-08 00:00:00',
    is_collapsed: false,
    is_pending: false,
    is_pinned: false,
    is_verified: false,
    is_allow_reply: true,
    page_key: '/test',
    visible: true,
    site_name: 'test',
    vote_up: 0,
    vote_down: 0,
  }
  const opts = {
    ...Defaults,
    flatMode: false,
    getApi: vi.fn(),
    replyComment: vi.fn(),
    editComment: vi.fn(),
    dateFormatter: (date: Date) => date.toISOString(),
  }

  // Unix seconds take precedence; UTC strings must also preserve the instant.
  expect(+parseDate({ ...base, date_unix: instant / 1000, date_utc: '2025-01-09T00:00:00Z' })).toBe(
    instant,
  )
  for (const date_utc of [
    '2025-01-02T03:04:05Z',
    '2025-01-02T11:04:05+08:00',
    '2025-01-01T22:04:05-05:00',
  ]) {
    expect(+parseDate({ ...base, date_utc })).toBe(instant)
  }

  const comments: CommentData[] = [
    base,
    { ...base, id: 3, rid: 1, date_utc: '2025-01-03T03:04:05Z' },
    { ...base, id: 4, rid: 1, date: legacy },
    { ...base, id: 2, rid: 1, date_unix: instant / 1000, date_utc: '2025-01-09T00:00:00Z' },
  ]
  const expected = [instant, Date.parse('2025-01-03T03:04:05Z'), +new Date(2025, 0, 1, 3, 4, 5)]
  for (const [index, id] of [2, 3, 4].entries()) {
    const data = comments.find((comment) => comment.id === id)!
    const node = new CommentNode(data, opts)
    expect(+parseDate(data)).toBe(expected[index])
    expect(+node.getCreatedAt()).toBe(expected[index])
    expect(node.getDateFormatted()).toBe(new Date(expected[index]).toISOString())
  }
  expect(makeNestCommentNodeList(comments, 'DATE_ASC')[0].children.map((node) => node.id)).toEqual([
    4, 2, 3,
  ])
  expect(makeNestCommentNodeList(comments, 'DATE_DESC')[0].children.map((node) => node.id)).toEqual(
    [3, 2, 4],
  )
})

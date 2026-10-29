import { expect, it, vi } from 'vitest'
import CommentNode from '@/comment/comment-node'
import { Defaults } from '@/defaults'
import { parseDate } from '@/lib/utils'
import { makeNestCommentNodeList } from '@/list/nest'
import type { CommentData } from '@/types'

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

it('parseDate prefers date_unix and falls back to legacy date', () => {
  const instant = Date.parse('2025-01-02T03:04:05Z')
  expect(+parseDate({ ...base, date_unix: instant / 1000 })).toBe(instant)
  expect(+parseDate({ date: '2025-01-01 03:04:05' })).toBe(+new Date(2025, 0, 1, 3, 4, 5))
})

it('comment display and date sorting use parsed timestamps', () => {
  const t = (iso: string) => Date.parse(iso) / 1000
  const comments: CommentData[] = [
    base,
    { ...base, id: 2, rid: 1, date_unix: t('2025-01-03T00:00:00Z') },
    { ...base, id: 3, rid: 1, date_unix: t('2025-01-02T00:00:00Z') },
    { ...base, id: 4, rid: 1, date: '2025-01-01 00:00:00' },
  ]

  const node = new CommentNode(comments[1], {
    ...Defaults,
    flatMode: false,
    getApi: vi.fn(),
    replyComment: vi.fn(),
    editComment: vi.fn(),
    avatarURLBuilder: undefined,
    scrollRelativeTo: undefined,
    dateFormatter: (date: Date) => date.toISOString(),
  })
  expect(node.getDateFormatted()).toBe('2025-01-03T00:00:00.000Z')

  const ids = (sort: 'DATE_ASC' | 'DATE_DESC') =>
    makeNestCommentNodeList(comments, sort)[0].children.map((n) => n.id)
  expect(ids('DATE_ASC')).toEqual([4, 3, 2])
  expect(ids('DATE_DESC')).toEqual([2, 3, 4])
})

import { vi } from 'vitest'
import Context from '@/context'
import { Api } from '@/api'
import type { EntityCommentDTO, EntityPageDTO, HttpResponse } from '@/api/v2'
import { EventsService } from '@/services/events'
import { ConfigService } from '@/services/config'
import { DataService } from '@/services/data'
import { UserService } from '@/services/user'
import { CommentsService } from '@/services/comments'
import { Fetch } from '@/plugins/list/fetch'
import type { Context as IContext } from '@/types'

export const version = { app: 'Artalk', version: '2.10.1', commit_hash: 'test' }

export const page: EntityPageDTO = {
  date: '2026-01-01 00:00:00',
  pv: 0,
  id: 1,
  key: '/',
  title: 'Test',
  url: 'https://example.com/',
  site_name: 'Test',
  admin_only: false,
  vote_up: 0,
  vote_down: 0,
}

export function comment(id: number, content = `Comment ${id}`): EntityCommentDTO {
  return {
    id,
    content,
    content_marked: '',
    user_id: 1,
    badge_name: '',
    badge_color: '',
    page_url: 'https://example.com/',
    nick: 'Visitor',
    email_encrypted: '',
    link: '',
    rid: 0,
    ua: '',
    date: '2026-01-01 00:00:00',
    is_collapsed: false,
    is_pending: false,
    is_pinned: false,
    is_verified: false,
    is_allow_reply: true,
    page_key: '/',
    site_name: 'Test',
    visible: true,
    vote_up: 0,
    vote_down: 0,
  }
}

export function list(comments: EntityCommentDTO[] = []) {
  return { comments, count: comments.length, roots_count: comments.length, page }
}

export function response<T>(data: T): HttpResponse<T, any> {
  return { data } as HttpResponse<T, any>
}

/** Mock the existing API methods; keep their v2 parameters and response shape. */
export function mockApi() {
  const api = new Api({
    baseURL: 'https://example.com/api/v2',
    siteName: 'Test',
    pageKey: '/',
    pageTitle: 'Test',
  })
  vi.spyOn(api.conf, 'conf').mockResolvedValue(response({ frontend_conf: {}, version }))
  vi.spyOn(api.comments, 'getComments').mockResolvedValue(response(list()))
  vi.spyOn(api.comments, 'getComment').mockImplementation(async (id) =>
    response({ comment: comment(id), reply_comment: comment(0) }),
  )
  vi.spyOn(api.comments, 'createComment').mockImplementation(async (params) =>
    response({
      ...comment(params.rid ? 2 : 1, params.content),
      rid: params.rid || 0,
      page_key: params.page_key,
      site_name: params.site_name,
    }),
  )
  vi.spyOn(api.comments, 'updateComment').mockImplementation(async (id, params) =>
    response({ ...comment(id), ...params }),
  )
  vi.spyOn(api.comments, 'deleteComment').mockResolvedValue(response({}))
  vi.spyOn(api.notifies, 'getNotifies').mockResolvedValue(response({ notifies: [], count: 0 }))
  vi.spyOn(api.notifies, 'markNotifyRead').mockResolvedValue(response({}))
  vi.spyOn(api.user, 'getUser').mockResolvedValue(
    response({
      is_login: false,
      notifies: [],
      notifies_count: 0,
      user: {
        id: 1,
        name: 'Visitor',
        email: 'visitor@example.com',
        link: '',
        badge_color: '',
        badge_name: '',
        is_admin: false,
        receive_email: false,
      },
    }),
  )
  vi.spyOn(api.pages, 'logPv').mockResolvedValue(response({ pv: 1 }))
  vi.spyOn(api.pages, 'updatePage').mockImplementation(async (_id, params) =>
    response({ ...page, ...params }),
  )
  vi.spyOn(api.stats, 'getStats').mockResolvedValue(response({ data: {} }))
  vi.spyOn(api.votes, 'getVote').mockResolvedValue(
    response({ up: 0, down: 0, is_up: false, is_down: false }),
  )
  vi.spyOn(api.votes, 'createVote').mockResolvedValue(
    response({ up: 1, down: 0, is_up: true, is_down: false }),
  )
  vi.spyOn(api.upload, 'upload').mockResolvedValue(
    response({
      public_url: 'https://example.com/upload.png',
      file_name: 'image.png',
      file_type: 'image/png',
    }),
  )
  return api
}

export function bindApi(ctx: IContext, api: Api) {
  vi.spyOn(ctx, 'getApi').mockReturnValue(api)
  vi.spyOn(api, 'getUserFields').mockImplementation(() => {
    const { name, email } = ctx.getUser().getData()
    return name && email ? { name, email } : undefined
  })
}

export function setup(api = mockApi()) {
  const ctx = new Context(document.createElement('div'))
  EventsService(ctx)
  ConfigService(ctx)
  ctx.updateConf({ server: 'https://example.com', pageKey: '/', site: 'Test', flatMode: false })
  UserService(ctx)
  DataService(ctx)
  bindApi(ctx, api)
  CommentsService(ctx)
  Fetch(ctx)
  // eslint-disable-next-line artalk/artalk-plugin -- Test composition root.
  return { ctx, data: ctx.getData(), service: ctx.inject('comments'), api }
}

export function deferred<T>() {
  let resolve!: (value: T) => void
  let reject!: (reason: Error) => void
  const promise = new Promise<T>((done, fail) => {
    resolve = done
    reject = fail
  })
  return { promise, resolve, reject }
}

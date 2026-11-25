import { describe, expect, it, vi } from 'vitest'
import { AnnotationApi } from '../src/api'

function mockCtx() {
  const getComments = vi.fn().mockResolvedValue({ data: { comments: [] } })
  const createComment = vi.fn().mockResolvedValue({ data: { id: 1 } })
  const updateCommentAnchor = vi.fn().mockResolvedValue({ data: { id: 1 } })
  const deleteComment = vi.fn().mockResolvedValue({ data: {} })
  const ctx: any = {
    getConf: () => ({ pageKey: '/doc', site: 'Test' }),
    getUser: () => ({
      getData: () => ({ name: 'alice', email: 'a@example.com', link: '' }),
    }),
    getApi: () => ({
      comments: { getComments, createComment, updateCommentAnchor, deleteComment },
    }),
  }
  return { ctx, getComments, createComment, updateCommentAnchor, deleteComment }
}

describe('api', () => {
  it('lists annotation threads', async () => {
    const { ctx, getComments } = mockCtx()
    await new AnnotationApi(ctx).listThreads()
    expect(getComments).toHaveBeenCalledWith(
      expect.objectContaining({ has_anchor: 'true', page_key: '/doc', site_name: 'Test' }),
    )
  })

  it('paginates annotation threads', async () => {
    const { ctx, getComments } = mockCtx()
    getComments.mockResolvedValue({ data: { comments: [], roots_count: 3 } })
    const page = await new AnnotationApi(ctx).listThreads(10, 20)
    expect(getComments).toHaveBeenCalledWith(
      expect.objectContaining({ offset: 10, limit: 20, has_anchor: 'true' }),
    )
    expect(page.rootsCount).toBe(3)
  })

  it('creates an annotation with an anchor', async () => {
    const { ctx, createComment } = mockCtx()
    const anchor = { quote: { type: 'TextQuoteSelector', exact: 'Hello' } }
    await new AnnotationApi(ctx).createAnnotation('note', anchor)
    expect(createComment).toHaveBeenCalledWith(
      expect.objectContaining({
        name: 'alice',
        email: 'a@example.com',
        content: 'note',
        anchor,
      }),
    )
  })

  it('replies with rid and updates anchors', async () => {
    const { ctx, createComment, updateCommentAnchor } = mockCtx()
    const api = new AnnotationApi(ctx)
    await api.reply(7, 'reply')
    expect(createComment).toHaveBeenCalledWith(expect.objectContaining({ rid: 7 }))

    await api.updateAnchor(7, { resolved: true })
    expect(updateCommentAnchor).toHaveBeenCalledWith(7, { resolved: true })
  })
})

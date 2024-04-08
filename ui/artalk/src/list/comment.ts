import type { ContextApi, CommentData } from '@/types'
import { CommentNode } from '@/comment'

interface CreateCommentNodeOptions {
  forceFlatMode?: boolean
}

export function createCommentNode(
  ctx: ContextApi,
  comment: CommentData,
  replyComment?: CommentData,
  opts?: CreateCommentNodeOptions,
): CommentNode {
  const instance = new CommentNode(comment, {
    onAfterRender: () => {
      ctx.trigger('comment-rendered', instance)
    },
    onDelete: (c: CommentNode) => {
      ctx.getData().deleteComment(c.getID())
    },

    replyTo: replyComment,

    // TODO simplify reference
    flatMode:
      typeof opts?.forceFlatMode === 'boolean'
        ? opts?.forceFlatMode
        : (ctx.conf.flatMode as boolean),
    gravatar: ctx.conf.gravatar,
    nestMax: ctx.conf.nestMax,
    heightLimit: ctx.conf.heightLimit,
    avatarURLBuilder: ctx.conf.avatarURLBuilder,
    scrollRelativeTo: ctx.conf.scrollRelativeTo,
    vote: ctx.conf.vote,
    voteDown: ctx.conf.voteDown,
    uaBadge: ctx.conf.uaBadge,

    // TODO: move to plugin folder and remove from core
    getApi: () => ctx.getApi(),
    replyComment: (c, $el) => ctx.replyComment(c, $el),
    editComment: (c, $el) => ctx.editComment(c, $el),
  })

  // 渲染元素
  instance.render()

  return instance
}

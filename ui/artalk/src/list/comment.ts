import type { ConfigManager, EventManager, CommentData, DataManager } from '@/types'
import { CommentNode } from '@/comment'
import type { Api } from '@/api'

interface CreateCommentNodeOptions {
  getApi: () => Api
  getEvents: () => EventManager
  getData: () => DataManager
  getConf: () => ConfigManager
  replyComment: (c: CommentData, $el: HTMLElement) => void
  editComment: (c: CommentData, $el: HTMLElement) => void
  forceFlatMode?: boolean
}

export function createCommentNode(
  opts: CreateCommentNodeOptions,
  comment: CommentData,
  replyComment?: CommentData,
): CommentNode {
  const conf = opts.getConf().get()

  const instance = new CommentNode(comment, {
    onAfterRender: () => {
      opts.getEvents().trigger('comment-rendered', instance)
    },
    onDelete: (c: CommentNode) => {
      opts.getData().deleteComment(c.getID())
    },

    replyTo: replyComment,

    // TODO simplify reference
    flatMode:
      typeof opts?.forceFlatMode === 'boolean' ? opts?.forceFlatMode : (conf.flatMode as boolean),
    gravatar: conf.gravatar,
    nestMax: conf.nestMax,
    heightLimit: conf.heightLimit,
    avatarURLBuilder: conf.avatarURLBuilder,
    scrollRelativeTo: conf.scrollRelativeTo,
    vote: conf.vote,
    voteDown: conf.voteDown,
    uaBadge: conf.uaBadge,
    dateFormatter: conf.dateFormatter,

    // TODO: move to plugin folder and remove from core
    getApi: opts.getApi,
    replyComment: opts.replyComment,
    editComment: opts.editComment,
  })

  // Render comment
  instance.render()

  return instance
}

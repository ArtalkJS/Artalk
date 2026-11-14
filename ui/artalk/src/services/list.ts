import type { ArtalkPlugin } from '@/types'
import { List } from '@/list/list'

export const ListService: ArtalkPlugin = (ctx) => {
  ctx.provide(
    'list',
    (events, config, data, editor, markdown) => {
      const list = new List({
        getApi: () => ctx.getApi(),
        getEvents: () => events,
        getConf: () => config,
        renderMarkdown: markdown.render,
        getData: () => data,

        replyComment: (c, $el) => editor.setReplyComment(c, $el),
        editComment: (c, $el) => editor.setEditComment(c, $el),
        resetEditorState: () => editor.resetState(),
        onListGotoFirst: () => ctx.listGotoFirst(),
      })
      ctx.getEl().appendChild(list.getEl())
      return list
    },
    ['events', 'config', 'data', 'editor', 'markdown'] as const,
  )
}

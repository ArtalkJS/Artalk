import type { ArtalkPlugin } from '@/types'
import { List } from '@/list/list'

export const ListService: ArtalkPlugin = (ctx) => {
  ctx.provide(
    'list',
    (api, events, config, data, editor) => {
      const list = new List({
        getApi: () => api,
        getEvents: () => events,
        getConf: () => config,
        getData: () => data,

        replyComment: (c, $el) => editor.setReplyComment(c, $el),
        editComment: (c, $el) => editor.setEditComment(c, $el),
        resetEditorState: () => editor.resetState(),
        onListGotoFirst: () => ctx.listGotoFirst(),
      })
      ctx.getEl().appendChild(list.getEl())
      return list
    },
    ['api', 'events', 'config', 'data', 'editor'] as const,
  )
}

import type { ArtalkPlugin } from '@/types'
import { Editor } from '@/editor/editor'

export const EditorService: ArtalkPlugin = (ctx) => {
  ctx.provide(
    'editor',
    (events, config) => {
      const editor = new Editor({
        getEvents: () => events,
        getConf: () => config,
      })
      ctx.getEl().appendChild(editor.getEl())
      return editor
    },
    ['events', 'config'] as const,
  )
}

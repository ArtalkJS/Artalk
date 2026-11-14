import type { ArtalkPlugin } from '@/types'
import { createMarkdown } from '@/lib/marked'

export const MarkdownService: ArtalkPlugin = (ctx) => {
  ctx.provide(
    'markdown',
    (config) => {
      const markdown = createMarkdown(() => config.get().markedReplacers)
      config.watchConf(['imgLazyLoad', 'markedOptions'], (conf) => markdown.configure(conf))
      return markdown
    },
    ['config'] as const,
  )
}

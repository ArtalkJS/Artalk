import type { ArtalkPlugin } from '@/types'

export const Markdown: ArtalkPlugin = (ctx) => {
  ctx.inject('markdown')
}

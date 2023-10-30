import type { ArtalkPlugin } from '~/types'

export const DarkMode: ArtalkPlugin = (ctx) => {
  ctx.on('inited', () => {
    ctx.setDarkMode(ctx.conf.darkMode)
  })

  ctx.on('conf-loaded', () => {
    ctx.setDarkMode(ctx.conf.darkMode)
  })
}

import { syncDarkModeConf } from '@/lib/dark-mode'
import type { ArtalkPlugin } from '~/types'

export const DarkMode: ArtalkPlugin = (ctx) => {
  ctx.on('inited', () => {
    syncDarkModeConf(ctx)
  })

  ctx.on('conf-loaded', () => {
    syncDarkModeConf(ctx)
  })
}

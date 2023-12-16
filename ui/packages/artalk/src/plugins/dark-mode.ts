import type { ArtalkPlugin } from '~/types'
import { applyDarkMode } from '@/components/dark-mode'

export const DarkMode: ArtalkPlugin = (ctx) => {
  const sync = () => {
    applyDarkMode([ctx.$root, ctx.get('layerManager').getEl()], ctx.conf.darkMode)
  }

  ctx.on('inited', () => sync())
  ctx.on('conf-loaded', () => sync())
  ctx.on('dark-mode-changed', () => sync())
}

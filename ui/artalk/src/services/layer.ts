import type { ArtalkPlugin } from '@/types'
import { LayerManager } from '@/layer'

export const LayersService: ArtalkPlugin = (ctx) => {
  ctx.provide('layers', () => {
    const layerManager = new LayerManager()
    document.body.appendChild(layerManager.getEl())
    ctx.on('unmounted', () => {
      layerManager?.destroy()
    })
    return layerManager
  })
}

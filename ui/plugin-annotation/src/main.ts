import type { ArtalkPlugin } from 'artalk'
import { AnnotationController } from './controller'
import type { AnnotationPluginOptions } from './types'

export const ArtalkAnnotationPlugin: ArtalkPlugin<AnnotationPluginOptions> = (
  ctx,
  options = {},
) => {
  const controller = new AnnotationController(ctx, options)

  ctx.on('mounted', () => controller.mount())
  ctx.on('unmounted', () => controller.destroy())
  ctx.on('user-changed', () => controller.onUserChange())
  ctx.watchConf(['pageKey'], () => controller.onPageChange())
  ctx.watchConf(['darkMode'], () => controller.onThemeChange())
  ctx.watchConf(['locale'], () => controller.onLocaleChange())
}

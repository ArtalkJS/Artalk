import ArtalkPlugin from '~/types/plugin'
import { version as ARTALK_VERSION } from '~/package.json'

export const ListCopyright: ArtalkPlugin = (ctx) => {
  ctx.on('conf-loaded', () => {
    const list = ctx.get('list')
    if (!list) return

    const $copyright = list.$el.querySelector<HTMLElement>('.atk-copyright')
    if (!$copyright) return

    $copyright.innerHTML = (
      `Powered By <a href="https://artalk.js.org" ` +
      `target="_blank" title="Artalk v${ARTALK_VERSION}">` +
      `Artalk</a>`)
  })
}

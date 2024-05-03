import type { ContextApi } from 'artalk'
import { JSX } from 'solid-js'
import { render } from 'solid-js/web'

export const createLayer = (ctx: ContextApi) => {
  const layer = ctx.get('layerManager').create('login')
  const show = (el: (l: typeof layer) => JSX.Element) => {
    const $el = document.createElement('div')
    render(() => el(layer), $el)
    layer.getEl().append($el.firstChild!)

    layer.show()
  }

  return {
    show,
  }
}

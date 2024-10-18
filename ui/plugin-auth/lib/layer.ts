import { JSX } from 'solid-js'
import { render } from 'solid-js/web'
import type { Layer } from 'artalk'
import type { AuthContext } from '../types'

export const createLayer = (ctx: AuthContext) => {
  const layer = ctx.getLayers().create('login')
  const show = (el: (l: Layer) => JSX.Element) => {
    const $el = document.createElement('div')
    render(() => el(layer), $el)
    layer.getEl().append($el.firstChild!)

    layer.show()
  }

  return {
    show,
  }
}

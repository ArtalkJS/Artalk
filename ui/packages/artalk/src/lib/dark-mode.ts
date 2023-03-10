import type Context from '~/types/context'
import { getLayerWrap } from '../layer'

const darkModeMedia = window.matchMedia('(prefers-color-scheme: dark)')
let darkModeAutoFunc: (evt: MediaQueryListEvent) => void

export function syncDarkModeConf(ctx: Context) {
  setDarkMode(ctx, ctx.conf.darkMode, false) // ref `conf.darkMode` data, so no need to `alterConf`
}

export function setDarkMode(ctx: Context, darkMode: boolean|'auto', alterConf = true) {
  const apply = (d: boolean) => {
    updateClassName(ctx, d)
    if (alterConf) alterCtxConf(ctx, d)
  }

  if (darkMode === 'auto') {
    if (!darkModeAutoFunc) { // 自动切换夜间模式，事件监听
      darkModeAutoFunc = (evt) => apply(evt.matches)
      darkModeMedia.addEventListener('change', darkModeAutoFunc)
    }

    apply(darkModeMedia.matches)
  } else {
    // if the type of darkMode is boolean
    if (darkModeAutoFunc) // 解除事件监听绑定
      darkModeMedia.removeEventListener('change', darkModeAutoFunc)

    apply(darkMode)
  }
}

function alterCtxConf(ctx: Context, darkMode: boolean) {
  ctx.conf.darkMode = darkMode
}

const DarkModeClassName = 'atk-dark-mode'
export function updateClassName(ctx: Context, darkMode: boolean) {
  if (darkMode) ctx.$root.classList.add(DarkModeClassName)
  else ctx.$root.classList.remove(DarkModeClassName)

  // for Layer
  const { $wrap: $layerWrap } = getLayerWrap()
  if ($layerWrap) {
    if (darkMode) $layerWrap.classList.add(DarkModeClassName)
    else $layerWrap.classList.remove(DarkModeClassName)
  }
}

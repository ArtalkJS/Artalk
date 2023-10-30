import { getLayerWrap } from '../layer'

const darkModeMedia = window.matchMedia('(prefers-color-scheme: dark)')
let darkModeAutoFunc: (evt: MediaQueryListEvent) => void

export function setDarkMode($el: HTMLElement, darkMode: boolean|'auto') {
  const apply = (d: boolean) => {
    updateClassName($el, d)
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

const DarkModeClassName = 'atk-dark-mode'
function updateClassName($el: HTMLElement, darkMode: boolean) {
  if (darkMode) $el.classList.add(DarkModeClassName)
  else $el.classList.remove(DarkModeClassName)

  // for Layer
  const { $wrap: $layerWrap } = getLayerWrap()
  if ($layerWrap) {
    if (darkMode) $layerWrap.classList.add(DarkModeClassName)
    else $layerWrap.classList.remove(DarkModeClassName)
  }
}

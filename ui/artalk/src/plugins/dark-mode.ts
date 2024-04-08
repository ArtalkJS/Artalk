import type { ArtalkPlugin } from '@/types'

// Notice: Singleton pattern needs to be loaded as lazy as possible,
//         because the SSG application does not have a `window` context.
let darkModeMedia: MediaQueryList | undefined

function updateClassnames($els: HTMLElement[], darkMode: boolean) {
  const DarkModeClassName = 'atk-dark-mode'

  $els.forEach(($el) => {
    if (darkMode) $el.classList.add(DarkModeClassName)
    else $el.classList.remove(DarkModeClassName)
  })
}

export const DarkMode: ArtalkPlugin = (ctx) => {
  // the handler bind to Artalk instance, don't forget to remove it when Artalk instance destroyed
  let darkModeAutoHandler: ((evt: MediaQueryListEvent) => void) | undefined

  const sync = (darkMode: boolean | 'auto') => {
    // the elements that classnames need to be updated when darkMode changed
    const $els = [ctx.$root, ctx.get('layerManager').getEl()]

    // init darkModeMedia if not exists, and only create once
    if (!darkModeMedia) {
      darkModeMedia = window.matchMedia('(prefers-color-scheme: dark)')
    }

    if (darkMode === 'auto') {
      // if darkMode is 'auto', add handler
      if (!darkModeAutoHandler) {
        // the handler that will be called when darkModeMedia changed
        darkModeAutoHandler = (evt) => updateClassnames($els, evt.matches)
        darkModeMedia.addEventListener('change', darkModeAutoHandler)
      }

      // update classnames immediately
      updateClassnames($els, darkModeMedia.matches)
    } else {
      // if darkMode is boolean, remove handler
      if (darkModeAutoHandler) {
        darkModeMedia.removeEventListener('change', darkModeAutoHandler)
        darkModeAutoHandler = undefined
      }

      // update classnames immediately
      updateClassnames($els, darkMode)
    }
  }

  ctx.watchConf(['darkMode'], (conf) => sync(conf.darkMode))
  ctx.on('created', () => sync(ctx.conf.darkMode))
  ctx.on('unmounted', () => {
    // if handler exists, don't forget to remove it, or it will cause memory leak
    darkModeAutoHandler && darkModeMedia?.removeEventListener('change', darkModeAutoHandler)
    darkModeAutoHandler = undefined
  })
}

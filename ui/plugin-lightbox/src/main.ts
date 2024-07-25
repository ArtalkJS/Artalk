import { ArtalkPlugin } from 'artalk'

import { getDeps } from './helper'

const LOADED_ATTR = 'atk-lightbox-loaded'
const IMG_LINK_EL_CLASS = 'atk-lightbox-img'
const IMG_LINK_EL_SEL = `.${IMG_LINK_EL_CLASS}`

export interface ArtalkLightboxPluginOptions {
  lightGallery?: {
    lib?: () => Promise<any>
  }
  lightBox?: {
    lib?: () => Promise<any>
  }
  fancyBox?: {
    lib?: () => Promise<any>
  }
  photoSwipe?: {
    lib?: () => Promise<any>
    pswpModule?: () => Promise<any>
  }

  /** Config for all lightbox plugins */
  config?: any
}

export const ArtalkLightboxPlugin: ArtalkPlugin<ArtalkLightboxPluginOptions> = (ctx, opts) => {
  // add artalk event listener
  ctx.on('list-loaded', async () => {
    const deps = await getDeps(opts)

    const $imgLinks: HTMLAnchorElement[] = []
    const $contents = new Set<HTMLElement>()

    ctx.getCommentNodes().forEach((comment) => {
      const $content = comment.getRender().$content
      $content
        .querySelectorAll<HTMLImageElement>(`img:not([atk-emoticon]):not([${LOADED_ATTR}])`)
        .forEach(($img) => {
          $img.setAttribute(LOADED_ATTR, '') // mark as loaded

          const $imgLink = document.createElement('a')
          $imgLink.setAttribute('class', IMG_LINK_EL_CLASS)
          $imgLink.setAttribute('href', $img.src)
          $imgLink.setAttribute('data-src', $img.src)

          // set image size for PhotoSwipe
          if (deps.photoSwipe) {
            // @see https://photoswipe.com/getting-started/#required-html-markup
            const updateImageSize = () => {
              $imgLink.setAttribute('data-pswp-width', $img.naturalWidth.toString())
              $imgLink.setAttribute('data-pswp-height', $img.naturalHeight.toString())
              $img.removeEventListener('load', updateImageSize)
            }
            updateImageSize()
            $img.addEventListener('load', updateImageSize)
          }

          $imgLink.append($img.cloneNode())

          $img.replaceWith($imgLink)
          $imgLinks.push($imgLink)
        })

      $contents.add($content)
    })

    // lightgallery
    if (deps.lightGallery) {
      // lightGallery
      // @see https://github.com/sachinchoolur/lightGallery
      $contents.forEach((el) => {
        deps.lightGallery!(el, {
          selector: IMG_LINK_EL_SEL,
          ...(opts?.config || window.ATK_LIGHTBOX_CONF || {}),
        })
      })
    }

    if (deps.lightBox) {
      // lightbox2
      // @see https://github.com/lokesh/lightbox2
      $imgLinks.forEach((el) => {
        el.setAttribute('data-title', el.querySelector('img')!.alt)
        el.setAttribute('rel', 'lightbox')
        el.onclick = (evt) => {
          evt.preventDefault()
          deps.lightBox!.start(window.$(el))
        }
      })
    }

    if (deps.photoSwipe) {
      // PhotoSwipe
      // @see https://github.com/dimsemenov/photoswipe
      const lightbox = new deps.photoSwipe({
        gallery: `.atk-content`,
        showHideAnimationType: 'fade',
        thumbSelector: `${IMG_LINK_EL_SEL}`,
        children: `${IMG_LINK_EL_SEL}`,
        pswpModule: opts?.photoSwipe!.pswpModule,
        ...(opts?.config || window.ATK_LIGHTBOX_CONF || {}),
      })
      lightbox.init()
    }

    if (deps.fancyBox) {
      // Fancybox
      // @see https://github.com/fancyapps/fancybox
      if (!window.$)
        throw new Error('Fancybox requires jQuery which is available in window global scope')
      if (!window.$.fancybox) deps.fancyBox(window.$)
      window
        .$(`.artalk .atk-list ${IMG_LINK_EL_SEL}`, opts?.config || window.ATK_LIGHTBOX_CONF)
        .fancybox()
    }
  })
}

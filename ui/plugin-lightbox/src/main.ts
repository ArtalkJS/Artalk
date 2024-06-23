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
  }
}

export const ArtalkLightboxPlugin: ArtalkPlugin<ArtalkLightboxPluginOptions> = (ctx, opts) => {
  // add artalk event listener
  ctx.on('list-loaded', async () => {
    const $imgLinks: HTMLAnchorElement[] = []
    const $contents = new Set<HTMLElement>()

    ctx.getCommentNodes().forEach((comment) => {
      const $content = comment.getRender().$content
      $content
        .querySelectorAll<HTMLImageElement>(`img:not([atk-emoticon]):not([${LOADED_ATTR}])`)
        .forEach(($img) => {
          $img.setAttribute(LOADED_ATTR, '') // 初始化标记

          const $imgLink = document.createElement('a')
          $imgLink.setAttribute('class', IMG_LINK_EL_CLASS)
          $imgLink.setAttribute('href', $img.src)
          $imgLink.setAttribute('data-src', $img.src)
          $imgLink.append($img.cloneNode())

          $img.replaceWith($imgLink)
          $imgLinks.push($imgLink)
        })

      $contents.add($content)
    })

    const deps = await getDeps(opts)

    // lightgallery
    if (deps.lightGallery) {
      // lightGallery
      // @link https://github.com/sachinchoolur/lightGallery
      $contents.forEach((el) => {
        deps.lightGallery!(el, {
          selector: IMG_LINK_EL_SEL,
          ...(window.ATK_LIGHTBOX_CONF || {}),
        })
      })
    }

    if (deps.lightBox) {
      // lightbox2
      // @link https://github.com/lokesh/lightbox2
      $imgLinks.forEach((el) => {
        window.$(el).attr('data-title', window.$(el).find('img').attr('alt'))
        window.$(el).attr('rel', 'lightbox')
        el.onclick = (evt) => {
          evt.preventDefault()
          deps.lightBox!.start(window.$(el))
        }
      })
    }

    if (deps.photoSwipe) {
      // PhotoSwipe
      // @link https://github.com/dimsemenov/photoswipe
      const dataSource: any[] = []
      const lightbox = new deps.photoSwipe({
        dataSource,
        ...(window.ATK_LIGHTBOX_CONF || {}),
      })

      $imgLinks.forEach((el) => {
        const $img = el.querySelector<HTMLImageElement>('img')!
        dataSource.push({ src: $img.src })

        el.onclick = (evt) => {
          evt.preventDefault()
          // @ts-expect-error missing type but works
          lightbox.loadAndOpen(0)
        }
      })

      lightbox.init()
    }

    if (deps.fancyBox) {
      // Fancybox Event bind
      // @link https://github.com/fancyapps/fancybox
      deps.fancyBox.then(({ FancyBox }) => {
        if (!FancyBox) return
        FancyBox.bind(`.artalk .atk-list ${IMG_LINK_EL_SEL}`, window.ATK_LIGHTBOX_CONF)
      })
    }
  })
}

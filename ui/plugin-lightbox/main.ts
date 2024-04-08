import Artalk from 'artalk'

Artalk.use((ctx) => {
  const LOADED_ATTR = 'atk-lightbox-loaded'
  const IMG_LINK_EL_CLASS = 'atk-lightbox-img'
  const IMG_LINK_EL_SEL = `.${IMG_LINK_EL_CLASS}`

  const typeSet = window.ATK_LIGHTBOX_TYPE
  const typeIs = (t: string) =>
    (!typeSet && window[t]) || (typeSet || '').toLowerCase() === t.toLowerCase()

  // add artalk event listener
  ctx.on('list-loaded', () => {
    const $imgLinks: HTMLAnchorElement[] = []

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

      // lightGallery
      // @link https://github.com/sachinchoolur/lightGallery
      if (typeIs('lightGallery')) {
        window.lightGallery($content, {
          selector: IMG_LINK_EL_SEL,
          ...(window.ATK_LIGHTBOX_CONF || {}),
        })
      }
    })

    // lightbox2
    // @link https://github.com/lokesh/lightbox2
    if (typeIs('lightbox')) {
      $imgLinks.forEach((el) => {
        window.$(el).attr('data-title', window.$(el).find('img').attr('alt'))
        window.$(el).attr('rel', 'lightbox')
        el.onclick = (evt) => {
          evt.preventDefault()
          window.lightbox.start(window.$(el))
        }
      })
    }

    // PhotoSwipe
    // @link https://github.com/dimsemenov/photoswipe
    /*
    if (typeIs('PhotoSwipeLightbox') || typeIs('PhotoSwipe')) {
      const dataSource: any[] = []
      const lightbox = new window.PhotoSwipeLightbox({
        dataSource, ...(window.ATK_LIGHTBOX_CONF || {})
      })

      $imgLinks.forEach((el) => {
        const $img = el.querySelector<HTMLImageElement>('img')!
        dataSource.push({ src: $img.src })

        el.onclick = (evt) => {
          evt.preventDefault()
          lightbox.loadAndOpen(0);
        }
      })

      lightbox.init()
    }
    */
  })

  // Fancybox Event bind
  // @link https://github.com/fancyapps/fancybox
  if (typeIs('Fancybox')) {
    window.Fancybox.bind(`.artalk .atk-list ${IMG_LINK_EL_SEL}`, window.ATK_LIGHTBOX_CONF)
  }
})

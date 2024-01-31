import type { ArtalkPlugin } from '@/types'

export const ReachBottom: ArtalkPlugin = (ctx) => {
  let observer: IntersectionObserver|null = null

  const setupObserver = ($target: HTMLElement) => {
    const scrollEvtAt = (ctx.conf.scrollRelativeTo && ctx.conf.scrollRelativeTo()) || null

    // eslint-disable-next-line compat/compat
    observer = new IntersectionObserver(([entries]) => {
      if (entries.isIntersecting) {
        clearObserver() // clear before trigger event to avoid trigger twice `list-reach-bottom`

        // note that this event will be triggered only once
        // until the next list is loaded
        ctx.trigger('list-reach-bottom')
      }
    }, {
      threshold: 0.9, // when the target is 90% visible

      // @see https://developer.mozilla.org/en-US/docs/Web/API/IntersectionObserver/root
      // If the root is null, then the bounds of the actual document viewport are used.
      root: scrollEvtAt
    })
    observer.observe($target)
  }

  const clearObserver = () => {
    observer?.disconnect()
    observer = null
  }

  ctx.on('list-loaded', () => {
    clearObserver()

    const list = ctx.get('list')

    // get the second last child
    const children = list.getCommentsWrapEl().childNodes
    const $target = children.length > 2 ? children[children.length - 2] as HTMLElement : null
    if (!$target) return

    // check IntersectionObserver support
    if (!('IntersectionObserver' in window)) {
      console.warn('IntersectionObserver api not supported')
      return
    }

    // use IntersectionObserver to detect reach bottom
    setupObserver($target)
  })

  ctx.on('unmounted', () => {
    clearObserver()
  })
}

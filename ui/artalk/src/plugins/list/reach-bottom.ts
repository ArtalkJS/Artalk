import type { ArtalkPlugin } from '@/types'

export const ReachBottom: ArtalkPlugin = (ctx) => {
  const scrollEvtAt = document // TODO: support ref ctx.conf
  let observer: IntersectionObserver|null = null

  const setupObserver = ($target: HTMLElement) => {
    // eslint-disable-next-line compat/compat
    observer = new IntersectionObserver(([entries]) => {
      if (entries.isIntersecting) {
        clearObserver()
        ctx.trigger('list-reach-bottom')
      }
    }, {
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

    // use IntersectionObserver to detect reach bottom
    const $target = list.getCommentsWrapEl().querySelector<HTMLElement>('.atk-comment-wrap:nth-last-child(3)')
    if (!$target) return

    // check IntersectionObserver support
    if (!('IntersectionObserver' in window)) {
      console.warn('IntersectionObserver api not supported')
      return
    }

    setupObserver($target)
  })

  ctx.on('destroy', () => {
    clearObserver()
  })
}

import type { ArtalkPlugin } from '~/types'

export const ReachBottom: ArtalkPlugin = (ctx) => {
  const scrollEvtAt = document.documentElement // TODO support ref ctx.conf
  let observer: IntersectionObserver|null = null

  ctx.on('inited', () => {
    const list = ctx.get('list')
    if (!list) return

    // use IntersectionObserver to detect reach bottom
    const $target = list.$el.querySelector<HTMLElement>('.atk-list-comments-wrap')

    // check IntersectionObserver support
    if (!('IntersectionObserver' in window)) {
      console.warn('IntersectionObserver api not supported')
      return
    }

    // eslint-disable-next-line compat/compat
    observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (entry.intersectionRatio > 0) {
          ctx.trigger('list-reach-bottom')
        }
      })
    }, {
      root: scrollEvtAt,
      threshold: 0,
    })
    observer.observe($target!)
  })

  ctx.on('destroy', () => {
    observer?.disconnect()
  })
}

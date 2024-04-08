import * as Ui from '@/lib/ui'

let bodyOrgOverflow: string
let bodyOrgPaddingRight: string

export function getScrollbarHelper() {
  return {
    init() {
      bodyOrgOverflow = document.body.style.overflow
      bodyOrgPaddingRight = document.body.style.paddingRight
    },

    unlock() {
      document.body.style.overflow = bodyOrgOverflow
      document.body.style.paddingRight = bodyOrgPaddingRight
    },

    lock() {
      document.body.style.overflow = 'hidden'
      const barPaddingRight = parseInt(
        window.getComputedStyle(document.body, null).getPropertyValue('padding-right'),
        10,
      )
      document.body.style.paddingRight = `${Ui.getScrollBarWidth() + barPaddingRight || 0}px`
    },
  }
}

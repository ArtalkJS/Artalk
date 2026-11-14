import * as Ui from '@/lib/ui'

const owners = new Set<object>()
let originalOverflow = ''
let originalPadding = ''

export function getScrollbarHelper() {
  const owner = {}
  return {
    // Kept for callers of the old helper. Styles are captured on the first lock.
    init() {},
    unlock() {
      if (!owners.delete(owner) || owners.size) return
      document.body.style.overflow = originalOverflow
      document.body.style.paddingRight = originalPadding
    },
    lock() {
      if (owners.has(owner)) return
      if (!owners.size) {
        originalOverflow = document.body.style.overflow
        originalPadding = document.body.style.paddingRight
        const padding = parseInt(window.getComputedStyle(document.body).paddingRight, 10) || 0
        document.body.style.paddingRight = `${Ui.getScrollBarWidth() + padding}px`
        document.body.style.overflow = 'hidden'
      }
      owners.add(owner)
    },
  }
}

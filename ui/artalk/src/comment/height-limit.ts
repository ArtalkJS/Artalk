import * as Utils from '../lib/utils'
import $t from '../i18n'

export interface IHeightLimitConf {
  /** After expand btn click */
  afterExpandBtnClick?: () => void
  /** Allow Scroll */
  scrollable?: boolean
}

export interface IHeightLimitRule {
  /** Target element need to check */
  el: HTMLElement | null | undefined

  /** Max height (unit: px) */
  max: number

  /** Whether or not the element contains `<img />` */
  imgCheck?: boolean
}

export type THeightLimitRuleSet = IHeightLimitRule[]

/** Check all elements below the max height limit */
export function check(conf: IHeightLimitConf, rules: THeightLimitRuleSet) {
  rules.forEach(({ el, max, imgCheck }) => {
    if (!el) return

    // set max height for avoiding img exceed the limit while loading
    if (imgCheck) el.style.maxHeight = `${max + 1}px` // allow 1px more for next detecting

    let lock = false
    const _check = () => {
      if (lock) return
      if (Utils.getHeight(el) <= max) return // if not exceed the limit, do nothing

      const afterExpandBtnClick = () => {
        lock = true // add lock to prevent collapse again after expand when image lazy loaded
        conf.afterExpandBtnClick?.()
      }

      !conf.scrollable
        ? applyHeightLimit({ el, max, afterExpandBtnClick })
        : applyScrollableHeightLimit({ el, max })
    }

    // check immediately
    _check()

    // check images after loaded
    if (imgCheck) {
      // check again when image loaded
      const imgs = el.querySelectorAll<HTMLImageElement>('.atk-content img')
      if (imgs.length === 0) el.style.maxHeight = ''
      imgs.forEach((img) => {
        img.onload = () => _check()
      })
    }
  })
}

/** Height limit CSS class name */
const HEIGHT_LIMIT_CSS = 'atk-height-limit'

/** Apply height limit on an element and add expand btn */
export function applyHeightLimit(obj: {
  el: HTMLElement
  max: number
  afterExpandBtnClick?: (e: MouseEvent) => void
}) {
  if (!obj.el) return
  if (!obj.max) return
  if (obj.el.classList.contains(HEIGHT_LIMIT_CSS)) return

  obj.el.classList.add(HEIGHT_LIMIT_CSS)
  obj.el.style.height = `${obj.max}px`
  obj.el.style.overflow = 'hidden'

  /* Expand button */
  const $expandBtn = Utils.createElement(
    `<div class="atk-height-limit-btn">${$t('readMore')}</span>`,
  )
  $expandBtn.onclick = (e) => {
    e.stopPropagation()
    disposeHeightLimit(obj.el)

    if (obj.afterExpandBtnClick) obj.afterExpandBtnClick(e)
  }
  obj.el.append($expandBtn)
}

/** Dispose height limit on an element and remove expand btn */
export function disposeHeightLimit($el: HTMLElement) {
  if (!$el) return
  if (!$el.classList.contains(HEIGHT_LIMIT_CSS)) return

  $el.classList.remove(HEIGHT_LIMIT_CSS)
  Array.from($el.children).forEach((e) => {
    if (e.classList.contains('atk-height-limit-btn')) e.remove()
  })
  $el.style.height = ''
  $el.style.maxHeight = ''
  $el.style.overflow = ''
}

/** Height limit scrollable CSS class name */
const HEIGHT_LIMIT_SCROLL_CSS = 'atk-height-limit-scroll'

/** Apply scrollable height limit */
export function applyScrollableHeightLimit(opt: { el: HTMLElement; max: number }) {
  if (!opt.el) return
  if (opt.el.classList.contains(HEIGHT_LIMIT_SCROLL_CSS)) return
  opt.el.classList.add(HEIGHT_LIMIT_SCROLL_CSS)
  opt.el.style.height = `${opt.max}px`
}

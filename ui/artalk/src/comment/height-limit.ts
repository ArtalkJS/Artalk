import * as Utils from '../lib/utils'
import $t from '../i18n'

export interface IHeightLimitConf {
  /** Post expand btn click */
  postExpandBtnClick?: (e: MouseEvent) => void
  /** Allow Scroll */
  scrollable?: boolean
}

export interface IHeightLimitRule {
  /** Target element need to check */
  el: HTMLElement|null|undefined

  /** Max height (unit: px) */
  max: number

  /** Whether or not the element contains `<img />` */
  imgContains?: boolean
}

export type THeightLimitRuleSet = IHeightLimitRule[]

/** Check all elements below the max height limit */
export function check(conf: IHeightLimitConf, rules: THeightLimitRuleSet) {
  rules.forEach(({
    el, max: maxHeight, imgContains
  }) => {
    const _apply = () => {
      if (!el) return
      if (!conf.scrollable)
        applyHeightLimit({ el, maxHeight, postBtnClick: conf.postExpandBtnClick })
      else
        applyScrollableHeightLimit({ el, maxHeight })
    }

    // checking
    const _check = () => {
      if (el && Utils.getHeight(el) > maxHeight) _apply() // 是否超过高度
    }
    _check() // check now
    if (imgContains && el) // check again if img contains
      Utils.onImagesLoaded(el, () => _check())
  })
}

/** Height limit CSS class name */
const HEIGHT_LIMIT_CSS = 'atk-height-limit'

/** Apply height limit on an element and add expand btn */
export function applyHeightLimit(obj: {
  el: HTMLElement,
  maxHeight: number,
  postBtnClick?: (e: MouseEvent) => void
}) {
  if (!obj.el) return
  if (!obj.maxHeight) return
  if (obj.el.classList.contains(HEIGHT_LIMIT_CSS)) return

  obj.el.classList.add(HEIGHT_LIMIT_CSS)
  obj.el.style.height = `${obj.maxHeight}px`
  obj.el.style.overflow = 'hidden'

  /* Expand button */
  const $expandBtn = Utils.createElement(`<div class="atk-height-limit-btn">${$t('readMore')}</span>`)
  $expandBtn.onclick = (e) => {
    e.stopPropagation()
    disposeHeightLimit(obj.el)

    if (obj.postBtnClick) obj.postBtnClick(e)
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
  $el.style.overflow = ''
}

/** Height limit scrollable CSS class name */
const HEIGHT_LIMIT_SCROLL_CSS = 'atk-height-limit-scroll'

/** Apply scrollable height limit */
export function applyScrollableHeightLimit(obj: {
  el: HTMLElement,
  maxHeight: number
}) {
  if (!obj.el) return
  if (obj.el.classList.contains(HEIGHT_LIMIT_SCROLL_CSS)) return
  obj.el.classList.add(HEIGHT_LIMIT_SCROLL_CSS)
  obj.el.style.height = `${obj.maxHeight}px`
}

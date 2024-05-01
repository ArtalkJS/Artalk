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

    // set max height to limit the height
    if (imgCheck) el.style.maxHeight = `${max + 1}px` // allow 2px more for next detecting

    const _apply = () => {
      const postBtnClick = conf.postExpandBtnClick
      !conf.scrollable
        ? applyHeightLimit({ el, max, postBtnClick })
        : applyScrollableHeightLimit({ el, max })
    }

    // checking
    const _check = () => {
      if (Utils.getHeight(el) > max) _apply() // 是否超过高度
    }

    _check() // check immediately

    // image check
    if (imgCheck) {
      // check again when image loaded
      el.querySelectorAll<HTMLImageElement>('.atk-content img').forEach((img) => {
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
  postBtnClick?: (e: MouseEvent) => void
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

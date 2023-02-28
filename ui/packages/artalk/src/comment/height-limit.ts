import * as Utils from '../lib/utils'
import $t from '../i18n'

export interface IHeightLimitConf {
  /** Post expand btn click */
  postExpandBtnClick?: (e: MouseEvent) => void
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
  rules.forEach(rule => {
    const _check = () => {
      // 是否超过高度
      if (!rule.el) return
      if (Utils.getHeight(rule.el) > rule.max)
        applyHeightLimit({
          el: rule.el, maxHeight: rule.max,
          postBtnClick: conf.postExpandBtnClick
        })
    }

    _check() // check now
    if (rule.imgContains && rule.el) // check again if img contains
      Utils.onImagesLoaded(rule.el, () => _check())
  })
}

/** Apply height limit on an element and add expand btn */
export function applyHeightLimit(obj: {
  el: HTMLElement,
  maxHeight: number,
  postBtnClick?: (e: MouseEvent) => void
}) {
  if (!obj.el) return
  if (obj.el.classList.contains('atk-height-limit')) return

  obj.el.classList.add('atk-height-limit')
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
  if (!$el.classList.contains('atk-height-limit')) return

  $el.classList.remove('atk-height-limit')
  Array.from($el.children).forEach((e) => {
    if (e.classList.contains('atk-height-limit-btn')) e.remove()
  })
  $el.style.height = ''
  $el.style.overflow = ''
}

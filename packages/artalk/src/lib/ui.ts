import * as Utils from './utils'

/** 显示加载 */
export function showLoading(parentElem: HTMLElement, conf?: { transparentBg?: boolean }) {
  // Use :scope https://developer.mozilla.org/en-US/docs/Web/CSS/:scope
  let $loading = parentElem.querySelector<HTMLElement>(':scope > .atk-loading')
  if (!$loading) {
    $loading = Utils.createElement(
    `<div class="atk-loading atk-fade-in" style="display: none;">
      <div class="atk-loading-spinner">
        <svg viewBox="25 25 50 50"><circle cx="50" cy="50" r="20" fill="none" stroke-width="2" stroke-miterlimit="10"></circle></svg>
      </div>
    </div>`)
    if (conf?.transparentBg) $loading.style.background = 'transparent'
    parentElem.appendChild($loading)
  }
  $loading.style.display = ''

  // spinner 延迟显示，若加载等待时间太短，没必要显示（闪一下即可）
  const $spinner = $loading.querySelector<HTMLElement>('.atk-loading-spinner')
  if ($spinner) {
    $spinner.style.display = 'none'
    window.setTimeout(() => {
      $spinner.style.display = ''
    }, 500)
  }
}

/** 隐藏加载 */
export function hideLoading(parentElem: HTMLElement) {
  const $loading = parentElem.querySelector<HTMLElement>(':scope > .atk-loading')
  if ($loading) $loading.style.display = 'none'
}

/** 元素加载动画 */
export function setLoading(val: boolean, parentElem: HTMLElement) {
  if (val) showLoading(parentElem)
  else hideLoading(parentElem)
}

/** 元素是否用户可见 */
export function isVisible(el: HTMLElement, viewport: HTMLElement = document.documentElement) {
  const viewportHeight = viewport.clientHeight

  const docViewTop = viewport.scrollTop
  const docViewBottom = docViewTop + viewportHeight

  const rect = el.getBoundingClientRect()
  const elemTop = rect.top + docViewTop
  const elemBottom = elemTop + rect.height

  return (elemBottom <= docViewBottom) //&& (elemTop >= docViewTop) 注释因为假如 el 比 viewport 还高就会失效
}

/** 滚动到元素中心 */
export function scrollIntoView(elem: HTMLElement, enableAnim: boolean = true) {
  const top =
    Utils.getOffset(elem).top +
    Utils.getHeight(elem) / 2 -
    document.documentElement.clientHeight / 2

  if (enableAnim) {
    window.scroll({
      top: top > 0 ? top : 0,
      left: 0,
      // behavior: 'smooth',
    })
  } else {
    // 无动画
    window.scroll(0, top > 0 ? top : 0)
  }
}

/** 显示消息 */
export function showNotify(
  wrapElem: HTMLElement,
  msg: string,
  type: 's' | 'e' | 'w' | 'i',
) {
  const colors = { s: '#57d59f', e: '#ff6f6c', w: '#ffc721', i: '#2ebcfc' }
  const timeout = 3000 // 持续显示时间 ms

  const notifyElem = Utils.createElement(
    `<div class="atk-notify atk-fade-in" style="background-color: ${colors[type]}"><span class="atk-notify-content"></span></div>`
  )
  const notifyContentEl = notifyElem.querySelector<HTMLElement>('.atk-notify-content')!
  notifyContentEl.innerHTML = Utils.htmlEncode(msg).replace('\n', '<br/>')

  wrapElem.appendChild(notifyElem)

  const notifyRemove = () => {
    notifyElem.classList.add('atk-fade-out')
    setTimeout(() => {
      notifyElem.remove()
    }, 200)
  }

  let timeoutFn: number
  if (timeout > 0) {
    timeoutFn = window.setTimeout(() => {
      notifyRemove()
    }, timeout)
  }

  notifyElem.addEventListener('click', () => {
    notifyRemove()
    window.clearTimeout(timeoutFn)
  })
}

/** fade 动画 */
export function playFadeAnim(
  elem: HTMLElement,
  after?: () => void,
  type: 'in' | 'out' = 'in'
) {
  elem.classList.add(`atk-fade-${type}`)
  // 动画结束清除 class
  const onAnimEnded = () => {
    elem.classList.remove(`atk-fade-${type}`)
    elem.removeEventListener('animationend', onAnimEnded)
    if (after) after()
  }
  elem.addEventListener('animationend', onAnimEnded)
}

/** 渐入动画 */
export function playFadeInAnim(elem: HTMLElement, after?: () => void) {
  playFadeAnim(elem, after, 'in')
}

/** 渐出动画 */
export function playFadeOutAnim(elem: HTMLElement, after?: () => void) {
  playFadeAnim(elem, after, 'out')
}

/** 设置全局错误 */
export function setError(parentElem: HTMLElement, html: string | HTMLElement | null, title: string = '<span class="atk-error-title">Artalk Error</span>') {
  let elem = parentElem.querySelector<HTMLElement>('.atk-error-layer')
  if (html === null) {
    if (elem !== null) elem.remove()
    return
  }
  if (!elem) {
    elem = Utils.createElement(
      `<div class="atk-error-layer">${title}<span class="atk-error-text"></span></div>`
    )
    parentElem.appendChild(elem)
  }

  const errorTextEl = elem.querySelector<HTMLElement>('.atk-error-text')!
  errorTextEl.innerHTML = ''
  if (html === null) return

  if (html instanceof HTMLElement) {
    errorTextEl.appendChild(html)
  } else {
    errorTextEl.innerText = html
  }
}

export function getScrollBarWidth() {
  const inner = document.createElement('p')
  inner.style.width = '100%'
  inner.style.height = '200px'

  const outer = document.createElement('div')
  outer.style.position = 'absolute'
  outer.style.top = '0px'
  outer.style.left = '0px'
  outer.style.visibility = 'hidden'
  outer.style.width = '200px'
  outer.style.height = '150px'
  outer.style.overflow = 'hidden'
  outer.appendChild(inner)

  document.body.appendChild(outer)
  const w1 = inner.offsetWidth
  outer.style.overflow = 'scroll'
  let w2 = inner.offsetWidth
  if (w1 === w2) w2 = outer.clientWidth

  document.body.removeChild(outer)

  return (w1 - w2)
}

import * as Utils from './utils'
import Context from '../Context'

/** 显示加载 */
export function showLoading(parentElem: HTMLElement|Context) {
  if (parentElem instanceof Context) parentElem = parentElem.rootEl

  let loadingEl = parentElem.querySelector<HTMLElement>('.atk-loading')
  if (!loadingEl) {
    loadingEl = Utils.createElement(`
        <div class="atk-loading" style="display: none;">
          <div class="atk-loading-spinner">
          <svg viewBox="25 25 50 50"><circle cx="50" cy="50" r="20" fill="none" stroke-width="2" stroke-miterlimit="10"></circle></svg>
          </div>
          </div>`)
    parentElem.appendChild(loadingEl)
  }
  loadingEl.style.display = ''
}

/** 隐藏加载 */
export function hideLoading(parentElem: HTMLElement|Context) {
  if (parentElem instanceof Context) parentElem = parentElem.rootEl

  const loadingEl = parentElem.querySelector<HTMLElement>('.atk-loading')
  if (loadingEl) loadingEl.style.display = 'none'
}

/** 元素是否用户可见 */
export function isVisible(elem: HTMLElement) {
  const docViewTop = window.scrollY
  const docViewBottom = docViewTop + document.documentElement.clientHeight

  const elemTop = Utils.getOffset(elem).top
  const elemBottom = elemTop + elem.offsetHeight

  return elemBottom <= docViewBottom && elemTop >= docViewTop
}

/** 滚动到元素中心 */
export function scrollIntoView(elem: HTMLElement, enableAnim: boolean = true) {
  if (isVisible(elem)) return
  const top =
    Utils.getOffset(elem).top +
    Utils.getHeight(elem) / 2 -
    document.documentElement.clientHeight / 2

  if (enableAnim) {
    window.scroll({
      top: top > 0 ? top : 0,
      left: 0,
      behavior: 'smooth',
    })
  } else {
    // 无动画
    window.scroll(0, top > 0 ? top : 0)
  }
}

/** 显示对话框 */
export function buildDialog (
  html: HTMLElement,
  onConfirm?: (btnElem: HTMLElement) => boolean | void,
  onCancel?: () => boolean | void,
): HTMLElement {
  const dialogElem = Utils.createElement(
    `<div class="atk-layer-dialog-wrap">
      <div class="atk-layer-dialog">
      <div class="atk-layer-dialog-content"></div>
      <div class="atk-layer-dialog-action">
      </div>`
  )

  // 按钮
  const actionElem = dialogElem.querySelector<HTMLElement>('.atk-layer-dialog-action')!
  const onclick =
    (f: (btnElem: HTMLElement) => boolean | void) =>
    (evt: Event) => {
      const returnVal = f(evt.currentTarget as HTMLElement)
      if (returnVal === undefined || returnVal === true) {
        dialogElem.remove()
      }
    }

  if (typeof onConfirm === 'function') {
    const btn = Utils.createElement<HTMLButtonElement>(
      '<button data-action="confirm">确定</button>'
    )
    btn.onclick = onclick(onConfirm)
    actionElem.appendChild(btn)
  }

  if (typeof onCancel === 'function') {
    const btn = Utils.createElement<HTMLButtonElement>(
      '<button data-action="cancel">取消</button>'
    )
    btn.onclick = onclick(onCancel)
    actionElem.appendChild(btn)
  }

  // 内容
  dialogElem.querySelector('.atk-layer-dialog-content')!.appendChild(html)

  return dialogElem
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
export function setError(parentElem: HTMLElement|Context, html: string | HTMLElement | null) {
  if (parentElem instanceof Context) parentElem = parentElem.rootEl

  let elem = parentElem.querySelector<HTMLElement>('.atk-error-layer')
  if (html === null) {
    if (elem !== null) elem.remove()
    return
  }
  if (!elem) {
    elem = Utils.createElement(
      '<div class="atk-error-layer"><span class="atk-error-title">Artalk Error</span><span class="atk-error-text"></span></div>'
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

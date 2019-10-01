export default class Utils {
  public static createElement (htmlStr: string = ''): HTMLElement {
    const div = document.createElement('div')
    div.innerHTML = htmlStr.trim()
    return div.firstElementChild as HTMLElement || div
  }

  public static getHeight (el: HTMLElement) {
    return parseFloat(getComputedStyle(el, null).height.replace('px', ''))
  }

  public static htmlEncode (str: string) {
    let temp = document.createElement('div')
    temp.innerText = str
    const output = temp.innerHTML
    temp = null
    return output
  }

  public static htmlDecode (str: string) {
    let temp = document.createElement('div')
    temp.innerHTML = str
    const output = temp.innerText
    temp = null
    return output
  }

  public static getOffset (el: HTMLElement) {
    const rect = el.getBoundingClientRect()
    return {
      top: rect.top + window.scrollY,
      left: rect.left + window.scrollX
    }
  }

  public static animate (elem: HTMLElement, style: string, unit: string, from: number, to: number, time: number, prop: boolean) {
    if (!elem) return
    const start = new Date().getTime()
    const timer = setInterval(() => {
      const step = Math.min(1, (new Date().getTime() - start) / time);
      if (prop) {
        elem[style] = (from + step * (to - from))+unit;
      } else {
        elem.style[style] = (from + step * (to - from))+unit;
      }
      if (step === 1) {
        clearInterval(timer);
      }
    }, 25)

    if (prop) {
      elem[style] = from+unit;
    } else {
      elem.style[style] = from+unit;
    }
  }
}

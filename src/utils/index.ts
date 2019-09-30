export default class Utils {
  public static createElement (htmlStr: string = ''): HTMLElement {
    const div = document.createElement('div')
    div.innerHTML = htmlStr.trim()
    return div.firstElementChild as HTMLElement || div
  }

  public static getElementHeight (el: HTMLElement) {
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
      top: rect.top + document.body.scrollTop,
      left: rect.left + document.body.scrollLeft
    }
  }
}

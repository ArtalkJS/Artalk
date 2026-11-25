import CSS from './styles.css?inline'

export class Host {
  readonly el: HTMLElement
  readonly root: ShadowRoot

  constructor(container: HTMLElement = document.body) {
    this.el = document.createElement('div')
    this.el.className = 'atk-annotation-host'
    this.el.setAttribute('data-artalk-annotation', '')
    container.appendChild(this.el)

    this.root = this.el.attachShadow({ mode: 'open' })
    const style = document.createElement('style')
    style.textContent = CSS
    this.root.appendChild(style)
  }

  setDarkMode(dark: boolean): void {
    this.el.dataset.theme = dark ? 'dark' : 'light'
  }

  destroy(): void {
    this.el.remove()
  }
}

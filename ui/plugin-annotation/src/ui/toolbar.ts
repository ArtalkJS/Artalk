export class Toolbar {
  readonly el: HTMLButtonElement

  constructor(root: ShadowRoot, label: string) {
    this.el = document.createElement('button')
    this.el.className = 'atk-annotation-toolbar'
    this.el.type = 'button'
    this.el.textContent = label
    this.el.hidden = true
    root.appendChild(this.el)
  }

  show(rect: DOMRect, onClick: () => void): void {
    this.el.style.top = `${window.scrollY + rect.top - 40}px`
    this.el.style.left = `${window.scrollX + rect.left + rect.width / 2}px`
    this.el.onclick = (evt) => {
      evt.preventDefault()
      evt.stopPropagation()
      onClick()
    }
    this.el.hidden = false
  }

  hide(): void {
    this.el.hidden = true
  }

  destroy(): void {
    this.el.remove()
  }
}

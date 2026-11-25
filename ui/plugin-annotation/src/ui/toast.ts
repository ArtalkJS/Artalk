export type ToastKind = 'success' | 'error' | 'info'

export class Toast {
  private readonly el: HTMLDivElement
  private timer?: number

  constructor(root: ShadowRoot) {
    this.el = document.createElement('div')
    this.el.className = 'atk-annotation-toast'
    this.el.hidden = true
    root.appendChild(this.el)
  }

  show(message: string, kind: ToastKind = 'info'): void {
    this.el.textContent = message
    this.el.dataset.kind = kind
    this.el.hidden = false
    window.clearTimeout(this.timer)
    this.timer = window.setTimeout(() => {
      this.el.hidden = true
    }, 3000)
  }

  destroy(): void {
    window.clearTimeout(this.timer)
    this.el.remove()
  }
}

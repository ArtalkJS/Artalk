import type { Messages } from '../i18n'

export interface ConfirmOptions {
  title: string
  body: string
  onConfirm: () => void
  onCancel?: () => void
}

export class ConfirmDialog {
  readonly el: HTMLDivElement
  private readonly titleEl: HTMLElement
  private readonly bodyEl: HTMLElement
  private readonly confirmBtn: HTMLButtonElement
  private readonly cancelBtn: HTMLButtonElement
  private options?: ConfirmOptions

  constructor(root: ShadowRoot, msg: Messages) {
    this.el = document.createElement('div')
    this.el.className = 'atk-annotation-confirm'
    this.el.hidden = true
    this.el.innerHTML = `
      <div class="confirm-card" role="dialog" aria-modal="true">
        <strong class="confirm-title"></strong>
        <p class="confirm-body"></p>
        <div class="confirm-actions">
          <button type="button" class="btn cancel"></button>
          <button type="button" class="btn danger confirm"></button>
        </div>
      </div>
    `
    root.appendChild(this.el)
    this.titleEl = this.el.querySelector('.confirm-title')!
    this.bodyEl = this.el.querySelector('.confirm-body')!
    this.confirmBtn = this.el.querySelector('.confirm')!
    this.cancelBtn = this.el.querySelector('.cancel')!
    this.cancelBtn.textContent = msg.cancel

    this.confirmBtn.addEventListener('click', () => {
      const opts = this.options
      this.close()
      opts?.onConfirm()
    })
    this.cancelBtn.addEventListener('click', () => {
      const opts = this.options
      this.close()
      opts?.onCancel?.()
    })
    this.el.addEventListener('click', (evt) => {
      if (evt.target === this.el) this.close()
    })
  }

  open(options: ConfirmOptions): void {
    this.options = options
    this.titleEl.textContent = options.title
    this.bodyEl.textContent = options.body
    this.el.hidden = false
    this.confirmBtn.focus()
  }

  close(): void {
    this.el.hidden = true
    this.options = undefined
  }

  destroy(): void {
    this.el.remove()
  }
}

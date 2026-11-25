import type { Messages } from '../i18n'

export interface ComposerOpenOptions {
  title: string
  quote: string
  onSubmit: (text: string) => Promise<void> | void
  onCancel: () => void
}

export class Composer {
  readonly el: HTMLDivElement
  private readonly titleEl: HTMLElement
  private readonly quoteEl: HTMLElement
  private readonly textarea: HTMLTextAreaElement
  private readonly errorEl: HTMLElement
  private readonly submitBtn: HTMLButtonElement
  private readonly cancelBtn: HTMLButtonElement
  private readonly countEl: HTMLElement
  private readonly msg: Messages
  private options?: ComposerOpenOptions
  private busy = false

  constructor(root: ShadowRoot, msg: Messages) {
    this.msg = msg
    this.el = document.createElement('div')
    this.el.className = 'atk-annotation-composer'
    this.el.hidden = true
    this.el.innerHTML = `
      <div class="composer-card" role="dialog" aria-modal="true">
        <div class="composer-head">
          <strong class="composer-title"></strong>
          <button type="button" class="icon-btn close" aria-label="close">×</button>
        </div>
        <blockquote class="composer-quote"></blockquote>
        <textarea rows="4" maxlength="4000"></textarea>
        <div class="composer-meta"><span class="composer-error"></span><span class="composer-count"></span></div>
        <div class="composer-actions">
          <button type="button" class="btn cancel"></button>
          <button type="button" class="btn primary submit"></button>
        </div>
      </div>
    `
    root.appendChild(this.el)
    this.titleEl = this.el.querySelector('.composer-title')!
    this.quoteEl = this.el.querySelector('.composer-quote')!
    this.textarea = this.el.querySelector('textarea')!
    this.errorEl = this.el.querySelector('.composer-error')!
    this.submitBtn = this.el.querySelector('.submit')!
    this.cancelBtn = this.el.querySelector('.cancel')!
    this.countEl = this.el.querySelector('.composer-count')!
    this.cancelBtn.textContent = msg.cancel
    this.submitBtn.textContent = msg.submit

    this.textarea.addEventListener('input', () => this.updateCount())
    this.textarea.addEventListener('keydown', (evt) => {
      if (evt.key === 'Enter' && (evt.metaKey || evt.ctrlKey)) {
        evt.preventDefault()
        void this.submit()
      }
    })
    this.submitBtn.addEventListener('click', () => void this.submit())
    this.cancelBtn.addEventListener('click', () => this.cancel())
    this.el.querySelector('.close')!.addEventListener('click', () => this.cancel())
    this.el.addEventListener('click', (evt) => {
      if (evt.target === this.el) this.cancel()
    })
  }

  open(options: ComposerOpenOptions): void {
    this.options = options
    this.titleEl.textContent = options.title
    this.quoteEl.textContent = options.quote
    this.textarea.value = ''
    this.errorEl.textContent = ''
    this.busy = false
    this.submitBtn.disabled = false
    this.el.hidden = false
    this.updateCount()
    this.textarea.focus()
  }

  close(): void {
    this.el.hidden = true
    this.options = undefined
  }

  setError(message: string): void {
    this.errorEl.textContent = message
  }

  private updateCount(): void {
    this.countEl.textContent = `${this.textarea.value.length}/4000`
  }

  private cancel(): void {
    if (this.busy) return
    const opts = this.options
    this.close()
    opts?.onCancel()
  }

  private async submit(): Promise<void> {
    if (this.busy || !this.options) return
    const text = this.textarea.value.trim()
    if (!text) {
      this.setError(this.msg.required)
      return
    }
    this.busy = true
    this.submitBtn.disabled = true
    this.errorEl.textContent = ''
    try {
      await this.options.onSubmit(text)
      this.close()
    } catch (err: any) {
      this.setError(err?.message || String(err))
    } finally {
      this.busy = false
      this.submitBtn.disabled = false
    }
  }

  destroy(): void {
    this.el.remove()
  }
}

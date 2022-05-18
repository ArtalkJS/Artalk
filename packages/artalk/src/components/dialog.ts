import Context from '~/types/context'
import * as Utils from '../lib/utils'

type BtnClickHandler = (btnEl: HTMLElement, dialog: Dialog) => boolean|void

/**
 * 对话框
 */
export default class Dialog {
  public ctx: Context
  public $el: HTMLElement
  public $content: HTMLElement
  public $actions: HTMLElement

  constructor (ctx: Context, contentEl: HTMLElement) {
    this.ctx = ctx

    this.$el = Utils.createElement(
      `<div class="atk-layer-dialog-wrap">
        <div class="atk-layer-dialog">
          <div class="atk-layer-dialog-content"></div>
          <div class="atk-layer-dialog-actions"></div>
        </div>
      </div>`
    )

    // 按钮
    this.$actions = this.$el.querySelector<HTMLElement>('.atk-layer-dialog-actions')!

    // 内容
    this.$content = this.$el.querySelector('.atk-layer-dialog-content')!
    this.$content.appendChild(contentEl)

    return this
  }

  /** 按钮 · 确定 */
  public setYes(handler: BtnClickHandler) {
    const btn = Utils.createElement<HTMLButtonElement>(
      `<button data-action="confirm">${this.ctx.$t('confirm')}</button>`
    )
    btn.onclick = this.onBtnClick(handler)
    this.$actions.appendChild(btn)

    return this
  }

  /** 按钮 · 取消 */
  public setNo(handler: BtnClickHandler) {
    const btn = Utils.createElement<HTMLButtonElement>(
      `<button data-action="cancel">${this.ctx.$t('cancel')}</button>`
    )
    btn.onclick = this.onBtnClick(handler)
    this.$actions.appendChild(btn)

    return this
  }

  private onBtnClick(handler: BtnClickHandler) {
    return (evt: Event) => {
      const re = handler(evt.currentTarget as HTMLElement, this)
      if (re === undefined || re === true) {
        this.$el.remove()
      }
    }
  }
}

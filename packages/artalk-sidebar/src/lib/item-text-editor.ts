import '../style/item-text-editor.less'

import * as Utils from 'artalk/src/lib/utils'

interface ItemTextEditorConf {
  initValue?: string
  validator?: (value: string) => boolean
  onYes?: (value: string) => boolean|void|Promise<boolean|void>
  onNo?: () => boolean|void|Promise<boolean|void>
  placeholder?: string
}

/**
 * 通用列表项目文本编辑器
 */
export default class ItemTextEditor {
  private conf: ItemTextEditorConf
  public $el: HTMLElement

  public $input: HTMLInputElement
  public $yesBtn: HTMLElement
  public $noBtn: HTMLElement

  public value: string = ''
  public allowSubmit = true

  constructor(conf: ItemTextEditorConf) {
    this.conf = conf

    this.$el = Utils.createElement(
    `<div class="atk-item-text-editor-layer">
      <div class="atk-edit-form">
        <input class="atk-main-input" type="text" placeholder="输入内容..." autocomplete="off" autofocus>
      </div>
      <div class="atk-actions">
        <div class="atk-item atk-yes-btn">
          <i class="atk-icon atk-icon-yes"></i>
        </div>
        <div class="atk-item atk-no-btn">
          <i class="atk-icon atk-icon-no"></i>
        </div>
      </div>
    </div>`)

    this.$input = this.$el.querySelector('.atk-main-input')!
    this.$yesBtn = this.$el.querySelector('.atk-yes-btn')!
    this.$noBtn = this.$el.querySelector('.atk-no-btn')!

    this.$input.value = conf.initValue || ''
    this.value = conf.initValue || ''
    if (this.conf.placeholder) this.$input.placeholder = this.conf.placeholder

    this.$input.oninput = () => this.onInput()
    this.$input.onkeyup = (evt) => {
      if (evt.key === 'Enter' || evt.keyCode === 13) { // 按下回车键
        evt.preventDefault()
        this.submit()
      }
    }

    window.setTimeout(() => this.$input.focus(), 80)

    this.$yesBtn.onclick = () => { this.submit() }
    this.$noBtn.onclick = () => { this.cancel() }
  }

  public appendTo(parentDOM: HTMLElement) {
    parentDOM.append(this.$el)
    return this
  }

  private onInput() {
    this.value = this.$input.value

    // 验证器
    if (this.conf.validator) {
      const ok = this.conf.validator(this.value)
      this.setAllowSubmit(ok)
      if (!ok) {
        this.$input.classList.add('atk-invalid')
      } else {
        this.$input.classList.remove('atk-invalid')
      }
    }
  }

  public setAllowSubmit(allow: boolean) {
    if (this.allowSubmit === allow) return
    this.allowSubmit = allow
    if (!allow) {
      this.$yesBtn.classList.add('.atk-disabled')
    } else {
      this.$yesBtn.classList.remove('.atk-disabled')
    }
  }

  public async submit() {
    if (!this.allowSubmit) return

    if (this.conf.onYes) {
      let isContinue: any
      if (this.conf.onYes instanceof (async () => {}).constructor) {
        isContinue = await this.conf.onYes(this.value)
      } else {
        isContinue = this.conf.onYes(this.value)
      }

      if (isContinue === undefined || isContinue === true) {
        this.closeEditor()
      }
    } else {
      this.closeEditor()
    }
  }

  public async cancel() {
    if (this.conf.onNo) {
      let isContinue: any
      if (this.conf.onNo instanceof (async () => {}).constructor) {
        isContinue = await this.conf.onNo()
      } else {
        isContinue = this.conf.onNo()
      }

      if (isContinue === undefined || isContinue === true) {
        this.closeEditor()
      }
    } else {
      this.closeEditor()
    }
  }

  public closeEditor() {
    this.$el.remove()
  }
}

import * as Utils from '../lib/utils'

export interface ReadMoreBtnOptions {
  /** 回调函数 */
  onClick: (offset: number) => void

  pageSize: number

  text?: string
}

/**
 * 阅读更多按钮
 */
export default class ReadMoreBtn {
  public opts: ReadMoreBtnOptions
  public $el: HTMLElement
  private $loading: HTMLElement
  private $text: HTMLElement
  private offset: number = 0
  private total: number = 0
  private origText: string = 'Load More'

  /** 是否有更多内容 */
  get hasMore() {
    return this.total > this.offset + this.opts.pageSize
  }

  get hasError() {
    return this.$el.classList.contains('atk-err')
  }

  public constructor(opts: ReadMoreBtnOptions) {
    this.opts = opts
    this.origText = this.opts.text || this.origText

    this.$el = Utils.createElement(
      `<div class="atk-list-read-more" style="display: none;">
      <button type="button" tabindex="0" class="atk-list-read-more-inner">
        <span class="atk-loading-icon" style="display: none;"></span>
        <span class="atk-text">${this.origText}</span>
      </button>
    </div>`,
    )

    this.$loading = this.$el.querySelector<HTMLElement>('.atk-loading-icon')!
    this.$text = this.$el.querySelector<HTMLElement>('.atk-text')!

    this.$el.onclick = () => {
      this.click()
    }
  }

  click() {
    if (this.hasMore) this.opts.onClick(this.offset + this.opts.pageSize)
    this.checkDisabled()
  }

  /** 显示 */
  show() {
    this.$el.style.display = ''
  }

  /** 隐藏 */
  hide() {
    this.$el.style.display = 'none'
  }

  /** 加载 */
  setLoading(isLoading: boolean) {
    const button = this.$el.querySelector('button')!
    button.disabled = isLoading
    button.setAttribute('aria-busy', String(isLoading))
    if (isLoading) {
      this.$el.classList.remove('atk-err')
      this.$text.innerText = this.origText
    }
    this.$loading.style.display = isLoading ? '' : 'none'
    this.$text.style.display = isLoading ? 'none' : ''
  }

  /** 错误提示 */
  showErr(errMsg: string) {
    this.setLoading(false)

    this.$text.innerText = errMsg
    this.$el.classList.add('atk-err')
  }

  /** @internal */
  dispose() {
    this.$el.onclick = null
    this.$el.remove()
  }

  /** 更新数据 */
  public update(offset: number, total: number) {
    this.offset = offset
    this.total = total

    this.checkDisabled()
  }

  public checkDisabled() {
    if (this.hasMore) this.show()
    else this.hide()
  }
}

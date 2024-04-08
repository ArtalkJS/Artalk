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

  public constructor(opts: ReadMoreBtnOptions) {
    this.opts = opts
    this.origText = this.opts.text || this.origText

    this.$el = Utils.createElement(
      `<div class="atk-list-read-more" style="display: none;">
      <div class="atk-list-read-more-inner">
        <div class="atk-loading-icon" style="display: none;"></div>
        <span class="atk-text">${this.origText}</span>
      </div>
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
    this.$loading.style.display = isLoading ? '' : 'none'
    this.$text.style.display = isLoading ? 'none' : ''
  }

  /** 错误提示 */
  showErr(errMsg: string) {
    this.setLoading(false)

    this.$text.innerText = errMsg
    this.$el.classList.add('atk-err')
    window.setTimeout(() => {
      this.$text.innerText = this.origText
      this.$el.classList.remove('atk-err')
    }, 2000) // 2s后错误提示复原
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

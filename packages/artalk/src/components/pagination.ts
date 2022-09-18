import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'

export interface PaginationConf {
  /** 每页条数 */
  pageSize: number

  /** 回调函数 */
  onChange: (offset: number) => void
}

export default class Pagination {
  private conf: PaginationConf
  public total: number
  public $el: HTMLElement
  public $input: HTMLInputElement
  public inputTimer?: number
  public $prevBtn: HTMLElement
  public $nextBtn: HTMLElement

  public page: number = 1
  get pageSize(): number {
    return this.conf.pageSize
  }
  get offset(): number {
    return this.pageSize * (this.page - 1)
  }
  get maxPage(): number {
    return Math.ceil(this.total / this.pageSize)
  }

  public constructor(total: number, conf: PaginationConf) {
    this.total = total
    this.conf = conf

    this.$el = Utils.createElement(
      `<div class="atk-pagination-wrap">
        <div class="atk-pagination">
          <div class="atk-btn atk-btn-prev">Prev</div>
          <input type="text" class="atk-input" aria-label="Enter the number of page" />
          <div class="atk-btn atk-btn-next">Next</div>
        </div>
      </div>`)
    this.$input = this.$el.querySelector('.atk-input')!
    this.$input.value = `${this.page}`

    this.$input.oninput = () => this.input()
    this.$input.onkeydown = (e) => this.keydown(e)

    this.$prevBtn = this.$el.querySelector('.atk-btn-prev')!
    this.$nextBtn = this.$el.querySelector('.atk-btn-next')!

    this.$prevBtn.onclick = () => this.prev()
    this.$nextBtn.onclick = () => this.next()

    this.checkDisabled()
  }

  public update(offset: number, total: number) {
    this.page = Math.ceil(offset / this.pageSize) + 1
    this.total = total

    this.setInput(this.page)
    this.checkDisabled()
  }

  public setInput(page: number) {
    this.$input.value = `${page}`
  }

  public input(now: boolean = false) {
    window.clearTimeout(this.inputTimer)

    const value = this.$input.value.trim()

    const modify = () => {
      if (value === '') { this.setInput(this.page);return }
      let page = Number(value)
      if (Number.isNaN(page)) { this.setInput(this.page);return }
      if (page < 1) { this.setInput(this.page);return }
      if (page > this.maxPage) { this.setInput(this.maxPage);page = this.maxPage }
      this.change(page)
    }

    // 延迟 800ms 执行
    if (!now) this.inputTimer = window.setTimeout(() => modify(), 800)
    else modify()
  }

  public prev() {
    const page = this.page - 1
    if (page < 1) { return }
    this.change(page)
  }

  public next() {
    const page = this.page + 1
    if (page > this.maxPage) { return }
    this.change(page)
  }

  public change(page: number) {
    this.page = page
    this.conf.onChange(this.offset)
    this.setInput(page)
    this.checkDisabled()
  }

  public checkDisabled() {
    if (this.page + 1 > this.maxPage) {
      this.$nextBtn.classList.add('atk-disabled')
    } else {
      this.$nextBtn.classList.remove('atk-disabled')
    }

    if (this.page - 1 < 1) {
      this.$prevBtn.classList.add('atk-disabled')
    } else {
      this.$prevBtn.classList.remove('atk-disabled')
    }
  }

  public keydown(e: KeyboardEvent) {
    const keyCode = e.keyCode || e.which

    if (keyCode === 38) {
      // 上键
      const page = Number(this.$input.value) + 1
      if (page > this.maxPage) { return }
      this.setInput(page);
      this.input()
    } else if (keyCode === 40) {
      // 下键
      const page = Number(this.$input.value) - 1
      if (page < 1) { return }
      this.setInput(page)
      this.input()
    } else if (keyCode === 13) {
      // 回车
      this.input(true)
    }
  }

  /** 加载 */
  setLoading (isLoading: boolean) {
    if (isLoading) Ui.showLoading(this.$el)
    else Ui.hideLoading(this.$el)
  }
}

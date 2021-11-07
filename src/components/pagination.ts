import '../style/pagination.less'
import * as Utils from '../lib/utils'

interface PaginationConf {
  /** 每页条数 */
  pageSize?: number

  /** 数据总数 */
  total: number

  /** 回调函数 */
  onChange: (offset: number) => void
}

export default class Pagination {
  private conf: PaginationConf
  public $el: HTMLElement
  public $input: HTMLInputElement
  public inputTimer?: number
  public $prevBtn: HTMLElement
  public $nextBtn: HTMLElement

  public page: number = 1
  get pageSize(): number {
    return this.conf.pageSize || 15
  }
  get offset(): number {
    return this.pageSize * (this.page - 1)
  }
  get maxPage(): number {
    return Math.ceil(this.conf.total / this.pageSize)
  }

  public constructor(conf: PaginationConf) {
    this.conf = conf
    this.$el = Utils.createElement(
      `<div class="atk-pagination">
        <div class="atk-btn atk-btn-prev">Prev</div>
        <input type="text" class="atk-input" />
        <div class="atk-btn atk-btn-next">Next</div>
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

  public input(now: boolean = false) {
    window.clearTimeout(this.inputTimer)

    const value = this.$input.value.trim()
    if (value === '') { return }

    const reset = (p: number) => { this.$input.value = `${p}` }
    const modify = () => {
      const page = Number(value)
      if (Number.isNaN(page)) { reset(this.page);return }
      if (page < 1) { reset(1);return }
      if (page > this.maxPage) { reset(this.maxPage);return }
      this.change(page)
    }

    // 延迟 800ms 执行
    if (!now) this.inputTimer = window.setTimeout(() => modify(), 800)
    else modify()
  }

  public keydown(e: KeyboardEvent) {
    const keyCode = e.keyCode || e.which

    if (keyCode === 38) {
      // 上键
      this.next()
    } else if (keyCode === 40) {
      // 下键
      this.prev()
    } else if (keyCode === 13) {
      // 回车
      this.input(true)
    }
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
    this.$input.value = String(page)
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
}

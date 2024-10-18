import type { Paginator, IPgHolderOpt } from '.'
import ReadMoreBtn from '@/components/read-more-btn'
import $t from '@/i18n'

/**
 * 阅读更多形式的分页
 */
export default class ReadMorePaginator implements Paginator {
  private instance!: ReadMoreBtn
  private onReachedBottom: (() => void) | null = null
  private opt!: IPgHolderOpt

  create(opt: IPgHolderOpt) {
    this.opt = opt

    this.instance = new ReadMoreBtn({
      pageSize: opt.pageSize,
      onClick: async (o) => {
        opt.getData().fetchComments({
          offset: o,
        })
      },
      text: $t('loadMore'),
    })

    // 滚动到底部自动加载
    if (opt.readMoreAutoLoad) {
      this.onReachedBottom = () => {
        if (!this.instance.hasMore || this.opt.getData().getLoading()) return
        this.instance.click()
      }
      this.opt.getEvents().on('list-reach-bottom', this.onReachedBottom)
    }

    return this.instance.$el
  }

  setLoading(val: boolean) {
    this.instance.setLoading(val)
  }

  update(offset: number, total: number) {
    this.instance.update(offset, total)
  }

  showErr(msg: string) {
    this.instance.showErr(msg)
  }

  next() {
    this.instance.click()
  }

  getHasMore(): boolean {
    return this.instance.hasMore
  }

  getIsClearComments(params: { offset?: number }): boolean {
    return params.offset === 0
  }

  dispose(): void {
    this.onReachedBottom && this.opt.getEvents().off('list-reach-bottom', this.onReachedBottom)
    this.instance.$el.remove()
  }
}

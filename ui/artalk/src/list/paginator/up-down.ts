import type { Paginator, IPgHolderOpt } from '.'
import PaginationComponent from '@/components/pagination'

/**
 * 翻页形式的分页
 */
export default class UpDownPaginator implements Paginator {
  private instance!: PaginationComponent

  create(opt: IPgHolderOpt) {
    this.instance = new PaginationComponent(opt.total, {
      pageSize: opt.pageSize,
      onChange: async (o) => {
        opt.resetEditorState()

        opt.getData().fetchComments({
          offset: o,
          onSuccess: () => {
            opt.onListGotoFirst?.()
          },
        })
      },
    })

    return this.instance.$el
  }

  setLoading(val: boolean) {
    this.instance.setLoading(val)
  }

  update(offset: number, total: number) {
    this.instance.update(offset, total)
  }

  next() {
    this.instance.next()
  }

  getHasMore(): boolean {
    return this.instance.getHasMore()
  }

  getIsClearComments(): boolean {
    return true
  }

  dispose(): void {
    this.instance.$el.remove()
  }
}

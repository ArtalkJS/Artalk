import Pagination from '@/components/pagination'
import * as Utils from '@/lib/utils'
import { IPgAdaptor } from '.'

interface IPaginationAdaptor extends IPgAdaptor<Pagination> {}

/**
 * 翻页形式的分页
 */
export default <IPaginationAdaptor>{
  createInstance(conf) {
    const instance = new Pagination(conf.total, {
      pageSize: conf.pageSize,
      onChange: async (o) => {
        if (conf.list.conf.editorTravel === true)
          conf.list.ctx.editorTravelBack() // 防止评论框被吞

        await conf.list.fetchComments(o)

        // 滚动到第一个评论的位置
        if (conf.list.repositionAt) {
          const at = conf.list.scrollListenerAt || window
          at.scroll({
            top: conf.list.repositionAt ? Utils.getOffset(conf.list.repositionAt).top : 0,
            left: 0,
          })
        }
      }
    })

    return [instance, instance.$el]
  },
  setLoading(val) {
    this.instance.setLoading(val)
  },
  update(offset, total) {
    this.instance.update(offset, total)
  },
  next() {
    this.instance.next()
  }
}

import ReadMoreBtn from '@/components/read-more-btn'
import $t from '@/i18n'
import * as Ui from '@/lib/ui'
import { IPgAdaptor } from '.'

interface IReadMoreAdaptor extends IPgAdaptor<ReadMoreBtn> {
  autoLoadScrollEvent?: () => void
}

/**
 * 阅读更多形式的分页
 */
export default <IReadMoreAdaptor>{
  createInstance(conf) {
    const readMoreBtn = new ReadMoreBtn({
      pageSize: conf.pageSize,
      onClick: async (o) => {
        await conf.list.fetchComments(o)
      },
      text: $t('loadMore'),
    })

    // 滚动到底部自动加载
    if (conf.list.conf.pagination.autoLoad) {
      // 添加滚动事件监听
      const at = conf.list.scrollListenerAt || document
      if (this.autoLoadScrollEvent) at.removeEventListener('scroll', this.autoLoadScrollEvent) // 解除原有
      this.autoLoadScrollEvent = () => {
        if (conf.mode !== 'read-more'
          || !readMoreBtn?.hasMore
          || conf.list.getLoading()
        ) return

        const $target = conf.list.$el.querySelector<HTMLElement>('.atk-list-comments-wrap > .atk-comment-wrap:nth-last-child(3)') // 获取倒数第3个评论元素
        if ($target && Ui.isVisible($target, conf.list.scrollListenerAt)) {
          readMoreBtn.click() // 自动点击加载更多按钮
        }
      }
      at.addEventListener('scroll', this.autoLoadScrollEvent)
    }

    return [readMoreBtn, readMoreBtn.$el]
  },
  setLoading(val) {
    this.instance.setLoading(val)
  },
  update(offset, total) {
    this.instance.update(offset, total)
  },
  showErr(msg) {
    this.instance.showErr(msg)
  },
  next() {
    this.instance.click()
  }
}

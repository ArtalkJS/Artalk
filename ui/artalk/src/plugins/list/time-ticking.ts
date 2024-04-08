import type { ArtalkPlugin } from '@/types'
import * as Utils from '@/lib/utils'

/** 评论时间自动更新 */
export const TimeTicking: ArtalkPlugin = (ctx) => {
  let timer: number | null = null

  ctx.on('mounted', () => {
    timer = window.setInterval(() => {
      const list = ctx.get('list')

      list.$el.querySelectorAll<HTMLElement>('[data-atk-comment-date]').forEach((el) => {
        const date = el.getAttribute('data-atk-comment-date')
        el.innerText = Utils.timeAgo(new Date(Number(date)), ctx.$t)
      })
    }, 30 * 1000) // 30s 更新一次
  })

  ctx.on('unmounted', () => {
    timer && window.clearInterval(timer)
  })
}

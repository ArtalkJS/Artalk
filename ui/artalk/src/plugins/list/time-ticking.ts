import type { ArtalkPlugin } from '@/types'
import * as Utils from '@/lib/utils'

/** 评论时间自动更新 */
export const TimeTicking: ArtalkPlugin = (ctx) => {
  const list = ctx.inject('list')
  const conf = ctx.inject('config')
  let timer: number | null = null

  ctx.on('mounted', () => {
    timer = window.setInterval(() => {
      list
        .getEl()
        .querySelectorAll<HTMLElement>('[data-atk-comment-date]')
        .forEach((el) => {
          const date = new Date(Number(el.getAttribute('data-atk-comment-date')))
          el.innerText = conf.get().dateFormatter?.(date) || Utils.timeAgo(date, ctx.$t)
        })
    }, 30 * 1000) // 30s 更新一次
  })

  ctx.on('unmounted', () => {
    timer && window.clearInterval(timer)
  })
}

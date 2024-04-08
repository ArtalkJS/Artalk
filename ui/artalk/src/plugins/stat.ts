import type { ContextApi, ArtalkPlugin, ArtalkConfig } from '@/types'
import { Api } from '@/api'

export interface CountOptions {
  getApi(): Api

  siteName: string
  pageKey?: string
  pageTitle?: string
  countEl: string
  pvEl: string

  /** 是否增加当前页面 PV 数 */
  pvAdd?: boolean
}

export const PvCountWidget: ArtalkPlugin = (ctx: ContextApi) => {
  ctx.watchConf(['site', 'pageKey', 'pageTitle', 'countEl', 'pvEl'], (conf) => {
    initCountWidget({
      getApi: () => ctx.getApi(),
      siteName: conf.site,
      pageKey: conf.pageKey,
      pageTitle: conf.pageTitle,
      countEl: conf.countEl,
      pvEl: conf.pvEl,
      pvAdd: typeof ctx.conf.pvAdd === 'boolean' ? ctx.conf.pvAdd : true,
    })
  })
}

/** 初始化评论数和 PV 数量展示元素 */
export async function initCountWidget(opt: CountOptions) {
  // 评论数
  if (opt.countEl && document.querySelector(opt.countEl)) {
    refreshStatCount(opt, { query: 'page_comment', numEl: opt.countEl })
  }

  // PV
  const initialData =
    opt.pvAdd && opt.pageKey
      ? {
          [opt.pageKey]: (
            await opt.getApi().pages.logPv({
              page_key: opt.pageKey,
              page_title: opt.pageTitle,
              site_name: opt.siteName,
            })
          ).data.pv, // pv+1 and get pv count
        }
      : undefined

  if (opt.pvEl && document.querySelector(opt.pvEl)) {
    refreshStatCount(opt, {
      query: 'page_pv',
      numEl: opt.pvEl,
      data: initialData,
    })
  }
}

type CountData = { [pageKey: string]: number }

async function refreshStatCount(
  opt: CountOptions,
  args: {
    query: 'page_pv' | 'page_comment'
    numEl: string
    data?: CountData
  },
) {
  let data: CountData = args.data || {}

  // Get page keys which will be queried
  let queryPageKeys = Array.from(document.querySelectorAll(args.numEl))
    .map((e) => e.getAttribute('data-page-key') || opt.pageKey)
    .filter((k) => k && typeof data[k] !== 'number') // filter out keys that already have data

  queryPageKeys = [...new Set(queryPageKeys)] // deduplicate

  // Fetch count data from server
  if (queryPageKeys.length > 0) {
    const res = (
      await opt.getApi().stats.getStats(args.query, {
        page_keys: queryPageKeys.join(','),
        site_name: opt.siteName,
      })
    ).data.data as CountData
    data = { ...data, ...res }
  }

  const defaultCount = opt.pageKey ? data[opt.pageKey] : 0
  applyCountData(args.numEl, data, defaultCount)
}

function applyCountData(selector: string, data: CountData, defaultCount: number) {
  document.querySelectorAll(selector).forEach((el) => {
    const pageKey = el.getAttribute('data-page-key')
    const count = Number(pageKey ? data[pageKey] : defaultCount)
    el.innerHTML = `${count}`
  })
}

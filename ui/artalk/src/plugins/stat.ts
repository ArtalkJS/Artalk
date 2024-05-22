import type { ContextApi, ArtalkPlugin, ArtalkConfig } from '@/types'
import { Api } from '@/api'

export interface CountOptions {
  getApi(): Api

  siteName: string
  pageKey?: string
  pageTitle?: string
  countEl: string
  pvEl: string
  pageKeyAttr: string

  /** 是否增加当前页面 PV 数 */
  pvAdd?: boolean
}

export const PvCountWidget: ArtalkPlugin = (ctx: ContextApi) => {
  ctx.watchConf(['site', 'pageKey', 'pageTitle', 'countEl', 'pvEl', 'statPageKeyAttr'], (conf) => {
    initCountWidget({
      getApi: () => ctx.getApi(),
      siteName: conf.site,
      pageKey: conf.pageKey,
      pageTitle: conf.pageTitle,
      countEl: conf.countEl,
      pvEl: conf.pvEl,
      pageKeyAttr: conf.statPageKeyAttr,
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
  const cacheData =
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
      cacheData,
    })
  }
}

type CountData = { [pageKey: string]: number }

async function refreshStatCount(
  opt: CountOptions,
  args: {
    query: 'page_pv' | 'page_comment'
    numEl: string
    cacheData?: CountData
  },
) {
  let data: CountData = args.cacheData || {}

  // Retrieve elements
  const els = Array.from(document.querySelectorAll<HTMLElement>(args.numEl))

  // Get page keys which will be queried
  let pageKeys = els
    .map((el) => el.getAttribute(opt.pageKeyAttr) || opt.pageKey)
    .filter((pageKey) => pageKey && typeof data[pageKey] !== 'number') // filter out keys that already have data

  pageKeys = [...new Set(pageKeys)] // deduplicate

  // Fetch count data from server
  if (pageKeys.length > 0) {
    const res = (
      await opt.getApi().stats.getStats(args.query, {
        page_keys: pageKeys.join(','),
        site_name: opt.siteName,
      })
    ).data.data as CountData
    data = { ...data, ...res }
  }

  const defaultCount = opt.pageKey ? data[opt.pageKey] : 0
  applyCountData(els, data, defaultCount, opt.pageKeyAttr)
}

function applyCountData(elements: HTMLElement[], data: CountData, defaultCount: number, pageKeyAttr: string) {
  elements.forEach((el) => {
    const pageKey = el.getAttribute(pageKeyAttr)
    const count = Number(pageKey ? data[pageKey] : defaultCount) // if pageKey is not set, use defaultCount
    el.innerText = `${count}`
  })
}

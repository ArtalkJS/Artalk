import type { Context, ArtalkPlugin } from '@/types'
import { Api } from '@/api'

type CountCache = { [pageKey: string]: number }

export interface CountOptions {
  getApi(): Api

  siteName: string
  pageKey?: string
  pageTitle?: string
  countEl: string
  pvEl: string
  pageKeyAttr: string

  /** Whether to add PV count when initializing */
  pvAdd?: boolean
}

export const PvCountWidget: ArtalkPlugin = (ctx: Context) => {
  ctx.watchConf(
    ['site', 'pageKey', 'pageTitle', 'countEl', 'pvEl', 'statPageKeyAttr', 'pvAdd'],
    (conf) => {
      initCountWidget({
        getApi: () => ctx.getApi(),
        siteName: conf.site,
        pageKey: conf.pageKey,
        pageTitle: conf.pageTitle,
        countEl: conf.countEl,
        pvEl: conf.pvEl,
        pageKeyAttr: conf.statPageKeyAttr,
        pvAdd: conf.pvAdd,
      })
    },
  )
}

/** Initialize count widgets */
export async function initCountWidget(opt: CountOptions) {
  // Load comment count
  await loadCommentCount(opt)

  // Increment PV count
  const cacheData = await incrementPvCount(opt)

  // Load PV count
  await loadPvCount(opt, cacheData)
}

/** Increment PV count and get cache data that contains PV count */
async function incrementPvCount(opt: CountOptions) {
  if (!opt.pvAdd || !opt.pageKey) return undefined

  const pvCount = (
    await opt.getApi().pages.logPv({
      page_key: opt.pageKey,
      page_title: opt.pageTitle,
      site_name: opt.siteName,
    })
  ).data.pv // pv+1 and get pv count

  return {
    [opt.pageKey]: pvCount,
  }
}

/** Load comment count */
async function loadCommentCount(opt: CountOptions) {
  await loadStatCount({ opt, query: 'page_comment', containers: [opt.countEl, '#ArtalkCount'] })
}

/** Load PV count */
async function loadPvCount(opt: CountOptions, cache?: CountCache) {
  await loadStatCount({ opt, query: 'page_pv', containers: [opt.pvEl, '#ArtalkPV'], cache })
}

async function loadStatCount(args: {
  opt: CountOptions
  query: 'page_pv' | 'page_comment'
  containers: string[]
  cache?: CountCache
}) {
  const { opt } = args
  let cache: CountCache = args.cache || {}

  // Retrieve elements
  const els = retrieveElements(args.containers)

  // Get page keys which will be queried
  const pageKeys = getPageKeys(els, opt.pageKeyAttr, opt.pageKey, cache)

  // Fetch count data from server
  if (pageKeys.length > 0) {
    const res = (
      await opt.getApi().stats.getStats(args.query, {
        page_keys: pageKeys.join(','),
        site_name: opt.siteName,
      })
    ).data.data as CountCache
    cache = { ...cache, ...res }
  }

  updateElementsText(els, cache, opt.pageKey)
}

/** Retrieve elements based on selectors */
function retrieveElements(containers: string[]): Set<HTMLElement> {
  const els = new Set<HTMLElement>()
  new Set(containers).forEach((selector) => {
    document.querySelectorAll<HTMLElement>(selector).forEach((el) => els.add(el))
  })
  return els
}

/** Get page keys to be queried */
function getPageKeys(
  els: Set<HTMLElement>,
  pageKeyAttr: string,
  pageKey: string | undefined,
  cache: CountCache,
): string[] {
  const pageKeys = Array.from(els)
    .map((el) => el.getAttribute(pageKeyAttr) || pageKey)
    .filter((key) => key && typeof cache[key] !== 'number') // filter out keys that already have data

  return [...new Set(pageKeys as string[])] // deduplicate
}

/** Update elements text content with the count data */
function updateElementsText(els: Set<HTMLElement>, data: CountCache, defaultPageKey?: string) {
  els.forEach((el) => {
    const pageKey = el.getAttribute('data-page-key')
    const count = (pageKey && data[pageKey]) || (defaultPageKey && data[defaultPageKey]) || 0 // if pageKey is not set, use defaultCount
    el.innerText = `${Number(count)}`
  })
}

export const exportedForTesting = {
  incrementPvCount,
  loadCommentCount,
  loadPvCount,
  retrieveElements,
  getPageKeys,
  updateElementsText,
}

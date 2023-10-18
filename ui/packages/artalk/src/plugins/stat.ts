import ContextApi from '~/types/context'
import ArtalkPlug from '~/types/plug'

export interface CountConf {
  ctx: ContextApi

  /** 是否增加当前页面 PV 数 */
  pvAdd?: boolean
}

export const PvCountWidget: ArtalkPlug = (ctx: ContextApi) => {
  if (!ctx.conf.useBackendConf) {
    // 不使用后端配置，在 Artalk 实例被创建后，立刻初始化
    initCountWidget({ ctx, pvAdd: true })
  } else {
    // 若使用后端配置，需待配置成功获取后 (来自后端设定的 pvEl 等)，再初始化
    ctx.on('list-loaded', () => {
      initCountWidget({ ctx, pvAdd: true })
    })
  }
}

/** 初始化评论数和 PV 数量展示元素 */
export async function initCountWidget(p: CountConf) {
  // 评论数
  const countEl = p.ctx.conf.countEl
  if (countEl && document.querySelector(countEl)) {
    handleStatCount(p, { api: 'page_comment', countEl })
  }

  // PV
  const curtPagePvNum = p.pvAdd ? await p.ctx.getApi().page.pv() : undefined
  const pvEl = p.ctx.conf.pvEl
  if (pvEl && document.querySelector(pvEl)) {
    handleStatCount(p, {
      api: 'page_pv',
      countEl: pvEl,
      curtPageCount: curtPagePvNum,
    })
  }
}

export async function handleStatCount(
  p: CountConf,
  args: {
    api: 'page_pv' | 'page_comment'
    countEl: string
    curtPageCount?: number
  }
) {
  let pageCounts: { [key: string]: number } = {}

  // 当前页面的统计数
  const curtPageKey = p.ctx.conf.pageKey
  if (args.curtPageCount) pageCounts[curtPageKey] = args.curtPageCount

  // 查询其他页面的统计数
  let queryPageKeys = Array.from(document.querySelectorAll(args.countEl))
    .map((e) => e.getAttribute('data-page-key') || curtPageKey)
    .filter((pageKey) => pageCounts[pageKey] === undefined)
  queryPageKeys = [...new Set(queryPageKeys)] // 去重
  if (queryPageKeys.length > 0) {
    const counts: any = await p.ctx.getApi().page.stat(args.api, queryPageKeys)
    pageCounts = { ...pageCounts, ...counts }
  }

  document.querySelectorAll(args.countEl).forEach((el) => {
    const pageKey = el.getAttribute('data-page-key') || curtPageKey
    el.innerHTML = `${Number(pageCounts[pageKey] || 0)}`
  })
}

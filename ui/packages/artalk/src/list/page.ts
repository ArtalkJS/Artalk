import type { ArtalkConfig, ContextApi } from '~/types'
import $t from '@/i18n'
import { Paginator } from './paginator'
import ReadMorePaginator from './paginator/read-more'
import UpDownPaginator from './paginator/up-down'

function createPaginatorByConf(conf: ArtalkConfig): Paginator {
  if (conf.pagination.readMore) return new ReadMorePaginator()
  return new UpDownPaginator()
}

function getPageDataByLastData(ctx: ContextApi): { offset: number, total: number } {
  const last = ctx.getData().getListLastFetch()
  const r = { offset: 0, total: 0 }
  if (!last) return r

  r.offset = last.params.offset
  if (last.data) r.total = last.params.flatMode ? last.data.total : last.data.total_roots

  return r
}

export const initListPaginatorFunc = (ctx: ContextApi) => {
  let paginator: Paginator|null = null

  // Init paginator when conf loaded
  ctx.on('conf-loaded', (conf) => {
    const list = ctx.get('list')

    if (paginator) paginator.dispose() // if had been init, dispose it

    // create paginator instance
    paginator = createPaginatorByConf(conf)

    // create paginator dom
    const { offset, total } = getPageDataByLastData(ctx)
    const $paginator = paginator.create({
      ctx, pageSize: conf.pagination.pageSize, total,

      readMoreAutoLoad: conf.pagination.autoLoad,
    })

    // mount paginator dom
    list.$el.append($paginator)
  })

  // When list loaded
  ctx.on('list-loaded', (comments) => {
    // update paginator info
    const { offset, total } = getPageDataByLastData(ctx)
    paginator?.update(offset, total)
  })

  // When list fetch
  ctx.on('list-fetch', (params) => {
    // if clear comments when fetch new page data
    if (ctx.getData().getComments().length > 0 && paginator?.getIsClearComments(params)) {
      ctx.getData().clearComments()
    }
  })

  // When list error
  ctx.on('list-error', () => {
    paginator?.showErr?.($t('loadFail'))
  })

  // List goto auto next page when comment not found
  const autoSwitchPageForFindComment = (commentID: number) => {
    const comment = ctx.getData().findComment(commentID)
    if (!!comment || !paginator?.getHasMore()) return

    // TODO 自动范围改为直接跳转到计算后的页面
    paginator?.next()

    // wait for list loaded
    ctx.on('list-loaded', () => {
      autoSwitchPageForFindComment(commentID) // recursive, until comment found or no more page
    }, { once: true })
  }

  ctx.on('list-goto', (commentID) => {
    autoSwitchPageForFindComment(commentID)
  })

  // loading
  ctx.on('list-fetch', (params) => {
    paginator?.setLoading(true)
  })
  ctx.on('list-fetched', ({ params }) => {
    paginator?.setLoading(false)
  })
}

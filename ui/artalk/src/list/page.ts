import type { ArtalkConfig, ContextApi } from '@/types'
import $t from '@/i18n'
import { Paginator } from './paginator'
import ReadMorePaginator from './paginator/read-more'
import UpDownPaginator from './paginator/up-down'

function createPaginatorByConf(conf: Pick<ArtalkConfig, 'pagination'>): Paginator {
  if (conf.pagination.readMore) return new ReadMorePaginator()
  return new UpDownPaginator()
}

function getPageDataByLastData(ctx: ContextApi): { offset: number, total: number } {
  const last = ctx.getData().getListLastFetch()
  const r = { offset: 0, total: 0 }
  if (!last) return r

  r.offset = last.params.offset
  if (last.data) r.total = last.params.flatMode ? last.data.count : last.data.roots_count

  return r
}

export const initListPaginatorFunc = (ctx: ContextApi) => {
  let paginator: Paginator|null = null

  // Init paginator when conf loaded
  ctx.watchConf(['pagination', 'locale'], (conf) => {
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
    list.$commentsWrap.after($paginator)

    // update paginator info
    paginator?.update(offset, total)
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
  ctx.on('list-failed', () => {
    paginator?.showErr?.($t('loadFail'))
  })

  // loading
  ctx.on('list-fetch', (params) => {
    paginator?.setLoading(true)
  })
  ctx.on('list-fetched', ({ params }) => {
    paginator?.setLoading(false)
  })
}

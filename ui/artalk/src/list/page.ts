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

  // List goto auto next page when comment not found
  // autoJumpToNextPage(ctx, paginator)

  // loading
  ctx.on('list-fetch', (params) => {
    paginator?.setLoading(true)
  })
  ctx.on('list-fetched', ({ params }) => {
    paginator?.setLoading(false)
  })
}

function autoJumpToNextPage(ctx: ContextApi, paginator: Paginator|null) {
  // TODO: Disable this feature temporarily. Because it may cause the page memory leak.
  // if the comments is too much and the comment still not found.
  // Consider to refactor to a better solution.
  // Such as: calculate in backend and jump to the specific page directly.
  return

  const autoSwitchPageForFindComment = (commentID: number) => {
    const comment = ctx.getData().findComment(commentID)
    if (!!comment || !paginator?.getHasMore()) return

    // wait for list loaded
    ctx.on('list-loaded', () => {
      autoSwitchPageForFindComment(commentID) // recursive, until comment found or no more page
    }, { once: true })

    // TODO: 自动范围改为直接跳转到计算后的页面
    setTimeout(() => {
      paginator?.next()
    }, 80)
  }

  ctx.on('list-goto', (commentID) => {
    autoSwitchPageForFindComment(commentID)
  })
}

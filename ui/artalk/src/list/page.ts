import type { ListOptions } from './list'
import { Paginator } from './paginator'
import ReadMorePaginator from './paginator/read-more'
import UpDownPaginator from './paginator/up-down'
import $t from '@/i18n'
import type { Config, List, ListLastFetchData } from '@/types'

function createPaginatorByConf(conf: Pick<Config, 'pagination'>): Paginator {
  if (conf.pagination.readMore) return new ReadMorePaginator()
  return new UpDownPaginator()
}

function getPageDataByLastData(last: ListLastFetchData | undefined): {
  offset: number
  total: number
} {
  const r = { offset: 0, total: 0 }
  if (!last) return r

  r.offset = last.params.offset
  if (last.data) r.total = last.params.flatMode ? last.data.count : last.data.roots_count

  return r
}

export interface PaginatorInitOptions extends ListOptions {
  getList: () => List
}

export const initListPaginatorFunc = (opts: PaginatorInitOptions) => {
  let paginator: Paginator | null = null
  let pageData = getPageDataByLastData(opts.getData().getListLastFetch())
  let failed = false
  const showError = () => paginator?.showErr?.(`${$t('loadFail')} · ${$t('listRetry')}`)
  opts.getEvents().on('unmounted', () => paginator?.dispose())

  // Init paginator when conf loaded
  opts.getConf().watchConf(['pagination', 'locale'], (conf) => {
    const list = opts.getList()

    if (paginator) paginator.dispose() // if had been init, dispose it

    // create paginator instance
    paginator = createPaginatorByConf(conf)

    // create paginator dom
    // Keep the last loaded page when a locale change follows a failed request.
    const { offset, total } = pageData
    const $paginator = paginator.create({
      pageSize: conf.pagination.pageSize,
      total,

      readMoreAutoLoad: conf.pagination.autoLoad,

      ...opts,
    })

    // mount paginator dom
    list.getCommentsWrapEl().after($paginator)

    // update paginator info
    paginator?.update(offset, total)
    if (failed) showError()
  })

  // When list loaded
  opts.getEvents().on('list-loaded', (comments) => {
    failed = false
    // update paginator info
    pageData = getPageDataByLastData(opts.getData().getListLastFetch())
    const { offset, total } = pageData
    paginator?.update(offset, total)
  })

  // When list fetch
  opts.getEvents().on('list-fetch', (params) => {
    failed = false
    // if clear comments when fetch new page data
    if (opts.getData().getComments().length > 0 && paginator?.getIsClearComments(params)) {
      opts.getData().clearComments()
    }
  })

  // When list error
  opts.getEvents().on('list-failed', () => {
    failed = true
    showError()
  })

  // loading
  opts.getEvents().on('list-fetch', (params) => {
    paginator?.setLoading(true)
  })
  opts.getEvents().on('list-fetched', ({ params }) => {
    paginator?.setLoading(false)
  })
}

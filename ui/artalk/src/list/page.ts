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

  // Init paginator when conf loaded
  opts.getConf().watchConf(['pagination', 'locale'], (conf) => {
    const list = opts.getList()
    const data = opts.getData()

    if (paginator) paginator.dispose() // if had been init, dispose it

    // create paginator instance
    paginator = createPaginatorByConf(conf)

    // create paginator dom
    const { offset, total } = getPageDataByLastData(data.getListLastFetch())
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
  })

  // When list loaded
  opts.getEvents().on('list-loaded', (comments) => {
    // update paginator info
    const { offset, total } = getPageDataByLastData(opts.getData().getListLastFetch())
    paginator?.update(offset, total)
  })

  // When list fetch
  opts.getEvents().on('list-fetch', (params) => {
    // if clear comments when fetch new page data
    if (opts.getData().getComments().length > 0 && paginator?.getIsClearComments(params)) {
      opts.getData().clearComments()
    }
  })

  // When list error
  opts.getEvents().on('list-failed', () => {
    paginator?.showErr?.($t('loadFail'))
  })

  // loading
  opts.getEvents().on('list-fetch', (params) => {
    paginator?.setLoading(true)
  })
  opts.getEvents().on('list-fetched', ({ params }) => {
    paginator?.setLoading(false)
  })
}

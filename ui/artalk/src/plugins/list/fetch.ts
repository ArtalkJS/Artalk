import type { ListData, ListFetchParams, ArtalkPlugin } from '@/types'

export const Fetch: ArtalkPlugin = (ctx) => {
  const conf = ctx.inject('config')

  ctx.on('list-fetch', (_params) => {
    if (ctx.getData().getLoading()) return
    ctx.getData().setLoading(true)

    const params: ListFetchParams = {
      // default params
      offset: 0,
      limit: conf.get().pagination.pageSize,
      flatMode: conf.get().flatMode as boolean, // always be boolean because had been handled in Artalk.init
      paramsModifier: conf.get().listFetchParamsModifier,
      ..._params,
    }

    // must before other function call
    ctx.getData().setListLastFetch({
      params,
    })

    // prepare params for request
    const reqParams = {
      limit: params.limit,
      offset: params.offset,
      flat_mode: params.flatMode,
      page_key: conf.get().pageKey,
      site_name: conf.get().site,
    }

    // call the modifier function
    if (params.paramsModifier) params.paramsModifier(reqParams)

    // start request
    ctx
      .getApi()
      .comments.getComments({
        ...reqParams,
        ...ctx.getApi().getUserFields(),
      })
      .then(({ data }) => {
        const scope = (reqParams as { scope?: 'page' | 'user' | 'site' }).scope
        if ((!scope || scope === 'page') && !data.page) {
          throw new Error('Page data is missing from comment list response')
        }

        // Keep the public ListData contract for page-scoped consumers. Site and user scopes
        // intentionally omit page data and are used by the sidebar and message center.
        const listData = data as ListData

        // Must before all other function call and event trigger,
        // because it will depend on the lastData
        // TODO: this is global variable, easy to use, but not good, consider to refactor.
        // refactor work is hard, because it is used in many places.
        ctx.getData().setListLastFetch({ params, data: listData })

        // 装置评论
        ctx.getData().loadComments(listData.comments)

        // 更新页面数据
        if (data.page) ctx.getData().updatePage(data.page)

        // trigger events when success
        params.onSuccess && params.onSuccess(listData)

        ctx.trigger('list-fetched', { params, data: listData })
      })
      .catch((e) => {
        // 显示错误对话框
        const error = {
          msg: e.msg || String(e),
          data: e.data,
        }

        params.onError && params.onError(error)

        // trigger events when error
        ctx.trigger('list-failed', error)
        ctx.trigger('list-fetched', { params, error })

        throw e
      })
      .finally(() => {
        ctx.getData().setLoading(false)
      })
  })
}

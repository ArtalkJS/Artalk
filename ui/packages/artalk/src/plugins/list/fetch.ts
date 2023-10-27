import type ArtalkConfig from '~/types/artalk-config'
import type { ListFetchParams } from '~/types/artalk-data'
import type ContextApi from '~/types/context'
import type ArtalkPlugin from '~/types/plugin'
import { handleBackendRefConf } from '../../config'

export const Fetch: ArtalkPlugin = (ctx) => {
  ctx.on('list-fetch', (_params) => {
    if (ctx.getData().getLoading()) return

    const params: ListFetchParams = {
      // default params
      offset: 0,
      limit: ctx.conf.pagination.pageSize,
      flatMode: ctx.conf.flatMode as boolean, // always be boolean because had been handled in Artalk.init
      ..._params
    }

    // must before other function call
    ctx.getData().setListLastFetch({
      params
    })

    ctx.getApi().comment
      .get(params.offset, params.limit, params.flatMode)
      .then((data) => {
        // Must before all other function call and event trigger,
        // because it will depend on the lastData
        // TODO this is global variable, easy to use, but not good, consider to refactor
        ctx.getData().setListLastFetch({ params, data })

        // 装载后端提供的配置
        loadConf(ctx, {
          useBackendConf: ctx.conf.useBackendConf,
          conf: data.conf.frontend_conf,
          apiVersion: data.api_version.version
        })

        // 装置评论
        ctx.getData().loadComments(data.comments)

        // 更新页面数据
        ctx.getData().updatePage(data.page)

        // 未读消息提示功能
        ctx.getData().updateUnreads(data.unread || [])

        params.onSuccess && params.onSuccess(data)

        ctx.trigger('list-fetched', { params, data })
      })
      .catch((e) => {
        // 显示错误对话框
        const error = {
          msg: e.msg || String(e),
          data: e.data
        }

        params.onError && params.onError(error)

        ctx.trigger('list-error', error)
        ctx.trigger('list-fetched', { params, error })

        throw e
      })
  })

  // When list load error
  ctx.on('list-error', (err) => {
    if (!confLoaded) {
      ctx.updateConf({})
    }
  })
}

let confLoaded = false

function loadConf(ctx: ContextApi, apiRes: { useBackendConf: boolean, conf: any, apiVersion: string }) {
  if (!confLoaded) { // 仅应用一次配置
    let conf: Partial<ArtalkConfig> = {
      apiVersion: apiRes.apiVersion
    }

    // reference conf from backend
    if (ctx.conf.useBackendConf) {
      if (!apiRes.conf) throw new Error('The remote backend does not respond to the frontend conf, but `useBackendConf` conf is enabled')
      conf = { ...conf, ...handleBackendRefConf(apiRes.conf) }
    }

    ctx.updateConf(conf)
    confLoaded = true
  }
}

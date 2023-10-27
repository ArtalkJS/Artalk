import type ArtalkConfig from '~/types/artalk-config'
import type ArtalkPlugin from '~/types/plugin'
import { handleBackendRefConf } from '@/config'

export const ConfRemoter: ArtalkPlugin = (ctx) => {
  let confLoaded = false

  ctx.on('inited', () => {
    if (confLoaded) return

    ctx.getApi().system.conf().then((data) => {
      let conf: Partial<ArtalkConfig> = {
        apiVersion: data.version.version, // version info
      }

      // reference conf from backend
      if (ctx.conf.useBackendConf) {
        if (!data.frontend_conf) throw new Error('The remote backend does not respond to the frontend conf, but `useBackendConf` conf is enabled')
        conf = { ...conf, ...handleBackendRefConf(data.frontend_conf) }
      }

      ctx.updateConf(conf)

      // 评论获取
      ctx.fetch({ offset: 0 })
    }).catch((err) => {
      // TODO show error dialog
      window.alert(`[Artalk Error] ${err.msg || String(err)}`)

      ctx.updateConf({})
    }).finally(() => {
      confLoaded = true
    })
  })
}

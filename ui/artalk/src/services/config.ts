import type { Config, ArtalkPlugin, ConfigManager } from '@/types'
import { mergeDeep } from '@/lib/merge-deep'
import { handelCustomConf } from '@/config'
import { watchConf } from '@/lib/watch-conf'

export const ConfigService: ArtalkPlugin = (ctx) => {
  let conf: Config = handelCustomConf({}, true)

  ctx.provide(
    'config',
    (events) => {
      let mounted = false
      let disposed = false
      events.on('mounted', () => (mounted = true))
      events.on('unmounted', () => (disposed = true))
      const instance: ConfigManager = {
        watchConf: (keys, effect) => {
          watchConf({
            keys,
            effect,
            getConf: () => conf,
            getEvents: () => events,
          })
        },
        get: () => conf,
        update: (config) => {
          if (disposed) return
          conf = mergeDeep<Config>(conf, handelCustomConf(config, false))
          mounted && events.trigger('updated', conf)
        },
      }
      return instance
    },
    ['events'] as const,
  )
}

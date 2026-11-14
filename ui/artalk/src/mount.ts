import { handleConfFormServer } from './config'
import { DefaultPlugins } from './plugins'
import { mergeDeep } from './lib/merge-deep'
import { MountError } from './plugins/mount-error'
import { Services } from './services'
import { setLocale } from './i18n'
import type { ConfigPartial, ArtalkPlugin, Context } from '@/types'

/**
 * Global Plugins for all Artalk instances
 */
export const GlobalPlugins: Set<ArtalkPlugin> = new Set([...DefaultPlugins])

/**
 * Plugin options for plugin initialization
 */
export const PluginOptions: WeakMap<ArtalkPlugin, any> = new WeakMap()

// Keep registrations across a failed initialization and its retry.
const instances = new WeakMap<Context, { loaded: Set<ArtalkPlugin>; mounted: boolean }>()

/**
 * Load plugins and the remote config, then mount Artalk.
 *
 * Returns early without touching the root element when the instance is
 * destroyed while the remote config or network plugins are still loading.
 * The root may already be reused by a new instance at that point.
 */
export async function mount(localConf: ConfigPartial, ctx: Context) {
  if (ctx.isDestroyed()) return

  let instance = instances.get(ctx)
  if (!instance) {
    instance = { loaded: new Set(), mounted: false }
    instances.set(ctx, instance)
  }
  if (instance.mounted) return
  const { loaded } = instance
  const loadPlugins = (plugins: Set<ArtalkPlugin>) => {
    plugins.forEach((plugin) => {
      if (ctx.isDestroyed()) return
      if (typeof plugin !== 'function') return
      if (loaded.has(plugin)) return
      plugin(ctx, PluginOptions.get(plugin))
      loaded.add(plugin)
    })
  }

  // Register service providers required to retrieve the remote config. The
  // providers are lazy, so feature and external plugins are still initialized
  // only after the final config has been applied.
  loadPlugins(Services)
  setLocale(ctx.getConf().locale)

  // Get conf from server
  // The deprecated useBackendConf flag remains ignored, as in v2.
  const res = await ctx
    .getApi()
    .conf.conf()
    .catch((err) => {
      if (ctx.isDestroyed()) return null
      MountError(ctx, { err, onRetry: () => mount(localConf, ctx) })
      throw err
    })
  if (!res || ctx.isDestroyed() || instance.mounted) return

  // Merge remote and local config
  let conf: ConfigPartial = {
    ...localConf,
    apiVersion: res.data.version?.version,
  }

  const remoteConf = handleConfFormServer({ ...res.data.frontend_conf })
  conf = conf.preferRemoteConf ? mergeDeep(conf, remoteConf) : mergeDeep(remoteConf, conf)

  // Apply local + remote conf
  ctx.updateConf(conf)
  setLocale(ctx.getConf().locale)

  // Download remote plugin scripts before initialization so every plugin sees
  // the same final local + remote config.
  let networkPlugins = new Set<ArtalkPlugin>()
  if (conf.pluginURLs) {
    networkPlugins = await loadNetworkPlugins(conf.pluginURLs, ctx.getConf().server).catch(
      (err) => {
        console.error('Failed to load plugin', err)
        return new Set<ArtalkPlugin>()
      },
    )
  }

  if (ctx.isDestroyed() || instance.mounted) return

  // Initialize built-in/local plugins first, followed by network plugins.
  loadPlugins(GlobalPlugins)
  loadPlugins(networkPlugins)
  if (ctx.isDestroyed()) return

  instance.mounted = true
  ctx.trigger('created')
  if (!ctx.isDestroyed()) ctx.trigger('mounted')
}

/**
 * Dynamically load plugins from Network
 */
async function loadNetworkPlugins(scripts: string[], apiBase: string): Promise<Set<ArtalkPlugin>> {
  const networkPlugins = new Set<ArtalkPlugin>()
  if (!scripts || !Array.isArray(scripts)) return networkPlugins

  const tasks: Promise<void>[] = []

  scripts.forEach((url) => {
    // check url valid
    if (!/^(http|https):\/\//.test(url))
      url = `${apiBase.replace(/\/$/, '')}/${url.replace(/^\//, '')}`

    tasks.push(
      new Promise<void>((resolve) => {
        // check if loaded
        if (document.querySelector(`script[src="${url}"]`)) {
          resolve()
          return
        }

        // load script
        const script = document.createElement('script')
        script.src = url
        document.head.appendChild(script)
        script.onload = () => resolve()
        script.onerror = (err) => {
          console.error('[artalk] Failed to load plugin', err)
          resolve()
        }
      }),
    )
  })

  await Promise.all(tasks)

  // Read ArtalkPlugins object from window
  Object.values(window.ArtalkPlugins || {}).forEach((plugin) => {
    if (typeof plugin === 'function') networkPlugins.add(plugin)
  })

  return networkPlugins
}

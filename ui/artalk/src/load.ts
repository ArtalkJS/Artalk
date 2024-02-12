import type { ArtalkConfig, ArtalkPlugin, ContextApi } from '@/types'
import { handleConfFormServer } from '@/config'
import { showErrorDialog } from '@/components/error-dialog'
import { DefaultPlugins } from './plugins'

/**
 * Global Plugins for all Artalk instances
 */
export const GlobalPlugins: ArtalkPlugin[] = [ ...DefaultPlugins ]

export async function load(ctx: ContextApi) {
  const loadedPlugins: ArtalkPlugin[] = []
  const loadPlugins = (plugins: ArtalkPlugin[]) => {
    plugins.forEach((plugin) => {
      if (typeof plugin === 'function' && !loadedPlugins.includes(plugin)) {
        plugin(ctx)
        loadedPlugins.push(plugin)
      }
    })
  }

  // Load local plugins
  loadPlugins(GlobalPlugins)

  // Get conf from server
  const { data } = await ctx.getApi().conf.conf().catch((err) => {
    onLoadErr(ctx, err)
    throw err
  })

  // Initial config
  let conf: Partial<ArtalkConfig> = {
    apiVersion: data.version?.version, // version info
  }

  // Reference conf from backend
  if (ctx.conf.useBackendConf) {
    if (!data.frontend_conf) throw new Error('The remote backend does not respond to the frontend conf, but `useBackendConf` conf is enabled')
    conf = { ...conf, ...handleConfFormServer(data.frontend_conf) }
  }

  // Apply conf modifier
  ctx.conf.remoteConfModifier && ctx.conf.remoteConfModifier(conf)

  // Dynamically load network plugins
  conf.pluginURLs && await loadNetworkPlugins(conf.pluginURLs, ctx.conf.server).then((plugins) => {
    loadPlugins(plugins)
  }).catch((err) => {
    console.error('Failed to load plugin', err)
  })

  // After all plugins are loaded
  ctx.trigger('created')

  // Apply conf updating
  ctx.updateConf(conf)

  // Trigger mounted event
  ctx.trigger('mounted')

  // Load comment list
  if (!ctx.conf.remoteConfModifier) {  // only auto fetch when no remoteConfModifier
    ctx.fetch({ offset: 0 })
  }
}

/**
 * Dynamically load plugins from Network
 */
async function loadNetworkPlugins(scripts: string[], apiBase: string): Promise<ArtalkPlugin[]> {
  if (!scripts || !Array.isArray(scripts)) return []

  const tasks: Promise<void>[] = []

  scripts.forEach((url) => {
    // check url valid
    if (!/^(http|https):\/\//.test(url))
      url = `${apiBase.replace(/\/$/, '')}/${url.replace(/^\//, '')}`

    tasks.push(new Promise<void>((resolve, reject) => {
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
      script.onerror = (err) => reject(err)
    }))
  })

  await Promise.all(tasks)

  return Object.values(window.ArtalkPlugins || {})
}

export function onLoadErr(ctx: ContextApi, err: any) {
  let sidebarOpenView = ''

  // if response err_no_site, modify the sidebar open view to create site
  if (err.data?.err_no_site) {
    const viewLoadParam = { create_name: ctx.conf.site, create_urls: `${window.location.protocol}//${window.location.host}` }
    sidebarOpenView = `sites|${JSON.stringify(viewLoadParam)}`
  }

  showErrorDialog({
    $err: ctx.get('list').$el,
    errMsg: err.msg || String(err),
    errData: err.data,
    retryFn: () => load(ctx),
    onOpenSidebar: ctx.get('user').getData().isAdmin ? () => ctx.showSidebar({
      view: sidebarOpenView as any
    }) : undefined // only show open sidebar button when user is admin
  })
}

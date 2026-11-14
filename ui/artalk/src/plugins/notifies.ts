import type { ArtalkPlugin } from '@/types'

export const Notifies: ArtalkPlugin = (ctx) => {
  let generation = 0
  const scope = () => {
    const { server, site, pageKey } = ctx.getConf()
    const { name, email, token } = ctx.getUser().getData()
    return JSON.stringify([server, site, pageKey, name, email, token])
  }
  ctx.on('list-fetch', (params) => {
    if (params.offset !== 0) return
    const current = ++generation
    const api = ctx.getApi()
    const user = api.getUserFields()
    if (!user) {
      ctx.getData().updateNotifies([])
      return
    }
    const snapshot = scope()
    const active = () => current === generation && snapshot === scope() && !ctx.isDestroyed()
    api.notifies
      .getNotifies(user)
      .then(({ data }) => {
        if (active()) ctx.getData().updateNotifies(data.notifies)
      })
      .catch((err) => {
        if (active()) console.error(err)
      })
  })
}

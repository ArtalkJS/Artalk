/* @refresh reload */
import './style.scss'
import { ArtalkPlugin } from 'artalk'
import { DialogMain } from './DialogMain'
import { createLayer } from './lib/layer'
import { RenderEditorUser } from './EditorUser'

export const ArtalkAuthPlugin: ArtalkPlugin = (ctx) => {
  ctx.getApiHandlers().add('need_auth_login', () => {
    openAuthDialog()
    throw new Error('Login required')
  })

  let anonymous = false
  const refreshBtn = () => {
    ctx.get('editor').getUI().$submitBtn.innerText =
      ctx.get('user').getData().token || anonymous
        ? ctx.conf.sendBtn || ctx.$t('send')
        : ctx.$t('signIn')
  }

  ctx.watchConf(['locale', 'sendBtn'], () => refreshBtn())
  ctx.on('user-changed', () => refreshBtn())

  ctx.on('mounted', () => {
    ctx.get('editor').getUI().$header.style.display = 'none'

    RenderEditorUser(ctx)
  })

  const onSkip = () => {
    ctx.get('editor').getUI().$header.style.display = ''
    ctx.get('editor').getUI().$nick.focus()
    ctx.updateConf({
      beforeSubmit: undefined,
    })

    anonymous = true
    refreshBtn()
  }

  const openAuthDialog = () => {
    createLayer(ctx).show((layer) => (
      <DialogMain ctx={ctx} onClose={() => layer.destroy()} onSkip={onSkip} />
    ))
  }

  ctx.updateConf({
    beforeSubmit: (editor, next) => {
      if (!ctx.get('user').getData().token) {
        openAuthDialog()
      } else {
        next()
      }
    },
  })
}

// Mount plugin to browser window global
if (window) {
  !window.ArtalkPlugins && (window.ArtalkPlugins = {})
  window.ArtalkPlugins.Auth = ArtalkAuthPlugin
  window.Artalk?.use(ArtalkAuthPlugin)
}

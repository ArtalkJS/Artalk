import type { ArtalkPlugin, List } from '@/types'
import { version as ARTALK_VERSION } from '~/package.json'
import * as Utils from '@/lib/utils'
import $t from '@/i18n'

export const VersionCheck: ArtalkPlugin = (ctx) => {
  const list = ctx.inject('list')
  let ignored = false
  let notice: HTMLElement | undefined
  const refresh = () => {
    notice?.remove()
    notice = undefined
    const conf = ctx.getConf()
    if (ctx.getUser().getData().is_admin && conf.apiVersion && conf.versionCheck && !ignored)
      notice = versionCheck(list, ARTALK_VERSION, conf.apiVersion, () => (ignored = true))
  }
  ctx.watchConf(['apiVersion', 'versionCheck', 'locale'], refresh)
  ctx.on('user-changed', refresh)
  ctx.on('unmounted', () => notice?.remove())
}

function versionCheck(list: List, feVer: string, beVer: string, ignore: () => void) {
  const comp = Utils.versionCompare(feVer, beVer)
  const sameVer = comp === 0
  if (sameVer) return

  const errEl = Utils.createElement(
    `<div class="atk-version-check-notice">${$t('updateMsg', {
      name: comp < 0 ? $t('client') : $t('server'),
    })} <span class="atk-info">` +
      `${$t('currentVersion')}: ${$t('client')} ${feVer} / ${$t('server')} ${beVer}` +
      `</span></div>`,
  )
  const ignoreBtn = Utils.createElement(
    `<button type="button" tabindex="0" class="atk-ignore-btn">${$t('ignore')}</button>`,
  )
  ignoreBtn.onclick = () => {
    errEl.remove()
    ignore()
  }
  errEl.append(ignoreBtn)
  list.getEl().parentElement!.prepend(errEl)
  return errEl
}

import type { ArtalkPlugin, List } from '@/types'
import { version as ARTALK_VERSION } from '~/package.json'
import * as Utils from '@/lib/utils'
import $t from '@/i18n'

let IgnoreVersionCheck = false

export const VersionCheck: ArtalkPlugin = (ctx) => {
  const list = ctx.inject('list')
  ctx.watchConf(['apiVersion', 'versionCheck'], (conf) => {
    if (conf.apiVersion && conf.versionCheck && !IgnoreVersionCheck)
      versionCheck(list, ARTALK_VERSION, conf.apiVersion)
  })
}

function versionCheck(list: List, feVer: string, beVer: string) {
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
  const ignoreBtn = Utils.createElement(`<span class="atk-ignore-btn">${$t('ignore')}</span>`)
  ignoreBtn.onclick = () => {
    errEl.remove()
    IgnoreVersionCheck = true
  }
  errEl.append(ignoreBtn)
  list.getEl().parentElement!.prepend(errEl)
}

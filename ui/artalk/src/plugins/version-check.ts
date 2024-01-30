import type { ArtalkPlugin } from '@/types'
import { version as ARTALK_VERSION } from '~/package.json'
import type List from '~/src/list/list'
import * as Ui from '@/lib/ui'
import * as Utils from '@/lib/utils'
import $t from '@/i18n'

let IgnoreVersionCheck = false

export const VersionCheck: ArtalkPlugin = (ctx) => {
  ctx.watchConf(['apiVersion', 'versionCheck'], (conf) => {
    const list = ctx.get('list')
    if (conf.apiVersion && conf.versionCheck && !IgnoreVersionCheck)
      versionCheck(list, ARTALK_VERSION, conf.apiVersion)
  })
}

function versionCheck(list: List, feVer: string, beVer: string) {
  const comp = Utils.versionCompare(feVer, beVer)
  const sameVer = (comp === 0)
  if (sameVer) return

  const errEl = Utils.createElement(
    `<div>${$t('updateMsg', { name: (comp < 0) ? $t('frontend') : $t('backend') })}` +
    `<br/><br/><span style="color: var(--at-color-meta);">` +
    `${$t('currentVersion')}: ${$t('frontend')} ${feVer} / ${$t('backend')} ${beVer}` +
    `</span><br/><br/></div>`)
  const ignoreBtn = Utils.createElement(`<span style="cursor:pointer">${$t('ignore')}</span>`)
  ignoreBtn.onclick = () => {
    Ui.setError(list.$el.parentElement!, null)
    IgnoreVersionCheck = true
    list.ctx.fetch({ offset: 0 })
  }
  errEl.append(ignoreBtn)
  Ui.setError(list.$el.parentElement!, errEl, '<span class="atk-warn-title">Artalk Warn</span>')
}

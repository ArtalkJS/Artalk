import type ArtalkPlugin from '~/types/plugin'
import { version as ARTALK_VERSION } from '~/package.json'
import type ListLite from '@/list/list-lite'
import * as Ui from '@/lib/ui'
import * as Utils from '@/lib/utils'
import $t from '@/i18n'

let IgnoreVersionCheck = false

export const VersionCheck: ArtalkPlugin = (ctx) => {
  ctx.on('conf-loaded', () => {
    const list = ctx.get('list')
    if (list && ctx.conf.apiVersion && ctx.conf.versionCheck && !IgnoreVersionCheck)
      versionCheck(list, ARTALK_VERSION, ctx.conf.apiVersion)
  })
}

function versionCheck(list: ListLite, feVer: string, beVer: string) {
  const comp = Utils.versionCompare(feVer, beVer)
  const sameVer = (comp === 0)
  if (sameVer) return

  const errEl = Utils.createElement(
    `<div>请更新 Artalk ${comp < 0 ? $t('frontend') : $t('backend')}以获得完整体验 ` +
    `(<a href="https://artalk.js.org/" target="_blank">帮助文档</a>)` +
    `<br/><br/><span style="color: var(--at-color-meta);">` +
    `当前版本：${$t('frontend')} ${feVer} / ${$t('backend')} ${beVer}` +
    `</span><br/><br/></div>`)
  const ignoreBtn = Utils.createElement('<span style="cursor:pointer">忽略</span>')
  ignoreBtn.onclick = () => {
    Ui.setError(list.$el.parentElement!, null)
    IgnoreVersionCheck = true
    list.fetchComments(0)
  }
  errEl.append(ignoreBtn)
  Ui.setError(list.$el.parentElement!, errEl, '<span class="atk-warn-title">Artalk Warn</span>')
}

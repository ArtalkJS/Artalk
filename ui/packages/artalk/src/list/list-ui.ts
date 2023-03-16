import ListLite from './list-lite'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import User from '../lib/user'
import $t from '../i18n'

/** 界面刷新 */
export function refreshUI(list: ListLite) {
    // 无评论
    const isNoComment = list.ctx.getCommentList().length <= 0
    let $noComment = list.getCommentsWrapEl().querySelector<HTMLElement>('.atk-list-no-comment')

    if (isNoComment) {
      if (!$noComment) {
        $noComment = Utils.createElement('<div class="atk-list-no-comment"></div>')
        $noComment.innerHTML = list.noCommentText || list.ctx.conf.noComment || list.ctx.$t('noComment')
        list.getCommentsWrapEl().appendChild($noComment)
      }
    } else {
      $noComment?.remove()
    }

    // 仅管理员显示控制
    list.ctx.checkAdminShowEl()
}

/** 构建错误提示界面 */
export function renderErrorDialog(list: ListLite, errMsg: string, errData?: any): HTMLElement {
  const errEl = Utils.createElement(`<span>${errMsg}，${$t('listLoadFailMsg')}<br/></span>`)

  const $retryBtn = Utils.createElement(`<span style="cursor:pointer;">${$t('listRetry')}</span>`)
  $retryBtn.onclick = () => (list.fetchComments(0))
  errEl.appendChild($retryBtn)

  const adminBtn = Utils.createElement('<span atk-only-admin-show> | <span style="cursor:pointer;">打开控制台</span></span>')
  errEl.appendChild(adminBtn)
  if (!User.data.isAdmin) adminBtn.classList.add('atk-hide')

  let sidebarView = ''

  // 找不到站点错误，打开侧边栏并填入创建站点表单
  if (errData?.err_no_site) {
    const viewLoadParam = {
      create_name: list.ctx.conf.site,
      create_urls: `${window.location.protocol}//${window.location.host}`
    }
    sidebarView = `sites|${JSON.stringify(viewLoadParam)}`
  }

  adminBtn.onclick = () => list.ctx.showSidebar({
    view: sidebarView as any
  })

  return errEl
}

/** 版本检测弹窗 */
export function versionCheckDialog(list: ListLite, feVer: string, beVer: string): boolean {
  const comp = Utils.versionCompare(feVer, beVer)
  const notSameVer = (comp !== 0)
  if (notSameVer) {
    const errEl = Utils.createElement(
      `<div>请更新 Artalk ${comp < 0 ? $t('frontend') : $t('backend')}以获得完整体验 ` +
      `(<a href="https://artalk.js.org/" target="_blank">帮助文档</a>)` +
      `<br/><br/><span style="color: var(--at-color-meta);">` +
      `当前版本：${$t('frontend')} ${feVer} / ${$t('backend')} ${beVer}` +
      `</span><br/><br/></div>`)
    const ignoreBtn = Utils.createElement('<span style="cursor:pointer">忽略</span>')
    ignoreBtn.onclick = () => {
      Ui.setError(list.$el.parentElement!, null)
      list.ctx.conf.versionCheck = false
      list.fetchComments(0)
    }
    errEl.append(ignoreBtn)
    Ui.setError(list.$el.parentElement!, errEl, '<span class="atk-warn-title">Artalk Warn</span>')
    return true
  }

  return false
}

/** 评论排序方式选择下拉菜单 */
export function renderDropdown(conf: {
  $dropdownWrap: HTMLElement,
  dropdownList: [string, () => void][]
}) {
  const { $dropdownWrap, dropdownList } = conf
  if ($dropdownWrap.querySelector('.atk-dropdown')) return

  // 修改 class
  $dropdownWrap.classList.add('atk-dropdown-wrap')

  // 插入图标
  $dropdownWrap.append(Utils.createElement(`<span class="atk-arrow-down-icon"></span>`))

  // 列表项点击事件
  let curtActive = 0 // 当前选中
  const onItemClick = (i: number, $item: HTMLElement, name: string, action: Function) => {
    action()

    // set active
    curtActive = i
    $dropdown.querySelectorAll('.active').forEach((e) => { e.classList.remove('active') })
    $item.classList.add('active')

    // 关闭层 (临时消失，取消 :hover)
    $dropdown.style.display = 'none'
    setTimeout(() => { $dropdown.style.display = '' }, 80)
  }

  // 生成列表元素
  const $dropdown = Utils.createElement(`<ul class="atk-dropdown atk-fade-in"></ul>`)
  dropdownList.forEach((item, i) => {
    const name = item[0] as string
    const action = item[1] as Function

    const $item = Utils.createElement(`<li class="atk-dropdown-item"><span></span></li>`)
    const $link = $item.querySelector<HTMLElement>('span')!
    $link.innerText = name
    $link.onclick = () => { onItemClick(i, $item, name, action) }
    $dropdown.append($item)

    if (i === curtActive) $item.classList.add('active') // 默认选中项
  })

  $dropdownWrap.append($dropdown)
}

/** 删除评论排序方式选择下拉菜单 */
export function removeDropdown(conf: {
  $dropdownWrap: HTMLElement
}) {
  const { $dropdownWrap } = conf
  $dropdownWrap.classList.remove('atk-dropdown-wrap')
  $dropdownWrap.querySelector('.atk-arrow-down-icon')?.remove()
  $dropdownWrap.querySelector('.atk-dropdown')?.remove()
}

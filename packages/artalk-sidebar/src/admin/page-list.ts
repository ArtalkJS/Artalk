import '../style/page-list.less'

import Context from 'artalk/src/context'
import * as Utils from 'artalk/src/lib/utils'
import * as Ui from 'artalk/src/lib/ui'
import { PageData } from 'artalk/types/artalk-data'
import Api from 'artalk/src/api'
import ActionBtn from 'artalk/src/components/action-btn'
import Component from '../sidebar-component'
import { SidebarCtx } from '../main'
import ItemTextEditor from '../lib/item-text-editor'

export default class PageList extends Component {
  $editor?: HTMLElement
  $inputer?: HTMLElement

  pages: PageData[] = []

  constructor(ctx: Context, sidebar: SidebarCtx) {
    super(ctx, sidebar)

    this.$el = Utils.createElement(`<div class="atk-page-list"></div>`)
  }

  /** 初始化 PageList (清空列表) */
  initPageList() {
    this.pages = []

    this.$el.innerHTML = `<div class="atk-header-action-bar">
    <span class="atk-update-all-title-btn"><i class="atk-icon atk-icon-sync"></i> <span class="atk-text">更新标题</span></span>
    <span class="atk-cache-flush-all-btn"><span class="atk-text">缓存清除</span></span>
    <span class="atk-cache-warm-up-btn"><span class="atk-text">缓存预热</span></span>
    </div>`

    // 缓存操作按钮
    const $cacheFlushBtn = this.$el.querySelector<HTMLElement>('.atk-cache-flush-all-btn')!
    $cacheFlushBtn.onclick = () => { new Api(this.ctx).cacheFlushAll().then((d: any) => alert(d.msg)).catch(() => alert('操作失败')) }
    const $cacheWarmBtn = this.$el.querySelector<HTMLElement>('.atk-cache-warm-up-btn')!
    $cacheWarmBtn.onclick = () => { new Api(this.ctx).cacheWarmUp().then((d: any) => alert(d.msg)).catch(() => alert('操作失败')) }

    // 更新全部页面标题按钮
    ;(async () => {
      const $updateAllTitleBtn = this.$el.querySelector<HTMLElement>('.atk-update-all-title-btn')!
      const $updateAllTitleIcon = $updateAllTitleBtn.querySelector<HTMLElement>('i')!
      const $updateAllTitleText = $updateAllTitleBtn.querySelector<HTMLElement>('.atk-text')!
      const btnTextOrg = $updateAllTitleText.innerText

      const done = () => {
        $updateAllTitleText.innerText = '更新完毕'
        $updateAllTitleIcon.classList.remove('atk-rotate')
        window.setTimeout(() => {
          $updateAllTitleText.innerText = btnTextOrg
        }, 1500)
      }
      const checkStatus = async () => {
        const status = await new Api(this.ctx).pageFetch(undefined, undefined, true)
        return status
      }
      const startWatchdog = () => {
        $updateAllTitleIcon.classList.add('atk-rotate')

        // 不完美的轮询更新状态
        const timerID = window.setInterval(async () => {
          const d = await new Api(this.ctx).pageFetch(undefined, undefined, true)

          if (d.is_progress === false) {
            clearInterval(timerID)
            done()
            return
          }

          $updateAllTitleText.innerText = d.msg
        }, 1000)
      }

      // 图标恢复
      const status = await checkStatus()
      if (status.is_progress === true) {
        startWatchdog()
        $updateAllTitleText.innerText = status.msg
      }

      // 点击发起任务
      $updateAllTitleBtn.onclick = async () => {
        if ($updateAllTitleIcon.classList.contains('atk-rotate')) return
        startWatchdog()
        $updateAllTitleText.innerText = '开始更新...'

        // 发起任务
        try {
          await new Api(this.ctx).pageFetch(undefined, this.sidebar.curtSite)
        } catch (err: any) {
          alert(err.msg)
          done()
        }
      }
    })()
  }

  /** 导入 Page 数据 */
  importPages(pages: PageData[]) {
    this.pages.push(...pages)

    pages.forEach((page) => {
      const $page = this.renderPage(page)
      this.$el.append($page)
    })
  }

  /** Page Element */
  renderPage(page: PageData) {
    const $page = Utils.createElement(
      `<div class="atk-page-item">
        <div class="atk-page-main">
          <div class="atk-title"></div>
          <div class="atk-sub"></div>
        </div>
        <div class="atk-page-actions">
          <div class="atk-item atk-edit-btn">
            <i class="atk-icon atk-icon-edit"></i>
          </div>
        </div>
      </div>`)

    const $main = $page.querySelector<HTMLElement>('.atk-page-main')!
    const $title = $main.querySelector<HTMLElement>('.atk-title')!
    const $sub = $main.querySelector<HTMLElement>('.atk-sub')!
    const $editBtn = $page.querySelector<HTMLElement>('.atk-edit-btn')!

    $title.innerText = page.title
    $sub.innerText = page.url || page.key
    $editBtn.onclick = () => this.showEditor(page, $page)

    $title.onclick = () => { if (page.url) window.open(page.url) }
    $sub.onclick = () => { if (page.url) window.open(page.url) }

    return $page
  }

  /** 显示编辑器 */
  showEditor(page: PageData, $page: HTMLElement) {
    this.closeEditor()

    this.$editor = Utils.createElement(
    `<div class="atk-page-edit-layer">
      <div class="atk-page-main-actions">
        <div class="atk-item atk-title-edit-btn">标题修改</div>
        <div class="atk-item atk-key-edit-btn">KEY 变更</div>
        <div class="atk-item atk-admin-only-btn"></div>
      </div>
      <div class="atk-page-actions">
        <div class="atk-item atk-sync-btn">
          <i class="atk-icon atk-icon-sync"></i>
        </div>
        <div class="atk-item atk-del-btn">
          <i class="atk-icon atk-icon-del"></i>
        </div>
        <div class="atk-item atk-close-btn">
          <i class="atk-icon atk-icon-close"></i>
        </div>
      </div>
    </div>`)
    $page.prepend(this.$editor)

    const $titleEditBtn = this.$editor.querySelector<HTMLElement>('.atk-title-edit-btn')!
    const $keyEditBtn = this.$editor.querySelector<HTMLElement>('.atk-key-edit-btn')!
    const $adminOnlyBtn = this.$editor.querySelector<HTMLElement>('.atk-admin-only-btn')!
    const $syncBtn = this.$editor.querySelector<HTMLElement>('.atk-sync-btn')!
    const $delBtn = this.$editor.querySelector<HTMLElement>('.atk-del-btn')!
    const $closeBtn = this.$editor.querySelector<HTMLElement>('.atk-close-btn')!
    const showLoading = () => { Ui.showLoading(this.$editor!) }
    const hideLoading = () => { Ui.hideLoading(this.$editor!) }
    const showError = (msg: string) => { window.alert(msg) }

    // 关闭编辑器
    $closeBtn.onclick = () => this.closeEditor()

    // 文本编辑
    const openTextEditor = (key: string) => {
      const textEditor = new ItemTextEditor({
        initValue: page[key] || '',
        onYes: async (val: string) => {
          Ui.showLoading(textEditor.$el)
          let p: PageData
          try {
            p = await new Api(this.ctx).pageEdit({ ...page, [key]: val })
          } catch (err: any) {
            showError(`修改失败：${err.msg || '未知错误'}`)
            console.error(err)
            return false
          } finally { Ui.hideLoading(textEditor.$el) }
          $page.replaceWith(this.renderPage(p))
          return true
        }
      })
      textEditor.appendTo(this.$editor!)
    }

    // 标题修改
    $titleEditBtn.onclick = () => openTextEditor('title')

    // Key 修改
    $keyEditBtn.onclick = () => openTextEditor('key')

    // 仅管理员可评
    const adminOnlyActionBtn = new ActionBtn({
      text: () => {
        $adminOnlyBtn.classList.remove('atk-green', 'atk-yellow')
        $adminOnlyBtn.classList.add(!page.admin_only ? 'atk-green' : 'atk-yellow')
        return (!page.admin_only) ? '所有人可评' : '管理员可评'
      }
    }).appendTo($adminOnlyBtn)

    $adminOnlyBtn.onclick = async () => {
      showLoading()
      let p: PageData
      try {
        p = await new Api(this.ctx).pageEdit({ ...page, admin_only: !page.admin_only })
      } catch (err: any) {
        showError(`修改失败：${err.msg || '未知错误'}`)
        console.log(err)
        return
      } finally { hideLoading() }
      page.admin_only = p.admin_only
      adminOnlyActionBtn.updateText()
    }

    // 同步操作
    $syncBtn.onclick = async () => {
      showLoading()
      let p: PageData
      try {
        p = (await new Api(this.ctx).pageFetch(page.id)).page
      } catch (err: any) {
        showError(`同步失败：${err.msg || '未知错误'}`)
        console.log(err)
        return
      } finally { hideLoading() }
      $page.replaceWith(this.renderPage(p))
    }

    // 删除
    $delBtn.onclick = () => {
      const del = async () => {
        showLoading()
        try {
          await new Api(this.ctx).pageDel(page.key, page.site_name)
        } catch (err: any) {
          console.log(err)
          showError(`删除失败 ${String(err)}`)
          return
        } finally { hideLoading() }
        $page.remove()
      }
      if (window.confirm(
        `确认删除页面 "${page.title || page.key}"？将会删除所有相关数据`
      )) del()
    }
  }

  /** 关闭编辑器 */
  public closeEditor() {
    if (!this.$editor) return

    this.$editor.remove()
  }
}

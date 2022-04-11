import '@/style/page-list.less'

import Context from '@/context'
import Component from '@/lib/component'
import * as Utils from '@/lib/utils'
import * as Ui from '@/lib/ui'
import { PageData } from '~/types/artalk-data'
import ItemTextEditor from '../item-text-editor'
import Api from '~/src/api'
import ActionBtn from '../action-btn'

export default class PageList extends Component {
  $editor?: HTMLElement
  $inputer?: HTMLElement

  pages: PageData[] = []

  constructor(ctx: Context) {
    super(ctx)

    this.$el = Utils.createElement(`<div class="atk-page-list"></div>`)
  }

  /** 清空所有 Pages */
  clearAll() {
    this.pages = []
    this.$el.innerHTML = ''
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
        p = await new Api(this.ctx).pageFetch(page.id)
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

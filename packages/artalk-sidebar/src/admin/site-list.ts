import '../style/site-list.less'

import Context from 'artalk/src/context'
import Component from 'artalk/src/lib/component'
import * as Utils from 'artalk/src/lib/utils'
import * as Ui from 'artalk/src/lib/ui'
import { SiteData } from 'artalk/types/artalk-data'
import Api from 'artalk/src/api'
import ItemTextEditor from '../item-text-editor'

export default class SiteList extends Component {
  sites: SiteData[] = []

  $header: HTMLElement
  $headerTitle: HTMLElement
  $headerActions: HTMLElement
  $rowsWrap: HTMLElement

  $editor?: HTMLElement
  activeSite: string = ''

  $add?: HTMLElement

  constructor(ctx: Context) {
    super(ctx)

    this.$el = Utils.createElement(
    `<div class="atk-site-list">
      <div class="atk-header">
        <div class="atk-title"></div>
        <div class="atk-actions">
          <div class="atk-item atk-site-add-btn"><i class="atk-icon atk-icon-plus"></i></div>
        </div>
      </div>
      <div class="atk-site-rows-wrap"></div>
    </div>`)

    this.$header = this.$el.querySelector('.atk-header')!
    this.$headerTitle = this.$header.querySelector('.atk-title')!
    this.$headerActions = this.$header.querySelector('.atk-actions')!
    this.$rowsWrap = this.$el.querySelector('.atk-site-rows-wrap')!

    // 标题
    this.$headerTitle.innerText = `共 0 个站点`

    // 新建站点按钮
    const $addBtn = this.$headerActions.querySelector<HTMLElement>('.atk-site-add-btn')!
    $addBtn.onclick = () => {
      this.closeEditor()
      this.showAdd()
    }
  }

  public loadSites(sites: SiteData[]) {
    this.sites = sites
    this.activeSite = ''
    this.$rowsWrap.innerHTML = ''
    this.$headerTitle.innerText = `共 0 个站点`

    let $row: HTMLElement

    for (let i = 0; i < sites.length; i++) {
      const site = sites[i]

      if (i % 4 === 0) { // 每 4 个 item 为一行
        $row = Utils.createElement('<div class="atk-site-row">')
        this.$rowsWrap.append($row)
      }

      // 创建 site item
      const $site = this.renderSite(site, $row!)
      $row!.append($site)
    }

    this.$headerTitle.innerText = `共 ${sites.length} 个站点`
  }

  public renderSite(site: SiteData, $row: HTMLElement) {
    const $site = Utils.createElement(
      `<div class="atk-site-item">
        <div class="atk-site-logo"></div>
        <div class="atk-site-name"></div>
      </div>`)

    const $siteLogo = $site.querySelector<HTMLElement>('.atk-site-logo')!
    const $siteName = $site.querySelector<HTMLElement>('.atk-site-name')!
    const setActive = () => { $site.classList.add('atk-active') }

    $siteLogo.innerText = site.name.substr(0, 1)
    $siteName.innerText = site.name

    // click
    $site.onclick = () => {
      this.closeEditor()
      this.closeAdd()
      setActive()
      this.showEditor(site, $site, $row)
    }

    if (this.activeSite === site.name) { setActive() }

    return $site
  }

  public showEditor(site: SiteData, $site: HTMLElement, $row: HTMLElement) {
    this.activeSite = site.name
    this.$editor = Utils.createElement(`
    <div class="atk-site-edit">
    <div class="atk-header">
      <div class="atk-site-info">
        <span class="atk-site-name"></span>
        <span class="atk-site-urls"></span>
      </div>
      <div class="atk-close-btn">
        <i class="atk-icon atk-icon-close"></i>
      </div>
    </div>
    <div class="atk-main">
      <div class="atk-site-text-actions">
        <div class="atk-item atk-rename-btn">重命名</div>
        <div class="atk-item atk-edit-url-btn">修改 URL</div>
        <!--<div class="atk-item atk-export-btn">导出</div>
        <div class="atk-item atk-import-btn">导入</div>-->
      </div>
      <div class="atk-site-btn-actions">
        <div class="atk-item atk-del-btn">
          <i class="atk-icon atk-icon-del"></i>
        </div>
      </div>
    </div>
    </div>`)

    // 插入
    $row.before(this.$editor)

    // header
    const $siteName = this.$editor.querySelector<HTMLElement>('.atk-site-name')!
    const $siteUrls = this.$editor.querySelector<HTMLElement>('.atk-site-urls')!
    const $closeBtn = this.$editor.querySelector<HTMLElement>('.atk-close-btn')!
    $closeBtn.onclick = () => this.closeEditor()

    const update = (s: SiteData) => {
      site = s
      $siteName.innerText = site.name
      $siteName.onclick = () => { if (site.first_url) window.open(site.first_url) }
      $siteUrls.innerHTML = ''
      site.urls?.forEach(u => {
        const $item = Utils.createElement('<span class="atk-url-item"></span>')
        $siteUrls.append($item)
        $item.innerText = (u || '').replace(/\/$/, '')
        $item.onclick = () => { window.open(u) }
      })
    }
    update(site)

    // actions
    const $main = this.$editor.querySelector<HTMLElement>('.atk-main')!
    const $actions = this.$editor.querySelector<HTMLElement>('.atk-site-text-actions')!
    const $renameBtn = $actions.querySelector<HTMLElement>('.atk-rename-btn')!
    const $editUrlBtn = $actions.querySelector<HTMLElement>('.atk-edit-url-btn')!
    // const $exportBtn = $actions.querySelector<HTMLElement>('.atk-export-btn')!
    // const $importBtn = $actions.querySelector<HTMLElement>(`.atk-import-btn`)!
    const $delBtn = this.$editor.querySelector<HTMLElement>('.atk-del-btn')!
    const showLoading = () => { Ui.showLoading(this.$editor!) }
    const hideLoading = () => { Ui.hideLoading(this.$editor!) }
    const showError = (msg: string) => { window.alert(msg) }

    // 文本编辑
    const openTextEditor = (key: string) => {
      let initValue = site[key] || ''
      if (key === 'urls') initValue = site.urls_raw || ''
      const textEditor = new ItemTextEditor({
        initValue,
        onYes: async (val: string) => {
          Ui.showLoading(textEditor.$el)
          let s: SiteData
          try {
            s = await new Api(this.ctx).siteEdit({ ...site, [key]: val })
          } catch (err: any) {
            showError(`修改失败：${err.msg || '未知错误'}`)
            console.error(err)
            return false
          } finally { Ui.hideLoading(textEditor.$el) }
          $site.replaceWith(this.renderSite(s, $row))
          update(s)
          return true
        }
      })
      textEditor.appendTo($main)
    }

    // 重命名
    $renameBtn.onclick = () => openTextEditor('name')

    // 修改 URL
    $editUrlBtn.onclick = () => openTextEditor('urls')

    // 导出
    // $exportBtn.onclick = () => {}

    // // 导入
    // $importBtn.onclick = () => {}

    // 删除
    $delBtn.onclick = () => {
      const del = async () => {
        showLoading()
        try {
          await new Api(this.ctx).siteDel(site.id, true)
        } catch (err: any) {
          console.log(err)
          showError(`删除失败 ${String(err)}`)
          return
        } finally { hideLoading() }
        this.closeEditor()
        $site.remove()
        this.sites = this.sites.filter(s => s.name !== site.name)
      }
      if (window.confirm(
        `确认删除站点 "${site.name}"？将会删除所有相关数据`
      )) del()
    }
  }

  public closeEditor() {
    if (!this.$editor) return

    this.$editor.remove()
    this.$rowsWrap.querySelectorAll('.atk-site-item').forEach((e) => e.classList.remove('atk-active'))
    this.activeSite = ''
  }

  public showAdd() {
    this.closeAdd()

    this.$add = Utils.createElement(`
    <div class="atk-site-add">
    <div class="atk-header">
      <div class="atk-title">新增站点</div>
      <div class="atk-close-btn">
        <i class="atk-icon atk-icon-close"></i>
      </div>
    </div>
    <div class="atk-form">
      <input type="text" name="AtkSiteName" placeholder="站点名称" autocomplete="off">
      <input type="text" name="AtkSiteUrls" placeholder="站点 URL（多个用逗号隔开）" autocomplete="off">
      <button class="atk-btn" name="AtkSubmit">创建</button>
    </div>
    </div>`)
    this.$header.after(this.$add)

    const $closeBtn = this.$add.querySelector<HTMLElement>('.atk-close-btn')!
    $closeBtn.onclick = () => this.closeAdd()

    const $siteName = this.$add.querySelector<HTMLInputElement>('[name="AtkSiteName"]')!
    const $siteUrls = this.$add.querySelector<HTMLInputElement>('[name="AtkSiteUrls"]')!
    const $submitBtn = this.$add.querySelector<HTMLButtonElement>('[name="AtkSubmit"]')!

    $submitBtn.onclick = async () => {
      const siteName = $siteName.value.trim()
      const siteUrls = $siteUrls.value.trim()

      if (siteName === '') { $siteName.focus(); return }

      Ui.showLoading(this.$add!)
      let s: SiteData
      try {
        s = await new Api(this.ctx).siteAdd(siteName, siteUrls)
      } catch (err: any) {
        window.alert(`创建失败：${err.msg || ''}`)
        console.error(err)
        return
      } finally {  Ui.hideLoading(this.$add!) }

      this.sites.push(s)
      this.loadSites(this.sites)
      this.closeAdd()
    }

    // 回车键提交
    const keyDown = (evt: KeyboardEvent) => { if (evt.key === 'Enter') { $submitBtn.click() } }
    $siteName.onkeyup = (evt) => keyDown(evt)
    $siteUrls.onkeyup = (evt) => keyDown(evt)
  }

  public closeAdd() {
    this.$add?.remove()
  }
}

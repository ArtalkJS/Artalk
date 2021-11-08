import '@/style/site-list.less'

import Context from '@/context'
import Component from '@/lib/component'
import * as Utils from '@/lib/utils'
import * as Ui from '@/lib/ui'
import { SiteData } from '~/types/artalk-data'

export default class SiteList extends Component {
  sites: SiteData[] = []

  $header: HTMLElement
  $headerTitle: HTMLElement
  $headerActions: HTMLElement
  $rowsWrap: HTMLElement

  $editor?: HTMLElement

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
    $addBtn.onclick = () => {}
  }

  public loadSites(sites: SiteData[]) {
    this.sites = sites
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
      const $site = Utils.createElement(
        `<div class="atk-site-item">
          <div class="atk-site-logo"></div>
          <div class="atk-site-name"></div>
        </div>`)
      $row!.append($site)

      const $siteLogo = $site.querySelector<HTMLElement>('.atk-site-logo')!
      const $siteName = $site.querySelector<HTMLElement>('.atk-site-name')!

      $siteLogo.innerText = site.name.substr(0, 1)
      $siteName.innerText = site.name

      // click
      $site.onclick = () => {
        this.closeEditor()
        $site.classList.add('atk-active')
        this.editSite(site, $site)
      }
    }

    this.$headerTitle.innerText = `共 ${sites.length} 个站点`
  }

  public editSite(site: SiteData, $site: HTMLElement) {
    this.$editor = Utils.createElement(`
    <div class="atk-site-edit">
    <div class="atk-header">
      <div class="atk-site-name">Site</div>
      <div class="atk-close-btn">
        <i class="atk-icon atk-icon-close"></i>
      </div>
    </div>
    <div class="atk-main">
      <div class="atk-site-text-actions">
        <div class="atk-item atk-rename-btn">重命名</div>
        <div class="atk-item atk-edit-url-btn">修改 URL</div>
        <div class="atk-item atk-export-btn">导出</div>
        <div class="atk-item atk-import-btn">导入</div>
      </div>
      <div class="atk-site-btn-actions">
        <div class="atk-item atk-del-btn">
          <i class="atk-icon atk-icon-del"></i>
        </div>
      </div>
    </div>
    </div>`)

    // 插入
    const $row = $site.parentElement!
    $row.before(this.$editor)

    // header
    this.$editor.querySelector<HTMLElement>('.atk-site-name')!.innerText = site.name
    this.$editor.querySelector<HTMLElement>('.atk-close-btn')!.onclick = () => {
      this.closeEditor()
    }

    // actions
    const $actions = this.$editor.querySelector<HTMLElement>('.atk-site-text-actions')!
    const $renameBtn = $actions.querySelector<HTMLElement>('.atk-rename-btn')!
    const $editUrlBtn = $actions.querySelector<HTMLElement>('.atk-edit-url-btn')!
    const $exportBtn = $actions.querySelector<HTMLElement>('.atk-export-btn')!
    const $importBtn = $actions.querySelector<HTMLElement>(`.atk-import-btn`)!
    const $delBtn = this.$editor.querySelector<HTMLElement>('.atk-del-btn')!

    // 重命名
    $renameBtn.onclick = () => {}

    // 修改 URL
    $editUrlBtn.onclick = () => {}

    // 导出
    $exportBtn.onclick = () => {}

    // 导入
    $importBtn.onclick = () => {}

    // 删除
    $delBtn.onclick = () => {}
  }

  public closeEditor() {
    if (!this.$editor) return

    this.$editor.remove()
    this.$rowsWrap.querySelectorAll('.atk-site-item').forEach((e) => e.classList.remove('atk-active'))
  }
}

import '@/style/page-list.less'

import Context from '@/context'
import Component from '@/lib/component'
import * as Utils from '@/lib/utils'
import * as Ui from '@/lib/ui'
import { PageData } from '~/types/artalk-data'

export default class PageList extends Component {
  pages: PageData[] = []

  $editor?: HTMLElement
  $inputer?: HTMLElement

  constructor(ctx: Context) {
    super(ctx)

    this.$el = Utils.createElement(`<div class="atk-page-list"></div>`)
  }

  importPages(pages: PageData[]) {
    this.pages.push(...pages)

    pages.forEach((page) => {
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
      this.$el.append($page)

      const $main = $page.querySelector<HTMLElement>('.atk-page-main')!
      const $title = $main.querySelector<HTMLElement>('.atk-title')!
      const $sub = $main.querySelector<HTMLElement>('.atk-sub')!
      const $editBtn = $page.querySelector<HTMLElement>('.atk-edit-btn')!

      $title.innerText = page.title
      $sub.innerText = page.url || page.key
      $editBtn.onclick = () => {
        this.editPage(page, $page)
      }
    })
  }

  clearAll() {
    this.pages = []
    this.$el.innerHTML = ''
  }

  editPage(page: PageData, $page: HTMLElement) {
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

    $titleEditBtn.onclick = () => {}
    $keyEditBtn.onclick = () => {}

    $adminOnlyBtn.classList.add(!page.admin_only ? 'atk-green' : 'atk-yellow')
    $adminOnlyBtn.innerText = (!page.admin_only) ? '所有人可评' : '仅管理员可评'
    $adminOnlyBtn.onclick = () => {}

    $syncBtn.onclick = () => {}
    $delBtn.onclick = () => {}
    $closeBtn.onclick = () => { this.closeEditor() }
  }

  public closeEditor() {
    if (!this.$editor) return

    this.$editor.remove()
  }
}

import Context from '../../Context'
import SidebarView from './SidebarView'
import * as Utils from '../../lib/utils'
import Comment from '../Comment'
import ListLite from '../ListLite'
import { CreateCommentList } from './MessageView'
import Api from '~/src/lib/api'

export default class AdminView extends SidebarView {
  name = 'admin'
  title = '控制台'
  actions = {
    comment: '评论',
    page: '页面',
    site: '站点',
    conf: '配置',
  }
  activeAction = ''

  cListInstance?: ListLite
  get cList () {
    if (!this.cListInstance) this.cListInstance = CreateCommentList(this.ctx)
    return this.cListInstance
  }

  constructor (ctx: Context) {
    super(ctx)

    this.el = Utils.createElement(`<div class="atk-msg-center"></div>`)
  }

  init () {
    this.activeAction = 'comment'
    this.switch('comment')
  }

  switch (action: string) {
    this.el.innerHTML = ''

    if (action === 'comment') {
      this.initCommentList()
    } else if (action === 'page') {
      this.initPageList()
    } else if (action === 'site') {
      this.initSiteList()
    } else if (action === 'setting') {
      this.initSetting()
    }
  }

  initCommentList() {
    this.el.innerHTML = ''
    this.el.append(this.cList.el)

    const reqComments = (type, siteName) => {
      if (!this.cList) return
      this.cList.type = `admin_${type}` as any
      this.cList.isFirstLoad = true
      this.cList.paramsEditor = (params) => {
        params.site_name = siteName
      }
      this.cList.reqComments()
    }

    let loaded = false
    const conf = { typeName: 'all', siteName: '' }
    this.showFilterBar({ all: '全部', pending: '待审' }, (typeName) => {
      if (!loaded) return
      conf.typeName = typeName
      reqComments(conf.typeName, conf.siteName)
    })

    // 初始化 site 列表
    new Api(this.ctx).siteGet()
      .then((sites) => {
        const items: {[name: string]: string} = {'': '所有站点'}
        sites.forEach((site) => { items[site.name] = site.name })

        this.showFilterBar(items, (s) => {
          conf.siteName = s
          reqComments(conf.typeName, conf.siteName)
        })
      })
      .finally(() => {
        loaded = true
      })
  }

  async initPageList () {
    this.el.innerHTML = ''

    const pListEl = Utils.createElement(`<div class="atk-sidebar-list"></div>`)
    this.el.append(pListEl)
    const pages = await new Api(this.ctx).pageGet()
    pages.forEach(page => {
      const pageItemEl = Utils.createElement(`
      <div class="atk-item">
      <div class="atk-title"></div>
      <div class="atk-sub"></div>
      <div class="atk-actions">
        <div class="atk-item" data-action="page-admin-only"></div>
        <div class="atk-item" data-action="page-del">删除</div>
      </div>
      </div>`)
      pListEl.append(pageItemEl)

      pageItemEl.setAttribute('data-page-id', String(page.id))

      const titleEl = pageItemEl.querySelector<HTMLElement>('.atk-title')!
      const keyEl = pageItemEl.querySelector<HTMLElement>('.atk-sub')!

      titleEl.innerText = page.title || page.key
      keyEl.innerText = page.key
      titleEl.onclick = () => {
        window.open(`${page.url}`)
      }
      keyEl.onclick = () => {
        window.open(`${page.url}`)
      }

      const adminOnlyBtn = pageItemEl.querySelector<HTMLElement>('[data-action="page-admin-only"]')!
      const renderAdminOnlyBtn = () => {
        adminOnlyBtn.innerText =  page.admin_only ? '仅管理员可评' : '所有人可评'
      }
      renderAdminOnlyBtn()
      adminOnlyBtn.onclick = (e) => {
        e.stopPropagation() // 防止穿透

        // TODO: loading ui
        new Api(this.ctx).pageEdit(page.key, {
          adminOnly: !page.admin_only,
        })
        .then((p) => {
          page = p
          renderAdminOnlyBtn()
        })
      }

      const delBtn = pageItemEl.querySelector<HTMLElement>('[data-action="page-del"]')!
      delBtn.onclick = () => {
        const del = () => {
          new Api(this.ctx).pageDel(page.key, page.site_name)
          .then(success => {
            if (success) pageItemEl.remove()
          })
        }
        if (window.confirm(`确认删除页面 "${page.title || page.key}"？将会删除所有相关数据`)) del()
      }
    })

  }

  async initSiteList () {
    // TODO: 可复用，特别是 actions
    this.el.innerHTML = ''
    const sListEl = Utils.createElement(`<div class="atk-sidebar-list"></div>`)
    this.el.append(sListEl)
    const sites = await new Api(this.ctx).siteGet()
    sites.forEach(site => {
      const siteItemEl = Utils.createElement(`
      <div class="atk-item">
      <div class="atk-title"></div>
      <div class="atk-sub"></div>
      <div class="atk-actions">
        <div class="atk-item" data-action="site-del">删除</div>
      </div>
      </div>`)
      sListEl.append(siteItemEl)

      siteItemEl.setAttribute('data-site-id', String(site.id))

      const nameEl = siteItemEl.querySelector<HTMLElement>('.atk-title')!
      const urlEl = siteItemEl.querySelector<HTMLElement>('.atk-sub')!

      nameEl.innerText = site.name || site.url
      urlEl.innerText = site.url
      nameEl.onclick = () => {
        window.open(`${site.url}`)
      }
      urlEl.onclick = () => {
        window.open(`${site.url}`)
      }

      const delBtn = siteItemEl.querySelector<HTMLElement>('[data-action="site-del"]')!
      delBtn.onclick = () => {
        const del = () => {
          new Api(this.ctx).siteDel(site.id, true)
          .then(success => {
            if (success) siteItemEl.remove()
          })
        }
        if (window.confirm(`确认删除站点 "${site.name || site.url}"？将会删除所有相关数据`))
          del()
      }
    })
  }

  initSetting () {

  }

  showFilterBar (items: {[name: string]: string}, clickEvt: (item) => void) {
    const filterBarEl = Utils.createElement(`<div class="atk-filter-bar"></div>`)
    this.el.prepend(filterBarEl)

    Object.entries(items).forEach(([name, label]) => {
      const itemEl = Utils.createElement(`<span class="atk-filter-item"></span>`)
      itemEl.innerText = label
      itemEl.addEventListener('click', () => {
        clickEvt(name)
        filterBarEl.querySelectorAll('.atk-active')
          .forEach(item => item.classList.remove('atk-active')) // 删除其他 active
        itemEl.classList.add('atk-active')
      })
      filterBarEl.append(itemEl)
    })

    ;(filterBarEl.firstChild as HTMLElement).click() // 默认打开第一个项目
  }
}

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
    setting: '配置',
  }
  activeAction = ''
  adminOnly = true

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

  async initCommentList() {
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
    const sites = await new Api(this.ctx).siteGet()
    const siteItems: {[name: string]: string} = {'': '所有站点'}
    if (sites) sites.forEach((site) => { siteItems[site.name] = site.name })
    this.showFilterBar(siteItems, (s) => {
      conf.siteName = s
      reqComments(conf.typeName, conf.siteName)
    })
    loaded = true
  }

  async initPageList () {
    this.el.innerHTML = ''

    const pListEl = Utils.createElement(`<div class="atk-sidebar-list"></div>`)
    this.el.append(pListEl)

    // 初始化 site 列表
    let loaded = false
    const sitesData = await new Api(this.ctx).siteGet()
    const siteItems: {[name: string]: string} = {'': '所有站点'}
    if (sitesData) sitesData.forEach((site) => { siteItems[site.name] = site.name })
    this.showFilterBar(siteItems, async (s) => {
      if (!loaded) return
      await this.reqPages(pListEl, s)
    })
    loaded = true

    await this.reqPages(pListEl)
  }

  async reqPages(pListEl: HTMLElement, siteName?: string) {
    pListEl.innerHTML = ''

    const pages = await new Api(this.ctx).pageGet(siteName)
    if (!pages) {
      return
    }

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
    const sListEl = Utils.createElement(`
    <div class="atk-sidebar-list">
      <div class="atk-site-add atk-form-inline-wrap">
        <input type="text" name="siteAdd_name" placeholder="Name..." />
        <input type="text" name="siteAdd_url" placeholder="URL..." />
        <button name="siteAdd_submit">Add</button>
      </div>
    </div>
    `)
    this.el.append(sListEl)
    const siteAddNameEl = sListEl.querySelector<HTMLInputElement>('[name="siteAdd_name"]')!
    const siteAddUrlEl = sListEl.querySelector<HTMLInputElement>('[name="siteAdd_url"]')!
    sListEl.querySelector<HTMLElement>('[name="siteAdd_submit"]')!.onclick = () => {
      const name = siteAddNameEl.value.trim()
      const url = siteAddUrlEl.value.trim()
      if (name === '') {
        siteAddNameEl.focus()
        return
      }
      new Api(this.ctx).siteAdd(name, url)
        .then(() => {
          this.initSiteList()
        })
        .catch((err) => {
          window.alert(`创建失败：${err.msg || ''}`)
        })
    }
    const sites = await new Api(this.ctx).siteGet()
    if (!sites) {
      return
    }
    sites.forEach(site => {
      const siteItemEl = Utils.createElement(`
      <div class="atk-item">
      <div class="atk-title"></div>
      <div class="atk-sub"></div>
      <div class="atk-actions">
        <div class="atk-item" data-action="site-edit">修改 URL</div>
        <div class="atk-item" data-action="site-del">删除</div>
      </div>
      </div>`)
      sListEl.append(siteItemEl)

      siteItemEl.setAttribute('data-site-id', String(site.id))

      const nameEl = siteItemEl.querySelector<HTMLElement>('.atk-title')!
      const urlEl = siteItemEl.querySelector<HTMLElement>('.atk-sub')!

      nameEl.innerText = site.name || site.first_url
      urlEl.innerText = site.urls_raw
      nameEl.onclick = () => {
        window.open(`${site.first_url}`)
      }
      urlEl.onclick = () => {
        window.open(`${site.first_url}`)
      }

      const delBtn = siteItemEl.querySelector<HTMLElement>('[data-action="site-del"]')!
      delBtn.onclick = () => {
        const del = () => {
          new Api(this.ctx).siteDel(site.id, true)
          .then(success => {
            if (success) siteItemEl.remove()
          })
        }
        if (window.confirm(`确认删除站点 "${site.name || site.first_url}"？将会删除所有相关数据`))
          del()
      }

      const editBtn = siteItemEl.querySelector<HTMLElement>('[data-action="site-edit"]')!
      editBtn.onclick = () => {
        const val = window.prompt('修改 URL (多个 URL 用逗号隔开):', site.urls_raw)
        if (val !== null) {
          new Api(this.ctx).siteEdit(site.id, { name: site.name, urls: val })
            .then(() => {
              this.initSiteList()
            })
            .catch((err) => {
              window.alert(`修改失败：${err.msg || '未知错误'}`)
            })
        }
      }
    })
  }

  initSetting () {
    this.el.innerHTML = ''

    const settingEl = Utils.createElement(`
    <div class="atk-setting">
    <div class="atk-group atk-importer atk-form-wrap">
    <div class="atk-title">导入评论数据</div>
    <input type="file" name="importer_dataFile" accept="text/plain,.json" />
    <div class="atk-label">数据类型：</div>
    <select name="importer_dataType">
      <option value="artalk_v1">Artalk v1 (PHP 旧版)</option>
    </select>
    <div class="atk-label">目标站点名：</div>
    <input type="text" name="importer_siteName" />
    <button class="atk-btn" name="importer_submit">导入</button>
    </div>
    </div>`)
    this.el.append(settingEl)

    const impDataFileEl = settingEl.querySelector<HTMLInputElement>('[name="importer_dataFile"]')!
    const impSiteNameEl = settingEl.querySelector<HTMLInputElement>('[name="importer_siteName"]')!
    const impDataTypeEl = settingEl.querySelector<HTMLSelectElement>('[name="importer_dataType"]')!
    const impSubmitEl = settingEl.querySelector<HTMLButtonElement>('[name="importer_submit"]')!
    const impSubmitTextOrg = impSubmitEl.innerText
    impSubmitEl.onclick = () => {
      const siteName = impSiteNameEl.value.trim()
      const dataType = impDataTypeEl.value.trim()

      if (!impDataFileEl.files || impDataFileEl.files.length === 0) {
        window.alert('请打开文件')
        return
      }
      if (dataType === '') {
        window.alert('请选择数据类型')
        return
      }

      const reader = new FileReader()
      reader.onload = () => {
        const data = String(reader.result)
        if (!data) return

        impSubmitEl.innerText = '请稍后...'
        new Api(this.ctx).importer(data, dataType, siteName)
        .then(() => {
          window.alert(`导入成功`)
        })
        .catch((err) => {
          console.error(err)
          window.alert(`导入失败：${err.msg || '未知错误'}`)
        })
        .finally(() => {
          impSubmitEl.innerText = impSubmitTextOrg
        })
      }
      reader.readAsText(impDataFileEl.files[0]);
    }
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

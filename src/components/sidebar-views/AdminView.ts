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

  /** site 筛选 · 初始化 */
  async initSiteFilterBar (clickEvt: (item: FilterBarItem) => void) {
    const sitesData = await new Api(this.ctx).siteGet()
    const siteItems: FilterBarItem[] = [
      {name: '__ATK_SITE_ALL', label: '全部站点'}
    ]
    if (sitesData) sitesData.forEach((site) => {
      siteItems.push({ name: site.name, label: site.name })
    })

    const filterBarEl = BuildFilterBar(siteItems, (item) => clickEvt(item))
    return filterBarEl
  }

  /** 评论列表 · 初始化 */
  async initCommentList() {
    const el = Utils.createElement('<div class="atk-admin-comment-list"></div>')
    el.append(this.cList.el) // TODO: 统一 loading 动画

    const reqComments = (type: string, siteName: string) => {
      if (!this.cList) return
      this.cList.type = `admin_${type}` as any
      this.cList.isFirstLoad = true
      this.cList.paramsEditor = (params) => {
        params.site_name = siteName
      }
      this.cList.reqComments()
    }

    let loaded = false
    const curt = { typeName: 'all', siteName: '' }

    // 初始化
    const typeFilterBarEl = BuildFilterBar([
      { name: 'all', label: '全部' },
      { name: 'pending', label: '待审' }
    ], (item) => {
      if (!loaded) return
      curt.typeName = item.name
      reqComments(curt.typeName, curt.siteName)
    })
    el.prepend(typeFilterBarEl)

    // 初始化 site 筛选
    const siteFilterBarEl  = await this.initSiteFilterBar((item) => {
      curt.siteName = item.name
      reqComments(curt.typeName, curt.siteName)
    })
    el.prepend(siteFilterBarEl)
    loaded = true

    this.el.innerHTML = ''
    this.el.append(el)
  }

  /** 页面列表 · 初始化 */
  async initPageList () {
    const el = Utils.createElement('<div class="atk-admin-page-list"></div>')

    const pListEl = Utils.createElement(`<div class="atk-sidebar-list"></div>`)
    el.append(pListEl)

    // 初始化 site 筛选
    const siteFilterBarEl = await this.initSiteFilterBar((item) => {
      this.reqPages(pListEl, item.name)
    })
    el.prepend(siteFilterBarEl)

    this.el.innerHTML = ''
    this.el.append(el)
  }

  async reqPages(pListEl: HTMLElement, siteName?: string) {
    pListEl.innerHTML = ''

    const pages = await new Api(this.ctx).pageGet(siteName)
    if (!pages) {
      pListEl.innerHTML = '<div class="atk-sidebar-no-content">无内容</div>'
      return
    }

    pages.forEach(page => {
      const pageItemEl = Utils.createElement(`
      <div class="atk-item">
      <div class="atk-title"></div>
      <div class="atk-sub"></div>
      <div class="atk-actions">
        <div class="atk-item" data-action="page-edit-title">修改标题</div>
        <div class="atk-item" data-action="page-fetch">获取标题</div>
        <div class="atk-item" data-action="page-edit-key">修改 KEY</div>
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
        page.admin_only = !page.admin_only
        new Api(this.ctx).pageEdit(page)
        .then((p) => {
          page = p
          renderAdminOnlyBtn()
        })
      }

      const editTitleBtn = pageItemEl.querySelector<HTMLElement>('[data-action="page-edit-title"]')!
      editTitleBtn.onclick = () => {
        const val = window.prompt('修改标题：', page.title)
        if (val !== null) {
          page.title = val
          new Api(this.ctx).pageEdit(page)
            .then(() => {
              this.reqPages(pListEl, siteName)
            })
            .catch((err) => {
              window.alert(`修改失败：${err.msg || '未知错误'}`)
            })
        }
      }

      const fetchBtn = pageItemEl.querySelector<HTMLElement>('[data-action="page-fetch"]')!
      fetchBtn.onclick = () => {
        const btnOrgTxt = fetchBtn.innerText
        fetchBtn.innerText = '获取中...'
        new Api(this.ctx).pageFetch(page.id)
        .then((p) => {
          page = p
          pageItemEl.querySelector<HTMLElement>('.atk-title')!.innerText = p.title
        })
        .catch((err) => {
          window.alert(`获取失败：${err.msg || '未知错误'}`)
        })
        .finally(() => {
          fetchBtn.innerText = btnOrgTxt
        })
      }

      const editKeyBtn = pageItemEl.querySelector<HTMLElement>('[data-action="page-edit-key"]')!
      editKeyBtn.onclick = () => {
        const val = window.prompt('修改 Key：', page.key)
        if (val !== null) {
          page.key = val
          new Api(this.ctx).pageEdit(page)
            .then(() => {
              this.reqPages(pListEl, siteName)
            })
            .catch((err) => {
              window.alert(`修改失败：${err.msg || '未知错误'}`)
            })
        }
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

  /** 站点列表 · 初始化 */
  async initSiteList () {
    const el = Utils.createElement('<div class="atk-site-list"></div>')
    // TODO: 可复用，特别是 actions
    const sListEl = Utils.createElement(`
    <div class="atk-sidebar-list">
      <div class="atk-site-add atk-form-inline-wrap">
        <input type="text" name="siteAdd_name" placeholder="Name..." />
        <input type="text" name="siteAdd_urls" placeholder="URL..." />
        <button name="siteAdd_submit">Add</button>
      </div>
    </div>
    `)
    el.append(sListEl)
    const siteAddNameEl = sListEl.querySelector<HTMLInputElement>('[name="siteAdd_name"]')!
    const siteAddUrlsEl = sListEl.querySelector<HTMLInputElement>('[name="siteAdd_urls"]')!
    sListEl.querySelector<HTMLElement>('[name="siteAdd_submit"]')!.onclick = () => {
      const name = siteAddNameEl.value.trim()
      const urls = siteAddUrlsEl.value.trim()
      if (name === '') {
        siteAddNameEl.focus()
        return
      }
      new Api(this.ctx).siteAdd(name, urls)
        .then(() => {
          this.initSiteList()
        })
        .catch((err) => {
          window.alert(`创建失败：${err.msg || ''}`)
        })
    }
    const sites = await new Api(this.ctx).siteGet()
    if (!sites) {
      el.append(Utils.createElement('<div class="atk-sidebar-no-content">无内容</div>'))
      return
    }
    sites.forEach(site => {
      const siteItemEl = Utils.createElement(`
      <div class="atk-item">
      <div class="atk-title"></div>
      <div class="atk-sub"></div>
      <div class="atk-actions">
        <div class="atk-item" data-action="site-rename">重命名</div>
        <div class="atk-item" data-action="site-edit-urls">修改 URL</div>
        <div class="atk-item" data-action="site-del">删除</div>
      </div>
      </div>`)
      sListEl.append(siteItemEl)

      siteItemEl.setAttribute('data-site-id', String(site.id))

      const nameEl = siteItemEl.querySelector<HTMLElement>('.atk-title')!
      const urlsEl = siteItemEl.querySelector<HTMLElement>('.atk-sub')!

      nameEl.innerText = site.name || site.first_url
      nameEl.onclick = () => {
        window.open(`${site.first_url}`)
      }

      if (site.urls) {
        site.urls.forEach((u) => {
          const urlItemEl = Utils.createElement('<span style="margin-right: 10px;" />')
          urlItemEl.innerText = u
          urlItemEl.onclick = () => {
            window.open(`${u}`)
          }
          urlsEl.append(urlItemEl)
        })
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

      const editUrlsBtn = siteItemEl.querySelector<HTMLElement>('[data-action="site-edit-urls"]')!
      editUrlsBtn.onclick = () => {
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

      const renameBtn = siteItemEl.querySelector<HTMLElement>('[data-action="site-rename"]')!
      renameBtn.onclick = () => {
        const val = window.prompt('编辑站点名称：', site.name)
        if (val !== null) {
          new Api(this.ctx).siteEdit(site.id, { name: val, urls: site.urls_raw })
            .then(() => {
              this.initSiteList()
            })
            .catch((err) => {
              window.alert(`修改失败：${err.msg || '未知错误'}`)
            })
        }
      }
    })

    this.el.innerHTML = ''
    this.el.append(el)
  }

  /** 配置页面 · 初始化 */
  initSetting () {
    const el = Utils.createElement('<div class="atk-admin-setting"></div>')

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
    el.append(settingEl)

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

    this.el.innerHTML = ''
    this.el.append(el)
  }
}

/** 通用筛选条 */
function BuildFilterBar (items: FilterBarItem[], clickEvt: (item: FilterBarItem) => void): HTMLElement {
  const filterBarEl = Utils.createElement(`<div class="atk-filter-bar"></div>`)

  items.forEach(item => {
    const itemEl = Utils.createElement(`<span></span>`)
    filterBarEl.append(itemEl)
    itemEl.innerText = item.label
    itemEl.addEventListener('click', () => {
      clickEvt(item)
      filterBarEl.querySelectorAll('.atk-active')
        .forEach(el => el.classList.remove('atk-active')) // 删除其他 active
      itemEl.classList.add('atk-active')
    })
  })

  ;(filterBarEl.firstChild as HTMLElement).click() // 默认打开第一个项目

  return filterBarEl
}

/** 筛选条中的项目 */
interface FilterBarItem {
  name: string
  label: string
}

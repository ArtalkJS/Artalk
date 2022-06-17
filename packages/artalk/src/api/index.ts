import { CommentData, ListData, UserData, PageData, SiteData, NotifyData } from '~/types/artalk-data'
import ArtalkConfig from '~/types/artalk-config'
import Context from '~/types/context'
import { Fetch, ToFormData, POST } from './request'
import * as Utils from '../lib/utils'
import 'abortcontroller-polyfill/dist/polyfill-patch-fetch'

export default class Api {
  private ctx: Context
  private baseURL: string

  constructor (ctx: Context) {
    this.ctx = ctx
    this.baseURL = `${ctx.conf.server}/api`
  }

  // ============================
  //  评论 Comment
  // ============================

  /** 评论 · 获取 */
  public get(offset: number, pageSize: number, flatMode?: boolean, paramsEditor?: (params: any) => void) {
    const params: any = {
      page_key: this.ctx.conf.pageKey,
      site_name: this.ctx.conf.site || '',
      limit: pageSize,
      offset,
    }

    if (flatMode) params.flat_mode = flatMode // 平铺模式
    if (this.ctx.user.checkHasBasicUserInfo()) {
      params.name = this.ctx.user.data.nick
      params.email = this.ctx.user.data.email
    }

    if (paramsEditor) paramsEditor(params)

    return POST<ListData>(this.ctx, `${this.baseURL}/get`, params)
  }

  /** 评论 · 创建 */
  public async add(comment: { nick: string, email: string, link: string, content: string, rid: number, page_key: string, page_title?: string, site_name?: string }) {
    const params: any = {
      name: comment.nick,
      email: comment.email,
      link: comment.link,
      content: comment.content,
      rid: comment.rid,
      page_key: comment.page_key,
      ua: await Utils.getCorrectUserAgent(), // 需要后端支持，获取修正后的 UA
    }

    if (comment.page_title) params.page_title = comment.page_title
    if (comment.site_name) params.site_name = comment.site_name

    const data = await POST<any>(this.ctx, `${this.baseURL}/add`, params)
    return (data.comment as CommentData)
  }

  /** 评论 · 修改 */
  public async commentEdit(data: Partial<CommentData>) {
    const params: any = {
      ...data,
    }

    const d = await POST<any>(this.ctx, `${this.baseURL}/admin/comment-edit`, params)
    return (d.comment as CommentData)
  }

  /** 评论 · 删除 */
  public commentDel(commentID: number, siteName?: string) {
    const params: any = {
      id: String(commentID),
      site_name: siteName || '',
    }

    return POST(this.ctx, `${this.baseURL}/admin/comment-del`, params)
  }

  // ============================
  //  用户 User
  // ============================

  /** 用户 · 登录 */
  public async login(name: string, email: string, password: string) {
    const params: any = {
      name, email, password
    }

    if (this.ctx.conf.site) params.site_name = this.ctx.conf.site

    const data = await POST<any>(this.ctx, `${this.baseURL}/login`, params)
    return (data.token as string)
  }

  /** 用户 · 获取  */
  public userGet(name: string, email: string) {
    const ctrl = new AbortController()
    const params: any = {
      name, email, site_name: (this.ctx.conf.site || '')
    }

    const req = Fetch(this.ctx, `${this.baseURL}/user-get`, {
      method: 'POST',
      body: ToFormData(params),
      signal: ctrl.signal,
    }).then((json) => ({
      user: json.data.user as UserData|null,
      is_login: json.data.is_login as boolean,
      unread: (json.data.unread || []) as NotifyData[],
      unread_count: json.data.unread_count || 0,
    }))

    return {
      req,
      abort: () => { ctrl.abort() },
    }
  }

  /** 用户 · 登录状态 */
  public async loginStatus() {
    const data = await POST<any>(this.ctx, `${this.baseURL}/login-status`, {
      name: this.ctx.user.data.nick,
      email: this.ctx.user.data.email
    })
    return (data || { is_login: false, is_admin: false }) as { is_login: boolean, is_admin: boolean }
  }

  /** 用户 · 注销 */
  public async logout() {
    return POST(this.ctx, `${this.baseURL}/logout`)
  }

  // ============================
  //  页面 Page
  // ============================

  /** 页面 · 获取 */
  public async pageGet(siteName?: string, offset?: number, limit?: number) {
    const params: any = {
      site_name: siteName || '',
      offset: offset || 0,
      limit: limit || 15,
    }

    const d = await POST<any>(this.ctx, `${this.baseURL}/admin/page-get`, params)
    return (d as { pages: PageData[], total: number })
  }

  /** 页面 · 修改 */
  public async pageEdit(data: PageData) {
    const params: any = {
      id: data.id,
      key: data.key,
      title: data.title,
      admin_only: data.admin_only,
      site_name: data.site_name || this.ctx.conf.site,
    }

    const d = await POST<any>(this.ctx, `${this.baseURL}/admin/page-edit`, params)
    return (d.page as PageData)
  }

  /** 页面 · 删除 */
  public pageDel(pageKey: string, siteName?: string) {
    const params: any = {
      key: String(pageKey),
      site_name: siteName || '',
    }

    return POST(this.ctx, `${this.baseURL}/admin/page-del`, params)
  }

  /** 页面 · 数据更新 */
  public async pageFetch(id?: number, siteName?: string, getStatus?: boolean) {
    const params: any = {}
    if (id) params.id = id
    if (siteName) params.site_name = siteName
    if (getStatus) params.get_status = getStatus

    const d = await POST<any>(this.ctx, `${this.baseURL}/admin/page-fetch`, params)
    return (d as any)
  }

  // ============================
  //  站点 Site
  // ============================

  /** 站点 · 获取 */
  public async siteGet() {
    const params: any = {}

    const d = await POST<any>(this.ctx, `${this.baseURL}/admin/site-get`, params)
    return (d.sites as SiteData[])
  }

  /** 站点 · 创建 */
  public async siteAdd(name: string, urls: string) {
    const params: any = {
      name, urls,
    }

    const d = await POST<any>(this.ctx, `${this.baseURL}/admin/site-add`, params)
    return (d.site as SiteData)
  }

  /** 站点 · 修改 */
  public async siteEdit(data: SiteData) {
    const params: any = {
      id: data.id,
      name: data.name || '',
      urls: data.urls || '',
    }

    const d = await POST<any>(this.ctx, `${this.baseURL}/admin/site-edit`, params)
    return (d.site as SiteData)
  }

  /** 站点 · 删除 */
  public siteDel(id: number, delContent = false) {
    const params: any = { id, del_content: delContent }

    return POST(this.ctx, `${this.baseURL}/admin/site-del`, params)
  }

  /** 导出 */
  public async export() {
    const d = await Fetch(this.ctx, `${this.baseURL}/admin/export`, { method: 'POST' }, 0)
    return (d.data?.data || '' as string)
  }

  // ============================
  //  杂项
  // ============================

  /** 投票 */
  public async vote(targetID: number, type: string) {
    const params: any = {
      site_name: this.ctx.conf.site || '',
      target_id: targetID,
      type,
    }

    if (this.ctx.user.checkHasBasicUserInfo()) {
      params.name = this.ctx.user.data.nick
      params.email = this.ctx.user.data.email
    }

    const data = await POST<any>(this.ctx, `${this.baseURL}/vote`, params)
    return (data as {up: number, down: number})
  }

  /** 已读标记 */
  public markRead(notifyKey: string, readAll = false) {
    const params: any = {
      site_name: this.ctx.conf.site || '',
      notify_key: notifyKey,
    }

    if (readAll) {
      delete params.notify_key
      params.read_all = true
      params.name = this.ctx.user.data.nick
      params.email = this.ctx.user.data.email
    }

    return POST(this.ctx, `${this.baseURL}/mark-read`, params)
  }

  /** PV */
  public async pv() {
    const params: any = {
      site_name: this.ctx.conf.site || '',
      page_key: this.ctx.conf.pageKey || '',
      page_title: this.ctx.conf.pageTitle || ''
    }

    const p = await POST<any>(this.ctx, `${this.baseURL}/pv`, params)
    return p.pv as number
  }

  /** 统计 */
  public async stat(
    type: 'latest_comments'|'latest_pages'|'pv_most_pages'|'comment_most_pages'|
          'page_pv'|'site_pv'|'page_comment'|'site_comment',
    pageKeys?: string|string[],
    limit?: number
  ) {
    const params: any = {
      type,
      site_name: this.ctx.conf.site || '',
    }

    if (pageKeys) params.page_keys = Array.isArray(pageKeys) ? pageKeys.join(',') : pageKeys
    if (limit) params.limit = limit

    const data = await POST<PageData[]|CommentData[]|object|number>(this.ctx, `${this.baseURL}/stat`, params)
    return data
  }

  /** 图片上传 */
  public async imgUpload(file: File) {
    const params: any = {
      name: this.ctx.user.data.nick,
      email: this.ctx.user.data.email,
      page_key: this.ctx.conf.pageKey,
    }

    if (this.ctx.conf.site) params.site_name = this.ctx.conf.site

    const form = ToFormData(params)
    form.set('file', file)

    const init: RequestInit = {
      method: 'POST',
      body: form
    }

    const json = await Fetch(this.ctx, `${this.baseURL}/img-upload`, init)
    return ((json.data || {}) as any) as { img_file: string, img_url: string }
  }

  /** 获取配置 */
  public async conf() {
    const data = await POST<any>(this.ctx, `${this.baseURL}/conf`)
    const conf = (data.frontend_conf || {}) as ArtalkConfig

    // Patch: `emoticons` config string to json
    if (conf.emoticons && typeof conf.emoticons === "string") {
      conf.emoticons = conf.emoticons.trim()
      if (conf.emoticons.startsWith("[") || conf.emoticons.startsWith("{")) {
        conf.emoticons = JSON.parse(conf.emoticons) // pase json
      } else if (conf.emoticons === "false") {
        conf.emoticons = false
      }
    }

    return conf
  }

  // ============================
  //  验证码
  // ============================

  /** 验证码 · 获取 */
  public async captchaGet() {
    const data = await POST<any>(this.ctx, `${this.baseURL}/captcha/refresh`)
    return (data.img_data || '') as string
  }

  /** 验证码 · 检验 */
  public async captchaCheck(value: string) {
    const data = await POST<any>(this.ctx, `${this.baseURL}/captcha/check`, { value })
    return (data.img_data || '') as string
  }

  /** 验证码 · 状态 */
  public async captchaStatus() {
    const data = await POST<any>(this.ctx, `${this.baseURL}/captcha/status`)
    return (data || { is_pass: false }) as { is_pass: boolean }
  }

  // ============================
  //  管理员缓存操作
  // ============================

  /** 缓存清除 */
  public cacheFlushAll() {
    const params: any = { flush_all: true }
    return POST(this.ctx, `${this.baseURL}/admin/cache-flush`, params)
  }

  /** 缓存预热 */
  public cacheWarmUp() {
    const params: any = {}
    return POST(this.ctx, `${this.baseURL}/admin/cache-warm`, params)
  }
}

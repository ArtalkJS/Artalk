import { CommentData, ListData, UserData, PageData, SiteData, NotifyData } from '~/types/artalk-data'
import Context from '../Context'

export default class Api {
  private ctx: Context
  private serverURL: string

  constructor (ctx: Context) {
    this.ctx = ctx
    this.serverURL = ctx.conf.server
  }

  public get(offset: number, type?: string, paramsEditor?: (params: any) => void): Promise<ListData> {
    const params: any = {
      page_key: this.ctx.conf.pageKey,
      limit: this.ctx.conf.readMore?.pageSize || 15,
      offset,
    }

    if (type) {
      params.type = type
    }

    if (this.ctx.user.checkHasBasicUserInfo()) {
      params.name = this.ctx.user.data.nick
      params.email = this.ctx.user.data.email
    }

    if (this.ctx.conf.site) params.site_name = this.ctx.conf.site

    if (paramsEditor) paramsEditor(params)

    return CommonFetch(this.ctx, `${this.serverURL}/get`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => (json.data as ListData))
  }

  public add(comment: { nick: string, email: string, link: string, content: string, rid: number }): Promise<CommentData> {
    const params: any = {
      name: comment.nick,
      email: comment.email,
      link: comment.link,
      content: comment.content,
      rid: comment.rid,
      page_key: this.ctx.conf.pageKey,
      page_url: this.ctx.conf.pageUrl || '',
      page_title: this.ctx.conf.pageTitle || '',
    }

    if (this.ctx.conf.site) params.site_name = this.ctx.conf.site

    return CommonFetch(this.ctx, `${this.serverURL}/add`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => (json.data.comment as CommentData))
  }

  public login(name: string, email: string, password: string): Promise<string> {
    const params: any = {
      name, email, password
    }

    if (this.ctx.conf.site) params.site_name = this.ctx.conf.site

    return CommonFetch(this.ctx, `${this.serverURL}/login`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => (json.data.token))
  }

  public userGet(name: string, email: string) {
    const ctrl = new AbortController()
    const { signal } = ctrl

    const params: any = {
      name, email
    }

    if (this.ctx.conf.site) params.site_name = this.ctx.conf.site

    const req = CommonFetch(this.ctx, `${this.serverURL}/user-get`, {
      method: 'POST',
      body: getFormData(params),
      signal,
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

  public vote(targetID: number, type: string) {
    const params: any = {
      target_id: targetID,
      type,
    }

    if (this.ctx.user.checkHasBasicUserInfo()) {
      params.name = this.ctx.user.data.nick
      params.email = this.ctx.user.data.email
    }
    if (this.ctx.conf.site) params.site_name = this.ctx.conf.site

    return CommonFetch(this.ctx, `${this.serverURL}/vote`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => ((json.data?.vote_num || 0) as number))
  }

  public markRead(notifyKey: string, readAll = false) {
    const params: any = {
      notify_key: notifyKey,
    }

    if (readAll) {
      delete params.notify_key
      params.read_all = true
      params.name = this.ctx.user.data.nick
      params.email = this.ctx.user.data.email
    }

    if (this.ctx.conf.site) params.site_name = this.ctx.conf.site

    return CommonFetch(this.ctx, `${this.serverURL}/mark-read`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => (json.success as boolean))
  }

  public captchaGet(): Promise<string> {
    return CommonFetch(this.ctx, `${this.serverURL}/captcha/refresh`, {
      method: 'GET',
    }).then((json) => {
      if (!!json.success && !!json.data.img_data) {
        return json.data.img_data
      }

      return ''
    })
  }

  public captchaCheck(value: string): Promise<string> {
    return CommonFetch(this.ctx, `${this.serverURL}/captcha/check?${new URLSearchParams({ value })}`, {
      method: 'GET',
    }).then((json) => {
      if (!json.success && !!json.data.img_data) {
        return json.data.img_data
      }

      return ''
    })
  }

  public commentEdit(data: CommentData) {
    const params: any = {
      ...data,
    }

    if (params.is_collapsed !== undefined) params.is_collapsed = params.is_collapsed ? '1' : '0'
    if (params.is_pending !== undefined) params.is_pending = params.is_pending ? '1' : '0'

    if (this.ctx.conf.site) params.site_name = this.ctx.conf.site

    return CommonFetch(this.ctx, `${this.serverURL}/admin/comment-edit`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => (json.data.comment as CommentData))
  }

  public pageEdit(data: PageData) {
    const params: any = {
      id: data.id,
      key: data.key,
      title: data.title,
      admin_only: data.admin_only ? '1' : '0',
    }

    if (this.ctx.conf.site) params.site_name = this.ctx.conf.site

    return CommonFetch(this.ctx, `${this.serverURL}/admin/page-edit`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => (json.data.page as PageData))
  }

  public pageFetch(id: number) {
    const params: any = {
      id,
    }

    return CommonFetch(this.ctx, `${this.serverURL}/admin/page-fetch`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => (json.data.page as PageData))
  }

  public siteGet(): Promise<SiteData[]> {
    const params: any = {}

    return CommonFetch(this.ctx, `${this.serverURL}/admin/site-get`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => (json.data.sites as SiteData[]))
  }

  public pageGet(siteName?: string): Promise<PageData[]> {
    const params: any = {
      site_name: siteName || ''
    }

    return CommonFetch(this.ctx, `${this.serverURL}/admin/page-get`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => (json.data.pages as PageData[]))
  }

  public pageDel(pageKey: string, siteName?: string) {
    const params: any = {
      key: String(pageKey),
      site_name: siteName || '',
    }

    return CommonFetch(this.ctx, `${this.serverURL}/admin/page-del`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => (json.success as boolean))
  }

  public siteDel(id: number, delContent = false) {
    const params: any = {
      id,
      del_content: delContent ? '1' : '0'
    }

    return CommonFetch(this.ctx, `${this.serverURL}/admin/site-del`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => (json.success as boolean))
  }

  public siteAdd(name: string, urls: string) {
    const params: any = {
      name, urls,
    }

    return CommonFetch(this.ctx, `${this.serverURL}/admin/site-add`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => (json.site as SiteData))
  }

  public siteEdit(id: number, data: {
    name: string
    urls: string
  }) {
    const params: any = {
      id,
      name: data.name,
      urls: data.urls,
    }

    return CommonFetch(this.ctx, `${this.serverURL}/admin/site-edit`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => (json.site as SiteData))
  }

  public commentDel(commentID: number, siteName?: string) {
    const params: any = {
      id: String(commentID),
      site_name: siteName || '',
    }

    return CommonFetch(this.ctx, `${this.serverURL}/admin/comment-del`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => (json.success as boolean))
  }

  public importer(data: string, type: string, siteName: string) {
    const params: any = {
      data,
      type,
      site_name: siteName || '',
    }

    return CommonFetch(this.ctx, `${this.serverURL}/admin/importer`, {
      method: 'POST',
      body: getFormData(params),
    }).then((json) => (json.success as boolean))
  }
}

/** 公共请求函数 */
function CommonFetch(ctx: Context, input: RequestInfo, init: RequestInit): Promise<any> {
  if (ctx.user.data.token) {
    const requestHeaders: HeadersInit = new Headers();
    requestHeaders.set('Authorization', `Bearer ${ctx.user.data.token}`);

    init.headers = requestHeaders
  }

  // 15s timeout
  return timeoutPromise(15000, fetch(input, init)).then(async (resp) => {
    // 解析获取响应的 json
    let json: any = await resp.json()

    // 重新发起请求
    const recall = (resolve, reject) => {
      CommonFetch(ctx, input, init).then(d => {
        resolve(d)
      }).catch(err => {
        reject(err)
      })
    }

    if (json.data && json.data.need_captcha) {
      // 请求需要验证码
      json = await (new Promise<any>((resolve, reject) => {
        ctx.dispatchEvent('checker-captcha', {
          imgData: json.data.img_data,
          onSuccess: () => {
            recall(resolve, reject)
          },
          onCancel: () => {
            reject(json)
          }
        })
      }))
    } else if ((json.data && json.data.need_login) || resp.status === 401) {
      // 请求需要管理员权限
      json = await (new Promise<any>((resolve, reject) => {
        ctx.dispatchEvent('checker-admin', {
          onSuccess: () => {
            recall(resolve, reject)
          },
          onCancel: () => {
            reject(json)
          }
        })
      }))
    }

    if (!json.success)
      throw json // throw 相当于 reject(json)
    else
      return json
  })
}

function getFormData (object: any): FormData {
  const formData = new FormData()
  Object.keys(object).forEach(key => formData.append(key, String(object[key])))
  return formData
}

/** TODO: 我靠，一个 timeout，都要丑陋的实现 */
function timeoutPromise<T>(ms: number, promise: Promise<T>): Promise<T> {
  return new Promise((resolve, reject) => {
    const timeoutId = setTimeout(() => {
      reject(new Error("promise timeout"))
    }, ms);
    promise.then(
      (res) => {
        clearTimeout(timeoutId);
        resolve(res);
      },
      (err) => {
        clearTimeout(timeoutId);
        reject(err);
      }
    );
  })
}

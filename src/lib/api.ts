import { CommentData, ListData } from '~/types/artalk-data'
import Context from '../Context'

export default class Api {
  private ctx: Context
  private serverURL: string

  constructor (ctx: Context) {
    this.ctx = ctx
    this.serverURL = ctx.conf.serverUrl
  }

  public async get(offset: number): Promise<ListData> {
    const params = getFormData({
      page_key: this.ctx.conf.pageKey,
      limit: this.ctx.conf.readMore?.pageSize || 15,
      offset,
    })

    return commonFetch(this.ctx, `${this.serverURL}/get`, {
      method: 'POST',
      body: params,
    }).then((json) => (json.data as ListData))
  }

  public async add(comment: { nick: string, email: string, link: string, content: string, rid: number }): Promise<CommentData> {
    const params = getFormData({
      name: comment.nick,
      email: comment.email,
      link: comment.link,
      content: comment.content,
      rid: comment.rid,
      page_key: this.ctx.conf.pageKey,
      token: this.ctx.user.data.token,
    })

    return commonFetch(this.ctx, `${this.serverURL}/add`, {
      method: 'POST',
      body: params,
    }).then((json) => (json.data.comment as CommentData))
  }

  public async login(name: string, email: string, password: string): Promise<string> {
    const params = getFormData({
      name, email, password
    })

    return commonFetch(this.ctx, `${this.serverURL}/login`, {
      method: 'POST',
      body: params,
    }).then((json) => (json.data.token))
  }

  public async captchaGet(): Promise<string> {
    return commonFetch(this.ctx, `${this.serverURL}/captcha/refresh`, {
      method: 'GET',
    }).then((json) => {
      if (!!json.success && !!json.data.img_data) {
        return json.data.img_data
      }

      return ''
    })
  }

  public async captchaCheck(value: string): Promise<string> {
    return commonFetch(this.ctx, `${this.serverURL}/captcha/check?${new URLSearchParams({ value })}`, {
      method: 'GET',
    }).then((json) => {
      if (!json.success && !!json.data.img_data) {
        return json.data.img_data
      }

      return ''
    })
  }
}

function commonFetch(ctx: Context, input: RequestInfo, init?: RequestInit | undefined): Promise<any> {
  return timeoutPromise(4000, fetch(input, init)).then(async (resp) => {
    let json: any = await resp.json()

    if (json.data && json.data.need_captcha) { // 请求需要验证码
      const nPromise = new Promise<any>((resolve, reject) => {
        ctx.dispatchEvent('checker-captcha', {
          imgData: json.data.img_data,
          onSuccess: () => {
            commonFetch(ctx, input, init).then(d => {
              resolve(d)
            }).catch(err => {
              reject(err)
            })
          }
        })
      })

      json = await nPromise
    } else if (json.data && json.data.need_login) { // 请求需要管理员权限
      const nPromise = new Promise<any>((resolve, reject) => {
        ctx.dispatchEvent('checker-admin', {
          onSuccess: () => {
            commonFetch(ctx, input, init).then(d => {
              resolve(d)
            }).catch(err => {
              reject(err)
            })
          }
        })
      })

      json = await nPromise
    }

    if (!json.success) {
      throw json
    }

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

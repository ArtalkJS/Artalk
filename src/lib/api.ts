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

    return timeoutPromise(4000, fetch(`${this.serverURL}/api/get`, {
      method: 'POST',
      body: params,
    })).then(async (resp) => {
      const json: any = await resp.json()
      if (!json.success) {
        throw new Error(json.msg)
      }

      return json.data as ListData
    })
  }

  public async add(comment: CommentData): Promise<CommentData> {
    const params = getFormData({
      name: comment.nick,
      email: comment.email,
      link: comment.link,
      content: comment.content,
      rid: comment.rid,
      page_key: this.ctx.conf.pageKey,
      token: this.ctx.user.data.token,
    })

    return timeoutPromise(4000, fetch(`${this.serverURL}/api/add`, {
      method: 'POST',
      body: params,
    })).then(async (resp) => {
      const json: any = await resp.json()
      if (!json.success) {
        throw new Error(json.msg)
      }

      return json.data.comment as CommentData
    })
  }

  public async login(): Promise<string> {
    const params = getFormData({
      user: this, password：
    })

    return timeoutPromise(4000, fetch(`${this.serverURL}/login`, {
      method: 'POST',
      body: params,
    })).then(async (resp) => {
      const json: any = await resp.json()
      if (!json.success) {
        throw new Error(json.msg)
      }

      return json.data.token
    })
  }

  public async captchaGet(): Promise<string> {
    return timeoutPromise(4000, fetch(`${this.serverURL}/captcha/get`, {
      method: 'GET',
    })).then(async (resp) => {
      const json: any = await resp.json()
      if (!json.success && !!json.data.img_data) {
        return json.data.img_data
      }

      return ''
    })
  }

  public async captchaCheck(value: string): Promise<string> {
    return timeoutPromise(4000, fetch(`${this.serverURL}/captcha/get?${new URLSearchParams({ value })}`, {
      method: 'GET',
    })).then(async (resp) => {
      const json: any = await resp.json()
      if (!json.success && !!json.data.img_data) {
        return json.data.img_data
      }

      return ''
    })
  }
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

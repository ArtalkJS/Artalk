import Context from '~/types/context'
import User from '../lib/user'
import $t from '../i18n'

/** 公共请求函数 */
export async function Fetch(ctx: Context, input: RequestInfo, init: RequestInit, timeout?: number): Promise<any> {
  // JWT
  if (User.data.token) {
    const headers = new Headers(init.headers) // 保留原有 headers
    headers.set('Authorization', `Bearer ${User.data.token}`)
    init.headers = headers
  }

  // 跨域请求携带 cookie (无需)
  // init.credentials = 'include'

  // 请求操作
  const resp = await timeoutFetch(ctx, input, timeout || ctx.conf.reqTimeout || 15000, init)

  const respHttpCode = resp.status
  const noAccessCodes = [401, 400]
  const isNoAccess = noAccessCodes.includes(respHttpCode)

  if (!resp.ok && !isNoAccess)
    throw new Error(`${$t('reqGot')} ${respHttpCode}`)

  // 解析获取响应的 json
  let json: any = await resp.json()

  // 重新发起请求
  const recall = (resolve, reject) => {
    Fetch(ctx, input, init)
      .then(d => { resolve(d) })
      .catch(e => { reject(e) })
  }

  // 请求弹出层验证
  if (json.data?.need_captcha) {
    // 请求需要验证码
    json = await (new Promise<any>((resolve, reject) => {
      ctx.checkCaptcha({
        imgData: json.data.img_data,
        iframe: json.data.iframe,
        onSuccess: () => { recall(resolve, reject) },
        onCancel: () => { reject(json) }
      })
    }))
  } else if (json.data?.need_login || isNoAccess) {
    // 请求需要管理员权限
    json = await (new Promise<any>((resolve, reject) => {
      ctx.checkAdmin({
        onSuccess: () => { recall(resolve, reject) },
        onCancel: () => { reject(json) }
      })
    }))
  }

  if (!json.success) throw json // throw 相当于 reject(json)

  return json
}

/** 公共 POST 请求 */
export async function POST<T>(ctx: Context, url: string, data?: {[key: string]: any}) {
  const init: RequestInit = {
    method: 'POST',
  }
  data = { site_name: ctx.conf.site || '', ...(data || {}) } // 总是携带 site_name
  init.body = ToFormData(data)

  const json = await Fetch(ctx, url, init)
  return ((json.data || {}) as T)
}

/** 公共 GET 请求 (无 GET 方式 API) */
// export async function GET<T>(ctx: Context, url: string, data?: {[key: string]: any}) {
//   const json = await Fetch(ctx, url + (data ? (`?${new URLSearchParams(data)}`) : ''), {
//     method: 'GET',
//   })
//   return ((json.data || {}) as T)
// }

/** 对象转 FormData */
export function ToFormData(object: {[key: string]: any}): FormData {
  const formData = new FormData()
  Object.keys(object).forEach(key => formData.append(key, String(object[key])))
  return formData
}

/** 我靠，fetch 一个 timeout，都要丑陋的实现 */
function timeoutFetch(ctx: Context, url: RequestInfo, ms: number, opts: RequestInit) {
  const controller = new AbortController()
  opts.signal?.addEventListener('abort', () => controller.abort()) // 保留原有 signal 功能
  let promise = fetch(url, { ...opts, signal: controller.signal })
  if (ms > 0) {
    const timer = setTimeout(() => controller.abort(), ms)
    promise.finally(() => { clearTimeout(timer) })
  }
  promise = promise.catch((err) => {
    throw ((err || {}).name === 'AbortError') ? new Error($t('reqAborted')) : err
  })
  return promise
}

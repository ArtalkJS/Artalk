import Context from '../context'

/** 公共请求函数 */
export async function Fetch(ctx: Context, input: RequestInfo, init: RequestInit, timeout?: number): Promise<any> {
  // JWT
  if (ctx.user.data.token) {
    const requestHeaders: HeadersInit = new Headers()
    requestHeaders.set('Authorization', `Bearer ${ctx.user.data.token}`)

    init.headers = requestHeaders
  }

  // 请求操作
  try {

    let resp: Response
    if ((typeof timeout !== 'number' && ctx.conf.reqTimeout === 0) || timeout === 0) {
      resp = await fetch(input, init)
    } else {
      // 请求超时检测
      resp = await timeoutPromise(timeout || ctx.conf.reqTimeout || 15000, fetch(input, init))
    }

    const noAccessCodes = [401, 400]

    if (!resp.ok && !noAccessCodes.includes(resp.status)) throw new Error(`请求响应 ${resp.status}`)

    // 解析获取响应的 json
    let json: any = await resp.json()

    // 重新发起请求
    const recall = (resolve, reject) => {
      Fetch(ctx, input, init).then(d => {
        resolve(d)
      }).catch(err => {
        reject(err)
      })
    }

    // 请求弹出层验证
    if (json.data && json.data.need_captcha) {
      // 请求需要验证码
      json = await (new Promise<any>((resolve, reject) => {
        ctx.trigger('checker-captcha', {
          imgData: json.data.img_data,
          iframe: json.data.iframe,
          onSuccess: () => {
            recall(resolve, reject)
          },
          onCancel: () => {
            reject(json)
          }
        })
      }))
    } else if ((json.data && json.data.need_login) || noAccessCodes.includes(resp.status)) {
      // 请求需要管理员权限
      json = await (new Promise<any>((resolve, reject) => {
        ctx.trigger('checker-admin', {
          onSuccess: () => {
            recall(resolve, reject)
          },
          onCancel: () => {
            reject(json)
          }
        })
      }))
    }

    if (!json.success) throw json // throw 相当于 reject(json)

    return json

  } catch (err) {
    // 错误处理
    console.error(err)

    if (err instanceof TypeError)
      throw new Error(`网络错误`)

    throw err
  }
}

/** 公共 POST 请求 */
export async function POST<T>(ctx: Context, url: string, data?: {[key: string]: any}) {
  const init: RequestInit = {
    method: 'POST',
  }
  if (data) init.body = ToFormData(data)

  const json = await Fetch(ctx, url, init)
  return ((json.data || {}) as T)
}

/** 公共 GET 请求 */
export async function GET<T>(ctx: Context, url: string, data?: {[key: string]: any}) {
  const json = await Fetch(ctx, url + (data ? (`?${new URLSearchParams(data)}`) : ''), {
    method: 'GET',
  })
  return ((json.data || {}) as T)
}

/** 对象转 FormData */
export function ToFormData(object: {[key: string]: any}): FormData {
  const formData = new FormData()
  Object.keys(object).forEach(key => formData.append(key, String(object[key])))
  return formData
}

/** 我靠，fetch 一个 timeout，都要丑陋的实现 */
function timeoutPromise<T>(ms: number, promise: Promise<T>): Promise<T> {
  return new Promise((resolve, reject) => {
    const timeoutId = setTimeout(() => {
      reject(new Error("请求超时"))
    }, ms)

    promise.then(
      (res) => {
        clearTimeout(timeoutId)
        resolve(res)
      },
      (err) => {
        clearTimeout(timeoutId)
        reject(err)
      }
    )
  })
}

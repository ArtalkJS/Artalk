import 'abortcontroller-polyfill/dist/polyfill-patch-fetch'

import { FetchError } from '@/types'
import { ApiOptions } from './_options'
import $t from '../i18n'

/** 公共请求函数 */
export async function Fetch<T = any>(opts: ApiOptions, url: string, init: RequestInit, timeout?: number): Promise<T> {
  const headers = new Headers({
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    'X-Requested-With': 'XMLHttpRequest',
    'Authorization': opts.apiToken ? `Bearer ${opts.apiToken}` : '',
    ...init.headers,
  })

  if (!headers.get('Authorization')) headers.delete('Authorization')

  // 请求操作
  const resp = await timeoutFetch(
    url,
    timeout || opts.timeout || 15000,
    { ...init, headers })

  // 解析获取响应的 json
  let json: any = await resp.json().catch(() => {})

  if (!resp.ok) {
    // 重新发起请求
    const recall = (resolve, reject) => {
      Fetch(opts, url, init)
        .then(d => { resolve(d) })
        .catch(e => { reject(e) })
    }

    // 请求弹出层验证
    if (json.need_captcha) {
      // 请求需要验证码
      json = await (new Promise<any>((resolve, reject) => {
        opts.onNeedCheckCaptcha && opts.onNeedCheckCaptcha({
          data: { imgData: json.data.img_data, iframe: json.data.iframe },
          recall: () => { recall(resolve, reject) },
          reject: () => { reject(json) }
        })
      }))
    } else if (json.need_login) {
      // 请求需要管理员权限
      json = await (new Promise<any>((resolve, reject) => {
        opts.onNeedCheckAdmin && opts.onNeedCheckAdmin({
          recall: () => { recall(resolve, reject) },
          reject: () => { reject(json) }
        })
      }))
    } else {
      throw await createError(resp.status, json)
    }
  }

  return json
}

/** 对象转 FormData */
export function ToFormData(object: {[key: string]: any}): FormData {
  const formData = new FormData()
  Object.keys(object).forEach(key => formData.append(key, String(object[key])))
  return formData
}

/** 我靠，fetch 一个 timeout，都要丑陋的实现 */
function timeoutFetch(url: RequestInfo, ms: number, opts: RequestInit) {
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

export class FetchException extends Error implements FetchError {
  code: number = 0
  message: string = 'fetch error'
  data?: any
}

async function createError(code: number, data: any): Promise<FetchException> {
  const err = new FetchException()
  err.message = data.msg || data.message || 'fetch error'
  err.code = code
  err.data = data
  console.error(err)
  return err
}

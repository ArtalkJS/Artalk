import 'abortcontroller-polyfill/dist/polyfill-patch-fetch'

import { FetchError } from '@/types'
import { ApiOptions } from './options'

export const Fetch = async (opts: ApiOptions, input: string | URL | Request, init?: RequestInit) => {
  const apiToken = opts.getApiToken && opts.getApiToken()

  const headers = new Headers({
    'Authorization': apiToken ? `Bearer ${apiToken}` : '',
    ...init?.headers,
  })

  if (!headers.get('Authorization')) headers.delete('Authorization')

  // 请求操作
  const resp = await fetch(input, {
    ...init,
    headers,
  })

  if (!resp.ok) {
    // Deserialize response body (if it is JSON, otherwise returns `{}`)
    const json: any = (await resp.json().catch(() => {})) || {}

    // 请求弹出层验证
    if (json.need_captcha) {
        // 请求需要验证码
        opts.onNeedCheckCaptcha && await opts.onNeedCheckCaptcha({
          data: { imgData: json.img_data, iframe: json.iframe }
        })
        return Fetch(opts, input, init) // retry
    }

    if (json.need_login) {
        // 请求需要管理员权限
        opts.onNeedCheckAdmin && await opts.onNeedCheckAdmin({})
        return Fetch(opts, input, init)
    }

    throw createError(resp.status, json)
  }

  return resp
}

export class FetchException extends Error implements FetchError {
  code: number = 0
  message: string = 'fetch error'
  data?: any
}

function createError(code: number, data: any): FetchException {
  const err = new FetchException()
  err.message = data.msg || data.message || 'fetch error'
  err.code = code
  err.data = data
  console.error(err)
  return err
}

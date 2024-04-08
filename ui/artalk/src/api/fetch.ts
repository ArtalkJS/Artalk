import 'abortcontroller-polyfill/dist/polyfill-patch-fetch'

import { FetchError } from '@/types'
import { ApiOptions } from './options'

export const Fetch = async (
  opts: ApiOptions,
  input: string | URL | Request,
  init?: RequestInit,
) => {
  const apiToken = opts.getApiToken && opts.getApiToken()

  const headers = new Headers({
    Authorization: apiToken ? `Bearer ${apiToken}` : '',
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

    let retry = false
    opts.handlers &&
      (await opts.handlers.get().reduce(async (promise, item) => {
        await promise
        if (json[item.action] === true) {
          await item.handler(json)
          retry = true
        }
      }, Promise.resolve()))

    if (retry) return Fetch(opts, input, init)
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

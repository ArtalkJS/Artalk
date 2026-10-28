import { ApiOptions } from './options'
import { FetchError } from '@/types'

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

  // Only comment/notification endpoints need this capability. In particular,
  // PV, statistics, assets and authentication requests never carry it.
  const url = new URL(input instanceof Request ? input.url : String(input), window.location.href)
  const base = new URL(opts.baseURL, window.location.href)
  const path = url.pathname.slice(base.pathname.replace(/\/$/, '').length)
  if (
    !headers.has('Authorization') &&
    url.origin === base.origin &&
    url.pathname.startsWith(`${base.pathname.replace(/\/$/, '')}/`) &&
    /^(\/comments(?:\/\d+)?|\/notifies(?:\/read)?|\/user)$/.test(path)
  ) {
    const create = path === '/comments' && init?.method?.toUpperCase() === 'POST'
    const guest = opts.getGuestToken?.(create)
    if (guest) headers.set('X-Artalk-Guest', guest)
  }

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

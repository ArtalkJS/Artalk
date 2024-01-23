/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface CommonApiVersionData {
  app: string
  commit_hash: string
  version: string
}

export interface CommonConfData {
  frontend_conf: CommonMap
  version: CommonApiVersionData
}

export interface CommonJSONResult {
  /** Data */
  data: any
  /** Message */
  msg: string
}

export type CommonMap = Record<string, any>

export interface EntityCookedComment {
  badge_color: string
  badge_name: string
  content: string
  content_marked: string
  date: string
  email_encrypted: string
  id: number
  ip_region: string
  is_allow_reply: boolean
  is_collapsed: boolean
  is_pending: boolean
  is_pinned: boolean
  link: string
  nick: string
  page_key: string
  page_url: string
  rid: number
  site_name: string
  ua: string
  user_id: number
  visible: boolean
  vote_down: number
  vote_up: number
}

export interface EntityCookedNotify {
  comment_id: number
  id: number
  is_emailed: boolean
  is_read: boolean
  read_link: string
  user_id: number
}

export interface EntityCookedPage {
  admin_only: boolean
  id: number
  key: string
  pv: number
  site_name: string
  title: string
  url: string
  vote_down: number
  vote_up: number
}

export interface EntityCookedSite {
  first_url: string
  id: number
  name: string
  urls: string[]
  urls_raw: string
}

export interface EntityCookedUser {
  badge_color: string
  badge_name: string
  email: string
  id: number
  is_admin: boolean
  link: string
  name: string
  receive_email: boolean
}

export interface EntityCookedUserForAdmin {
  badge_color: string
  badge_name: string
  comment_count: number
  email: string
  id: number
  is_admin: boolean
  is_in_conf: boolean
  last_ip: string
  last_ua: string
  link: string
  name: string
  receive_email: boolean
}

export type HandlerMap = Record<string, any>

export interface HandlerParamsCaptchaVerify {
  /** The captcha value to check */
  value: string
}

export interface HandlerParamsCommentCreate {
  /** The comment content */
  content: string
  /** The comment email */
  email: string
  /** The comment link */
  link?: string
  /** The comment name */
  name: string
  /** The comment page_key */
  page_key: string
  /** The comment page_title */
  page_title?: string
  /** The comment rid */
  rid?: number
  /** The site name of your content scope */
  site_name: string
  /** The comment ua */
  ua?: string
}

export interface HandlerParamsCommentUpdate {
  /** The comment content */
  content: string
  /** The comment email */
  email?: string
  /** The comment ip */
  ip?: string
  /** The comment is_collapsed */
  is_collapsed: boolean
  /** The comment is_pending */
  is_pending: boolean
  /** The comment is_pinned */
  is_pinned: boolean
  /** The comment link */
  link?: string
  /** The comment nick */
  nick?: string
  /** The comment page_key */
  page_key: string
  /** The comment rid */
  rid: number
  /** The site name of your content scope */
  site_name: string
  /** The comment ua */
  ua: string
}

export interface HandlerParamsEmailSend {
  /** The body of email */
  body: string
  /** The subject of email */
  subject: string
  /** The email address of the receiver */
  to_addr: string
}

export interface HandlerParamsNotifyReadAll {
  /** The user email */
  email: string
  /** The username */
  name: string
}

export interface HandlerParamsPageFetchAll {
  /** If not empty, only fetch pages of this site */
  site_name?: string
}

export interface HandlerParamsPagePV {
  /** The page key */
  page_key: string
  /** The page title */
  page_title?: string
  /** The site name of your content scope */
  site_name?: string
}

export interface HandlerParamsPageUpdate {
  /** Updated page admin_only option */
  admin_only: boolean
  /** Updated page key */
  key: string
  /** The site name of your content scope */
  site_name: string
  /** Updated page title */
  title: string
}

export interface HandlerParamsSettingApply {
  /** The content of the config file in YAML format */
  yaml: string
}

export interface HandlerParamsSiteCreate {
  /** The site name */
  name: string
  /** The site urls */
  urls: string[]
}

export interface HandlerParamsSiteUpdate {
  /** Updated site name */
  name: string
  /** Updated site urls */
  urls: string[]
}

export interface HandlerParamsTransferImport {
  /** Automatically answer yes for all questions. */
  assumeyes?: boolean
  /** The JSON data */
  json_data?: string
  /** The JSON file path */
  json_file?: string
  /** The target site name */
  target_site_name?: string
  /** The target site url */
  target_site_url?: string
  /** Enable URL resolver */
  url_resolver?: boolean
}

export interface HandlerParamsUserCreate {
  /** The user badge color (hex format) */
  badge_color?: string
  /** The user badge name */
  badge_name?: string
  /** The user email */
  email: string
  /** The user is an admin */
  is_admin: boolean
  /** The user link */
  link?: string
  /** The user name */
  name: string
  /** The user password */
  password?: string
  /** The user receive email */
  receive_email: boolean
}

export interface HandlerParamsUserLogin {
  /** The user email */
  email: string
  /** The username */
  name?: string
  /** The user password */
  password: string
}

export interface HandlerParamsUserUpdate {
  /** The user badge color (hex format) */
  badge_color?: string
  /** The user badge name */
  badge_name?: string
  /** The user email */
  email: string
  /** The user is an admin */
  is_admin: boolean
  /** The user link */
  link?: string
  /** The user name */
  name: string
  /** The user password */
  password?: string
  /** The user receive email */
  receive_email: boolean
}

export interface HandlerParamsVote {
  /** The user email */
  email?: string
  /** The username */
  name?: string
}

export interface HandlerResponseAdminUserList {
  count: number
  users: EntityCookedUserForAdmin[]
}

export interface HandlerResponseCaptchaGet {
  img_data: string
}

export interface HandlerResponseCaptchaStatus {
  is_pass: boolean
}

export interface HandlerResponseCommentCreate {
  badge_color: string
  badge_name: string
  content: string
  content_marked: string
  date: string
  email_encrypted: string
  id: number
  ip_region: string
  is_allow_reply: boolean
  is_collapsed: boolean
  is_pending: boolean
  is_pinned: boolean
  link: string
  nick: string
  page_key: string
  page_url: string
  rid: number
  site_name: string
  ua: string
  user_id: number
  visible: boolean
  vote_down: number
  vote_up: number
}

export interface HandlerResponseCommentList {
  comments: EntityCookedComment[]
  count: number
  page: EntityCookedPage
  roots_count: number
}

export interface HandlerResponseCommentUpdate {
  badge_color: string
  badge_name: string
  content: string
  content_marked: string
  date: string
  email_encrypted: string
  id: number
  ip_region: string
  is_allow_reply: boolean
  is_collapsed: boolean
  is_pending: boolean
  is_pinned: boolean
  link: string
  nick: string
  page_key: string
  page_url: string
  rid: number
  site_name: string
  ua: string
  user_id: number
  visible: boolean
  vote_down: number
  vote_up: number
}

export interface HandlerResponseNotifyList {
  count: number
  notifies: EntityCookedNotify[]
}

export interface HandlerResponsePageFetch {
  admin_only: boolean
  id: number
  key: string
  pv: number
  site_name: string
  title: string
  url: string
  vote_down: number
  vote_up: number
}

export interface HandlerResponsePageFetchStatus {
  /** The number of pages that have been fetched */
  done: number
  /** If the task is in progress */
  is_progress: boolean
  /** The message of the task status */
  msg: string
  /** The total number of pages */
  total: number
}

export interface HandlerResponsePageList {
  count: number
  pages: EntityCookedPage[]
}

export interface HandlerResponsePagePV {
  pv: number
}

export interface HandlerResponsePageUpdate {
  admin_only: boolean
  id: number
  key: string
  pv: number
  site_name: string
  title: string
  url: string
  vote_down: number
  vote_up: number
}

export interface HandlerResponseSettingGet {
  yaml: string
}

export interface HandlerResponseSettingTemplate {
  yaml: string
}

export interface HandlerResponseSiteCreate {
  first_url: string
  id: number
  name: string
  urls: string[]
  urls_raw: string
}

export interface HandlerResponseSiteList {
  count: number
  sites: EntityCookedSite[]
}

export interface HandlerResponseSiteUpdate {
  first_url: string
  id: number
  name: string
  urls: string[]
  urls_raw: string
}

export interface HandlerResponseTransferExport {
  /** The exported data which is a JSON string */
  artrans: string
}

export interface HandlerResponseTransferUpload {
  /** The uploaded file name which can be used to import */
  filename: string
}

export interface HandlerResponseUpload {
  file_name: string
  file_type: string
  public_url: string
}

export interface HandlerResponseUserCreate {
  badge_color: string
  badge_name: string
  comment_count: number
  email: string
  id: number
  is_admin: boolean
  is_in_conf: boolean
  last_ip: string
  last_ua: string
  link: string
  name: string
  receive_email: boolean
}

export interface HandlerResponseUserInfo {
  is_login: boolean
  notifies: EntityCookedNotify[]
  notifies_count: number
  user: EntityCookedUser
}

export interface HandlerResponseUserLogin {
  token: string
  user: EntityCookedUser
}

export interface HandlerResponseUserStatus {
  is_admin: boolean
  is_login: boolean
}

export interface HandlerResponseUserUpdate {
  badge_color: string
  badge_name: string
  comment_count: number
  email: string
  id: number
  is_admin: boolean
  is_in_conf: boolean
  last_ip: string
  last_ua: string
  link: string
  name: string
  receive_email: boolean
}

export interface HandlerResponseVote {
  down: number
  up: number
}

export type QueryParamsType = Record<string | number, any>
export type ResponseFormat = keyof Omit<Body, 'body' | 'bodyUsed'>

export interface FullRequestParams extends Omit<RequestInit, 'body'> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean
  /** request path */
  path: string
  /** content type of request body */
  type?: ContentType
  /** query params */
  query?: QueryParamsType
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseFormat
  /** request body */
  body?: unknown
  /** base url */
  baseUrl?: string
  /** request cancellation token */
  cancelToken?: CancelToken
}

export type RequestParams = Omit<FullRequestParams, 'body' | 'method' | 'query' | 'path'>

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string
  baseApiParams?: Omit<RequestParams, 'baseUrl' | 'cancelToken' | 'signal'>
  securityWorker?: (securityData: SecurityDataType | null) => Promise<RequestParams | void> | RequestParams | void
  customFetch?: typeof fetch
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown> extends Response {
  data: D
  error: E
}

type CancelToken = Symbol | string | number

export enum ContentType {
  Json = 'application/json',
  FormData = 'multipart/form-data',
  UrlEncoded = 'application/x-www-form-urlencoded',
  Text = 'text/plain',
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = '/api/v2'
  private securityData: SecurityDataType | null = null
  private securityWorker?: ApiConfig<SecurityDataType>['securityWorker']
  private abortControllers = new Map<CancelToken, AbortController>()
  private customFetch = (...fetchParams: Parameters<typeof fetch>) => fetch(...fetchParams)

  private baseApiParams: RequestParams = {
    credentials: 'same-origin',
    headers: {},
    redirect: 'follow',
    referrerPolicy: 'no-referrer',
  }

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig)
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data
  }

  protected encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key)
    return `${encodedKey}=${encodeURIComponent(typeof value === 'number' ? value : `${value}`)}`
  }

  protected addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key])
  }

  protected addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key]
    return value.map((v: any) => this.encodeQueryParam(key, v)).join('&')
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {}
    const keys = Object.keys(query).filter((key) => 'undefined' !== typeof query[key])
    return keys
      .map((key) => (Array.isArray(query[key]) ? this.addArrayQueryParam(query, key) : this.addQueryParam(query, key)))
      .join('&')
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery)
    return queryString ? `?${queryString}` : ''
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === 'object' || typeof input === 'string') ? JSON.stringify(input) : input,
    [ContentType.Text]: (input: any) => (input !== null && typeof input !== 'string' ? JSON.stringify(input) : input),
    [ContentType.FormData]: (input: any) =>
      Object.keys(input || {}).reduce((formData, key) => {
        const property = input[key]
        formData.append(
          key,
          property instanceof Blob
            ? property
            : typeof property === 'object' && property !== null
            ? JSON.stringify(property)
            : `${property}`,
        )
        return formData
      }, new FormData()),
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  }

  protected mergeRequestParams(params1: RequestParams, params2?: RequestParams): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    }
  }

  protected createAbortSignal = (cancelToken: CancelToken): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken)
      if (abortController) {
        return abortController.signal
      }
      return void 0
    }

    const abortController = new AbortController()
    this.abortControllers.set(cancelToken, abortController)
    return abortController.signal
  }

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken)

    if (abortController) {
      abortController.abort()
      this.abortControllers.delete(cancelToken)
    }
  }

  public request = async <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format,
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams =
      ((typeof secure === 'boolean' ? secure : this.baseApiParams.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {}
    const requestParams = this.mergeRequestParams(params, secureParams)
    const queryString = query && this.toQueryString(query)
    const payloadFormatter = this.contentFormatters[type || ContentType.Json]
    const responseFormat = format || requestParams.format

    return this.customFetch(`${baseUrl || this.baseUrl || ''}${path}${queryString ? `?${queryString}` : ''}`, {
      ...requestParams,
      headers: {
        ...(requestParams.headers || {}),
        ...(type && type !== ContentType.FormData ? { 'Content-Type': type } : {}),
      },
      signal: (cancelToken ? this.createAbortSignal(cancelToken) : requestParams.signal) || null,
      body: typeof body === 'undefined' || body === null ? null : payloadFormatter(body),
    }).then(async (response) => {
      const r = response as HttpResponse<T, E>
      r.data = null as unknown as T
      r.error = null as unknown as E

      const data = !responseFormat
        ? r
        : await response[responseFormat]()
            .then((data) => {
              if (r.ok) {
                r.data = data
              } else {
                r.error = data
              }
              return r
            })
            .catch((e) => {
              r.error = e
              return r
            })

      if (cancelToken) {
        this.abortControllers.delete(cancelToken)
      }

      if (!response.ok) throw data
      return data
    })
  }
}

/**
 * @title Artalk API
 * @version 2.0
 * @license MIT (https://github.com/ArtalkJS/Artalk/blob/master/LICENSE)
 * @baseUrl /api/v2
 * @contact API Support <artalkjs@gmail.com> (https://artalk.js.org)
 *
 * Artalk is a modern comment system based on Golang.
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  cache = {
    /**
 * @description Flush all cache on the server
 *
 * @tags Cache
 * @name FlushCache
 * @summary Flush Cache
 * @request POST:/cache/flush
 * @secure
 * @response `200` `(HandlerMap & {
    msg?: string,

})` OK
 * @response `400` `(HandlerMap & {
    msg?: string,

})` Bad Request
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 */
    flushCache: (params: RequestParams = {}) =>
      this.request<
        HandlerMap & {
          msg?: string
        },
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/cache/flush`,
        method: 'POST',
        secure: true,
        format: 'json',
        ...params,
      }),

    /**
 * @description Cache warming helps you to pre-load the cache to improve the performance of the first request
 *
 * @tags Cache
 * @name WarmUpCache
 * @summary Warm-Up Cache
 * @request POST:/cache/warm_up
 * @secure
 * @response `200` `(HandlerMap & {
    msg?: string,

})` OK
 * @response `400` `(HandlerMap & {
    msg?: string,

})` Bad Request
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 */
    warmUpCache: (params: RequestParams = {}) =>
      this.request<
        HandlerMap & {
          msg?: string
        },
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/cache/warm_up`,
        method: 'POST',
        secure: true,
        format: 'json',
        ...params,
      }),
  }
  captcha = {
    /**
 * @description Get a base64 encoded captcha image or a HTML page to verify for user
 *
 * @tags Captcha
 * @name GetCaptcha
 * @summary Get Captcha
 * @request GET:/captcha
 * @response `200` `HandlerResponseCaptchaGet` OK
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    getCaptcha: (params: RequestParams = {}) =>
      this.request<
        HandlerResponseCaptchaGet,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/captcha`,
        method: 'GET',
        format: 'json',
        ...params,
      }),

    /**
     * @description Get the status of the user's captcha verification
     *
     * @tags Captcha
     * @name GetCaptchaStatus
     * @summary Get Captcha Status
     * @request GET:/captcha/status
     * @response `200` `HandlerResponseCaptchaStatus` OK
     */
    getCaptchaStatus: (params: RequestParams = {}) =>
      this.request<HandlerResponseCaptchaStatus, any>({
        path: `/captcha/status`,
        method: 'GET',
        format: 'json',
        ...params,
      }),

    /**
 * @description Verify user enters correct captcha code
 *
 * @tags Captcha
 * @name VerifyCaptcha
 * @summary Verify Captcha
 * @request POST:/captcha/verify
 * @response `200` `HandlerMap` OK
 * @response `403` `(HandlerMap & {
    img_data?: string,

})` Forbidden
 */
    verifyCaptcha: (data: HandlerParamsCaptchaVerify, params: RequestParams = {}) =>
      this.request<
        HandlerMap,
        HandlerMap & {
          img_data?: string
        }
      >({
        path: `/captcha/verify`,
        method: 'POST',
        body: data,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),
  }
  comments = {
    /**
 * @description Get a list of comments by some conditions
 *
 * @tags Comment
 * @name GetComments
 * @summary Get Comment List
 * @request GET:/comments
 * @secure
 * @response `200` `HandlerResponseCommentList` OK
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    getComments: (
      query: {
        /** The user email */
        email?: string
        /** Enable flat_mode */
        flat_mode?: boolean
        /** The limit for pagination */
        limit?: number
        /** The username */
        name?: string
        /** The offset for pagination */
        offset?: number
        /** The comment page_key */
        page_key: string
        /** The scope of comments */
        scope?: 'page' | 'user' | 'site'
        /** Search keywords */
        search?: string
        /** The site name of your content scope */
        site_name?: string
        /** Sort by condition */
        sort_by?: 'date_asc' | 'date_desc' | 'vote'
        /** Message center show type */
        type?: 'all' | 'mentions' | 'mine' | 'pending'
        /** Only show comments by admin */
        view_only_admin?: boolean
      },
      params: RequestParams = {},
    ) =>
      this.request<
        HandlerResponseCommentList,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/comments`,
        method: 'GET',
        query: query,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
 * @description Create a new comment
 *
 * @tags Comment
 * @name CreateComment
 * @summary Create Comment
 * @request POST:/comments
 * @secure
 * @response `200` `HandlerResponseCommentCreate` OK
 * @response `400` `(HandlerMap & {
    msg?: string,

})` Bad Request
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    createComment: (comment: HandlerParamsCommentCreate, params: RequestParams = {}) =>
      this.request<
        HandlerResponseCommentCreate,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/comments`,
        method: 'POST',
        body: comment,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
 * @description Update a specific comment
 *
 * @tags Comment
 * @name UpdateComment
 * @summary Update Comment
 * @request PUT:/comments/{id}
 * @secure
 * @response `200` `HandlerResponseCommentUpdate` OK
 * @response `400` `(HandlerMap & {
    msg?: string,

})` Bad Request
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `404` `(HandlerMap & {
    msg?: string,

})` Not Found
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    updateComment: (id: number, comment: HandlerParamsCommentUpdate, params: RequestParams = {}) =>
      this.request<
        HandlerResponseCommentUpdate,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/comments/${id}`,
        method: 'PUT',
        body: comment,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
 * @description Delete a specific comment
 *
 * @tags Comment
 * @name DeleteComment
 * @summary Delete Comment
 * @request DELETE:/comments/{id}
 * @secure
 * @response `200` `HandlerMap` OK
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `404` `(HandlerMap & {
    msg?: string,

})` Not Found
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    deleteComment: (id: number, params: RequestParams = {}) =>
      this.request<
        HandlerMap,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/comments/${id}`,
        method: 'DELETE',
        secure: true,
        format: 'json',
        ...params,
      }),
  }
  conf = {
    /**
     * @description Get System Configs for UI
     *
     * @tags System
     * @name Conf
     * @summary Get System Configs
     * @request GET:/conf
     * @response `200` `CommonConfData` OK
     */
    conf: (params: RequestParams = {}) =>
      this.request<CommonConfData, any>({
        path: `/conf`,
        method: 'GET',
        format: 'json',
        ...params,
      }),
  }
  notifies = {
    /**
 * @description Get a list of notifies for user
 *
 * @tags Notify
 * @name GetNotifies
 * @summary Get Notifies
 * @request GET:/notifies
 * @response `200` `HandlerResponseNotifyList` OK
 * @response `400` `(HandlerMap & {
    msg?: string,

})` Bad Request
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    getNotifies: (
      query: {
        /** The user email */
        email: string
        /** The user name */
        name: string
      },
      params: RequestParams = {},
    ) =>
      this.request<
        HandlerResponseNotifyList,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/notifies`,
        method: 'GET',
        query: query,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
 * @description Mark all notifies as read for user
 *
 * @tags Notify
 * @name MarkAllNotifyRead
 * @summary Mark All Notifies as Read
 * @request POST:/notifies/read
 * @response `200` `HandlerMap` OK
 * @response `400` `(HandlerMap & {
    msg?: string,

})` Bad Request
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    markAllNotifyRead: (options: HandlerParamsNotifyReadAll, params: RequestParams = {}) =>
      this.request<
        HandlerMap,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/notifies/read`,
        method: 'POST',
        body: options,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
 * @description Mark specific notification as read for user
 *
 * @tags Notify
 * @name MarkNotifyRead
 * @summary Mark Notify as Read
 * @request POST:/notifies/{comment_id}/{notify_key}
 * @response `200` `HandlerMap` OK
 * @response `400` `(HandlerMap & {
    msg?: string,

})` Bad Request
 * @response `404` `(HandlerMap & {
    msg?: string,

})` Not Found
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    markNotifyRead: (commentId: number, notifyKey: string, params: RequestParams = {}) =>
      this.request<
        HandlerMap,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/notifies/${commentId}/${notifyKey}`,
        method: 'POST',
        format: 'json',
        ...params,
      }),
  }
  pages = {
    /**
 * @description Get a list of pages by some conditions
 *
 * @tags Page
 * @name GetPages
 * @summary Get Page List
 * @request GET:/pages
 * @secure
 * @response `200` `HandlerResponsePageList` OK
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 */
    getPages: (
      query?: {
        /** The limit for pagination */
        limit?: number
        /** The offset for pagination */
        offset?: number
        /** The site name of your content scope */
        site_name?: string
      },
      params: RequestParams = {},
    ) =>
      this.request<
        HandlerResponsePageList,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/pages`,
        method: 'GET',
        query: query,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
 * @description Fetch the data of all pages
 *
 * @tags Page
 * @name FetchAllPages
 * @summary Fetch All Pages Data
 * @request POST:/pages/fetch
 * @secure
 * @response `200` `HandlerMap` OK
 * @response `400` `(HandlerMap & {
    msg?: string,

})` Bad Request
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    fetchAllPages: (options: HandlerParamsPageFetchAll, params: RequestParams = {}) =>
      this.request<
        HandlerMap,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/pages/fetch`,
        method: 'POST',
        body: options,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description Get the status of the task of fetching all pages
     *
     * @tags Page
     * @name GetPageFetchStatus
     * @summary Get Pages Fetch Status
     * @request GET:/pages/fetch/status
     * @secure
     * @response `200` `HandlerResponsePageFetchStatus` OK
     */
    getPageFetchStatus: (params: RequestParams = {}) =>
      this.request<HandlerResponsePageFetchStatus, any>({
        path: `/pages/fetch/status`,
        method: 'GET',
        secure: true,
        format: 'json',
        ...params,
      }),

    /**
     * @description Increase and get the number of page views
     *
     * @tags Page
     * @name LogPv
     * @summary Increase Page Views (PV)
     * @request POST:/pages/pv
     * @response `200` `HandlerResponsePagePV` OK
     */
    logPv: (page: HandlerParamsPagePV, params: RequestParams = {}) =>
      this.request<HandlerResponsePagePV, any>({
        path: `/pages/pv`,
        method: 'POST',
        body: page,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
 * @description Update a specific page
 *
 * @tags Page
 * @name UpdatePage
 * @summary Update Page
 * @request PUT:/pages/{id}
 * @secure
 * @response `200` `HandlerResponsePageUpdate` OK
 * @response `400` `(HandlerMap & {
    msg?: string,

})` Bad Request
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `404` `(HandlerMap & {
    msg?: string,

})` Not Found
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    updatePage: (id: number, page: HandlerParamsPageUpdate, params: RequestParams = {}) =>
      this.request<
        HandlerResponsePageUpdate,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/pages/${id}`,
        method: 'PUT',
        body: page,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
 * @description Delete a specific page
 *
 * @tags Page
 * @name DeletePage
 * @summary Delete Page
 * @request DELETE:/pages/{id}
 * @secure
 * @response `200` `HandlerMap` OK
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `404` `(HandlerMap & {
    msg?: string,

})` Not Found
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    deletePage: (id: number, params: RequestParams = {}) =>
      this.request<
        HandlerMap,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/pages/${id}`,
        method: 'DELETE',
        secure: true,
        format: 'json',
        ...params,
      }),

    /**
 * @description Fetch the data of a specific page
 *
 * @tags Page
 * @name FetchPage
 * @summary Fetch Page Data
 * @request POST:/pages/{id}/fetch
 * @secure
 * @response `200` `HandlerResponsePageFetch` OK
 * @response `404` `(HandlerMap & {
    msg?: string,

})` Not Found
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    fetchPage: (id: number, params: RequestParams = {}) =>
      this.request<
        HandlerResponsePageFetch,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/pages/${id}/fetch`,
        method: 'POST',
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),
  }
  sendEmail = {
    /**
 * @description Send an email to test the email sender
 *
 * @tags System
 * @name SendEmail
 * @summary Send Email
 * @request POST:/send_email
 * @secure
 * @response `200` `HandlerMap` OK
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `500` `HandlerMap` Internal Server Error
 */
    sendEmail: (email: HandlerParamsEmailSend, params: RequestParams = {}) =>
      this.request<
        HandlerMap,
        | (HandlerMap & {
            msg?: string
          })
        | HandlerMap
      >({
        path: `/send_email`,
        method: 'POST',
        body: email,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),
  }
  settings = {
    /**
 * @description Get settings from app config file
 *
 * @tags System
 * @name GetSettings
 * @summary Get Settings
 * @request GET:/settings
 * @secure
 * @response `200` `HandlerResponseSettingGet` OK
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    getSettings: (params: RequestParams = {}) =>
      this.request<
        HandlerResponseSettingGet,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/settings`,
        method: 'GET',
        secure: true,
        format: 'json',
        ...params,
      }),

    /**
 * @description Apply settings and restart the server
 *
 * @tags System
 * @name ApplySettings
 * @summary Save and apply Settings
 * @request PUT:/settings
 * @secure
 * @response `200` `HandlerMap` OK
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    applySettings: (settings: HandlerParamsSettingApply, params: RequestParams = {}) =>
      this.request<
        HandlerMap,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/settings`,
        method: 'PUT',
        body: settings,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description Get config templates in different languages for rendering the settings page in the frontend
     *
     * @tags System
     * @name GetSettingsTemplate
     * @summary Get Settings Template
     * @request GET:/settings/template/{locale}
     * @secure
     * @response `200` `HandlerResponseSettingTemplate` OK
     */
    getSettingsTemplate: (locale: string, params: RequestParams = {}) =>
      this.request<HandlerResponseSettingTemplate, any>({
        path: `/settings/template/${locale}`,
        method: 'GET',
        secure: true,
        format: 'json',
        ...params,
      }),
  }
  sites = {
    /**
     * @description Get a list of sites by some conditions
     *
     * @tags Site
     * @name GetSites
     * @summary Get Site List
     * @request GET:/sites
     * @secure
     * @response `200` `HandlerResponseSiteList` OK
     */
    getSites: (params: RequestParams = {}) =>
      this.request<HandlerResponseSiteList, any>({
        path: `/sites`,
        method: 'GET',
        secure: true,
        format: 'json',
        ...params,
      }),

    /**
 * @description Create a new site
 *
 * @tags Site
 * @name CreateSite
 * @summary Create Site
 * @request POST:/sites
 * @secure
 * @response `200` `HandlerResponseSiteCreate` OK
 * @response `400` `(HandlerMap & {
    msg?: string,

})` Bad Request
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    createSite: (site: HandlerParamsSiteCreate, params: RequestParams = {}) =>
      this.request<
        HandlerResponseSiteCreate,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/sites`,
        method: 'POST',
        body: site,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description Update a specific site
     *
     * @tags Site
     * @name UpdateSite
     * @summary Update Site
     * @request PUT:/sites/{id}
     * @secure
     * @response `200` `HandlerResponseSiteUpdate` OK
     */
    updateSite: (id: number, site: HandlerParamsSiteUpdate, params: RequestParams = {}) =>
      this.request<HandlerResponseSiteUpdate, any>({
        path: `/sites/${id}`,
        method: 'PUT',
        body: site,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
 * @description Delete a specific site
 *
 * @tags Site
 * @name DeleteSite
 * @summary Site Delete
 * @request DELETE:/sites/{id}
 * @secure
 * @response `200` `HandlerMap` OK
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `404` `(HandlerMap & {
    msg?: string,

})` Not Found
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    deleteSite: (id: number, params: RequestParams = {}) =>
      this.request<
        HandlerMap,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/sites/${id}`,
        method: 'DELETE',
        secure: true,
        format: 'json',
        ...params,
      }),
  }
  stats = {
    /**
 * @description Get the statistics of various data analysis
 *
 * @tags Statistic
 * @name GetStats
 * @summary Statistic
 * @request GET:/stats/{type}
 * @response `200` `CommonJSONResult` OK
 * @response `400` `(HandlerMap & {
    msg?: string,

})` Bad Request
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `404` `(HandlerMap & {
    msg?: string,

})` Not Found
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    getStats: (
      type:
        | 'latest_comments'
        | 'latest_pages'
        | 'pv_most_pages'
        | 'comment_most_pages'
        | 'page_pv'
        | 'site_pv'
        | 'page_comment'
        | 'site_comment'
        | 'rand_comments'
        | 'rand_pages',
      query?: {
        /** The limit for pagination */
        limit?: number
        /** multiple page keys separated by commas */
        page_keys?: string
        /** The site name of your content scope */
        site_name?: string
      },
      params: RequestParams = {},
    ) =>
      this.request<
        CommonJSONResult,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/stats/${type}`,
        method: 'GET',
        query: query,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),
  }
  transfer = {
    /**
 * @description Export data from Artalk
 *
 * @tags Transfer
 * @name ExportArtrans
 * @summary Export Artrans
 * @request GET:/transfer/export
 * @secure
 * @response `200` `HandlerResponseTransferExport` OK
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    exportArtrans: (params: RequestParams = {}) =>
      this.request<
        HandlerResponseTransferExport,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/transfer/export`,
        method: 'GET',
        secure: true,
        format: 'json',
        ...params,
      }),

    /**
     * @description Import data to Artalk
     *
     * @tags Transfer
     * @name ImportArtrans
     * @summary Import Artrans
     * @request POST:/transfer/import
     * @secure
     * @response `200` `string` OK
     */
    importArtrans: (data: HandlerParamsTransferImport, params: RequestParams = {}) =>
      this.request<string, any>({
        path: `/transfer/import`,
        method: 'POST',
        body: data,
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
 * @description Upload a file to prepare to import
 *
 * @tags Transfer
 * @name UploadArtrans
 * @summary Upload Artrans
 * @request POST:/transfer/upload
 * @secure
 * @response `200` `(HandlerResponseTransferUpload & {
    filename?: string,

})` OK
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    uploadArtrans: (
      data: {
        /**
         * Upload file in preparation for import task
         * @format binary
         */
        file: File
      },
      params: RequestParams = {},
    ) =>
      this.request<
        HandlerResponseTransferUpload & {
          filename?: string
        },
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/transfer/upload`,
        method: 'POST',
        body: data,
        secure: true,
        type: ContentType.FormData,
        format: 'json',
        ...params,
      }),
  }
  upload = {
    /**
 * @description Upload file from this endpoint
 *
 * @tags Upload
 * @name Upload
 * @summary Upload
 * @request POST:/upload
 * @secure
 * @response `200` `HandlerResponseUpload` OK
 * @response `400` `(HandlerMap & {
    msg?: string,

})` Bad Request
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    upload: (
      data: {
        /**
         * Upload file
         * @format binary
         */
        file: File
      },
      params: RequestParams = {},
    ) =>
      this.request<
        HandlerResponseUpload,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/upload`,
        method: 'POST',
        body: data,
        secure: true,
        type: ContentType.FormData,
        format: 'json',
        ...params,
      }),
  }
  user = {
    /**
 * @description Get user info to prepare for login or check current user status
 *
 * @tags Account
 * @name GetUser
 * @summary Get User Info
 * @request GET:/user
 * @secure
 * @response `200` `HandlerResponseUserInfo` OK
 * @response `400` `(HandlerMap & {
    msg?: string,

})` Bad Request
 */
    getUser: (
      query?: {
        /** The user email */
        email?: string
        /** The username */
        name?: string
      },
      params: RequestParams = {},
    ) =>
      this.request<
        HandlerResponseUserInfo,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/user`,
        method: 'GET',
        query: query,
        secure: true,
        format: 'json',
        ...params,
      }),

    /**
 * @description Login user by name or email
 *
 * @tags Account
 * @name Login
 * @summary Get Access Token
 * @request POST:/user/access_token
 * @response `200` `HandlerResponseUserLogin` OK
 * @response `400` `(HandlerMap & {
    " data"?: {
    need_name_select?: (string)[],

},
    msg?: string,

})` Multiple users with the same email address are matched
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    login: (user: HandlerParamsUserLogin, params: RequestParams = {}) =>
      this.request<
        HandlerResponseUserLogin,
        | (HandlerMap & {
            ' data'?: {
              need_name_select?: string[]
            }
            msg?: string
          })
        | (HandlerMap & {
            msg?: string
          })
      >({
        path: `/user/access_token`,
        method: 'POST',
        body: user,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description Get user login status by header Authorization
     *
     * @tags Account
     * @name GetUserStatus
     * @summary Get Login Status
     * @request GET:/user/status
     * @secure
     * @response `200` `HandlerResponseUserStatus` OK
     */
    getUserStatus: (
      query?: {
        /** The user email */
        email?: string
        /** The username */
        name?: string
      },
      params: RequestParams = {},
    ) =>
      this.request<HandlerResponseUserStatus, any>({
        path: `/user/status`,
        method: 'GET',
        query: query,
        secure: true,
        format: 'json',
        ...params,
      }),
  }
  users = {
    /**
 * @description Create a new user
 *
 * @tags User
 * @name CreateUser
 * @summary Create User
 * @request POST:/users
 * @secure
 * @response `200` `HandlerResponseUserCreate` OK
 * @response `400` `(HandlerMap & {
    msg?: string,

})` Bad Request
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    createUser: (user: HandlerParamsUserCreate, params: RequestParams = {}) =>
      this.request<
        HandlerResponseUserCreate,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/users`,
        method: 'POST',
        body: user,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
 * @description Update a specific user
 *
 * @tags User
 * @name UpdateUser
 * @summary Update User
 * @request PUT:/users/{id}
 * @secure
 * @response `200` `HandlerResponseUserUpdate` OK
 * @response `400` `(HandlerMap & {
    msg?: string,

})` Bad Request
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `404` `(HandlerMap & {
    msg?: string,

})` Not Found
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    updateUser: (id: number, user: HandlerParamsUserUpdate, params: RequestParams = {}) =>
      this.request<
        HandlerResponseUserUpdate,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/users/${id}`,
        method: 'PUT',
        body: user,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
 * @description Delete a specific user
 *
 * @tags User
 * @name DeleteUser
 * @summary Delete User
 * @request DELETE:/users/{id}
 * @secure
 * @response `200` `HandlerMap` OK
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `404` `(HandlerMap & {
    msg?: string,

})` Not Found
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    deleteUser: (id: number, params: RequestParams = {}) =>
      this.request<
        HandlerMap,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/users/${id}`,
        method: 'DELETE',
        secure: true,
        format: 'json',
        ...params,
      }),

    /**
 * @description Get a list of users by some conditions
 *
 * @tags User
 * @name GetUsers
 * @summary Get User List
 * @request GET:/users/{type}
 * @secure
 * @response `200` `HandlerResponseAdminUserList` OK
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 */
    getUsers: (
      type?: 'all' | 'admin' | 'in_conf',
      query?: {
        /** The limit for pagination */
        limit?: number
        /** The offset for pagination */
        offset?: number
      },
      params: RequestParams = {},
    ) =>
      this.request<
        HandlerResponseAdminUserList,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/users/${type}`,
        method: 'GET',
        query: query,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),
  }
  version = {
    /**
     * @description Get the version of Artalk
     *
     * @tags System
     * @name GetVersion
     * @summary Get Version Info
     * @request GET:/version
     * @response `200` `CommonApiVersionData` OK
     */
    getVersion: (params: RequestParams = {}) =>
      this.request<CommonApiVersionData, any>({
        path: `/version`,
        method: 'GET',
        format: 'json',
        ...params,
      }),
  }
  votes = {
    /**
 * @description Sync the number of votes in the `comments` or `pages` data tables to keep them the same as the `votes` table
 *
 * @tags Vote
 * @name SyncVotes
 * @summary Sync Vote Data
 * @request POST:/votes/sync
 * @secure
 * @response `200` `HandlerMap` OK
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 */
    syncVotes: (params: RequestParams = {}) =>
      this.request<
        HandlerMap,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/votes/sync`,
        method: 'POST',
        secure: true,
        format: 'json',
        ...params,
      }),

    /**
 * @description Vote for a specific comment or page
 *
 * @tags Vote
 * @name Vote
 * @summary Vote
 * @request POST:/votes/{type}/{target_id}
 * @response `200` `HandlerResponseVote` OK
 * @response `403` `(HandlerMap & {
    msg?: string,

})` Forbidden
 * @response `404` `(HandlerMap & {
    msg?: string,

})` Not Found
 * @response `500` `(HandlerMap & {
    msg?: string,

})` Internal Server Error
 */
    vote: (
      type: 'comment_up' | 'comment_down' | 'page_up' | 'page_down',
      targetId: number,
      vote: HandlerParamsVote,
      params: RequestParams = {},
    ) =>
      this.request<
        HandlerResponseVote,
        HandlerMap & {
          msg?: string
        }
      >({
        path: `/votes/${type}/${targetId}`,
        method: 'POST',
        body: vote,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),
  }
}

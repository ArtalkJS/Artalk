import Artalk from '../Artalk'
import Utils from '../utils/index'
import Comment from "../components/Comment"

export interface ActionConf {
  data: object
  before: () => void
  after: () => void
  success: (msg: string, data: any) => void
  error: (msg: string, data: any) => void
}

export interface ActionAddConf extends ActionConf {
  data: {
    content: string
    nick: string
    email: string
    link: string
    rid: number
    page_key: string
    password: string
    captcha: string
  }
}

export interface ActionGetConf extends ActionConf {
  data: {
    page_key: string
    limit: number
    offset: number
  }
}

export default abstract class ApiDriver {
  public artalk: Artalk

  constructor (artalk: Artalk) {
    this.artalk = artalk
  }

  abstract actionAdd(conf: ActionAddConf): void
  abstract actionGet(conf: ActionGetConf): void

  public request (url: string, method: 'GET'|'POST', data: any, actionConf: ActionConf) {
    actionConf.before()

    const formData = new FormData()
    if (method === 'POST') {
      Object.keys(data).forEach(key => formData.set(key, data[key]))
    }

    let queryStr = ""
    if (method === 'GET') {
      queryStr = `?${Utils.urlQuerySerialize(data)}`
    }

    const xhr = new XMLHttpRequest()
    xhr.timeout = 5000
    xhr.open(method, this.artalk.conf.serverUrl + url + queryStr, true)

    xhr.onload = () => {
      actionConf.after()
      if (xhr.status >= 200 && xhr.status < 400) {
        const respData = JSON.parse(xhr.response)
        if (respData.success) {
          actionConf.success(respData.msg, respData.data)
        } else {
          actionConf.error(respData.msg, respData.data)
        }
      } else {
        actionConf.error(`服务器响应错误 Code: ${xhr.status}`, {})
      }
    };

    xhr.onerror = () => {
      actionConf.after()
      actionConf.error('网络错误', {})
    };

    if (method === 'POST') {
      xhr.send(formData)
    } else if (method === 'GET') {
      xhr.send()
    }
  }
}

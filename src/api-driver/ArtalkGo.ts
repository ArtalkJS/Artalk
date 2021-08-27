import Artalk from '../Artalk'
import ApiDriver, { ActionAddConf, ActionGetConf } from './ApiDriver';

export default class GoApi extends ApiDriver {
  constructor(artalk: Artalk) {
    super(artalk)
    console.log(artalk)
  }

  actionAdd(conf: ActionAddConf): void {
    const data = {
      name: conf.data.nick,
      email: conf.data.email,
      link: conf.data.link,
      content: conf.data.content,
      rid: conf.data.rid,
      page_key: conf.data.page_key,
      token: conf.data.password,
    }

    this.request("/api/add", "POST", data, conf)
  }
  actionGet(conf: ActionGetConf): void {
    const data = {
      page_key: conf.data.page_key,
      limit: conf.data.limit,
      offset: conf.data.offset,
    }

    this.request("/api/get", "GET", data, conf)
  }
}

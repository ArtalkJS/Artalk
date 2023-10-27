import type { ContextApi } from '~/types'
import 'abortcontroller-polyfill/dist/polyfill-patch-fetch'

import comment from './comment'
import page from './page'
import site from './site'
import user from './user'
import system from './system'
import captcha from './captcha'
import admin from './admin'
import upload from './upload'

const ApiComponents = {
  comment, page, site,
  user, system, captcha,
  admin, upload
}

class Api {
  protected ctx: ContextApi
  public get baseURL() {
    return `${this.ctx.conf.server}/api`
  }

  constructor (ctx: ContextApi) {
    this.ctx = ctx

    Object.entries(ApiComponents).forEach(([key, ApiComponent]) => {
      this[key] = new ApiComponent(this, this.ctx)
    })
  }
}

type TC = typeof ApiComponents
type AC = { [K in keyof TC]: InstanceType<TC[K]> }
interface Api extends AC {}

export default Api

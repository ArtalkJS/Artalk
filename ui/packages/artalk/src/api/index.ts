import { ApiOptions } from './_options'
import comment from './comment'
import page from './page'
import site from './site'
import user from './user'
import system from './system'
import captcha from './captcha'
import admin from './admin'
import upload from './upload'

const ApiComponents = {
  comment,
  page,
  site,
  user,
  system,
  captcha,
  admin,
  upload,
}

class Api {
  constructor(opts: ApiOptions) {
    Object.entries(ApiComponents).forEach(([key, ApiComponent]) => {
      this[key] = new ApiComponent(opts)
    })
  }
}

type TC = typeof ApiComponents
type AC = { [K in keyof TC]: InstanceType<TC[K]> }
interface Api extends AC {}

export default Api

import * as Utils from '../utils'
import { Checker } from '.'

const AdminChecker: Checker = {
  inputType: 'password',

  request(that, ctx, inputVal) {
    const data = {
      name: that.ctx.user.data.nick,
      email: that.ctx.user.data.email,
      password: inputVal
    }

    return (async () => {
      const resp = await that.ctx.getApi().user.login(data.name, data.email, data.password)
      return resp.token
    })()
  },

  body(that, ctx) {
    return Utils.createElement(`<span>${that.ctx.$t('adminCheck')}</span>`)
  },

  onSuccess(that, ctx, userToken, inputVal, formEl) {
    that.ctx.user.update({
      isAdmin: true,
      token: userToken
    })
    that.ctx.trigger('user-changed', that.ctx.user.data)
    that.ctx.listReload()
  },

  onError(that, ctx, err, inputVal, formEl) {}
}

export default AdminChecker

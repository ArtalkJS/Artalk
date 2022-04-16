import Api from '@/api'
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

    return new Api(that.ctx).login(data.name, data.email, data.password)
  },

  body(that, ctx) {
    return Utils.createElement('<span>敲入密码来验证管理员身份：</span>')
  },

  onSuccess(that, ctx, userToken, inputVal, formEl) {
    that.ctx.user.data.isAdmin = true
    that.ctx.user.data.token = userToken
    that.ctx.user.save()
    that.ctx.trigger('user-changed', that.ctx.user.data)
    that.ctx.trigger('list-reload')
  },

  onError(that, ctx, err, inputVal, formEl) {}
}

export default AdminChecker

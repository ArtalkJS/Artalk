import * as Utils from '../utils'
import User from '../user'
import { Checker } from '.'

const AdminChecker: Checker = {
  inputType: 'password',

  request(checker, inputVal) {
    const data = {
      name: User.data.nick,
      email: User.data.email,
      password: inputVal
    }

    return (async () => {
      const resp = await checker.getApi().user.login(data.name, data.email, data.password)
      return resp.token
    })()
  },

  body(checker) {
    return Utils.createElement(`<span>${checker.getCtx().$t('adminCheck')}</span>`)
  },

  onSuccess(checker, userToken, inputVal, formEl) {
    User.update({
      isAdmin: true,
      token: userToken
    })
    checker.getCtx().trigger('user-changed', User.data)
    checker.getCtx().listReload()
  },

  onError(checker, err, inputVal, formEl) {}
}

export default AdminChecker

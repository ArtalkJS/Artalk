import * as Utils from '../utils'
import { Checker } from '.'

const AdminChecker: Checker = {
  inputType: 'password',

  request(checker, inputVal) {
    const data = {
      name: checker.getUser().data.nick,
      email: checker.getUser().data.email,
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
    checker.getUser().update({
      isAdmin: true,
      token: userToken
    })
    checker.getCtx().trigger('user-changed', checker.getUser().data)
    checker.getCtx().listReload()
  },

  onError(checker, err, inputVal, formEl) {}
}

export default AdminChecker

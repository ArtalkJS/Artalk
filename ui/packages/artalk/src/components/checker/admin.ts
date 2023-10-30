import $t from '@/i18n'
import * as Utils from '@/lib/utils'
import User from '@/lib/user'
import type { Checker } from '.'

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
    return Utils.createElement(`<span>${$t('adminCheck')}</span>`)
  },

  onSuccess(checker, userToken, inputVal, formEl) {
    User.update({
      isAdmin: true,
      token: userToken
    })
    checker.getOpts().onReload()
  },

  onError(checker, err, inputVal, formEl) {}
}

export default AdminChecker

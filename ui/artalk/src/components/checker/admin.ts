import $t from '@/i18n'
import * as Utils from '@/lib/utils'
import type { Checker } from '.'

const AdminChecker: Checker<{ token: string }> = {
  inputType: 'password',

  request(checker, inputVal) {
    const data = {
      name: checker.getUser().getData().nick,
      email: checker.getUser().getData().email,
      password: inputVal
    }

    return checker.getApi().user.login(data.name, data.email, data.password)
  },

  body(checker) {
    return Utils.createElement(`<span>${$t('adminCheck')}</span>`)
  },

  onSuccess(checker, res, inputVal, formEl) {
    checker.getUser().update({
      isAdmin: true,
      token: res.token
    })
    checker.getOpts().onReload()
  },

  onError(checker, err, inputVal, formEl) {}
}

export default AdminChecker

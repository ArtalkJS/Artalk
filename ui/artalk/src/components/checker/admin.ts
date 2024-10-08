import type { Checker } from '.'
import $t from '@/i18n'
import * as Utils from '@/lib/utils'

const AdminChecker: Checker<{ token: string }> = {
  inputType: 'password',

  async request(checker, inputVal) {
    return (
      await checker.getApi().user.login({
        name: checker.getUser().getData().name,
        email: checker.getUser().getData().email,
        password: inputVal,
      })
    ).data
  },

  body(checker) {
    return Utils.createElement(`<span>${$t('adminCheck')}</span>`)
  },

  onSuccess(checker, res, inputVal, formEl) {
    checker.getUser().update({
      is_admin: true,
      token: res.token,
    })
    checker.getOpts().onReload()
  },

  onError(checker, err, inputVal, formEl) {},
}

export default AdminChecker

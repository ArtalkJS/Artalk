import { createStore } from 'solid-js/store'
import { loginByApiRes } from './lib/token-login'
import type { AuthContext } from './types'

export interface DialogPageLoginProps {
  ctx: AuthContext
  onRegisterNowClick: () => void
  changeTitle: (title: string) => void
  onComplete: () => void
}

export const DialogPageLogin = (props: DialogPageLoginProps) => {
  const { ctx, onRegisterNowClick } = props

  props.changeTitle(ctx.$t('signIn'))

  const [fields, setFields] = createStore({
    email: '',
    password: '',
  })

  const submitHandler = (e: SubmitEvent) => {
    e.preventDefault()
    ctx
      .getApi()
      .auth.loginByEmail(fields)
      .then((res) => {
        // console.log(res.data)
        loginByApiRes(ctx, res.data)
        props.onComplete()
      })
      .catch((err) => {
        alert(err.message)
      })
  }

  return (
    <div class="atk-login-page">
      <form class="atk-form" onSubmit={submitHandler}>
        <input
          type="email"
          placeholder={ctx.$t('email')}
          name="email"
          onInput={(e) => setFields('email', e.target.value.trim())}
          required
        />
        <input
          type="password"
          placeholder={ctx.$t('password')}
          name="password"
          onInput={(e) => setFields('password', e.target.value.trim())}
          required
        />
        <button type="submit">{ctx.$t('signIn')}</button>
        <div class="atk-form-bottom">
          {ctx.$t('noAccountPrompt')}{' '}
          <span class="atk-link" onClick={() => onRegisterNowClick()}>
            {ctx.$t('signUp')}
          </span>
        </div>
      </form>
    </div>
  )
}

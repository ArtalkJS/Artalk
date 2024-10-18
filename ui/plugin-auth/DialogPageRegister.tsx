import { Show, createSignal, createEffect } from 'solid-js'
import { createStore, reconcile } from 'solid-js/store'
import { VerifyButton } from './VerifyButton'
import { loginByApiRes } from './lib/token-login'
import type { AuthContext } from './types'

export interface DialogPageRegisterProps {
  ctx: AuthContext
  onLoginNowClick: () => void
  changeTitle: (title: string) => void
  onComplete: () => void
}

export const DialogPageRegister = (props: DialogPageRegisterProps) => {
  const { ctx, onLoginNowClick } = props

  const [mode, setMode] = createSignal<'register' | 'forget'>('register')

  const [fields, setFields] = createStore({
    email: '',
    code: '',
    username: '',
    password: '',
  })

  createEffect(() => {
    if (mode() === 'register') props.changeTitle(ctx.$t('signUp'))
    else props.changeTitle(ctx.$t('resetPassword'))
    setFields({ email: '', code: '', username: '', password: '' })
  })

  const onSubmit = (e: SubmitEvent) => {
    e.preventDefault()

    const data =
      mode() === 'register'
        ? {
            email: fields.email,
            code: fields.code,
            name: fields.username,
            password: fields.password,
          }
        : {
            email: fields.email,
            code: fields.code,
            password: fields.password,
          }

    ctx
      .getApi()
      .auth.registerByEmail(data)
      .then((res) => {
        if (mode() === 'register') {
          loginByApiRes(ctx, res.data)
          props.onComplete()
        } else {
          onLoginNowClick()
        }
      })
      .catch((err) => {
        alert(err.message)
      })
  }

  return (
    <div class="atk-register-page">
      <form class="atk-form" onSubmit={onSubmit}>
        <Show when={mode() === 'register'}>
          <input
            type="text"
            placeholder={ctx.$t('username')}
            name="username"
            value={fields.username}
            onInput={(e) => setFields('username', e.target.value.trim())}
            required
          />
          <input
            type="password"
            placeholder={ctx.$t('password')}
            name="password"
            value={fields.password}
            onInput={(e) => setFields('password', e.target.value.trim())}
            required
          />
        </Show>
        <input
          type="email"
          placeholder={ctx.$t('email')}
          name="email"
          value={fields.email}
          onInput={(e) => setFields('email', e.target.value.trim())}
          required
        />
        <div class="atk-input-grp">
          <input
            type="text"
            placeholder={ctx.$t('verificationCode')}
            name="code"
            value={fields.code}
            onInput={(e) => setFields('code', e.target.value.trim())}
            required
          />
          <VerifyButton ctx={ctx} getEmail={() => fields.email} />
        </div>
        <Show when={mode() === 'forget'}>
          <input
            type="password"
            placeholder={ctx.$t('resetPassword')}
            name="re_password"
            value={fields.password}
            onInput={(e) => setFields('password', e.target.value.trim())}
            required
          />
        </Show>
        <button type="submit">{ctx.$t('nextStep')}</button>
        <div class="atk-form-bottom">
          {ctx.$t('haveAccountPrompt')}{' '}
          <span class="atk-link" onclick={() => onLoginNowClick()}>
            {ctx.$t('signIn')}
          </span>{' '}
          |{' '}
          <span>
            {mode() === 'register' && (
              <span class="atk-link" onClick={() => setMode('forget')}>
                {ctx.$t('forgetPassword')}
              </span>
            )}
            {mode() === 'forget' && (
              <span class="atk-link" onClick={() => setMode('register')}>
                {ctx.$t('signUp')}
              </span>
            )}
          </span>
        </div>
      </form>
    </div>
  )
}

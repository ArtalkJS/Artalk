import { LocalUser } from 'artalk'
import { createResource, createSignal, Resource, Show } from 'solid-js'
import { createStore } from 'solid-js/store'
import type { AuthContext } from './types'
import { Dialog } from './Dialog'
import { VerifyButton } from './VerifyButton'

export interface UserProfileDialogProps {
  ctx: AuthContext
  onClose: () => void
}

export const UserProfileDialog = (props: UserProfileDialogProps) => {
  const { ctx, onClose } = props

  const [title, setTitle] = createSignal('')

  const [userInfo] = createResource<LocalUser>(async () => {
    const { data } = await ctx.getApi().user.getUser()
    if (!data.is_login) {
      ctx.getUser().logout()
      onClose()
      throw new Error('No login')
    }
    return { ...data.user, token: '' }
  })

  const pages = {
    basic: () => UserBasicProfileForm(ctx, () => setPage('password'), setTitle, userInfo, onClose),
    password: () => UserChangePasswordForm(ctx, setTitle, userInfo, onClose),
  }

  const homePage = 'basic'
  const [page, setPage] = createSignal<keyof typeof pages>(homePage)

  const showBackBtn = () => page() !== homePage
  const backHome = () => setPage(homePage)

  return (
    <Dialog showBackBtn={showBackBtn} onBack={backHome} onClose={onClose} title={title}>
      {() => pages[page()]()}
    </Dialog>
  )
}

const UserBasicProfileForm = (
  ctx: AuthContext,
  onChangePassword: () => void,
  setTitle: (t: string) => void,
  user: Resource<LocalUser>,
  onClose: () => void,
) => {
  setTitle(ctx.$t('userProfile'))

  const [fields, setFields] = createStore({
    name: user()?.name || '',
    email: user()?.email || '',
    link: user()?.link || '',
    code: '',
  })

  const needVerify = () => !user.loading && fields.email !== user()?.email

  const submitHandler = (e: SubmitEvent) => {
    e.preventDefault()

    ctx
      .getApi()
      .user.updateProfile({
        ...fields,
      })
      .then(({ data: { user } }) => {
        ctx.getUser().update(user)
        onClose()
      })
      .catch((e) => {
        window.alert(e.message)
      })
  }

  return (
    <div>
      <form class="atk-form" onSubmit={submitHandler}>
        <label>{ctx.$t('nick')}</label>
        <input
          type="text"
          name="name"
          required
          value={fields.name}
          onInput={(e) => setFields('name', e.target.value.trim())}
        />

        <label>{ctx.$t('email')}</label>
        <input
          type="email"
          name="email"
          required
          value={fields.email}
          onInput={(e) => {
            setFields('code', '')
            setFields('email', e.target.value.trim())
          }}
        />

        {needVerify() && (
          <>
            <label>{ctx.$t('verificationCode')}</label>
            <div class="atk-input-grp">
              <input
                type="text"
                name="code"
                value={fields.code}
                onInput={(e) => setFields('code', e.target.value.trim())}
                required
              />
              <VerifyButton ctx={ctx} getEmail={() => fields.email || ''} />
            </div>
          </>
        )}

        <label>{ctx.$t('link')}</label>
        <input
          type="text"
          name="link"
          value={fields.link}
          onInput={(e) => setFields('link', e.target.value.trim())}
        />
        <button type="submit">{ctx.$t('save')}</button>
        <div class="atk-form-bottom">
          <span class="atk-link" onclick={onChangePassword}>
            {ctx.$t('changePassword')}
          </span>
        </div>
      </form>
    </div>
  )
}

const UserChangePasswordForm = (
  ctx: AuthContext,
  setTitle: (t: string) => void,
  user: Resource<LocalUser>,
  onClose: () => void,
) => {
  setTitle(ctx.$t('changePassword'))

  const [fields, setFields] = createStore({
    code: '',
    newPassword: '',
    newPassword2: '',
  })

  const submitHandler = (e: SubmitEvent) => {
    e.preventDefault()

    if (fields.newPassword !== fields.newPassword2) {
      window.alert(ctx.$t('passwordMismatch'))
      return
    }

    ctx
      .getApi()
      .auth.registerByEmail({
        name: user()?.name || '',
        email: user()?.email || '',
        password: fields.newPassword,
        code: fields.code,
      })
      .then((res) => {
        ctx.getUser().logout()
        onClose()
      })
      .catch((err) => {
        alert(err.message)
      })
  }

  return (
    <div>
      <form class="atk-form" onSubmit={submitHandler}>
        <label>{ctx.$t('verificationCode')}</label>
        <div class="atk-input-grp">
          <input
            type="text"
            name="code"
            value={fields.code}
            onInput={(e) => setFields('code', e.target.value.trim())}
            required
          />
          <VerifyButton ctx={ctx} getEmail={() => user()!.email} />
        </div>

        <label>{ctx.$t('password')}</label>
        <input
          type="password"
          name="newPassword"
          onInput={(e) => setFields('newPassword', e.target.value.trim())}
          required
        />

        <label>{ctx.$t('confirmPassword')}</label>
        <input
          type="password"
          name="newPassword2"
          onInput={(e) => setFields('newPassword2', e.target.value.trim())}
          required
        />

        <button type="submit">{ctx.$t('nextStep')}</button>
      </form>
    </div>
  )
}

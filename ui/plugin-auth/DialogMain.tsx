import type { ContextApi } from 'artalk'
import { Show, createResource, createSignal } from 'solid-js'
import { Dialog } from './Dialog'
import { DialogMethods } from './DialogMethods'
import { DialogPageLogin } from './DialogPageLogin'
import { DialogPageRegister } from './DialogPageRegister'
import { createLayer } from './lib/layer'
import { DialogMerge } from './merge/DialogMerge'

interface DialogMainProps {
  ctx: ContextApi
  onClose: () => void
  onSkip: () => void
}

export const DialogMain = (props: DialogMainProps) => {
  const { ctx, onClose, onSkip, ...others } = props

  const [allowAnonymous, setAllowAnonymous] = createSignal<boolean>(false)
  const [methods] = createResource(async () => {
    const { data } = await ctx.getApi().conf.getSocialLoginProviders()
    setAllowAnonymous(data.anonymous)
    return data.providers
      .map<LoginMethod>(({ name, label, icon, path }) => {
        const mm: LoginMethod = { name, label, icon, link: path }
        if (name === 'email') mm.onClick = () => setPage('login')
        return mm
      })
      .sort((a, b) => {
        // email always on top
        if (a.name === 'email') return -1
        if (b.name === 'email') return 1
        // others by label
        return a.label.localeCompare(b.label)
      })
  })

  const [title, setTitle] = createSignal<string>('Login')
  const onComplete = () => {
    onClose()

    ctx.get('editor').getUI().$header.style.display = 'none'

    // Check need to merge
    ctx
      .getApi()
      .auth.checkDataMerge()
      .then(({ data }) => {
        if (data.need_merge) {
          setTimeout(() => {
            createLayer(ctx).show((layer) => (
              <DialogMerge ctx={ctx} onClose={() => layer.destroy()} usernames={data.user_names} />
            ))
          }, 500)
        }
      })
  }

  const pages = {
    methods: () => (
      <DialogMethods ctx={ctx} methods={methods} changeTitle={setTitle} onComplete={onComplete} />
    ),
    login: () => (
      <DialogPageLogin
        ctx={ctx}
        onRegisterNowClick={() => setPage('register')}
        changeTitle={setTitle}
        onComplete={onComplete}
      />
    ),
    register: () => (
      <DialogPageRegister
        ctx={ctx}
        onLoginNowClick={() => setPage('login')}
        changeTitle={setTitle}
        onComplete={onComplete}
      />
    ),
  }

  const homePage = 'methods'
  const [page, setPage] = createSignal<keyof typeof pages>(homePage)
  const showBackBtn = () => page() !== homePage && !(page() == 'login' && methods()?.length === 1)
  const backHome = () => {
    setPage(homePage)
  }

  return (
    <Dialog
      showBackBtn={showBackBtn}
      onBack={backHome}
      onClose={onClose}
      title={title}
      extras={() => (
        <Show when={allowAnonymous()}>
          <div
            class="atk-skip-btn"
            onClick={() => {
              onSkip()
              onClose()
            }}
          >
            {ctx.$t('skipNotVerify')}
          </div>
        </Show>
      )}
    >
      {() => pages[page()]()}
    </Dialog>
  )
}

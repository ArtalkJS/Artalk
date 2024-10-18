import { createResource, createSignal } from 'solid-js'
import { Dialog } from './Dialog'
import { DialogMethods } from './DialogMethods'
import { DialogPageLogin } from './DialogPageLogin'
import { DialogPageRegister } from './DialogPageRegister'
import { createLayer } from './lib/layer'
import { DialogMerge } from './merge/DialogMerge'
import { fetchMethods, LoginMethod } from './lib/methods'
import type { AuthContext } from './types'

interface DialogMainProps {
  ctx: AuthContext
  onClose: () => void
  onSkip: () => void
}

export const DialogMain = (props: DialogMainProps) => {
  const { ctx, onClose, onSkip, ...others } = props

  const [methods] = createResource(async () => {
    return fetchMethods(ctx.getApi())
      .then((m) =>
        m.map<LoginMethod>((mm) => {
          if (mm.name === 'email') mm.onClick = () => setPage('login')
          if (mm.name === 'skip') {
            mm.label = ctx.$t('skipVerify')
            mm.onClick = () => {
              onSkip()
              onClose()
            }
          }
          return mm
        }),
      )
      .catch((e) => {
        alert('Failed to fetch login methods: ' + e.message)
        onClose()
        return []
      })
  })

  const [title, setTitle] = createSignal<string>('Login')
  const onComplete = () => {
    onClose()

    ctx.getEditor().getUI().$header.style.display = 'none'

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
    <Dialog showBackBtn={showBackBtn} onBack={backHome} onClose={onClose} title={title}>
      {() => pages[page()]()}
    </Dialog>
  )
}

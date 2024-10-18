import { createEffect, createMemo, createSignal, For, Resource } from 'solid-js'
import { startOAuthLogin } from './lib/oauth-login'
import { loginByToken } from './lib/token-login'
import { LoginMethod } from './lib/methods'
import type { AuthContext } from './types'

export interface DialogMethodsProps {
  ctx: AuthContext
  methods: Resource<LoginMethod[]>
  changeTitle: (title: string) => void
  onComplete: () => void
}

export const DialogMethods = (props: DialogMethodsProps) => {
  const { ctx, methods } = props
  const [isLoading, setIsLoading] = createSignal(true)

  props.changeTitle(ctx.$t('signIn'))

  const clickHandler = (m: LoginMethod) => {
    if (m.onClick) m.onClick()
    else if (m.link)
      (async () => {
        const url = /^(http|https):\/\//.test(m.link!)
          ? m.link!
          : `${ctx.getConf().get().server}${m.link}`
        const { token } = await startOAuthLogin(ctx, url)
        loginByToken(ctx, token)
        props.onComplete()
      })()
  }

  createEffect(() => {
    if (methods()?.length) setIsLoading(false)
    else setIsLoading(true)

    if (methods()?.length === 1) {
      const m = methods()![0]
      clickHandler(m)
    }
  })

  const $methods = (
    <div class="atk-methods">
      <For each={methods()}>
        {(m) => (
          <div
            class="atk-method-item atk-fade-in"
            data-method={m.name}
            onClick={() => clickHandler(m)}
          >
            <div
              class="atk-method-icon"
              style={{
                'background-image': `url('${m.icon}')`,
              }}
            ></div>
            <div class="atk-method-text">{m.label}</div>
          </div>
        )}
      </For>
      {isLoading() &&
        [...Array(3)].map((_, i) => (
          <div class="atk-method-item atk-methods-loading atk-fade-in">
            <div class="atk-method-icon"></div>
            <div class="atk-method-text"></div>
          </div>
        ))}
    </div>
  )

  return <div class="atk-methods-page">{$methods}</div>
}

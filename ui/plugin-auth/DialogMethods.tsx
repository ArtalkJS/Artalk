import { createEffect, createMemo, createSignal, For, Resource } from 'solid-js'
import type { ContextApi } from 'artalk'
import { startOAuthLogin } from './lib/oauth-login'
import { loginByToken } from './lib/token-login'

export interface DialogMethodsProps {
  ctx: ContextApi
  methods: Resource<LoginMethod[]>
  changeTitle: (title: string) => void
  onComplete: () => void
}

export const DialogMethods = (props: DialogMethodsProps) => {
  const { ctx, methods } = props

  props.changeTitle(ctx.$t('signIn'))

  const clickHandler = (m: LoginMethod) => {
    if (m.onClick) m.onClick()
    else if (m.link)
      (async () => {
        const url = /^(http|https):\/\//.test(m.link!)
          ? m.link!
          : `${ctx.getConf().server}${m.link}`
        const { token } = await startOAuthLogin(ctx, url)
        loginByToken(ctx, token)
        props.onComplete()
      })()
  }

  createEffect(() => {
    if (methods()?.length === 1) {
      const m = methods()![0]
      clickHandler(m)
    }
  })

  const $methods = (
    <div class="atk-methods">
      <For each={methods()}>
        {(m) => (
          <div class="atk-method-item" data-method={m.name} onClick={() => clickHandler(m)}>
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
    </div>
  )

  return <div class="atk-methods-page">{$methods}</div>
}

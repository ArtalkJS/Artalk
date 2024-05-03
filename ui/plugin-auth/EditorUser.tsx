import type { ContextApi } from 'artalk'
import { Show, onCleanup, createSignal } from 'solid-js'
import { render } from 'solid-js/web'

const EditorUser = ({ ctx }: { ctx: ContextApi }) => {
  const logoutHandler = () => {
    window.confirm(ctx.$t('logoutConfirm')) &&
      ctx.get('user').update({
        token: '',
        nick: '',
        email: '',
        link: '',
        isAdmin: false,
      })
  }

  const getUser = () => ({ ...ctx.get('user').getData() }) // Must clone the object to avoid reactivity problem

  const [user, setUser] = createSignal(getUser())

  const userChangedHandler = (u: any) => {
    setUser(getUser())
  }
  ctx.on('user-changed', userChangedHandler)

  onCleanup(() => {
    ctx.off('user-changed', userChangedHandler)
  })

  return (
    <div class="atk-editor-user-wrap">
      <Show when={user().token}>
        <div class="atk-editor-user">
          <div class="atk-user-profile-btn atk-user-btn">
            {/* <div class="atk-avatar" style={{ "background-image": `url('https://avatars.githubusercontent.com/u/76841221?s=200&v=4')` }}></div> */}
            <div class="atk-name">{user().nick}</div>
          </div>
          <div class="atk-logout atk-user-btn" onClick={logoutHandler} aria-label="Logout">
            <svg
              stroke="currentColor"
              fill="currentColor"
              height="1em"
              width="1em"
              stroke-width="0"
              viewBox="-6 -6 36 36"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path d="M5 22C4.44772 22 4 21.5523 4 21V3C4 2.44772 4.44772 2 5 2H19C19.5523 2 20 2.44772 20 3V6H18V4H6V20H18V18H20V21C20 21.5523 19.5523 22 19 22H5ZM18 16V13H11V11H18V8L23 12L18 16Z"></path>
            </svg>
          </div>
        </div>
      </Show>
    </div>
  )
}

export const RenderEditorUser = (ctx: ContextApi) => {
  const editor = ctx.get('editor')
  const findEl = () => editor.getEl().querySelector('.atk-editor-user-wrap')

  if (!findEl()) {
    const el = document.createElement('div')
    render(() => <EditorUser ctx={ctx} />, el)
    editor.getUI().$header.after(el)
  }
}

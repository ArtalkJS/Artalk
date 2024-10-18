import { createSignal } from 'solid-js'
import { DialogFooter } from '../DialogFooter'
import { loginByToken } from '../lib/token-login'
import type { AuthContext } from '../types'

interface DialogMergePageConfirmProps {
  ctx: AuthContext
  usernames: string[]
  onClose: () => void
}

export const DialogMergePageConfirm = (props: DialogMergePageConfirmProps) => {
  const { usernames, ctx, onClose } = props
  const [targetName, setTargetName] = createSignal('')

  const selectUsername = (u: string) => {
    setTargetName(targetName() !== u ? u : '')
  }

  const confirmMerge = () => {
    ctx
      .getApi()
      .auth.applyDataMerge({
        user_name: targetName(),
      })
      .then(({ data }) => {
        loginByToken(ctx, data.user_token)
        // console.log(data)
        onClose()
      })
  }

  return (
    <div class="atk-merge-confirm-page">
      <div class="atk-text">{ctx.$t('accountMergeNotice')}</div>
      <div class="atk-usernames">
        {usernames.map((u) => (
          <div
            class="atk-item"
            classList={{ active: targetName() === u }}
            onClick={() => selectUsername(u)}
          >
            <div class="atk-username">{u}</div>
          </div>
        ))}
      </div>
      <div class="atk-text">
        {!targetName()
          ? ctx.$t('accountMergeSelectOne')
          : ctx.$t('accountMergeConfirm', { id: targetName() })}
      </div>
      <DialogFooter
        noText={ctx.$t('dismiss')}
        yesText={ctx.$t('merge')}
        onNo={() => {
          onClose()
        }}
        onYes={() => {
          confirmMerge()
        }}
        yesDisabled={() => !targetName()}
      />
    </div>
  )
}

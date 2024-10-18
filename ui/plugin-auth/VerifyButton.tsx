import { createSignal } from 'solid-js'
import type { AuthContext } from './types'

interface VerifyButtonProps {
  ctx: AuthContext
  getEmail: () => string
  onSend?: () => void
}

export const VerifyButton = (props: VerifyButtonProps) => {
  const { ctx, onSend, getEmail } = props
  const [btnText, setBtnText] = createSignal(ctx.$t('verifySend'))

  const TIMEOUT = 60

  let sent = false
  let timer: any = null
  let duration = TIMEOUT

  const reset = () => {
    timer && clearInterval(timer)
    timer = null
    sent = false
    duration = TIMEOUT
  }

  const clickHandler = async () => {
    if (sent) return
    sent = true

    try {
      await ctx.getApi().auth.sendVerifyEmail({
        email: getEmail(),
      })

      timer && clearInterval(timer)
      timer = setInterval(() => {
        if (duration <= 0) {
          reset()
          setBtnText(ctx.$t('verifyResend'))
          return
        }
        duration--
        setBtnText(ctx.$t('waitSeconds', { seconds: `${duration}` }))
      }, 1000)
      onSend?.()
    } catch (e: any) {
      sent = false
      // console.log(e.message)
      alert(e.message)
    }
  }

  return (
    <div class="atk-send-verify-btn atk-input-grp-btn" onClick={clickHandler}>
      {btnText()}
    </div>
  )
}

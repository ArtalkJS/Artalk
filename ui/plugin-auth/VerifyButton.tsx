import type { ContextApi } from 'artalk'
import { createSignal } from 'solid-js'

interface VerifyButtonProps {
  ctx: ContextApi
  getEmail: () => string
  onSend?: () => void
}

export const VerifyButton = (props: VerifyButtonProps) => {
  const { ctx, onSend, getEmail } = props
  const [btnText, setBtnText] = createSignal(ctx.$t('verifySend'))

  let sent = false
  let timer: any = null
  let duration = 60
  const clickHandler = () => {
    if (sent) return
    sent = true
    ctx
      .getApi()
      .auth.sendVerifyEmail({
        email: getEmail(),
      })
      .then(() => {
        timer && clearInterval(timer)
        timer = setInterval(() => {
          if (duration <= 0) {
            clearInterval(timer)
            setBtnText(ctx.$t('verifyResend'))
            sent = false
            return
          }
          duration--
          setBtnText(ctx.$t('waitSeconds', { seconds: `${duration}` }))
        }, 1000)
        onSend?.()
      })
      .catch((e) => {
        sent = false
        // console.log(e.message)
        alert(e.message)
      })
  }

  return (
    <div class="atk-send-verify-btn atk-input-grp-btn" onClick={clickHandler}>
      {btnText()}
    </div>
  )
}

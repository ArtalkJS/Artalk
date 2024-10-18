import type { AuthContext } from '../types'

let watchTimer: any = null
let messageHandler: ((evt: MessageEvent) => void) | null = null

const clearListener = () => {
  watchTimer && clearInterval(watchTimer)
  messageHandler && window.removeEventListener('message', messageHandler)
}

export const startOAuthLogin = (ctx: AuthContext, url: string) => {
  clearListener()

  const width = 1020
  const height = 618
  const left = (window.innerWidth - width) / 2
  const top = (window.innerHeight - height) / 2

  const handler = window.open(
    url,
    '_blank',
    `width=${width},height=${height},left=${left},top=${top},scrollbars=no,resizable=no,status=no,location=no,toolbar=no,menubar=no`,
  )

  return new Promise<{ token: string }>((resolve, reject) => {
    watchTimer = setInterval(() => {
      if (handler?.closed) {
        clearListener()
        reject(new Error('Login canceled'))
        return
      }
      handler?.postMessage({ type: 'ATK_LOGIN' }, '*')
    }, 1000)

    messageHandler = ({ data }: MessageEvent): void => {
      if (data.type === 'ATK_AUTH_CALLBACK' && data.payload) {
        clearListener()

        handler?.close()
        // console.log(data.payload)
        resolve({ token: data.payload })
      }
    }
    window.addEventListener('message', messageHandler)
  })
}

import './adsense.css'

let ready: Promise<string | undefined> | undefined

export function loadAdSense() {
  if (!import.meta.env.PROD || location.hostname !== 'artalk.js.org') return

  if (ready) return ready

  const client = document.querySelector<HTMLMetaElement>(
    'meta[name="google-adsense-account"]',
  )?.content
  if (!client) return

  ready = new Promise<string | undefined>((resolve) => {
    const script = document.createElement('script')
    script.async = true
    script.crossOrigin = 'anonymous'
    script.src = 'https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=' + client
    script.onload = () => resolve(client)
    script.onerror = () => resolve(undefined)
    document.head.append(script)
  })
  return ready
}

export function showDisplayAd(container: HTMLElement | null) {
  if (!container) return
  loadAdSense()?.then((client) => {
    if (!client || !container.isConnected || container.firstChild) return

    const ad = document.createElement('ins')
    ad.className = 'adsbygoogle'
    ad.dataset.adClient = client
    ad.dataset.adSlot = '7948324420'
    container.append(ad)

    const host = window as Window & { adsbygoogle?: object[] }
    host.adsbygoogle ??= []
    host.adsbygoogle.push({})
  })
}

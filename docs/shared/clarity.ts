export function loadClarity() {
  if (!import.meta.env.PROD || location.hostname !== 'artalk.js.org') return
  if ('globalPrivacyControl' in navigator && navigator.globalPrivacyControl) return

  const host = window as Window & { clarity?: (...args: unknown[]) => void }
  if (host.clarity) return

  const queue: unknown[][] = []
  host.clarity = Object.assign((...args: unknown[]) => queue.push(args), { q: queue })

  // The project's Cookie setting is off. Clarity reads Google Consent Mode when available.
  const script = document.createElement('script')
  script.async = true
  script.src = 'https://www.clarity.ms/tag/ziptdojfsg'
  document.head.append(script)
}

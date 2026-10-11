// A browser capability, scoped to one backend. It never represents a user login.
const memory = new Map<string, string>()
const valid = (token: unknown): token is string =>
  typeof token === 'string' && /^[a-f0-9]{64}$/.test(token)

function storageKey(server: string) {
  const url = new URL(server, window.location.href)
  return `ArtalkGuest:${url.origin}${url.pathname.replace(/\/+$/, '')}`
}

export function getGuestToken(server: string, create = false): string | undefined {
  const key = storageKey(server)
  // Keep the credential stable for this document, including when a trusted
  // sidebar handoff cannot replace an older value in read-only storage.
  if (memory.has(key)) return memory.get(key)
  try {
    const stored = window.localStorage.getItem(key)
    if (valid(stored)) {
      memory.set(key, stored)
      return stored
    }
  } catch {
    // Storage may be blocked in an iframe. Keep the current page usable.
  }
  if (!create || !window.crypto?.getRandomValues) return undefined
  const bytes = window.crypto.getRandomValues(new Uint8Array(32))
  const token = Array.from(bytes, (byte) => byte.toString(16).padStart(2, '0')).join('')
  setGuestToken(server, token)
  return token
}

export function setGuestToken(server: string, token: string) {
  if (!valid(token)) return
  const key = storageKey(server)
  memory.set(key, token)
  try {
    window.localStorage.setItem(key, token)
  } catch {
    // No persistence is possible, but the credential works for this page session.
  }
}

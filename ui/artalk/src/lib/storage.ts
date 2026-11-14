/** Ignore URL aliases that refer to the same backend. Keep subpath deployments separate. */
export function storageKey(name: string, server: string) {
  const url = new URL(server, window.location.href)
  return `${name}:${url.origin}${url.pathname.replace(/\/+$/, '')}`
}

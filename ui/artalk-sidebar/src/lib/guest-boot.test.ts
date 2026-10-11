// @vitest-environment jsdom
import { beforeEach, expect, it, vi } from 'vitest'

vi.mock('artalk', () => ({ default: class {} }))
beforeEach(() => {
  vi.resetModules()
  localStorage.clear()
})

it('consumes the sidebar fragment before routing and removes credentials from the address', async () => {
  const secret = 'ab'.repeat(32)
  history.replaceState(
    {},
    '',
    `/sidebar/?darkMode=1&user=${encodeURIComponent(JSON.stringify({ name: 'Guest', email: 'guest@example.com' }))}#guest_token=${secret}`,
  )
  const { bootParams } = await import('../global')
  expect(bootParams.guestToken).toBe(secret)
  expect(bootParams.user.token).toBe('')
  expect(bootParams.user.name).toBe('Guest')
  expect(location.hash).toBe('')
  expect(location.search).toBe('')
})

it('preserves direct sidebar routes and old clients without a guest credential', async () => {
  history.replaceState({}, '', '/sidebar/?darkMode=1#/pages')
  const { bootParams } = await import('../global')
  expect(bootParams.guestToken).toBe('')
  expect(location.hash).toBe('#/pages')
})

it('accepts the earlier sidebar fragment during a frontend upgrade', async () => {
  const secret = 'cd'.repeat(32)
  history.replaceState({}, '', `/sidebar/?darkMode=1#atk_guest=${secret}`)
  const { bootParams } = await import('../global')
  expect(bootParams.guestToken).toBe(secret)
  expect(location.hash).toBe('')
})

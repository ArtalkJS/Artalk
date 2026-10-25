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

it('restores login and guest credentials from the fragment before the router starts', async () => {
  const user = { name: 'Admin', email: 'admin@example.com', is_admin: true }
  const token = 'header.payload.signature'
  const guest = 'ab'.repeat(32)
  history.replaceState(
    {},
    '',
    `/sidebar/?darkMode=1&site=Test&view=comments&pageKey=%2Fpost&user=${encodeURIComponent(JSON.stringify(user))}#token=${token}&guest_token=${guest}`,
  )
  const { bootParams } = await import('../global')
  expect(bootParams.user).toMatchObject({ ...user, token })
  expect(bootParams.guestToken).toBe(guest)
  expect(bootParams).toMatchObject({
    site: 'Test',
    view: 'comments',
    pageKey: '/post',
    darkMode: true,
  })
  expect(location.search).toBe('')
  expect(location.hash).toBe('')
})

it.each([
  ['', 'legacy-token'],
  ['#token=new-token', 'new-token'],
  ['#token=', ''],
])(
  'accepts legacy query credentials but prefers an explicitly supplied fragment: %s',
  async (fragment, expected) => {
    history.replaceState(
      {},
      '',
      `/sidebar/?darkMode=1&user=${encodeURIComponent(JSON.stringify({ email: 'old@example.com', token: 'legacy-token' }))}${fragment}`,
    )
    const { bootParams } = await import('../global')
    expect(bootParams.user.token).toBe(expected)
    expect(location.search).toBe('')
    expect(location.hash).toBe('')
  },
)

it.each(['{', 'null', '{"email":123,"token":{},"is_admin":"false"}'])(
  'safely consumes malformed legacy user data: %s',
  async (user) => {
    history.replaceState({}, '', `/sidebar/?darkMode=1&user=${encodeURIComponent(user)}`)
    const { bootParams } = await import('../global')
    expect(bootParams.user).toEqual({ name: '', email: '', link: '', token: '', is_admin: false })
    expect(location.search).toBe('')
  },
)

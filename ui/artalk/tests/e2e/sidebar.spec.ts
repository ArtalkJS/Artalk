import { expect, test } from '@playwright/test'

const server = 'http://localhost:23366'
const admin = { name: 'Admin', email: 'admin@example.com', password: 'e2e-test-only' }

test('guest pending comments and notifications survive reload and the sidebar handoff', async ({
  page,
  request,
  browser,
}, testInfo) => {
  const errors: string[] = []
  page.on('pageerror', (error) => errors.push(error.message))
  await page.goto('/')
  await page.getByPlaceholder('Nickname').fill('Guest')
  await page.getByPlaceholder('Email').fill('guest@example.com')
  const content = 'My pending browser comment'
  await page.getByPlaceholder('Leave a comment').fill(content)
  const created = page.waitForResponse(
    (res) => res.url() === `${server}/api/v2/comments` && res.request().method() === 'POST',
  )
  await page.getByRole('button', { name: 'Send', exact: true }).click()
  const response = await created
  expect(response.ok()).toBe(true)
  const comment = await response.json()
  const guest = response.request().headers()['x-artalk-guest']
  expect(guest).toBeTruthy()

  // Set up moderation and a real reply through the backend API.
  const login = await request.post(`${server}/api/v2/user/access_token`, { data: admin })
  expect(login.ok()).toBe(true)
  const headers = { Authorization: `Bearer ${(await login.json()).token}` }
  const reply = await request.post(`${server}/api/v2/comments`, {
    headers,
    data: {
      ...admin,
      page_key: comment.page_key,
      site_name: 'ArtalkDocs',
      rid: comment.id,
      content: 'A reply for the guest',
    },
  })
  expect(reply.ok(), await reply.text()).toBe(true)
  const pending = await request.put(`${server}/api/v2/comments/${comment.id}`, {
    headers,
    data: { ...comment, is_pending: true },
  })
  expect(pending.ok(), await pending.text()).toBe(true)
  const forbidden = await request.delete(`${server}/api/v2/comments/${comment.id}`, {
    headers: { 'X-Artalk-Guest': guest },
  })
  expect(forbidden.status()).toBe(403)

  await page.reload()
  await expect(page.locator(`#atk-comment-${comment.id}`)).toContainText(content)
  await expect(page.locator('.atk-unread-badge')).toBeVisible()
  const handoff = page.waitForResponse(
    (res) =>
      res.url().includes('/api/v2/comments?') &&
      new URL(res.url()).searchParams.get('scope') === 'user',
  )
  await page.locator('[data-action="open-sidebar"]').click()
  const sidebarResponse = await handoff
  expect(sidebarResponse.ok()).toBe(true)
  expect(sidebarResponse.request().headers()['x-artalk-guest']).toBe(guest)
  const sidebar = page.frameLocator('.atk-sidebar-layer iframe')
  await expect(sidebar.locator('.atk-list')).toContainText('A reply for the guest')
  await expect(sidebar.locator('[atk-only-admin-show]:visible')).toHaveCount(0)
  await sidebar.locator('.top-navigation').getByText('Pending', { exact: true }).click()
  await expect(sidebar.locator('.atk-list')).toContainText(content)
  const frame = page.frames().find((frame) => frame.url().includes('/sidebar/'))!
  expect(new URL(frame.url()).search).toBe('')
  expect(frame.url()).not.toContain(guest)
  await frame.evaluate(() => location.reload())
  await expect(sidebar.locator('.atk-list')).toContainText(content)
  await expect(sidebar.locator('[atk-only-admin-show]:visible')).toHaveCount(0)
  await page.screenshot({
    path: testInfo.outputPath('guest-sidebar.png'),
    fullPage: true,
    animations: 'disabled',
  })
  await page.locator('.atk-sidebar-close').click()
  await expect(page.locator('body')).not.toHaveCSS('overflow', 'hidden')

  // Knowing the same nickname and email does not grant access to the pending comment.
  const stranger = await browser.newContext()
  try {
    const other = await stranger.newPage()
    await other.goto('http://localhost:5173')
    await other.getByPlaceholder('Nickname').fill('Guest')
    await other.getByPlaceholder('Email').fill('guest@example.com')
    const loaded = other.waitForResponse(
      (res) =>
        res.url().includes('/api/v2/comments?') &&
        new URL(res.url()).searchParams.get('page_key') === comment.page_key,
    )
    await other.evaluate((pageKey) => {
      const artalk = (window as any).artalk
      artalk.update({ pageKey })
      artalk.reload()
    }, comment.page_key)
    expect((await (await loaded).json()).comments).not.toEqual(
      expect.arrayContaining([expect.objectContaining({ id: comment.id })]),
    )
    await expect(other.locator('.atk-list')).not.toContainText(content)
  } finally {
    await stranger.close()
  }
  expect(errors).toEqual([])
})

test('the built sidebar supports login retry, moderation, refresh, and logout', async ({
  page,
  request,
}, testInfo) => {
  const errors: string[] = []
  page.on('pageerror', (error) => errors.push(error.message))
  const created = await request.post(`${server}/api/v2/comments`, {
    data: {
      name: 'Reader',
      email: 'reader@example.com',
      site_name: 'ArtalkDocs',
      page_key: `/admin-${Date.now()}`,
      content: 'Browser moderation test',
    },
  })
  expect(created.ok()).toBe(true)
  const { id } = await created.json()
  await page.goto(`${server}/sidebar/#/login`)
  await page.locator('.login-form').getByPlaceholder('Email').fill(admin.email)
  await page.getByPlaceholder('Password').fill('wrong-password')
  await page.getByRole('button', { name: 'Login', exact: true }).click()
  await expect(page.locator('.err-msg')).toBeVisible()
  await page.getByPlaceholder('Password').fill(admin.password)
  await page.getByRole('button', { name: 'Login', exact: true }).click()
  const comment = page.locator(`#atk-comment-${id}`)
  await expect(comment).toContainText('Browser moderation test')
  await comment.getByText('Edit', { exact: true }).click()
  await page.getByPlaceholder('Leave a comment').fill('Edited in the built sidebar')
  await page.getByRole('button', { name: 'Save', exact: true }).click()
  await expect(comment).toContainText('Edited in the built sidebar')
  await page.reload()
  await expect(comment).toContainText('Edited in the built sidebar')
  await page.screenshot({
    path: testInfo.outputPath('admin-sidebar.png'),
    fullPage: true,
    animations: 'disabled',
  })

  let statusRequests = 0
  await page.route(`${server}/api/v2/pages/fetch/status`, (route) => {
    statusRequests++
    return route.fulfill({ json: { is_progress: true, msg: 'Fetching titles' } })
  })
  await page.locator('.sidebar-navigation').getByText('Page', { exact: true }).click()
  await expect(page.locator('.atk-update-all-title-btn')).toContainText('Fetching titles')
  await expect(page.locator('.atk-page-list')).toContainText(`/admin-`)
  const pageNumber = page.getByLabel('Enter the number of page')
  await pageNumber.fill('1.5')
  await pageNumber.press('Enter')
  await expect(pageNumber).toHaveValue('1')
  let releaseUsers!: () => void
  let usersReady!: () => void
  let usersFinished!: () => void
  const usersGate = new Promise<void>((resolve) => (releaseUsers = resolve))
  const ready = new Promise<void>((resolve) => (usersReady = resolve))
  const finished = new Promise<void>((resolve) => (usersFinished = resolve))
  await page.route(`${server}/api/v2/users/all?*`, async (route) => {
    const response = await route.fetch()
    usersReady()
    await usersGate
    await route.fulfill({ response })
    usersFinished()
  })
  await page.locator('.sidebar-navigation').getByText('User', { exact: true }).click()
  await ready
  await page.locator('.top-tabbar').getByText('Admin', { exact: true }).click()
  await expect(page.locator('.user-list')).toContainText(admin.email)
  await expect(page.locator('.user-list')).not.toContainText('reader@example.com')
  releaseUsers()
  await finished
  await expect(page.locator('.user-list')).not.toContainText('reader@example.com')
  const previous = statusRequests
  // Two polling intervals must pass without requests from the route we left.
  await page.waitForTimeout(2200)
  expect(statusRequests).toBe(previous)
  await page.locator('.sidebar-navigation').getByText('Site', { exact: true }).click()
  await expect(page.locator('.main')).toContainText('ArtalkDocs')
  await page.locator('.sidebar-navigation').getByText('Comment', { exact: true }).click()
  await expect(comment).toContainText('Edited in the built sidebar')

  await comment.getByText('Delete', { exact: true }).click()
  await comment.getByText('Confirm', { exact: true }).click()
  await expect(comment).toHaveCount(0)
  await page.locator('.header .avatar').click()
  page.once('dialog', (dialog) => dialog.accept())
  await page.getByText('Logout', { exact: true }).click()
  await expect(page.getByPlaceholder('Password')).toBeVisible()
  await page.reload()
  await expect(page.getByPlaceholder('Password')).toBeVisible()
  expect(errors).toEqual([])
})

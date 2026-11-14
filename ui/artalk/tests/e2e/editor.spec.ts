import { expect, test } from '@playwright/test'

const png =
  'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+/l9sAAAAASUVORK5CYII='

test('mobile editor keeps drafts, renders emoticons, uploads, and retries a captcha', async ({
  page,
}, testInfo) => {
  const errors: string[] = []
  page.on('pageerror', (error) => errors.push(error.message))
  await page.setViewportSize({ width: 390, height: 844 })
  await page.goto('/')
  await expect(page.getByRole('button', { name: 'Send', exact: true })).toBeVisible()
  await page.evaluate(() =>
    (window as any).artalk.update({
      darkMode: true,
      emoticons: [{ name: 'Emoji', type: 'emoji', items: [{ key: 'smile', val: '🙂' }] }],
    }),
  )
  await page.getByPlaceholder('Nickname').fill('Mobile visitor')
  await page.getByPlaceholder('Email').fill('mobile@example.com')
  const input = page.getByPlaceholder('Leave a comment')
  await input.fill('**Mobile draft** ')
  await page.locator('[aria-label="Emoji"]').click()
  await page.locator('.atk-grp .atk-item').click()
  await expect(input).toHaveValue('**Mobile draft** 🙂')
  await page.locator('[aria-label="Preview"]').click()
  await expect(page.locator('.atk-editor-plug-preview strong')).toHaveText('Mobile draft')
  await expect(page.locator('#comments')).toHaveClass(/atk-dark-mode/)
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true)
  await page.screenshot({
    path: testInfo.outputPath('mobile-dark-preview.png'),
    fullPage: true,
    animations: 'disabled',
  })

  const chooseFile = page.waitForEvent('filechooser')
  await page.locator('[aria-label="Upload Image"]').click()
  await (
    await chooseFile
  ).setFiles({ name: 'pixel.png', mimeType: 'image/png', buffer: Buffer.from(png, 'base64') })
  await expect(input).toHaveValue(/!\[\]\(.*\/static\/images\/.*\.png\)/)
  const draft = await input.inputValue()
  await page.evaluate(() => (window as any).artalk.update({ pageKey: '/other-page' }))
  await expect(input).toHaveValue('')
  await page.evaluate(() =>
    (window as any).artalk.update({ pageKey: sessionStorage.getItem('test-page') }),
  )
  await expect(input).toHaveValue(draft)
  await page.reload()
  await expect(input).toHaveValue(draft)

  // Inject only the challenge responses; the final comment uses the real backend.
  let writes = 0
  let challenge = true
  await page.route('**/api/v2/comments', async (route) => {
    if (route.request().method() !== 'POST') return route.continue()
    writes++
    if (!challenge) return route.continue()
    await route.fulfill({
      status: 403,
      json: {
        msg: 'Captcha required',
        need_captcha: true,
        img_data: `data:image/png;base64,${png}`,
      },
    })
  })
  await page.getByRole('button', { name: 'Send', exact: true }).click()
  await expect(page.locator('.atk-captcha-img')).toBeVisible()
  await page.locator('[data-action="cancel"]').click()
  await expect(input).toHaveValue(draft)
  expect(writes).toBe(1)

  await page.route('**/api/v2/captcha', (route) =>
    route.fulfill({ json: { img_data: `data:image/png;base64,${png}` } }),
  )
  await page.route('**/api/v2/captcha/verify', async (route) => {
    if (route.request().postDataJSON().value !== 'test-answer') {
      return route.fulfill({ status: 400, json: { msg: 'Incorrect captcha' } })
    }
    challenge = false
    await route.fulfill({ json: {} })
  })
  await page.getByRole('button', { name: 'Send', exact: true }).click()
  await page.locator('#check').fill('wrong')
  await page.locator('[data-action="confirm"]').click()
  await expect(page.locator('[data-action="confirm"]')).toHaveText('Incorrect captcha')
  await page.locator('#check').fill('test-answer')
  const sent = page.waitForResponse(
    (res) =>
      res.url().endsWith('/api/v2/comments') && res.request().method() === 'POST' && res.ok(),
  )
  await page.locator('[data-action="confirm"]').click()
  const { id } = await (await sent).json()
  await expect(page.locator(`#atk-comment-${id}`)).toContainText('Mobile draft')
  await expect(input).toHaveValue('')
  expect(writes).toBe(3)
  expect(errors).toEqual([])
})

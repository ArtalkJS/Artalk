import { expect, test } from '@playwright/test'

test('widgets share a slow download of the built Auth plugin', async ({ page }) => {
  const errors: string[] = []
  page.on('pageerror', (error) => errors.push(error.message))
  const pluginURL = 'http://localhost:23366/dist/plugins/artalk-plugin-auth.js'
  let release!: () => void
  const gate = new Promise<void>((resolve) => (release = resolve))
  let downloads = 0
  await page.route(pluginURL, async (route) => {
    downloads++
    const response = await route.fetch()
    await gate
    await route.fulfill({ response })
  })
  await page.goto('/')
  await expect(page.getByRole('button', { name: 'Send', exact: true })).toBeVisible()
  await page.evaluate((pluginURL) => {
    const conf = (window as any).artalk.getConf()
    ;(window as any).artalk.destroy()
    ;(window as any).artalk = (window as any).Artalk.init({
      ...conf,
      pluginURLs: [pluginURL],
    })
  }, pluginURL)
  await expect.poll(() => downloads).toBe(1)
  await page.evaluate((pluginURL) => {
    const root = document.createElement('div')
    root.id = 'second-comments'
    document.body.append(root)
    ;(window as any).secondArtalk = (window as any).Artalk.init({
      el: root,
      server: 'http://localhost:23366',
      site: 'ArtalkDocs',
      pageKey: '/plugin-smoke',
      locale: 'en',
      emoticons: false,
      pluginURLs: [pluginURL],
    })
  }, pluginURL)
  await expect
    .poll(() => page.evaluate(() => (window as any).secondArtalk.getConf().apiVersion))
    .toBeTruthy()
  release()
  for (const root of ['#comments', '#second-comments']) {
    await expect(page.locator(root).getByRole('button', { name: 'Sign in' })).toBeVisible()
    await expect(page.locator(root).getByPlaceholder('Nickname')).toBeHidden()
  }
  expect(downloads).toBe(1)
  expect(errors).toEqual([])
})

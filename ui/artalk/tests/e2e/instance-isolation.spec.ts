import { expect, test } from '@playwright/test'

test('two real widgets share a language and keep Markdown, users, and cleanup separate', async ({
  page,
}) => {
  const errors: string[] = []
  page.on('pageerror', (error) => errors.push(error.message))
  await page.goto('/')
  await expect(page.locator('#comments .atk-send-btn')).toBeVisible()

  await page.evaluate(() => {
    const first = (window as any).artalk
    first.update({
      locale: 'en',
      markedOptions: { breaks: true },
      markedReplacers: [(text: string) => text.replaceAll('sample', 'first')],
    })
    const root = document.createElement('div')
    root.id = 'second-comments'
    document.body.append(root)
    ;(window as any).secondArtalk = (window as any).Artalk.init({
      el: root,
      server: first.getConf().server,
      site: first.getConf().site,
      pageKey: '/instance-isolation',
      locale: 'en',
      emoticons: false,
      markedOptions: { breaks: false },
    })
  })
  const first = page.locator('#comments')
  const second = page.locator('#second-comments')
  await expect(first.getByRole('button', { name: 'Send', exact: true })).toBeVisible()
  await expect(second.getByRole('button', { name: 'Send', exact: true })).toBeVisible()
  await first.locator('input[name="name"]').fill('First Visitor')
  await second.getByPlaceholder('Nickname').fill('Second Visitor')
  await expect(first.locator('input[name="name"]')).toHaveValue('First Visitor')

  await first.locator('textarea').fill('sample\nnext **bold**')
  await second.locator('textarea').fill('sample\nnext **bold**')
  await first.locator('[aria-label="Preview"]').click()
  await second.locator('[aria-label="Preview"]').click()
  await expect(first.locator('.atk-editor-plug-preview')).toContainText('first')
  await expect(first.locator('.atk-editor-plug-preview br')).toHaveCount(1)
  await expect(second.locator('.atk-editor-plug-preview')).toContainText('sample')
  await expect(second.locator('.atk-editor-plug-preview br')).toHaveCount(0)
  await expect(second.locator('.atk-editor-plug-preview strong')).toHaveText('bold')

  await page.evaluate(() => {
    ;(window as any).secondArtalk.update({ markedOptions: { breaks: true } })
  })
  await expect(second.locator('.atk-editor-plug-preview br')).toHaveCount(1)
  await expect(first.getByRole('button', { name: 'Send', exact: true })).toBeVisible()

  await page.evaluate(() => {
    for (const widget of [(window as any).artalk, (window as any).secondArtalk]) {
      widget.update({ locale: 'zh-CN' })
    }
  })
  await expect(first.getByRole('button', { name: '发送', exact: true })).toBeVisible()
  await expect(second.getByRole('button', { name: '发送', exact: true })).toBeVisible()
  await expect(first.locator('textarea')).toHaveValue('sample\nnext **bold**')
  await expect(second.locator('textarea')).toHaveValue('sample\nnext **bold**')

  await page.evaluate(() => {
    const first = (window as any).artalk
    const second = (window as any).secondArtalk
    first.ctx.inject('layers').create('first-test-layer').show()
    const layer = second.ctx.inject('layers').create('second-test-layer')
    layer.show()
    ;(window as any).secondTestLayer = layer
    first.destroy()
  })
  await expect(page.locator('body')).toHaveCSS('overflow', 'hidden')
  await expect(page.locator('.atk-layer-wrap')).toHaveCount(1)
  await page.evaluate(() => (window as any).secondTestLayer.hide())
  await expect(page.locator('body')).not.toHaveCSS('overflow', 'hidden')
  await expect(second.getByRole('button', { name: '发送', exact: true })).toBeVisible()
  await expect(second.locator('input[name="name"]')).toHaveValue('Second Visitor')
  await page.evaluate(() => (window as any).secondArtalk.destroy())
  await expect(page.locator('.atk-layer-wrap')).toHaveCount(0)
  expect(errors).toEqual([])
})

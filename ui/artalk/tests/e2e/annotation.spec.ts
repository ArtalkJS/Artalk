import { expect, test } from '@playwright/test'

test.beforeEach(async ({ page }) => {
  await page.goto('/')
  await page.getByPlaceholder('Nickname').fill('qwqcode')
  await page.getByPlaceholder('Email').fill('qwqcode@gmail.com')
  await expect(page.locator('.atk-annotation-fab')).toBeVisible()
})

test('creates, reloads, resolves, and deletes a highlight annotation', async ({ page }) => {
  await page.locator('#article p').evaluate((el) => {
    const text = el.firstChild as Text
    const range = document.createRange()
    range.setStart(text, 0)
    range.setEnd(text, 8)
    const selection = window.getSelection()!
    selection.removeAllRanges()
    selection.addRange(range)
    document.dispatchEvent(new Event('selectionchange'))
    document.dispatchEvent(new MouseEvent('mouseup', { bubbles: true }))
  })

  const annotate = page.locator('.atk-annotation-toolbar')
  await expect(annotate).toBeVisible()
  await annotate.click()

  const composer = page.locator('.atk-annotation-composer textarea')
  await composer.fill('E2E annotation')
  await page.locator('.atk-annotation-composer .submit').click()

  await expect(page.locator('mark.atk-annotation-mark')).toHaveCount(1)
  await expect(page.locator('.atk-annotation-drawer')).toBeVisible()
  await expect(page.locator('.atk-annotation-drawer')).toContainText('E2E annotation')

  await page.reload()
  await expect(page.locator('mark.atk-annotation-mark')).toHaveCount(1)

  await page.locator('.atk-annotation-fab').click()
  await page.locator('.atk-annotation-drawer [data-action="resolve"]').first().click()
  await expect(
    page.locator('.atk-annotation-drawer .badge').filter({ hasText: 'Resolved' }),
  ).toBeVisible()

  await page.locator('.atk-annotation-drawer [data-action="delete"]').first().click()
  await page.locator('.atk-annotation-confirm .confirm').click()
  await expect(page.locator('mark.atk-annotation-mark')).toHaveCount(0)
})

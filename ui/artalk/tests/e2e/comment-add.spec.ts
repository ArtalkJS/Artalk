import { expect, test } from '@playwright/test'

test.beforeEach(async ({ page }) => {
  await page.goto('/')

  await page.getByPlaceholder('Nickname').fill('qwqcode')
  await page.getByPlaceholder('Email').fill('qwqcode@gmail.com')
  await page.getByPlaceholder('Website').fill('https://github.com/qwqcode')
})

test('Comment Add', async ({ page }) => {
  const commentContent = '[A New Test Comment Content Here]'
  await page.getByPlaceholder('Leave a comment').fill(commentContent)

  await Promise.all([
    page.waitForResponse(
      (response) =>
        response.url() === 'http://localhost:23366/api/v2/comments' &&
        response.request().method() === 'POST',
    ),
    page.getByRole('button', { name: 'Send' }).click(),
  ])

  await expect(page.locator('.atk-list').getByText(commentContent)).toBeVisible()
})

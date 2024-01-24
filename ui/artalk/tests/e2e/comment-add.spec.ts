import { expect, test } from '@playwright/test'

test.beforeEach(async ({ page }) => {
  await page.goto('http://localhost:5173/')

  await page.getByPlaceholder('昵称').click()
  await page.getByPlaceholder('昵称').fill('qwqcode')
  await page.getByPlaceholder('昵称').press('Tab')
  await page.getByPlaceholder('邮箱').fill('qwqcode@gmail.com')
  await page.getByPlaceholder('邮箱').press('Tab')
  await page.getByPlaceholder('网址').fill('https://github.com/qwqcode')
})

test('Comment Add', async ({ page }) => {
  await page.getByPlaceholder('键入内容...').click()
  const CommentContent = '[A New Test Comment Content Here]'
  await page.getByPlaceholder('键入内容...').fill(CommentContent)
  await page.getByRole('button', { name: '发送评论' }).click()
  await page.waitForResponse('http://localhost:23366/api/v2/comments')
  expect(
    await page.locator('.atk-list').getByText(CommentContent).isVisible(),
  ).toBeTruthy()
})

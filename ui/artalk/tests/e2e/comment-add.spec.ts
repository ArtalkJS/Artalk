import { expect, test } from '@playwright/test'

test.beforeEach(async ({ page }) => {
  await page.goto('/')

  await page.getByPlaceholder('Nickname').fill('qwqcode')
  await page.getByPlaceholder('Email').fill('qwqcode@gmail.com')
  await page.getByPlaceholder('Website').fill('https://github.com/qwqcode')
})

test('Comment add, reply, sort, and reload', async ({ page }) => {
  const commentContent = '[A New Test Comment Content Here]'
  await page.getByPlaceholder('Leave a comment').fill(commentContent)

  const [created] = await Promise.all([
    page.waitForResponse(
      (response) =>
        response.url() === 'http://localhost:23366/api/v2/comments' &&
        response.request().method() === 'POST',
    ),
    page.getByRole('button', { name: 'Send' }).click(),
  ])

  expect(created.ok()).toBe(true)
  const { id } = await created.json()
  await expect(page.locator('.atk-list').getByText(commentContent)).toBeVisible()

  const replyContent = '[A Reply From The Same Guest]'
  await page.locator(`#atk-comment-${id} .atk-actions`).getByText('Reply', { exact: true }).click()
  await page.getByPlaceholder('Leave a comment').fill(replyContent)
  const [replied] = await Promise.all([
    page.waitForResponse(
      (response) =>
        response.url() === 'http://localhost:23366/api/v2/comments' &&
        response.request().method() === 'POST',
    ),
    page.getByRole('button', { name: 'Send' }).click(),
  ])
  expect(replied.ok()).toBe(true)
  expect(await replied.json()).toMatchObject({ rid: id })
  await expect(page.locator('.atk-list').getByText(replyContent)).toBeVisible()

  await page.reload()
  await expect(page.locator('.atk-list').getByText(commentContent)).toBeVisible()
  await expect(page.locator('.atk-list').getByText(replyContent)).toBeVisible()
  await expect(page.getByPlaceholder('Nickname')).toHaveValue('qwqcode')

  // Select oldest first through the public widget, then check the v2 query.
  await page.locator('.atk-dropdown-wrap').hover()
  const [sorted] = await Promise.all([
    page.waitForResponse(
      (response) => new URL(response.url()).searchParams.get('sort_by') === 'date_asc',
    ),
    page.locator('.atk-dropdown-item').nth(2).click(),
  ])
  expect(sorted.ok()).toBe(true)
  await expect(page.locator('.atk-list').getByText(replyContent)).toBeVisible()

  await page.evaluate(() => {
    const artalk = (window as any).artalk
    artalk.update({ flatMode: true, pagination: { pageSize: 1, readMore: false } })
    artalk.reload()
  })
  const pageNumber = page.getByLabel('Enter the number of page')
  await expect(pageNumber).toHaveValue('1')
  await page.getByLabel('Next page').click()
  await expect(pageNumber).toHaveValue('2')
  await expect(page.locator('.atk-list .atk-comment')).toHaveCount(1)
  await expect(page.locator('.atk-list')).toContainText(replyContent)
  for (const value of ['1.5', 'Infinity', '-1', 'text', '']) {
    await pageNumber.fill(value)
    await pageNumber.press('Enter')
    await expect(pageNumber).toHaveValue('2')
    await expect(page.locator('.atk-list')).toContainText(replyContent)
  }

  await page.evaluate(() => {
    const artalk = (window as any).artalk
    artalk.update({ pageKey: `/empty-${crypto.randomUUID()}` })
    artalk.reload()
  })
  await expect(page.locator('.atk-list .atk-comment')).toHaveCount(0)
  await expect(pageNumber).toHaveValue('1')
  await pageNumber.fill('99')
  await pageNumber.press('Enter')
  await expect(pageNumber).toHaveValue('1')
})

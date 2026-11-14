import { defineConfig, devices } from '@playwright/test'
import { copyFileSync } from 'node:fs'

// Build with `make build-frontend && make build-debug` before running test:e2e.
copyFileSync(
  new URL('./tests/e2e/index.html', import.meta.url),
  new URL('./dist/index.html', import.meta.url),
)

export default defineConfig({
  // Keep Playwright isolated from the Vitest suites in the parent tests directory.
  testDir: 'tests/e2e',

  // Run all tests in parallel.
  fullyParallel: true,

  // Fail the build on CI if you accidentally left test.only in the source code.
  forbidOnly: !!process.env.CI,

  // Retry on CI only.
  retries: process.env.CI ? 2 : 0,

  // Opt out of parallel tests on CI.
  workers: process.env.CI ? 1 : undefined,

  // Reporter to use
  reporter: 'html',

  use: {
    // Base URL to use in actions like `await page.goto('/')`.
    baseURL: 'http://localhost:5173',

    // Collect trace when retrying the failed test.
    trace: 'on-first-retry',
  },
  // Configure projects for major browsers.
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
  ],
  // Serve the built package and an isolated backend with the built sidebar.
  webServer: [
    {
      command: 'pnpm exec vite preview --port 5173 --strictPort',
      url: 'http://localhost:5173',
      reuseExistingServer: false,
    },
    {
      command: 'pnpm exec tsx tests/e2e/server.ts',
      url: 'http://localhost:23366',
      reuseExistingServer: false,
    },
  ],
})

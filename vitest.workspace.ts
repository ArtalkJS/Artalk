import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    coverage: {
      provider: 'v8',
      include: ['ui/artalk/src/**/*.ts'],
      exclude: ['**/*.d.ts'],
      reporter: ['text-summary', 'lcov'],
      reportsDirectory: './coverage/artalk',
    },
    projects: [
      './ui/artalk/vitest.config.ts',
      // "./ui/plugin-katex/vite.config.ts",
      './ui/plugin-auth/vitest.config.ts',
      // "./ui/plugin-lightbox/vite.config.ts",
      './ui/artalk-sidebar/vite.config.ts',
      // "./docs/landing/vite.config.ts",
    ],
  },
})

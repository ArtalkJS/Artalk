import { fileURLToPath } from 'node:url'
import { resolve } from 'path'
import { mergeConfig, defineConfig, configDefaults } from 'vitest/config'
import viteConfig from './vite.config'

export default mergeConfig(
  viteConfig,
  defineConfig({
    test: {
      environment: 'jsdom',
      exclude: [...configDefaults.exclude, 'tests/e2e/*'],
      root: fileURLToPath(new URL('./', import.meta.url)),
      setupFiles: ['tests/setup.ts'],
    },
  }),
)

import { defineConfig } from 'vite'
import tsconfigPaths from 'vite-tsconfig-paths'
import { resolve } from 'path'
import dts from 'vite-plugin-dts'
import checker from 'vite-plugin-checker'

export default defineConfig({
  root: __dirname,
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@import "./src/style/_variables.scss";@import "./src/style/_extend.scss";`
     },
    },
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '~': resolve(__dirname),
    }
  },
  plugins: [
    tsconfigPaths(),
    checker({
      typescript: true,
      eslint: {
        lintCommand: 'eslint "./src/**/*.{js,ts}"',
      },
    }),
    dts({
      // @see https://github.com/qmhc/vite-plugin-dts/blob/main/CHANGELOG.md#breaking-changes
      copyDtsFiles: true
    })
  ],
})

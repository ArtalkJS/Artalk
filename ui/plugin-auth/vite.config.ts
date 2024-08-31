import { defineConfig } from 'vite'
import { resolve, dirname } from 'node:path'
import solidPlugin from 'vite-plugin-solid'
import cssInjectedByJsPlugin from 'vite-plugin-css-injected-by-js'
import { fileURLToPath } from 'node:url'
import { ViteArtalkPluginKit } from '@artalk/plugin-kit'

const __dirname = dirname(fileURLToPath(import.meta.url))

export default defineConfig({
  build: {
    lib: {
      entry: resolve(__dirname, './main.tsx'),
    },
  },
  plugins: [ViteArtalkPluginKit(), solidPlugin(), cssInjectedByJsPlugin()],
})

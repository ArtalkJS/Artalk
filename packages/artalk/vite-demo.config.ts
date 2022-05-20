import { defineConfig } from 'vite'
import { resolve } from 'path'
import baseConf from './vite-base.config'

export default defineConfig({
  build: {
    target: 'esnext',
    outDir: resolve(__dirname, "deploy"),
    rollupOptions: {}
  },
  ...baseConf
})

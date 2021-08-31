import { defineConfig } from 'vite'
import { resolve } from 'path'
import commonConf from './vite-common.config'

export default defineConfig({
  build: {
    target: 'esnext',
    outDir: resolve(__dirname, "deploy"),
    rollupOptions: {}
  },
  ...commonConf
})

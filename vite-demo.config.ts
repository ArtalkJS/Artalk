import { defineConfig } from 'vite'
import { resolve } from 'path'

export default defineConfig({
  build: {
    target: 'esnext',
    outDir: resolve(__dirname, "deploy"),
    rollupOptions: {}
  }
})

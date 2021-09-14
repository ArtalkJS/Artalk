import { defineConfig } from 'vite'
import { resolve } from 'path'
import commonConf from './vite-common.config'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig({
  build: {
    target: 'esnext',
    outDir: resolve(__dirname, "deploy"),
    rollupOptions: {}
  },
  plugins: [tsconfigPaths()],
  ...commonConf
})

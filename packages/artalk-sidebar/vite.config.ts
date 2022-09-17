import { defineConfig } from 'vite'
import { resolve } from 'path'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig({
  base: './',
  build: {
    target: 'es2015',
    outDir: resolve(__dirname, "dist"),
    minify: 'terser'
  },
  plugins: [tsconfigPaths()],
  css: {
    preprocessorOptions: {
      less: {
         additionalData: `@import "./src/style/_variables.less";@import "./src/style/_extend.less";`
     },
    },
  },
})

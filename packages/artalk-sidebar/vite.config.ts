import { defineConfig } from 'vite'
import { resolve } from 'path'
import tsconfigPaths from 'vite-tsconfig-paths'
import { version as artalkVersion } from 'artalk/package.json'
import { version as sidebarVersion } from './package.json'

export default defineConfig({
  base: './',
  define: {
    ARTALK_VERSION: JSON.stringify(artalkVersion),
    ARTALK_SIDEBAR_VERSION: JSON.stringify(sidebarVersion)
  },
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

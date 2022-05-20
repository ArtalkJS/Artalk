import { defineConfig } from 'vite'
import tsconfigPaths from 'vite-tsconfig-paths'
import { version } from './package.json'

export default defineConfig({
  define: {
    ARTALK_VERSION: JSON.stringify(version),
  },
  css: {
    preprocessorOptions: {
      less: {
         additionalData: `@import "./src/style/_variables.less";@import "./src/style/_extend.less";`
     },
    },
  },
  plugins: [tsconfigPaths()],
})

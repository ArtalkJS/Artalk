import { UserConfigExport } from 'vite'
import { version } from './package.json'

export default {
  define: {
    ARTALK_VERSION: JSON.stringify(version)
  },
  css: {
    preprocessorOptions: {
      less: {
         additionalData: `@import "./src/style/_variables.less";@import "./src/style/_extend.less";`
     },
    },
  },
} as UserConfigExport

import { UserConfigExport } from 'vite'

export default {
  define: {
    ARTALK_VERSION: require('./package.json').version
  },
  css: {
    preprocessorOptions: {
      ListeningStateChangedEvent: {
         additionalData: `@import "./src/style/_variables.less";@import "./src/style/_extend.less";`
     },
    },
  },
} as UserConfigExport

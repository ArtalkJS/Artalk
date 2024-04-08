import DefaultTheme from 'vitepress/theme'
import Layout from './Layout.vue'
import Artalk from './Artalk.vue'
import Artransfer from './Artransfer.vue'
import { Theme } from 'vitepress'

export default {
  ...DefaultTheme,

  Layout,

  enhanceApp({ app, router, siteData }) {
    app.component('Artransfer', Artransfer)
    app.component('Artalk', Artalk)

    // app is the Vue 3 app instance from `createApp()`.
    // router is VitePress' custom router. `siteData` is
    // a `ref` of current site-level metadata.
  },
} as Theme

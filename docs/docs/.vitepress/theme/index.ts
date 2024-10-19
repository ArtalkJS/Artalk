import DefaultTheme from 'vitepress/theme'
import { Theme } from 'vitepress'
import { NolebaseGitChangelogPlugin } from '@nolebase/vitepress-plugin-git-changelog/client'
import Layout from './Layout.vue'
import Artalk from './Artalk.vue'
import Artransfer from './Artransfer.vue'

import '@nolebase/vitepress-plugin-git-changelog/client/style.css'

export default {
  ...DefaultTheme,

  Layout,

  enhanceApp({ app, router, siteData }) {
    app.component('Artransfer', Artransfer)
    app.component('Artalk', Artalk)
    app.use(NolebaseGitChangelogPlugin)

    // app is the Vue 3 app instance from `createApp()`.
    // router is VitePress' custom router. `siteData` is
    // a `ref` of current site-level metadata.
  },
} as Theme

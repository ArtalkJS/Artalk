import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Artalk from 'artalk'
import { createRouter, createWebHashHistory } from 'vue-router'
import { routes } from 'vue-router/auto-routes'
import messages from './i18n/messages'
import 'artalk/dist/Artalk.css'
import './style.scss'
import App from './App.vue'
import { setArtalk, initArtalk } from './global'

// i18n
const i18n = createI18n({
  legacy: false, // use i18n in Composition API
  locale: 'en',
  fallbackLocale: 'en',
  messages,
})

// Artalk extension
Artalk.use((ctx) => {
  // Sync config from artalk instance to sidebar
  ctx.watchConf(['locale'], (conf) => {
    if (typeof conf.locale === 'string' && conf.locale !== 'auto')
      i18n.global.locale.value = conf.locale as any
  })
})

const app = createApp(App)

app.use(i18n)

// Init Artalk
setArtalk(initArtalk())

// router
// @see https://github.com/posva/unplugin-vue-router
const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

app.use(router)

// pinia
const pinia = createPinia()
app.use(pinia)

app.mount('#app')

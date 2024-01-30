import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import { createRouter, createWebHashHistory } from 'vue-router/auto'
import messages from './i18n/messages'
import 'artalk/dist/Artalk.css'
import './style.scss'
import App from './App.vue'
import { setArtalk, initArtalk } from './global'
import Artalk from 'artalk'

// i18n
const i18n = createI18n({
  legacy: false, // use i18n in Composition API
  locale:  'en',
  fallbackLocale: 'en',
  messages
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
const router = createRouter({
  history: createWebHashHistory(),
})

app.use(router)

// pinia
const pinia = createPinia()
app.use(pinia)

app.mount('#app')

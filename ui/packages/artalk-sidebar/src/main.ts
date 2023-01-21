import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import { createRouter, createWebHashHistory } from 'vue-router/auto'
import messages from './i18n/messages'
import 'artalk/dist/Artalk.css'
import './style.scss'
import App from './App.vue'
import global, { createArtalkInstance, bootParams, getBootParams } from './global'

// 初始化 Artalk
global.setArtalk(createArtalkInstance())

const app = createApp(App)

// router
const router = createRouter({
  history: createWebHashHistory(),
})

app.use(router)

// i18n
const i18n = createI18n({
  legacy: false, // use i18n in Composition API
  locale: getBootParams().locale || 'en',
  fallbackLocale: 'en',
  messages
})

app.use(i18n)

// pinia
const pinia = createPinia()
app.use(pinia)

app.mount('#app')

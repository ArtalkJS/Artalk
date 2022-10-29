import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from 'vue-router/auto'
import 'artalk/dist/Artalk.css'
import './style.scss'
import App from './App.vue'
import global, { createArtalkInstance, bootParams } from './global'

// 初始化 Artalk
global.setArtalk(createArtalkInstance())

const app = createApp(App)

const router = createRouter({
  history: createWebHashHistory(),
})

app.use(router)

const pinia = createPinia()
app.use(pinia)

app.mount('#app')

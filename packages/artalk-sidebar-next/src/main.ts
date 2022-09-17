import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router/auto'
import 'artalk/dist/Artalk.css'
import Artalk from 'artalk'
import './style.css'
import App from './App.vue'
import global from './global'

// 启动参数
const p        = new URLSearchParams(document.location.search)
const pageKey  = p.get('pageKey') || ''
const site     = p.get('site') || ''
const user     = JSON.parse(p.get('user') || '{}')
const view     = p.get('view') || ''
const darkMode = p.get('darkMode') === '1'

// 初始化 Artalk
global.artalk = await (new Artalk({
  el: document.createElement('div'),
  server: (import.meta.env.DEV) ? 'http://localhost:23366' : '/',
  pageKey,
  site,
  darkMode,
  useBackendConf: true
}))

global.artalk.ctx.user.update({
  ...user
})

const app = createApp(App)

const router = createRouter({
  history: createWebHashHistory(),
})

app.use(router)
app.mount('#app')

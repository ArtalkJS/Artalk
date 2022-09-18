import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from 'vue-router/auto'
import 'artalk/dist/Artalk.css'
import Artalk from 'artalk'
import './style.css'
import App from './App.vue'
import global from './global'

// 启动参数
function getBootParameters() {
  const p = new URLSearchParams(document.location.search)
  return {
    pageKey:  p.get('pageKey') || '',
    site:     p.get('site') || '',
    user:     JSON.parse(p.get('user') || '{}'),
    view:     p.get('view') || '',
    darkMode: p.get('darkMode') === '1'
  }
}

function createArtalkInstance() {
  const p = getBootParameters()
  const artalkEl = document.createElement('div')
  artalkEl.style.display = 'none'
  document.body.append(artalkEl)

  Artalk.DisabledComponents = ['list']
  return new Artalk({
    el: artalkEl,
    server: (import.meta.env.DEV) ? 'http://localhost:23366' : '/',
    pageKey: 'https://artalk.js.org/guide/intro.html',
    site: p.site,
    darkMode: p.darkMode,
    useBackendConf: true
  }) as unknown as Promise<Artalk>
}

// 初始化 Artalk
global.artalk = await createArtalkInstance()
// note: 这里 await 会阻塞整个 vue app

const app = createApp(App)

const router = createRouter({
  history: createWebHashHistory(),
})

app.use(router)

const pinia = createPinia()
app.use(pinia)

app.mount('#app')

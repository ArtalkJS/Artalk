import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from 'vue-router/auto'
import 'artalk/dist/Artalk.css'
import Artalk from 'artalk'
import './style.scss'
import App from './App.vue'
import global, { bootParams } from './global'

function createArtalkInstance() {
  const artalkEl = document.createElement('div')
  artalkEl.style.display = 'none'
  document.body.append(artalkEl)

  Artalk.DisabledComponents = ['list']
  return new Artalk({
    el: artalkEl,
    server: (import.meta.env.DEV) ? 'http://localhost:23366' : '/',
    pageKey: 'https://artalk.js.org/guide/intro.html',
    site: bootParams.site,
    darkMode: bootParams.darkMode,
    useBackendConf: true
  }) as unknown as Promise<Artalk>
}

// 初始化 Artalk
global.setArtalk(await createArtalkInstance())
// note: 这里 await 会阻塞整个 vue app

// 更新用户资料
global.getArtalk()!.ctx.user.update(bootParams.user)

const app = createApp(App)

const router = createRouter({
  history: createWebHashHistory(),
})

app.use(router)

const pinia = createPinia()
app.use(pinia)

app.mount('#app')

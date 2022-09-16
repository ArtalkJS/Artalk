import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router/auto'
import './style.css'
import App from './App.vue'

const app = createApp(App)

const router = createRouter({
  history: createWebHashHistory(),
})

app.use(router)
app.mount('#app')

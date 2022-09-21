import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import VueRouter from 'unplugin-vue-router/vite'
import { VueRouterAutoImports } from 'unplugin-vue-router'
import Components from 'unplugin-vue-components/vite'
import { HeadlessUiResolver } from 'unplugin-vue-components/resolvers'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    VueRouter({ importMode: 'sync' }),
    vue(),
    Components({ resolvers: [HeadlessUiResolver()] }),
    AutoImport({
      imports: ['vue', VueRouterAutoImports],
    }),
  ],
  server: {
    port: 23367,
  },
  css: {
    preprocessorOptions: {
      scss: {
         additionalData: `@import "./src/style/_variables.scss";@import "./src/style/_extends.scss";`
     },
    },
  },
})

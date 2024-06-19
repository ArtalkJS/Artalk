import { defineConfig } from 'vite'

import { ViteArtalkPluginKit } from '@artalk/plugin-kit'

export default defineConfig({
  build: {
    rollupOptions: {
      external: ['lightgallery', 'lightbox2', 'photoswipe', 'fancybox'],
      output: {
        dynamicImportInCjs: false,
      },
    },
  },
  plugins: [ViteArtalkPluginKit()],
})

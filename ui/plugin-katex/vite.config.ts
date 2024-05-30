import { defineConfig } from 'vite'

import { ViteArtalkPluginKit } from '@artalk/plugin-kit'

export default defineConfig({
  build: {
    rollupOptions: {
      external: ['katex', 'katex/dist/katex.min.css'],
      output: {
        globals: {
          katex: 'katex',
        },
      },
    },
  },
  plugins: [ViteArtalkPluginKit()],
})

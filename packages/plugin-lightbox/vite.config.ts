import { defineConfig } from 'vite'
import { resolve } from 'path'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig({
  base: './',
  build: {
    target: 'es2015',
    outDir: resolve(__dirname, "dist"),
    minify: 'terser',
    lib: {
      entry: resolve(__dirname, './main.ts'),
      name: 'artalk-plugin-lightbox',
      fileName: (format) => ((format == "umd") ? 'artalk-plugin-lightbox.js' : `artalk-plugin-lightbox.${format}.js`),
      formats: ["es", "umd", "iife"]
    },
    rollupOptions: {
      external: ['artalk', 'lightgallery'],
      output: {
        globals: {
          artalk: 'Artalk',
          lightgallery: 'lightgallery',
        },
        extend: true
      }
    }
  },
  plugins: [tsconfigPaths()],
})

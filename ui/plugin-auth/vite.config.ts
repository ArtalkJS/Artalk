import { defineConfig } from 'vite'
import { resolve, dirname } from 'node:path'
import tsconfigPaths from 'vite-tsconfig-paths'
import solidPlugin from 'vite-plugin-solid'
import cssInjectedByJsPlugin from 'vite-plugin-css-injected-by-js'
import { fileURLToPath } from 'node:url'

const __dirname = dirname(fileURLToPath(import.meta.url))

function getFileName(name: string, format: string) {
  if (format == 'umd') return `${name}.js`
  else if (format == 'cjs') return `${name}.cjs`
  else if (format == 'es') return `${name}.mjs`
  return `${name}.${format}.js`
}

export default defineConfig({
  build: {
    target: 'es2015',
    minify: 'terser',
    lib: {
      entry: resolve(__dirname, './main.tsx'),
      name: 'artalk-plugin-auth',
      fileName: (format) => getFileName('artalk-plugin-auth', format),
      formats: ['es', 'umd', 'cjs', 'iife'],
    },
    rollupOptions: {
      external: ['artalk'],
      output: {
        globals: {
          artalk: 'Artalk',
        },
        extend: true,
      },
    },
  },
  plugins: [tsconfigPaths(), solidPlugin(), cssInjectedByJsPlugin()],
})

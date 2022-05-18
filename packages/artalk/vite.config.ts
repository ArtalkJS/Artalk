import { defineConfig } from 'vite'
import { resolve } from 'path'
import commonConf from './vite-common.config'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig({
  build: {
    target: 'es2015',
    outDir: resolve(__dirname, "dist"),
    minify: 'terser',
    emptyOutDir: false,
    lib: {
      entry: resolve(__dirname, 'src/main.ts'),
      name: 'Artalk',
      fileName: (format) => ((format == "umd") ? 'Artalk.js' : `Artalk.${format}.js`),
      formats: ["es", "umd", "iife"]
    },
    rollupOptions: {
      output: {
        assetFileNames: assetInfo => {
          if (/\.css$/.test(assetInfo.name)) {
            return 'Artalk.css'
          }
          return "[name].[ext]"
        }
      }
    }
  },
  plugins: [tsconfigPaths()],
  ...commonConf
})

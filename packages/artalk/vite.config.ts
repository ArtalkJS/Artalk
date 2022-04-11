import { defineConfig } from 'vite'
import { resolve } from 'path'
import commonConf from './vite-common.config'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig({
  build: {
    target: 'es2015',
    outDir: resolve(__dirname, "dist"),
    minify: 'terser',
    lib: {
      entry: resolve(__dirname, 'src/main.ts'),
      name: 'Artalk',
      fileName: (format) => ((format == "umd") ? 'Artalk.js' : `Artalk.${format}.js`),
      formats: ["es", "umd", "iife"]
    },
    rollupOptions: {
      // 确保外部化处理那些你不想打包进库的依赖
      external: ['vue'],
      output: {
        assetFileNames: assetInfo => {
          if (/\.css$/.test(assetInfo.name)) {
            return 'Artalk.css'
          }
          return "[name].[ext]"
        },
        // 在 UMD 构建模式下为这些外部化的依赖提供一个全局变量
        globals: {
          vue: 'Vue'
        }
      }
    }
  },
  plugins: [tsconfigPaths()],
  ...commonConf
})

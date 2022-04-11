import { defineConfig } from 'vite'
import { resolve } from 'path'
import tsconfigPaths from 'vite-tsconfig-paths'
import { version as artalkVersion } from 'artalk/package.json'
import { version as sidebarVersion } from './package.json'

export default defineConfig({
  base: './',
  define: {
    ARTALK_VERSION: JSON.stringify(artalkVersion),
    ARTALK_SIDEBAR_VERSION: JSON.stringify(sidebarVersion)
  },
  build: {
    target: 'es2015',
    outDir: resolve(__dirname, "dist"),
    minify: 'terser',
    // lib: {
    //   entry: resolve(__dirname, 'src/main.ts'),
    //   name: 'ArtalkSidebar',
    //   fileName: (format) => ((format == "umd") ? 'ArtalkSidebar.js' : `ArtalkSidebar.${format}.js`),
    //   formats: ["es", "umd", "iife"]
    // },
    rollupOptions: {
      // 确保外部化处理那些你不想打包进库的依赖
      // external: ['artalk'],
      output: {
        assetFileNames: assetInfo => {
          if (/\.css$/.test(assetInfo.name)) {
            return 'ArtalkSidebar.css'
          }
          return "[name].[ext]"
        },
        // 在 UMD 构建模式下为这些外部化的依赖提供一个全局变量
        globals: {
          artalk: 'Artalk'
        }
      }
    }
  },
  plugins: [tsconfigPaths()],
  css: {
    preprocessorOptions: {
      less: {
         additionalData: `@import "./src/style/_variables.less";@import "./src/style/_extend.less";`
     },
    },
  },
})

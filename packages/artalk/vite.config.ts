import { defineConfig } from 'vite'
import { resolve } from 'path'
import baseConf from './vite-base.config'
import * as Utils from './src/lib/utils'

export default Utils.mergeDeep(baseConf, defineConfig({
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
        assetFileNames: (assetInfo) => (/\.css$/.test(assetInfo.name || '') ? "Artalk.css" : "[name].[ext]"),

        // Entry module "src/main.ts" is using named and default exports together.
        // Consumers of your bundle will have to use `Artalk["default"]` to access the default export
        exports: 'named'
      }
    }
  },
  define: {
    ARTALK_LITE: false,
  },
}))

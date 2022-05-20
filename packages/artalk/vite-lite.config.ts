import fullVersionConf from './vite.config'
import * as Utils from './src/lib/utils'

export default Utils.mergeDeep(fullVersionConf, {
  build: {
    lib: {
      fileName: (format) => ((format == "umd") ? 'ArtalkLite.js' : `ArtalkLite.${format}.js`),
    },
    rollupOptions: {
      external: ['marked'],
      output: {
        globals: {
          marked: 'marked',
        },
        assetFileNames: (assetInfo) => (/\.css$/.test(assetInfo.name) ? "ArtalkLite.css" : "[name].[ext]")
      }
   }
  },
  define: {
    ARTALK_LITE: true,
  },
})

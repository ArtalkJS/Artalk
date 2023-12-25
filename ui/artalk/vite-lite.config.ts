import { mergeConfig } from 'vite'
import fullVersionConf from './vite.config'

export default mergeConfig(fullVersionConf, {
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

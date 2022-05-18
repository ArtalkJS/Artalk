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
      }
   }
  },
})

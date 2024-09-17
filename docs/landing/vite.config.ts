import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    rollupOptions: {
      output: {
        assetFileNames: (assetInfo) => {
          const extType = assetInfo.name.split('.')[1]

          if (/png|jpe?g|svg|gif|tiff|bmp|ico/i.test(extType)) {
            return `landing/img/[name]-[hash][extname]`
          }

          return `landing/[name]-[hash][extname]`
        },
        chunkFileNames: 'landing/[name]-[hash].js',
        entryFileNames: 'landing/[name]-[hash].js',
      },
    },
  },
})

import { defineConfig } from 'vite'
import { resolve } from 'path'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig({
  base: './',
  build: {
    target: 'es2015',
    outDir: resolve(__dirname, 'dist'),
    minify: 'terser',
    lib: {
      entry: resolve(__dirname, './main.ts'),
      name: 'artalk-plugin-katex',
      fileName: (format) =>
        format == 'umd' ? 'artalk-plugin-katex.js' : `artalk-plugin-katex.${format}.js`,
      formats: ['es', 'umd', 'iife'],
    },
    rollupOptions: {
      // 确保外部化处理那些你不想打包进库的依赖
      external: ['artalk', 'katex'],
      output: {
        // 在 UMD 构建模式下为这些外部化的依赖提供一个全局变量
        globals: {
          artalk: 'Artalk',
          katex: 'katex',
        },
        extend: true,
      },
    },
  },
  plugins: [tsconfigPaths()],
})

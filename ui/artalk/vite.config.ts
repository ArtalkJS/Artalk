import { defineConfig } from 'vite'
import { resolve, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'
import checker from 'vite-plugin-checker'
import dts from 'vite-plugin-dts'
import { copyFileSync } from 'node:fs'

const __dirname = dirname(fileURLToPath(import.meta.url))

export function getFileName(name: string, format: string) {
  if (format == 'umd') return `${name}.js`
  else if (format == 'cjs') return `${name}.cjs`
  else if (format == 'es') return `${name}.mjs`
  return `${name}.${format}.js`
}

const name = process.env.ARTALK_LITE ? 'ArtalkLite' : 'Artalk'

export default defineConfig({
  root: __dirname,
  build: {
    target: 'es2015',
    outDir: resolve(__dirname, 'dist'),
    minify: 'terser',
    sourcemap: true,
    emptyOutDir: name === 'Artalk', // wait for https://github.com/qmhc/vite-plugin-dts/pull/291
    lib: {
      name: 'Artalk',
      fileName: (format: string) => getFileName(name, format),
      entry: resolve(__dirname, 'src/main.ts'),
      formats: ['es', 'umd', 'cjs', 'iife'],
    },
    rollupOptions: {
      external: name === 'ArtalkLite' ? ['marked'] : [],
      output: {
        globals:
          name === 'ArtalkLite'
            ? {
                marked: 'marked',
              }
            : {},
        assetFileNames: (assetInfo) =>
          /\.css$/.test(assetInfo.names?.[0] ?? '') ? `${name}.css` : '[name].[ext]',
        // @see https://github.com/rollup/rollup/issues/587
        //  and https://github.com/rollup/rollup/pull/631/files
        exports: 'named',
      },
    },
  },
  css: {
    preprocessorOptions: {
      scss: {
        loadPaths: [resolve(__dirname, 'src/style')],
        silenceDeprecations: ['import'], // https://sass-lang.com/documentation/breaking-changes/import/
        additionalData: `@import "variables";@import "extends";`,
      },
    },
  },
  resolve: {
    tsconfigPaths: true,
    alias: {
      '@': resolve(__dirname, 'src'),
      '~': resolve(__dirname),
    },
  },
  define: {
    ARTALK_LITE: false,
  },
  plugins: [
    {
      ...checker({
        typescript: true,
        eslint: {
          useFlatConfig: true,
          lintCommand: 'eslint .',
        },
      }),
      apply: 'serve',
    },
    // @see https://github.com/qmhc/vite-plugin-dts
    name === 'Artalk'
      ? dts({
          include: ['src'],
          exclude: ['src/**/*.{spec,test}.ts', 'dist'],
          bundleTypes: true,
          compilerOptions: {
            composite: false,
          },
          afterBuild: () => {
            // @see https://github.com/arethetypeswrong/arethetypeswrong.github.io/tree/main/packages/cli
            // @fix https://github.com/arethetypeswrong/arethetypeswrong.github.io/blob/main/docs/problems/FalseESM.md#consequences
            copyFileSync('dist/main.d.ts', 'dist/main.d.cts')
            // node10 resolution of `artalk/dist/ArtalkLite` finds dist/ArtalkLite.js via filesystem,
            // but needs dist/ArtalkLite.d.ts alongside it for types
            copyFileSync('dist/main.d.ts', 'dist/ArtalkLite.d.ts')
            copyFileSync('dist/main.d.cts', 'dist/ArtalkLite.d.cts')
          },
        })
      : null,
  ],
})

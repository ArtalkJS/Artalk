import { defineConfig, Options } from 'tsup'
import RawPlugin from 'esbuild-plugin-raw'

const shared: Options = {
  format: ['esm', 'cjs'],
  target: 'node14',
  platform: 'node',
  shims: true,
  splitting: false,
  bundle: true,
  sourcemap: true,
  clean: false,
  dts: true,
}

export default defineConfig([
  {
    ...shared,
    outDir: 'dist',
    entry: ['src/plugin/main.ts'],
    external: ['artalk', 'typescript', 'picocolors', '@microsoft/api-extractor'],
    esbuildPlugins: [RawPlugin()],
  },
  {
    ...shared,
    outDir: 'dist/@runtime',
    entry: ['src/runtime/main.ts'],
    format: ['esm'],
    dts: false,
  },
])

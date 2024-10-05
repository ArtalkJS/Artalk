import { defineConfig, Options } from 'tsup'

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
    entry: ['src/main.ts'],
  },
])

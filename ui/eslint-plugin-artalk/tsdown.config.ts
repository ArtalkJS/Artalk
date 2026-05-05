import { defineConfig } from 'tsdown'

export default defineConfig({
  entry: ['src/main.ts'],
  format: ['esm', 'cjs'],
  target: 'node14',
  platform: 'node',
  shims: true,
  sourcemap: true,
  clean: true,
  dts: true,
})

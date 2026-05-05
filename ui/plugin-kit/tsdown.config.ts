import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { defineConfig, type TsdownPlugin } from 'tsdown'

const rawPlugin: TsdownPlugin = {
  name: 'raw-file',
  resolveId(id, importer) {
    if (!id.endsWith('?raw')) return null
    const cleanId = id.replace(/\?raw$/, '')
    const resolved = importer ? resolve(importer, '..', cleanId) : resolve(cleanId)
    return `\0raw:${resolved}`
  },
  load(id) {
    if (!id.startsWith('\0raw:')) return null
    const file = id.slice('\0raw:'.length)
    return `export default ${JSON.stringify(readFileSync(file, 'utf-8'))}`
  },
}

export default defineConfig([
  {
    entry: ['src/plugin/main.ts'],
    outDir: 'dist',
    format: ['esm', 'cjs'],
    target: 'node14',
    platform: 'node',
    shims: true,
    sourcemap: true,
    clean: true,
    dts: true,
    deps: { neverBundle: ['artalk', 'typescript', 'picocolors', '@microsoft/api-extractor'] },
    plugins: [rawPlugin],
  },
  {
    entry: ['src/runtime/main.ts'],
    outDir: 'dist/@runtime',
    format: ['esm'],
    target: 'node14',
    platform: 'node',
    shims: true,
    sourcemap: true,
    clean: false,
    dts: false,
  },
])

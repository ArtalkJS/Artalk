import fs from 'node:fs'
import path from 'node:path'
import { build, LibraryOptions } from 'vite'
import { fileURLToPath } from 'node:url'
import { getFileName } from '../vite.config'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const i18nPath = path.join(__dirname, '../src/i18n/')
const outDir = path.resolve(__dirname, '../dist/i18n')

// empty outDir before build
fs.rmSync(outDir, { recursive: true, force: true })

const libraries: LibraryOptions[] = []

fs.readdirSync(i18nPath).forEach((f) => {
  if (['index.ts', 'external.ts'].includes(f)) return

  const filename = path.join(i18nPath, f)
  const lang = path.parse(filename).name
  const content = fs.readFileSync(filename)

  // only compile the external locale
  if (!content.toString().includes('defineLocaleExternal')) return

  libraries.push({
    entry: filename,
    name: lang,
    fileName: (format) => getFileName(lang, format),
  })
})

libraries.forEach(async (lib) => {
  await build({
    build: {
      target: 'es2015',
      outDir,
      minify: 'terser',
      emptyOutDir: false,
      lib: {
        ...lib,
        formats: ['umd', 'cjs', 'es'],
      },
    },
    configFile: false, // prevent load any vite config file
  })

  // crete d.ts file
  const dts = `import type { I18n } from '../main.d.cts'\ndeclare const locale: I18n\nexport = locale\n`
  fs.writeFileSync(path.join(outDir, `${lib.name!}.d.ts`), dts)
  fs.writeFileSync(path.join(outDir, `${lib.name!}.d.cts`), dts)
})

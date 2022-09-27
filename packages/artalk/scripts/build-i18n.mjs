import fs from 'fs'
import path from 'path'
import { build } from 'vite'
import { fileURLToPath } from 'url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const i18nPath = path.join(__dirname, '../src/i18n/')
const outDir = path.resolve(__dirname, '../dist/i18n')

// empty outDir before build
fs.rmSync(outDir, { recursive: true, force: true })

const libraries = []

fs.readdirSync(i18nPath).forEach(f => {
  if (['index.ts'].includes(f)) return

  const filename = path.join(i18nPath, f)
  const lang = path.parse(filename).name
  const content = fs.readFileSync(filename);

  // only compile the external locale
  if (!content.toString().includes('defineLocaleExternal')) return

  libraries.push({
    entry: filename,
    name: lang,
    fileName: (format) => ((format == "umd") ? `${lang}.js` : `${lang}.${format}.js`),
  })
})

libraries.forEach(async (lib) => {
  await build({
    build: {
      target: 'es2015',
      outDir: outDir,
      minify: 'terser',
      emptyOutDir: false,
      lib: {
        ...lib,
        formats: ['umd'],
      },
    },
    configFile: false, // prevent load any vite config file
  })
})

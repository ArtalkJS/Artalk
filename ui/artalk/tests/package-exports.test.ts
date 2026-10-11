import fs from 'node:fs'
import path from 'node:path'
import { describe, expect, it } from 'vitest'
import pkg from '../package.json'

const i18nDir = path.resolve(__dirname, '../src/i18n')

// Same rule as `scripts/build-i18n.ts`: only external locales are built to `dist/i18n`
const externalLocales = fs
  .readdirSync(i18nDir)
  .filter((f) => !['index.ts', 'external.ts'].includes(f))
  .filter((f) => fs.readFileSync(path.join(i18nDir, f), 'utf-8').includes('defineLocaleExternal'))
  .map((f) => path.parse(f).name)

describe('package exports', () => {
  it('finds the external locales', () => {
    expect(externalLocales.length).toBeGreaterThan(0)
  })

  it.each(externalLocales)('exports the external locale `artalk/i18n/%s`', (lang) => {
    expect(pkg.exports).toHaveProperty([`./i18n/${lang}`], {
      import: {
        types: `./dist/i18n/${lang}.d.ts`,
        default: `./dist/i18n/${lang}.mjs`,
      },
      require: {
        types: `./dist/i18n/${lang}.d.cts`,
        default: `./dist/i18n/${lang}.cjs`,
      },
    })
  })
})

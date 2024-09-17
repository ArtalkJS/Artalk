import path from 'node:path'
import { writeFileSync, readFileSync } from 'node:fs'
import { getFeatures } from '../src/features'
import { initI18nSSR } from '../src/i18n'

const build = async (lang: string, filenames: string[]) => {
  const t = await initI18nSSR(lang)

  const __dirname = new URL('.', import.meta.url).pathname
  const features = getFeatures(t)
    .map((o) => `* [${o.name}](${o.link}): ${o.desc}`)
    .join('\n')

  filenames.forEach((p) => {
    p = path.resolve(__dirname, p)

    const newReadme = readFileSync(p, 'utf-8').replace(
      /<!-- features -->[\s\S]*<!-- \/features -->/,
      `<!-- features -->\n${features}\n<!-- /features -->`,
    )
    writeFileSync(p, newReadme)

    console.log(`Updated ${p}`)
  })

  console.log(`[DONE] locale "${lang}" finished.\n`)
}

await build('en', ['../../../README.md', '../../docs/en/guide/intro.md'])
await build('zh-CN', ['../../../README.zh.md', '../../docs/zh/guide/intro.md'])

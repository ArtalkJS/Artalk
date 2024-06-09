import path from 'node:path'
import { writeFileSync, readFileSync } from 'node:fs'
import { Features } from '../src/Features'

const __dirname = new URL('.', import.meta.url).pathname
const features = Features.map(o => `* [${o.name}](${o.link}): ${o.desc}`).join('\n')

const pathList = [
  '../../../README.md',
  '../../docs/guide/intro.md',
]

pathList.forEach(p => {
  p = path.resolve(__dirname, p)

  const newReadme = readFileSync(p, 'utf-8').replace(/<!-- features -->[\s\S]*<!-- \/features -->/, `<!-- features -->\n${features}\n<!-- /features -->`)
  writeFileSync(p, newReadme)

  console.log(`Updated ${p}`)
})

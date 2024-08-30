import fs from 'node:fs'
import path from 'node:path'
import { createRequire } from 'node:module'

const asDefaultRE = /export\s*\{.*\w+\s*\bas\s+default\b.*\}\s*from\s*['"].+['"]/

export function hasExportDefault(content: string) {
  return content.includes('export default') || asDefaultRE.test(content)
}

export function ensureFolderExist(folder: string) {
  folder = path.dirname(folder)
  if (!fs.existsSync(folder)) {
    fs.mkdirSync(folder, { recursive: true })
  }
}

const pluginNameRE = /^(@artalk\/|artalk-)(plugin-[a-z0-9-]+)$/

export function getPluginNameFromPackageJson(packageJsonPath: string) {
  const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf-8'))
  const nameMatch = packageJson.name.match(pluginNameRE)
  if (!nameMatch || nameMatch.length !== 3) {
    throw new Error(
      `Invalid package.json name: "${packageJson.name}" should match "${pluginNameRE}"`,
    )
  }
  return `artalk-${nameMatch[2]}`
}

export function kababCaseToPascalCase(str: string) {
  return str
    .split('-')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join('')
}

export function getDistFileName(name: string, format: string) {
  if (format === 'umd') return `${name}.js`
  if (format === 'cjs') return `${name}.cjs`
  if (format === 'es') return `${name}.mjs`
  return `${name}.${format}.js`
}

export function getTypescriptLibFolder() {
  return createRequire(import.meta.url)
    .resolve('typescript')
    .replace(/\\+/g, '/') // fix windows path
    .replace(/node_modules\/typescript.*/, 'node_modules/typescript')
}

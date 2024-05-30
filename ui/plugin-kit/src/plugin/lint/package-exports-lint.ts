import fs from 'node:fs'
import logger from '../logger'

export function runPackageExportsLint(
  packageJsonPath: string,
  artalkPluginName: string,
): LintResult {
  const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf-8'))
  let modified = false

  // check if package.json has exports field
  if (!packageJson.exports) {
    modified = true
    // auto generate exports field
    packageJson.exports = {
      '.': {
        require: {
          types: `./dist/${artalkPluginName}.d.cts`,
          default: `./dist/${artalkPluginName}.cjs`,
        },
        default: {
          types: `./dist/${artalkPluginName}.d.ts`,
          default: `./dist/${artalkPluginName}.mjs`,
        },
      },
    }

    logger.info(`Auto generate exports field in package.json`)
  }

  const filedMap: Record<string, string> = {
    main: `./dist/${artalkPluginName}.js`,
    module: `./dist/${artalkPluginName}.mjs`,
    types: `./dist/${artalkPluginName}.d.ts`,
    type: 'module',
  }

  Object.entries(filedMap).forEach(([field, value]) => {
    if (packageJson[field] !== value) {
      modified = true
      packageJson[field] = value
      logger.info(`Auto update "${field}" field to "${value}" in package.json`)
    }
  })

  modified && fs.writeFileSync(packageJsonPath, JSON.stringify(packageJson, null, 2))

  return {
    ok: true,
  }
}

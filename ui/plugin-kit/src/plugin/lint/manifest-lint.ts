import fs from 'node:fs'

export function runManifestLint(packageJsonPath: string): LintResult {
  const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf-8'))

  // Check minAppVersion
  const versionRE = /^(\d+)\.(\d+)\.(\d+)$/
  if (!packageJson.minAppVersion) {
    return {
      ok: false,
      level: 'warn',
      message: '`minAppVersion` is missing in package.json',
    }
  }
  if (!versionRE.test(packageJson.minAppVersion)) {
    return {
      ok: false,
      level: 'error',
      message: `\`minAppVersion\` is not match semver format ${versionRE} in package.json`,
    }
  }

  return {
    ok: true,
  }
}

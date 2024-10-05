import { exec } from 'child_process'
import { accessSync } from 'fs'
import { resolve } from 'path'

const __dirname = new URL('.', import.meta.url).pathname
const forceBuild = process.env.ATK_DEV_TOOLCHAIN_FORCE_BUILD === '1'

function runCommand(command) {
  return new Promise((resolve, reject) => {
    exec(command, (error, stdout, stderr) => {
      if (error) {
        reject(stderr || stdout)
      } else {
        resolve(stdout)
      }
    })
  })
}

function checkFolderExists(pathname) {
  if (forceBuild) return false
  try {
    accessSync(resolve(__dirname, '../', pathname))
    return true
  } catch (error) {
    return false
  }
}

const green = '\x1b[32m'

async function build() {
  try {
    // Build Artalk Plugin Kit for plugin development
    if (!checkFolderExists('ui/plugin-kit/dist')) {
      await runCommand('pnpm build:plugin-kit')
      console.log(green, '[ArtalkDev] build @artalk/plugin-kit success')
    }
    // Build Artalk eslint plugin for lint checking
    if (!checkFolderExists('ui/eslint-plugin-artalk/dist')) {
      await runCommand('pnpm build:eslint-plugin')
      console.log(green, '[ArtalkDev] build eslint-plugin-artalk success')
    }
  } catch (error) {
    console.error('[ArtalkDev] Artalk UI development environment setup failed:', error)
  }
}

build()

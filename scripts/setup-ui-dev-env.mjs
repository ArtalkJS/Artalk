import { exec } from 'child_process'

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

async function build() {
  try {
    // Build Artalk Plugin Kit for plugin development
    await runCommand('pnpm build:plugin-kit')
    const green = '\x1b[32m'
    console.log(green, '[ArtalkDev] build @artalk/plugin-kit success')
    // Build Artalk eslint plugin for lint checking
    await runCommand('pnpm build:eslint-plugin')
    console.log(green, '[ArtalkDev] build eslint-plugin-artalk success')
  } catch (error) {
    console.error('[ArtalkDev] Artalk UI development environment setup failed:', error)
  }
}

build()

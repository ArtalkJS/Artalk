import { execFileSync, spawn } from 'node:child_process'
import { mkdtempSync, rmSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { fileURLToPath } from 'node:url'

// Use a disposable database and uploads, never the developer's Artalk data.
const workdir = mkdtempSync(join(tmpdir(), 'artalk-e2e-'))
const binary = fileURLToPath(new URL('../../../../bin/artalk', import.meta.url))
const config = fileURLToPath(new URL('./artalk.yml', import.meta.url))
const args = ['--config', config, '--workdir', workdir]
try {
  execFileSync(
    binary,
    [
      ...args,
      'admin',
      '--name',
      'Admin',
      '--email',
      'admin@example.com',
      '--password',
      'e2e-test-only',
    ],
    { stdio: 'inherit' },
  )
} catch (error) {
  rmSync(workdir, { recursive: true, force: true })
  throw error
}
const server = spawn(binary, [...args, 'server'], { stdio: 'inherit' })
for (const signal of ['SIGINT', 'SIGTERM'] as const) process.on(signal, () => server.kill(signal))
server.on('exit', (code) => {
  rmSync(workdir, { recursive: true, force: true })
  process.exit(code ?? 0)
})

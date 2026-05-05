import fs from 'fs'
import { fileURLToPath } from 'url'
import { dirname, resolve } from 'path'

const __dirname = dirname(fileURLToPath(import.meta.url))

export const RUNTIME_PATH = '/@vite-plugin-artalk-plugin-kit-runtime'

export const wrapVirtualPrefix = (id: `/${string}`): `virtual:${string}` => `virtual:${id.slice(1)}`

export const getRuntimeCode = () =>
  `${fs.readFileSync(resolve(__dirname, './@runtime/main.mjs'), 'utf-8')};`

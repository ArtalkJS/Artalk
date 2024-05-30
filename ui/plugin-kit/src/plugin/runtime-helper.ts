import fs from 'fs'
import { createRequire } from 'module'

const _require = createRequire(import.meta.url)

export const RUNTIME_PATH = '/@vite-plugin-artalk-plugin-kit-runtime'

export const wrapVirtualPrefix = (id: `/${string}`): `virtual:${string}` => `virtual:${id.slice(1)}`

export const getRuntimeCode = () =>
  `${fs.readFileSync(_require.resolve('./@runtime/main.js'), 'utf-8')};`

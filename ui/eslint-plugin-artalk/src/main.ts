import type { ESLint } from 'eslint'
import { artalkPlugin } from './artalk-plugin'

const configs = {
  recommended: {
    rules: { 'artalk/artalk-plugin': 'warn' as const },
  },
}

const plugin: ESLint.Plugin & { configs: typeof configs } = {
  rules: { 'artalk-plugin': artalkPlugin },
  configs,
}

export default plugin

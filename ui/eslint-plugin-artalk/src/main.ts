import { artalkPlugin } from './artalk-plugin'

export const rules = {
  'artalk-plugin': artalkPlugin,
}

export default {
  rules,
  configs: {
    recommended: { plugins: ['artalk'], rules: { 'artalk/artalk-plugin': 'warn' } },
  },
}

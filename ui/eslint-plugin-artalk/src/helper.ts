import type { Rule } from 'eslint'

const docsBaseUrl =
  'https://github.com/ArtalkJS/Artalk/tree/master/ui/eslint-plugin-artalk/README.md#rule-'

export function createRule(def: {
  name: string
  meta: Rule.RuleMetaData
  defaultOptions?: unknown[]
  create: Rule.RuleModule['create']
}): Rule.RuleModule {
  return {
    meta: {
      ...def.meta,
      docs: { ...def.meta.docs, url: `${docsBaseUrl}${def.name}` },
    },
    create: def.create,
  }
}

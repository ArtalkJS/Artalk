import { ESLintUtils } from '@typescript-eslint/utils'

export const createRule = ESLintUtils.RuleCreator(
  (name) =>
    `https://github.com/ArtalkJS/Artalk/tree/master/ui/eslint-plugin-artalk/README.md#rule-${name}`,
)

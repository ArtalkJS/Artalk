const path = require('node:path')

module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  env: { browser: true },
  parserOptions: {
    project: ['./ui/*/tsconfig.json'],
    tsconfigRootDir: __dirname,
  },
  extends: [
    'airbnb-base',
    'airbnb-typescript/base',
    'prettier',
    'plugin:import/recommended',
    'plugin:import/typescript',
    'plugin:compat/recommended',
  ],
  plugins: ['@typescript-eslint', 'import'],
  rules: {
    'no-alert': 'warn',
    'no-unused-vars': 0,
    'no-plusplus': 0,
    'no-param-reassign': 0,
    'no-console': 0,
    'no-underscore-dangle': 0,
    'class-methods-use-this': 0,
    'spaced-comment': 0,
    'no-lonely-if': 0,
    'prefer-destructuring': 0,
    'import/prefer-default-export': 0,
    'import/no-extraneous-dependencies': 0,
    '@typescript-eslint/lines-between-class-members': 0,
    '@typescript-eslint/no-unused-vars': 'off',
    '@typescript-eslint/no-use-before-define': 0,
    '@typescript-eslint/naming-convention': 0,
    '@typescript-eslint/no-useless-constructor': 0,
    '@typescript-eslint/no-unused-expressions': 0, // for `func && func()` expressions
  },
  settings: {
    'import/resolver': {
      typescript: {
        project: ['ui/artalk/tsconfig.json'].map((p) => path.resolve(__dirname, p)),
      },
    },
    polyfills: ['AbortController'],
  },
}

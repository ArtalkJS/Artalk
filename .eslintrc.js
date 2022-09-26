module.exports = {
  root: true,
  parser: "@typescript-eslint/parser",
  env: { "browser": true },
  parserOptions: {
    project: ['./packages/*/tsconfig.json'],
    tsconfigRootDir: __dirname,
  },
  extends: [
    'airbnb-base',
    'airbnb-typescript/base',
    'prettier',
    'plugin:import/recommended',
    'plugin:import/typescript',
    "plugin:compat/recommended",
  ],
  plugins: ['@typescript-eslint', 'import'],
  rules: {
    'no-alert': 'off',
    'no-unused-vars': 'off',
    'no-plusplus': 0,
    'no-param-reassign': 0,
    'no-console': 0,
    'no-constructor-return': 0,
    'class-methods-use-this': 0,
    'spaced-comment': 0,
    'no-lonely-if': 0,
    'prefer-destructuring': 0,
    'import/no-cycle': 0,
    '@typescript-eslint/lines-between-class-members': 0,
    '@typescript-eslint/no-unused-vars': 'off',
    '@typescript-eslint/no-use-before-define': 0,
  },
  settings: {
    'import/resolver': {
      typescript: {
        "project": [__dirname + "/packages/*/tsconfig.json"]
      },
    },
    'polyfills': [
      'AbortController'
    ]
  },
}

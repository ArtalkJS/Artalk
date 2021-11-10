module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  extends: [
    'airbnb',
    'airbnb-typescript',
    'prettier',
    'plugin:import/recommended',
  ],
  plugins: ['@typescript-eslint', 'import'],
  rules: {
    '@typescript-eslint/lines-between-class-members': 0,
    'no-unused-vars': 'off',
    '@typescript-eslint/no-unused-vars': 'off',
    'no-plusplus': 0,
    'no-param-reassign': 0,
    'no-console': 0,
    'react/sort-comp': 0,
    '@typescript-eslint/no-use-before-define': 0,
    'class-methods-use-this': 0,
    'import/no-cycle': 0,
    'spaced-comment': 0,
    'no-lonely-if': 0,
  },
  settings: {
    'import/resolver': {
      typescript: {},
    },
  },
  parserOptions: {
    project: './tsconfig.json',
  },
}

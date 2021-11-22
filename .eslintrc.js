module.exports = {
  root: true,
  parser: "@typescript-eslint/parser",
  extends: [
    'airbnb-base',
    'airbnb-typescript/base',
    'prettier',
    'plugin:import/recommended',
    'plugin:import/typescript',
  ],
  plugins: ['@typescript-eslint', 'import'],
  rules: {
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
  parserOptions: {
    project: './tsconfig.json',
  },
  settings: {
    'import/resolver': {
      typescript: {},
    },
  },
}

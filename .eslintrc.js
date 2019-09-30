module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  env: {
    browser: true
  },
  extends: [
    'airbnb-typescript/base',
    'prettier',
    'prettier/@typescript-eslint',
  ],
  plugins: ['@typescript-eslint', 'prettier', 'import'],
  rules: {
    'no-console': 0,
    'arrow-parens': 0,
    'generator-star-spacing': 0,
    'no-unused-vars': 'off',
    '@typescript-eslint/no-unused-vars': 'off',
    'class-methods-use-this': 0,
    'no-param-reassign': 0,
    'no-plusplus': 0,
    'global-require': 0,
    'import/no-cycle': 0, /** TODO: 待优化（未来计划） */
    'import/no-extraneous-dependencies': 0,
    'lines-between-class-members': 0
  },
  settings: {
    'import/resolver': {
      // use <root>/tsconfig.json
      'typescript': {},
    }
  }
}

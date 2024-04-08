const path = require('node:path')

module.exports = {
  root: true,
  extends: [path.join(__dirname, '../../.eslintrc.cjs')],
}

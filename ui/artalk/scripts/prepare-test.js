const fs = require('fs')
const path = require('path')

const rootPath = path.resolve(__dirname, '../../../../')

const src = path.resolve(rootPath, 'data/ui_test.db')
const dest = path.resolve(rootPath, 'data/artalk.db')

fs.copyFileSync(src, dest)
console.log('ui_test.db copied')

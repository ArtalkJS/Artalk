const webpack = require('webpack')
const UglifyJsPlugin = require('uglifyjs-webpack-plugin')
const path = require('path')

const ROOT_PATH = path.resolve(__dirname)
const SRC_PATH = path.resolve(ROOT_PATH, 'src')
let BUILD_PATH = path.resolve(ROOT_PATH, 'dist')

const VERSION = require('./package.json').version
const BANNER =
  'Artalk v' + VERSION + '\n' +
  '(c) 2016-' + new Date().getFullYear() + ' qwqaq.com\n' +
  'Link: https://github.com/qwqcode/Artalk'

module.exports = (env, options) => {
  const dev = options.mode === 'development'
  // console.log(`MODE -> ${options.mode}`)
  /* if (!dev) {
    BUILD_PATH = path.resolve(ROOT_PATH, `dist/${VERSION}`)
  } */

  let conf = {
    entry: {
      'Artalk': [`${SRC_PATH}/Artalk.scss`, `${SRC_PATH}/Artalk.js`]
    },
    output: {
      path: BUILD_PATH,
      filename: '[name].js',
      library: '[name]',
      libraryTarget: 'umd'
    },
    module: {
      rules: [{
        test: /\.js$/,
        enforce: 'pre',
        loader: 'eslint-loader',
        include: [`${SRC_PATH}`]
      }, {
        test: /\.js$/,
        include: [SRC_PATH],
        exclude: /node_modules/,
        loader: 'babel-loader',
        options: {
          plugins: ['syntax-dynamic-import'],
          presets: [
            ['env', {
              modules: false
            }]
          ]
        }
      }, {
        test: /\.scss$/,
        use: [
          'style-loader',
          {
            loader: 'css-loader',
            options: {
              sourceMap: true
            }
          },
          'postcss-loader',
          {
            loader: 'sass-loader',
            options: {
              sourceMap: true
            }
          }
        ],
        include: SRC_PATH
      }, {
        test: /\.css$/,
        use: [
          'style-loader',
          'css-loader',
          'postcss-loader'
        ]
      }, {
        test: /\.(png|jpg|gif|svg)$/,
        use: ['url-loader?limit=1024*10']
      }]
    },
    devtool: 'cheap-module-source-map',
    devServer: {
      contentBase: path.resolve(ROOT_PATH, 'demo'),
      open: true,
      hot: true,
      inline: true,
      publicPath: '/dist/',
      compress: true,
      stats: 'errors-only',
      overlay: {
        errors: true,
        warnings: true
      }
    },
    plugins: []
  }

  if (!dev) {
    conf.plugins.push(new webpack.BannerPlugin(BANNER))
    conf.optimization = {
      minimizer: [
        new UglifyJsPlugin({
          uglifyOptions: {
            beautify: false,
            sourceMap: false,
            comments: false,
            mangle: true,
            compress: {
              drop_console: true,
              warnings: false,
              collapse_vars: true,
              reduce_vars: true
            }
          }
        })
      ]
    }
  } else {
    conf.plugins.push(new webpack.NamedModulesPlugin())
    conf.plugins.push(new webpack.HotModuleReplacementPlugin())
  }

  return conf
}

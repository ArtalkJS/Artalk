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
  const isDev = options.mode === 'development'
  // console.log(`MODE -> ${options.mode}`)
  /* if (!isDev) {
    BUILD_PATH = path.resolve(ROOT_PATH, `dist/${VERSION}`)
  } */

  let conf = {
    entry: {
      'Artalk': [`${SRC_PATH}/Artalk.js`]
    },
    output: {
      path: BUILD_PATH,
      filename: '[name].js',
      library: '[name]',
      libraryTarget: 'umd',
      libraryExport: 'default',
      umdNamedDefine: true
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
          presets: ['vue-app']
        }
      }, {
        test: /\.scss$/,
        use: [
          'style-loader',
          {
            loader: 'css-loader',
            options: {
              sourceMap: isDev
            }
          },
          'postcss-loader',
          {
            loader: 'sass-loader',
            options: {
              sourceMap: isDev,
              includePaths: [SRC_PATH],
              data: '@import "css/_variables.scss";'
            }
          }
        ],
        include: SRC_PATH
      }, {
        test: /\.css$/,
        use: [
          'style-loader',
          {
            loader: 'css-loader',
            options: {
              sourceMap: isDev
            }
          },
          'postcss-loader'
        ]
      }, {
        test: /\.(png|jpg|gif|svg)$/,
        use: ['url-loader?limit=1024*10']
      }, {
        test: /\.ejs$/,
        loader: 'ejs-compiled-loader',
        options: {
          htmlmin: !isDev,
          htmlminOptions: {
            removeComments: !isDev
          }
        }
      }]
    },
    devtool: 'cheap-module-source-map',
    devServer: {
      contentBase: path.resolve(ROOT_PATH, 'demo'),
      open: true,
      hot: true,
      inline: true,
      publicPath: '/',
      compress: true,
      stats: 'errors-only',
      overlay: {
        errors: true,
        warnings: true
      }
    },
    externals: {
      jquery: 'jQuery'
    },
    plugins: [
      new webpack.DefinePlugin({
        ARTALK_VERSION: `"${VERSION}"`
      })
    ]
  }

  if (!isDev) {
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

const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const UglifyJsPlugin = require('uglifyjs-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const OptimizeCSSAssetsPlugin = require('optimize-css-assets-webpack-plugin')
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
  const isDev = (options.mode === 'development')
  const isProd = !isDev
  // console.log(`MODE -> ${options.mode}`)
  /* if (isProd) {
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
    plugins: [
      new HtmlWebpackPlugin({
        title: `Artalk DEMO`,
        filename: `${ROOT_PATH}/index.html`,
        template: `${ROOT_PATH}/index-tpl.ejs`,
        inject: 'head',
        hash: true
      }),
      new webpack.DefinePlugin({
        ARTALK_VERSION: `"${VERSION}"`
      }),
      new MiniCssExtractPlugin({
        filename: '[name].css'
      })
    ],
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
          MiniCssExtractPlugin.loader,
          {
            loader: 'css-loader',
            options: {
              sourceMap: isDev,
              minimize: isProd
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
        test: /\.(png|jpg|gif|svg)$/,
        use: ['url-loader?limit=1024*10']
      }, {
        test: /\.ejs$/,
        loader: 'ejs-compiled-loader',
        options: {
          htmlmin: isProd,
          htmlminOptions: {
            removeComments: true,
            collapseWhitespace: true,
            preserveLineBreaks: true,
            conservativeCollapse: true
          }
        }
      }]
    },
    devtool: isDev ? 'cheap-module-source-map' : false,
    devServer: {
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
    externals: {
      jquery: 'jQuery'
    }
  }

  if (isProd) {
    conf.plugins.push(new webpack.BannerPlugin(BANNER))
    conf.optimization = {
      minimizer: [
        new UglifyJsPlugin({
          uglifyOptions: {
            sourceMap: false,
            beautify: false,
            comments: false,
            mangle: true,
            compress: {
              drop_console: false, // console.log
              warnings: false,
              collapse_vars: true,
              reduce_vars: true
            }
          }
        }),
        new OptimizeCSSAssetsPlugin({})
      ]
    }
  } else {
    conf.plugins.push(new webpack.NamedModulesPlugin())
    conf.plugins.push(new webpack.HotModuleReplacementPlugin())
  }

  return conf
}

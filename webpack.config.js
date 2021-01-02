/* eslint-disable global-require */
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const HtmlWebpackHarddiskPlugin = require('html-webpack-harddisk-plugin')
const UglifyJsPlugin = require('uglifyjs-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const OptimizeCSSAssetsPlugin = require('optimize-css-assets-webpack-plugin')
const path = require('path')

const ROOT_PATH = path.resolve(__dirname)
const SRC_PATH = path.resolve(ROOT_PATH, 'src')
const BUILD_PATH = path.resolve(ROOT_PATH, 'dist')

const VERSION = require('./package.json').version

const BANNER =
  `Artalk v${  VERSION  }\n` +
  `(c) 2016-${  new Date().getFullYear()  } qwqaq.com\n` +
  `Link: https://github.com/ArtalkJS/Artalk`

module.exports = (env, argv) => {
  const NODE_ENV = argv.mode || 'development'
  const IS_DEV = NODE_ENV === 'development'
  const IS_PROD = !IS_DEV
  /* if (IS_PROD) {
    BUILD_PATH = path.resolve(ROOT_PATH, `dist/${VERSION}`)
  } */

  const conf = {
    entry: {
      Artalk: [`${SRC_PATH}/Artalk.ts`]
    },
    output: {
      path: BUILD_PATH,
      filename: '[name].js',
      library: '[name]',
      libraryTarget: 'umd',
      libraryExport: '',
      umdNamedDefine: true,
      globalObject: 'this'
    },
    plugins: [
      new HtmlWebpackPlugin({
        title: 'Artalk DEMO',
        filename: `${ROOT_PATH}/index.html`,
        template: `${ROOT_PATH}/index-tpl.ejs`,
        inject: 'head',
        hash: true,
        alwaysWriteToDisk: true
      }),
      new HtmlWebpackHarddiskPlugin(),
      new webpack.DefinePlugin({
        ARTALK_VERSION: `"${VERSION}"`,
        'process.env.NODE_ENV': JSON.stringify(NODE_ENV)
      }),
      new MiniCssExtractPlugin({
        filename: '[name].css'
      })
    ],
    module: {
      rules: [{
        test: /\.(js|ts)$/,
        enforce: 'pre',
        exclude: /node_modules/,
        use: {
          loader: 'eslint-loader',
          options: {
            formatter: require('eslint-friendly-formatter')
          }
        }
      }, {
        test: /\.ts$/,
        exclude: /node_modules/,
        use: [
          {
            loader: 'babel-loader',
            options: {
              cacheDirectory: true,
            }
          },
          {
            loader: 'ts-loader',
          }
        ]
      }, {
        test: /\.js$/,
        exclude: /node_modules/,
        use: [
          {
            loader: 'babel-loader',
            options: {
              cacheDirectory: true,
            }
          }
        ]
      }, {
        test: /\.less$/,
        use: [
          MiniCssExtractPlugin.loader,
          {
            loader: 'css-loader',
            options: {
              sourceMap: IS_DEV
            }
          },
          'postcss-loader',
          {
            loader: 'less-loader',
            options: {
              sourceMap: IS_DEV,
              lessOptions: {
                paths: [SRC_PATH]
              },
              prependData: '@import "css/_variables.less";@import "css/_extend.less";'
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
          htmlmin: IS_PROD,
          htmlminOptions: {
            removeComments: true,
            collapseWhitespace: true,
            preserveLineBreaks: true,
            conservativeCollapse: true
          }
        }
      }]
    },
    resolve: {
      alias: {
        '@': path.join(__dirname, './src'),
        '~': path.join(__dirname, './')
      },
      extensions: ['.ts', '.js', '.vue', '.json', '.less', '.css', '.node']
    },
    devtool: IS_DEV ? 'cheap-module-source-map' : false,
    devServer: {
      open: true,
      hot: true,
      inline: true,
      publicPath: '/dist/', // 打包文件生成路径（内存中）
      contentBase: path.join(__dirname, '/'),
      stats: 'errors-only',
      overlay: {
        errors: true,
        warnings: true
      }
    }
  }

  if (IS_PROD) {
    conf.plugins.push(new webpack.BannerPlugin(BANNER))
    conf.optimization = {
      minimizer: [
        new UglifyJsPlugin({
          uglifyOptions: {
            sourceMap: false,
            beautify: false,
            comments: false,
            mangle: true,
            warnings: false,
            compress: {
              drop_console: false, // console.log
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

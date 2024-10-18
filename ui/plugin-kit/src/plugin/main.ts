import path from 'node:path'
import fs from 'node:fs'
import type { LibraryOptions, Plugin } from 'vite'
import ts from 'typescript'
import type { ConfigPartial } from 'artalk'

import { RUNTIME_PATH, getRuntimeCode, wrapVirtualPrefix } from './runtime-helper'
import { getInjectHTMLTags, hijackIndexPage } from './dev-page'
import { runManifestLint, runExportLint } from './lint'
import * as utils from './utils'
import logger from './logger'
import { generateDts, rollupDeclarationFiles } from './dts'
import { runPackageExportsLint } from './lint/package-exports-lint'
import { ViteArtalkPluginKitCtx } from './context'

const tsRE = /\.(m|c)?tsx?$/
const dtsRE = /\.d\.ts$/

export interface ViteArtalkPluginKitOptions {
  /**
   * The options for Artalk instance initialization in the dev page.
   */
  artalkInitOptions?: ConfigPartial
}

export const ViteArtalkPluginKit = (opts: ViteArtalkPluginKitOptions = {}): Plugin => {
  const ctx = new ViteArtalkPluginKitCtx()

  return {
    name: 'vite-plugin-artalk-dev-kit',
    enforce: 'pre',

    async config(_config) {
      return {
        // override vite user config
        build: {
          target: 'es2015',
          rollupOptions: {
            external: ['artalk'],
            output: {
              globals: {
                artalk: 'Artalk',
              },
              exports: 'named',
              extend: true,
            },
          },
        },
      }
    },

    async configResolved(config) {
      ctx.isDev = !config.isProduction && config.command === 'serve'
      ctx.rootDir = config.root
      ctx.outDir = path.resolve(ctx.rootDir, config.build.outDir || 'dist')
      ctx.dtsTempDir = path.resolve(ctx.rootDir, ctx.dtsTempDir)
      ctx.packageJsonPath = path.resolve(ctx.rootDir, 'package.json')
      ctx.pluginName = utils.getPluginNameFromPackageJson(ctx.packageJsonPath)
      ctx.entryExportName = `Artalk${utils.kababCaseToPascalCase(/plugin-(.*)/.exec(ctx.pluginName)![1])}Plugin`

      // entry
      const _lib = config.build.lib || {}
      config.build.lib = <LibraryOptions>{
        name: ctx.pluginName,
        entry: path.resolve(ctx.rootDir, 'src/main.ts'),
        formats: ['es', 'umd', 'cjs', 'iife'],
        fileName: (format) => utils.getDistFileName(ctx.pluginName, format),
        ..._lib,
      }
      if (typeof config.build.lib.entry !== 'string')
        throw new Error('Only support single entry file now')
      ctx.entryPath = path.resolve(ctx.rootDir, config.build.lib.entry)
    },

    configureServer(_server) {
      ctx.viteServer = _server
      _server.hot.on('artalk-plugin-kit:remote-add', (val, client) => {
        client.send('artalk-plugin-kit:update', { test: `Hello World at ${new Date()}` })
      })
      hijackIndexPage(_server.middlewares, _server.transformIndexHtml, opts.artalkInitOptions || {})
    },

    transformIndexHtml(html) {
      if (!ctx.isDev) return

      return {
        html,
        tags: getInjectHTMLTags(path.relative(ctx.rootDir, ctx.entryPath)),
      }
    },

    resolveId(id) {
      if (id === RUNTIME_PATH) return wrapVirtualPrefix(id)
      return undefined
    },

    load(id) {
      if (id === wrapVirtualPrefix(RUNTIME_PATH)) return getRuntimeCode()
      return undefined
    },

    async buildStart() {
      const tsConfigPath = ts.findConfigFile(ctx.rootDir, ts.sys.fileExists)
      if (!tsConfigPath) {
        logger.error(`Cannot find tsconfig.json in "${ctx.rootDir}"`)
        process.exit(1)
      }
      ctx.tsConfigPath = tsConfigPath

      // read tsconfig.json
      const readConfigFile = ts.readConfigFile(tsConfigPath, ts.sys.readFile)
      if (readConfigFile.error) {
        logger.error(`Cannot read tsconfig.json in "${tsConfigPath}"`)
        logger.error(
          ts.formatDiagnosticsWithColorAndContext([readConfigFile.error], {
            getCanonicalFileName: (fileName) => fileName,
            getCurrentDirectory: ts.sys.getCurrentDirectory,
            getNewLine: () => ts.sys.newLine,
          }),
        )
        process.exit(1)
      }

      // parse tsconfig.json
      const parsedConfig = ts.parseJsonConfigFileContent(
        readConfigFile.config,
        ts.sys,
        path.dirname(tsConfigPath),
      )
      if (parsedConfig.errors.length > 0) {
        logger.error(`Error parsing tsconfig.json in "${tsConfigPath}"`)
        logger.error(
          ts.formatDiagnosticsWithColorAndContext(parsedConfig.errors, {
            getCanonicalFileName: (fileName) => fileName,
            getCurrentDirectory: ts.sys.getCurrentDirectory,
            getNewLine: () => ts.sys.newLine,
          }),
        )
        process.exit(1)
      }
      const compilerOptions = parsedConfig.options

      ctx.tsCompilerOptionsRaw = compilerOptions
      ctx.tsCompilerOptions = {
        ...ctx.tsCompilerOptionsRaw,
        noEmit: false,
        declaration: true,
        emitDeclarationOnly: true,
        noUnusedParameters: false,
        checkJs: false,
        skipLibCheck: true,
        preserveSymlinks: false,
        noEmitOnError: undefined,
        target: ts.ScriptTarget.ESNext,
        moduleResolution: ts.ModuleResolutionKind.Bundler,
        outDir: '.',
        declarationDir: '.',
      }

      // Pre-build lints
      const lintTasks: (() => LintResult)[] = [
        () => runPackageExportsLint(ctx.packageJsonPath, ctx.pluginName),
        () => runManifestLint(ctx.packageJsonPath),
        () => runExportLint(ctx.entryPath, ctx.entryExportName),
      ]

      lintTasks.forEach((task) => {
        const result = task()
        if (!result.ok) {
          result.message && logger[result.level || 'info'](result.message)
          if (result.level === 'error') process.exit(1)
        }
      })
    },

    transform(code, id) {
      if (tsRE.test(id) && id === ctx.entryPath) {
        return {
          code: `${code}
          // Mount plugin to browser window global
          if (typeof window !== 'undefined') {
            !window.ArtalkPlugins && (window.ArtalkPlugins = {})
            window.ArtalkPlugins.${ctx.entryExportName} = ${ctx.entryExportName}
            window.Artalk?.use(${ctx.entryExportName})
          }`,
          map: null,
        }
      }

      return undefined
    },

    watchChange(id) {
      ctx.isBundled = false
    },

    async writeBundle() {
      if (ctx.isBundled) return
      ctx.isBundled = true

      // Cleanup dts temp folder
      if (fs.existsSync(ctx.dtsTempDir)) {
        fs.rmSync(ctx.dtsTempDir, { recursive: true })
      }

      // Generate dts files
      const dtsFiles = generateDts([ctx.entryPath], ctx.tsCompilerOptions!)

      // Save dts files to temp folder
      for (const [dtsFilename, dtsContent] of Object.entries(dtsFiles)) {
        const dtsPath = path.resolve(ctx.dtsTempDir, dtsFilename)
        utils.ensureFolderExist(dtsPath)
        fs.writeFileSync(dtsPath, dtsContent, 'utf-8')
      }

      // Prepare index dts
      const dtsIndexFileName = `${ctx.pluginName}.d.ts`
      const dtsIndexFilePath = path.normalize(path.resolve(ctx.dtsTempDir, dtsIndexFileName))

      let fromPath = path.normalize(path.basename(ctx.entryPath)).replace(tsRE, '')
      const relativePathRE = /^(\.\/|\.\.\/)/
      fromPath = !relativePathRE.test(fromPath) ? `./${fromPath}` : fromPath
      const content = `export * from '${fromPath}'\n`
      fs.writeFileSync(dtsIndexFilePath, content, 'utf-8')

      // Rollup
      rollupDeclarationFiles({
        entryPath: dtsIndexFilePath,
        fileName: dtsIndexFileName,
        rootDir: ctx.rootDir,
        outDir: ctx.dtsTempDir,
        tsConfigPath: ctx.tsConfigPath,
        libFolder: utils.getTypescriptLibFolder(),
      })

      // Copy dts files to dist
      const distDtsPath = path.resolve(ctx.outDir, dtsIndexFileName)
      fs.copyFileSync(dtsIndexFilePath, distDtsPath)

      // Copy .d.cts for dtsIndex
      fs.copyFileSync(distDtsPath, distDtsPath.replace(dtsRE, '.d.cts'))

      // Cleanup dts temp folder
      if (fs.existsSync(ctx.dtsTempDir)) {
        fs.rmSync(ctx.dtsTempDir, { recursive: true })
      }

      logger.info(`Generate dts files successfully!`)
    },
  }
}

import type { ViteDevServer } from 'vite'
import type ts from 'typescript'

export class ViteArtalkPluginKitCtx {
  /**
   * The artalk plugin name (kabab-case)
   */
  pluginName = ''

  /**
   * Whether the current environment is development
   */
  isDev = true

  /**
   * The root directory of the working project
   */
  rootDir = './'

  /**
   * The entry path of the plugin (main.ts)
   */
  entryPath = './src/main.ts'

  /**
   * The export name of entry code (PascalCase)
   */
  entryExportName = ''

  /**
   * The output directory of the plugin (dist)
   */
  outDir = './dist/'

  /**
   * Dts temporary directory
   *
   * (for generate dts, first store in this directory, then copy to outDir)
   */
  dtsTempDir = '.atk-vite-dts-temp' // TODO: node_modules/.atk-vite-dts-temp not working for rollup dts

  /**
   * The package.json path
   */
  packageJsonPath = './package.json'

  /**
   * The vite server instance
   */
  viteServer?: ViteDevServer

  /**
   * The TypeScript config path
   */
  tsConfigPath = 'tsconfig.json'

  /**
   * The TypeScript compiler options (for generate dts)
   */
  tsCompilerOptions?: ts.CompilerOptions

  /**
   * The TypeScript compiler options (raw)
   */
  tsCompilerOptionsRaw?: ts.CompilerOptions

  /**
   * Builded flag (for watchChange)
   */
  isBundled = false
}

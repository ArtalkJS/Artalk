import path from 'node:path'
import { Extractor, ExtractorConfig } from '@microsoft/api-extractor'

import ts from 'typescript'
import type { ExtractorLogLevel, IExtractorInvokeOptions } from '@microsoft/api-extractor'
import logger from './logger'

export function generateDts(
  fileNames: string[],
  options: ts.CompilerOptions,
): Record<string, string> {
  const createdFiles: Record<string, string> = {}
  const host = ts.createCompilerHost(options)
  host.writeFile = (fileName: string, contents: string) => {
    createdFiles[fileName] = contents
  }

  const program = ts.createProgram(fileNames, options, host)
  const diagnostics = [
    ...program.getDeclarationDiagnostics(),
    ...program.getSemanticDiagnostics(),
    ...program.getSyntacticDiagnostics(),
  ]

  if (diagnostics?.length) {
    logger.error(ts.formatDiagnosticsWithColorAndContext(diagnostics, host))
    process.exit(1)
  }

  program.emit()

  return createdFiles
}

export interface BundleOptions {
  rootDir: string
  tsConfigPath?: string
  outDir: string
  entryPath: string
  fileName: string
  libFolder?: string
  rollupOptions?: IExtractorInvokeOptions
}

const dtsRE = /\.d\.tsx?$/

export function rollupDeclarationFiles({
  rootDir,
  tsConfigPath,
  outDir,
  entryPath,
  fileName,
  libFolder,
  rollupOptions = {},
}: BundleOptions) {
  if (!dtsRE.test(fileName)) {
    fileName += '.d.ts'
  }

  const extractorConfig = ExtractorConfig.prepare({
    configObject: {
      projectFolder: rootDir,
      mainEntryPointFilePath: entryPath,
      compiler: {
        tsconfigFilePath: tsConfigPath,
      },
      apiReport: {
        enabled: false,
        reportFileName: '<unscopedPackageName>.api.md',
      },
      docModel: {
        enabled: false,
      },
      dtsRollup: {
        enabled: true,
        publicTrimmedFilePath: path.resolve(outDir, fileName),
      },
      tsdocMetadata: {
        enabled: false,
      },
      messages: {
        compilerMessageReporting: {
          default: {
            logLevel: 'none' as ExtractorLogLevel.None,
          },
        },
        extractorMessageReporting: {
          default: {
            logLevel: 'none' as ExtractorLogLevel.None,
          },
        },
      },
    },
    configObjectFullPath: path.resolve(rootDir, 'api-extractor.json'),
    packageJsonFullPath: path.resolve(rootDir, 'package.json'),
  })

  const result = Extractor.invoke(extractorConfig, {
    localBuild: false,
    showVerboseMessages: false,
    showDiagnostics: false,
    typescriptCompilerFolder: libFolder ? path.resolve(libFolder) : undefined,
    ...rollupOptions,
  })

  return result.succeeded
}

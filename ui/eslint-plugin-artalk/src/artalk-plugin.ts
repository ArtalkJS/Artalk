import type { Context } from '../../artalk'
import {
  isPluginName,
  checkPluginFunction,
  checkInjectCallOutsideArtalkPlugin,
} from './artalk-plugin-checkers'
import { createRule } from './helper'
type _ = Context // for IDE jump-to-definition

const depsMap: DepsStore = new Map()

export const artalkPlugin = createRule({
  name: 'artalk-plugin',
  meta: {
    type: 'problem',
    docs: {
      description:
        'Enforce best practices for ArtalkPlugin arrow functions, including Context usage.',
    },
    messages: {
      noLifeCycleEventInNestedBlocks:
        'The life-cycle event `{{ eventName }}` listeners should only be defined in the top-level scope of the ArtalkPlugin.',
      noEventInWatchConf: 'Avoid calling `{{ functionName }}` inside the `ctx.watchConf` effect.',
      noInjectInNestedBlocks:
        'The `ctx.inject` method should only be called in the top-level scope of the ArtalkPlugin.',
      noInjectOutsidePlugin:
        'The `ctx.inject` method should only be called inside the ArtalkPlugin arrow function.',
      noCycleDeps: 'Dependency cycle via `ctx.provide` ({{ route }}) in the ArtalkPlugin.',
      onePluginPerFile: 'There is more than one ArtalkPlugin in this file.',
    },
    schema: [],
  },
  defaultOptions: [],
  create(context) {
    const checkerContext: ArtalkPluginCheckerContext = {
      eslint: context,
      depsStore: depsMap,
    }

    let lastPluginFilePath = ''
    let lastPluginName = ''

    return {
      ImportDeclaration(node) {},

      TSTypeAnnotation(node) {
        const typeAnnotation = node.typeAnnotation
        if (typeAnnotation.type !== 'TSTypeReference') return
        if (typeAnnotation.typeName.type !== 'Identifier') return
        const typeName = typeAnnotation.typeName.name
        const identifier = node.parent

        if (isPluginName(typeName)) {
          if (identifier.type !== 'Identifier') return
          const pluginName = identifier.name

          // Get the variable declaration of the ArtalkPlugin
          const varDecl = identifier.parent
          if (
            varDecl.type === 'VariableDeclarator' &&
            varDecl.init &&
            varDecl.init?.type == 'ArrowFunctionExpression'
          ) {
            // console.log('Found ArtalkPlugin:', pluginName)
            checkPluginFunction(checkerContext, varDecl.init)

            // check for multiple ArtalkPlugins in the same file
            if (lastPluginFilePath === context.filename && lastPluginName !== pluginName) {
              context.report({
                node: identifier,
                messageId: 'onePluginPerFile',
              })
            }
            lastPluginFilePath = context.filename
            lastPluginName = pluginName
          }
        }
      },

      VariableDeclaration(fnNode) {},

      Identifier(node) {
        if (node.name === 'inject') {
          checkInjectCallOutsideArtalkPlugin(checkerContext, node)
        }
      },
    }
  },
})

type ArtalkPluginRuleContext = Parameters<(typeof artalkPlugin)['create']>[0]

export type DepsData = Map<string, Set<string>> // DepName -> ProviderNames
export type DepsStore = Map<string, DepsData> // FilePath -> DepsData

export interface ArtalkPluginCheckerContext {
  eslint: ArtalkPluginRuleContext
  depsStore: DepsStore
}

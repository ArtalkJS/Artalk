import { AST_NODE_TYPES, TSESLint, TSESTree } from '@typescript-eslint/utils'
import type { ContextApi } from '../../artalk/src/types/context'
import { createRule } from './helper'
type _ = ContextApi // for IDE jump-to-definition

/** Whether the given string is a ArtalkPlugin name */
function isPluginName(s: string) {
  return s === 'ArtalkPlugin' || /Artalk[A-Z0-9].*Plugin/.test(s)
}

/** The event function names in ContextApi */
const ctxEventFns = ['off', 'on', 'trigger']

/** The life-cycle event names in ContextApi */
const ctxLifeCycleEvents = ['mounted', 'destroyed', 'updated', 'list-fetched']

export const artalkPlugin = createRule({
  name: 'artalk-plugin',
  meta: {
    type: 'problem',
    docs: {
      description:
        'Enforce best practices for ArtalkPlugin arrow functions, including ContextApi usage.',
    },
    messages: {
      noLifeCycleEventInNestedBlocks:
        'The life-cycle event `{{ eventName }}` listeners should only be defined in the top-level scope of the ArtalkPlugin.',
      noEventInWatchConf: 'Avoid calling `{{ functionName }}` inside the `ctx.watchConf` effect.',
    },
    schema: [],
  },
  defaultOptions: [],
  create(context) {
    // Initialize the TypeScript parser services
    const parserServices = context.sourceCode.parserServices
    if (!parserServices || !parserServices.program) {
      console.error('[eslint-plugin-artalk] Missing typescript parser services.')
      return {}
    }
    const checker = parserServices.program.getTypeChecker()

    // -------------------------------------------------------------------
    //  Utility functions
    // -------------------------------------------------------------------
    const getTypeName = (node: TSESTree.Node) => {
      const tsNode = parserServices?.esTreeNodeToTSNodeMap?.get(node)
      const tsType = tsNode ? checker.getTypeAtLocation(tsNode) : null
      const typeName = tsType ? checker.typeToString(tsType) : ''
      return typeName
    }

    const getArrowFunctionType = (node: TSESTree.Node) => {
      if (node.type === 'ArrowFunctionExpression') return getTypeName(node)
      return ''
    }

    const isInsideArtalkPlugin = (node: TSESTree.Node) => {
      let curr: TSESTree.Node | undefined = node
      while (curr) {
        if (isPluginName(getArrowFunctionType(curr))) return true
        curr = curr.parent
      }
      return false
    }

    /**
     * Get the references to ContextApi in the top scope of the given scope
     */
    const getCtxRefNamesInTopScope = (ctxArgName: string, scope: TSESLint.Scope.Scope) => {
      const ctxRefs = new Map<TSESTree.Node, string>()

      const getFullMethodName = (node: TSESTree.Node) => {
        const methodNameArr: string[] = []
        let curr: TSESTree.Node | undefined = node
        while (curr) {
          if (curr.type === 'MemberExpression' && curr.property.type === 'Identifier')
            methodNameArr.push(curr.property.name)
          curr = curr.parent
        }
        return methodNameArr.join('.')
      }

      scope.references.forEach((reference) => {
        const identifier = reference.identifier
        if (identifier.name !== ctxArgName) return

        const methodName = getFullMethodName(identifier.parent)
        if (methodName) ctxRefs.set(identifier.parent, methodName)
      })

      return ctxRefs
    }

    /**
     * Get the references to ContextApi in the nested scopes of the given
     */
    const getCtxRefNamesInNestedScope = (
      ctxArgName: string,
      parentScope: TSESLint.Scope.Scope,
      keepTop = true,
    ) => {
      const ctxRefs = new Map<TSESTree.Node, string>()
      keepTop &&
        getCtxRefNamesInTopScope(ctxArgName, parentScope).forEach((v, k) => ctxRefs.set(k, v))
      parentScope.childScopes.forEach((childScope) => {
        getCtxRefNamesInNestedScope(ctxArgName, childScope).forEach((v, k) => ctxRefs.set(k, v))
      })
      return ctxRefs
    }

    // -------------------------------------------------------------------
    //  Checker functions
    // -------------------------------------------------------------------

    /**
     * Check the set of all function names in ContextApi
     *
     * (which is called in the top-level of ArtalkPlugin arrow-function scope)
     */
    const checkTopLevelCtxRefs = (m: Map<TSESTree.Node, string>) => {
      // console.debug('checkTopLevelCtxFnCalls', m.values())
      // ...
    }

    /**
     * Check the set of all function names in ContextApi
     *
     * (which is called in the nested scopes of ArtalkPlugin arrow-function scope)
     */
    const checkNestedCtxRefs = (m: Map<TSESTree.Node, string>) => {
      // console.debug('checkAllCtxFnCalls', m.values())
      // ...
      // TODO: Event Circular trigger Check
    }

    /**
     * Check the set of all function names in ContextApi
     *
     * (which is called in the nested scopes of ArtalkPlugin arrow-function scope, excluding the top-level)
     */
    const checkNestedCtxRefsNoTop = (m: Map<TSESTree.Node, string>) => {
      m.forEach((methodName, node) => {
        // Disallow life-cycle events in nested blocks
        if (methodName === 'on') {
          // Get the call arguments
          const parent = node.parent
          if (!parent || parent.type !== 'CallExpression') return
          if (parent.arguments.length == 0) return
          const eventNameArg = parent.arguments[0]
          if (eventNameArg.type !== 'Literal') return
          const eventName = eventNameArg.value
          if (typeof eventName !== 'string') return
          if (ctxLifeCycleEvents.includes(eventName)) {
            context.report({
              node: parent,
              messageId: 'noLifeCycleEventInNestedBlocks',
              data: {
                eventName,
              },
            })
          }
        }
      })
    }

    /**
     * Check the set of all function names in ContextApi
     *
     * (which is called in the watchConf effect function scope)
     */
    const checkWatchConfCalls = (m: Map<TSESTree.Node, string>) => {
      const disallowedMethods = [...ctxEventFns]
      m.forEach((methodName, node) => {
        if (disallowedMethods.includes(methodName)) {
          context.report({
            node: node.parent || node,
            messageId: 'noEventInWatchConf',
            data: { functionName: `ctx.${methodName}` },
          })
        }
      })
    }

    /**
     * Whether the ArtalkPlugin is imported
     *
     * (to enable the plugin checker)
     */
    let pluginCheckerEnabled = false

    return {
      ImportDeclaration(node) {
        // Check if contains ArtalkPlugin importing
        node.specifiers.forEach((specifier) => {
          if (specifier.type !== 'ImportSpecifier') return
          if (isPluginName(specifier.imported.name)) {
            pluginCheckerEnabled = true
          }
        })
      },

      VariableDeclaration(fnNode) {
        if (!pluginCheckerEnabled) return

        // Check if the variable declaration is ArtalkPlugin
        fnNode.declarations.forEach((decl) => {
          if (
            isPluginName(getTypeName(decl)) &&
            decl.init &&
            decl.init?.type == 'ArrowFunctionExpression'
          ) {
            // Is ArtalkPlugin arrow-function
            const pluginFn = decl.init

            // Get the first parameter name as the ContextApi reference
            if (pluginFn.params.length === 0) return // No ctx reference
            const ctxArg = pluginFn.params[0]
            if (ctxArg.type !== 'Identifier') return
            const ctxArgName = ctxArg.name

            // Visit the top-level scope of the ArtalkPlugin arrow-function
            const pluginFnScope = context.sourceCode.getScope(pluginFn.body)
            const topLevelCtxRefs = getCtxRefNamesInTopScope(ctxArgName, pluginFnScope)
            checkTopLevelCtxRefs(topLevelCtxRefs)

            // Visit all nested scopes (including the top-level) of the ArtalkPlugin arrow-function
            const nestedCtxRefsIncludeTop = getCtxRefNamesInNestedScope(
              ctxArgName,
              pluginFnScope,
              true,
            )
            checkNestedCtxRefs(nestedCtxRefsIncludeTop)

            // Visit all nested scopes (excluding the top-level) of the ArtalkPlugin arrow-function
            const nestedCtxRefsExcludeTop = getCtxRefNamesInNestedScope(
              ctxArgName,
              pluginFnScope,
              false,
            )
            checkNestedCtxRefsNoTop(nestedCtxRefsExcludeTop)

            // Visit watchConf effect function scope
            const watchConfCalls = new Map<TSESTree.Node, string>()
            topLevelCtxRefs.forEach((v, k) => {
              if (v === 'watchConf') {
                // Get the watchConf call expression
                const watchConfCall = k.parent
                if (!watchConfCall || watchConfCall.type !== AST_NODE_TYPES.CallExpression) return

                // Get the watchConf effect function
                if (watchConfCall.arguments.length < 2) return
                const watchConfEffectFn = watchConfCall.arguments[1]
                if (
                  watchConfEffectFn.type !== 'ArrowFunctionExpression' &&
                  watchConfEffectFn.type !== 'FunctionExpression'
                )
                  return

                // Get the references to ContextApi in the watchConf effect function top scope
                const scope = context.sourceCode.getScope(watchConfEffectFn.body)
                getCtxRefNamesInTopScope(ctxArgName, scope).forEach((v, k) =>
                  watchConfCalls.set(k, v),
                )
              }
            })
            checkWatchConfCalls(watchConfCalls)
          }
        })
      },

      Identifier(node) {},

      CallExpression(node) {},
    }
  },
})

import { AST_NODE_TYPES, TSESLint, TSESTree } from '@typescript-eslint/utils'
import type { ArtalkPluginCheckerContext, DepsData, DepsStore } from './artalk-plugin'
import { tarjan } from './scc'

/** The event function names in Context */
const ctxEventFns = ['off', 'on', 'trigger']

/** The life-cycle event names in Context */
const ctxLifeCycleEvents = ['mounted', 'destroyed', 'updated', 'list-fetched']

/** Whether the given string is a ArtalkPlugin name */
export function isPluginName(s: string) {
  return s === 'ArtalkPlugin' || /Artalk[A-Z0-9].*Plugin/.test(s)
}

/**
 * Get the references to Context in the top scope of the given scope
 */
const getCtxRefNamesInTopScope = (ctxArgName: string, scope: TSESLint.Scope.Scope) => {
  const ctxRefs = new Map<TSESTree.Node, string>()

  const getFullMethodName = (node: TSESTree.Node) => {
    const methodNameArr: string[] = []
    let curr: TSESTree.Node | undefined = node
    const visited = new Set<TSESTree.Node>()
    while (curr) {
      if (visited.has(curr)) break
      visited.add(curr)
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
 * Get the references to Context in the nested scopes of the given
 */
const getCtxRefNamesInNestedScope = (
  ctxArgName: string,
  parentScope: TSESLint.Scope.Scope,
  keepTop = true,
) => {
  const ctxRefs = new Map<TSESTree.Node, string>()
  keepTop && getCtxRefNamesInTopScope(ctxArgName, parentScope).forEach((v, k) => ctxRefs.set(k, v))
  parentScope.childScopes.forEach((childScope) => {
    getCtxRefNamesInNestedScope(ctxArgName, childScope).forEach((v, k) => ctxRefs.set(k, v))
  })
  return ctxRefs
}

/**
 * Check the set of all function names in Context
 *
 * (which is called in the top-level of ArtalkPlugin arrow-function scope)
 */
const checkTopLevelCtxRefs = (ctx: ArtalkPluginCheckerContext, m: Map<TSESTree.Node, string>) => {
  // console.debug('checkTopLevelCtxFnCalls', m.values())
  // ...
}

const getDepsMap = (node: TSESTree.CallExpression): DepsData | void => {
  if (node.arguments.length < 3) return

  // Get the dependencies
  let depsArg = node.arguments[2]
  if (depsArg.type === 'TSAsExpression') depsArg = depsArg.expression
  if (depsArg.type !== 'ArrayExpression') return
  const deps = depsArg.elements
    .map((e) => (e && e.type === 'Literal' && typeof e.value === 'string' ? e.value : ''))
    .filter((e) => !!e)
  if (deps.length === 0) return

  // Get the provider name
  const providerNameArg = node.arguments[0]
  if (providerNameArg.type !== 'Literal') return
  const providerName = providerNameArg.value
  if (!providerName || typeof providerName !== 'string') return

  // Record the dependency data for the file
  const depsMap: DepsData = new Map()
  deps.forEach((depName) => {
    if (!depsMap.has(depName)) depsMap.set(depName, new Set())
    depsMap.get(depName)!.add(providerName)
  })

  return depsMap
}

function removePath(graph: DepsData, path: string[]): DepsData {
  if (path.length < 2) {
    throw new Error('Path must have at least two nodes')
  }

  // Clone the data
  const updatedGraph = new Map(graph)
  updatedGraph.forEach((value, key) => {
    updatedGraph.set(key, new Set(value))
  })

  for (let i = 0; i < path.length - 1; i++) {
    const toNode = path[i]
    const fromNode = path[i + 1]

    // Remove the edge from fromNode to toNode
    if (updatedGraph.has(fromNode)) {
      const edges = updatedGraph.get(fromNode)
      if (edges) {
        // Remove the edge to toNode from the adjacency set of fromNode
        edges.delete(toNode)
        if (edges.size === 0) {
          // If the node no longer has any edges, remove the node
          updatedGraph.delete(fromNode)
        }
      }
    }
  }

  return updatedGraph
}

const checkCircularDependency = (
  depsStore: DepsStore,
  ctx: ArtalkPluginCheckerContext,
  _node: TSESTree.CallExpression,
) => {
  const args = _node.arguments
  const node = args.length >= 3 ? args[2] : _node

  // Merge all files' dependency data
  const depsGraph: DepsData = new Map()
  depsStore.forEach((depsData, filename) => {
    depsData.forEach((providers, depName) => {
      if (!depsGraph.has(depName)) depsGraph.set(depName, new Set())
      providers.forEach((provider) => depsGraph.get(depName)!.add(provider))
    })
  })

  // console.log('\n' + ctx.eslint.filename)
  // console.log('depsGraph', depsGraph)
  // console.log('tarjan', tarjan(depsGraph))

  // Basic check (self-reference, a->a)
  for (const [depName, providers] of depsGraph) {
    if (providers.has(depName)) {
      ctx.eslint.report({
        node,
        messageId: 'noCycleDeps',
        data: { route: `${depName}->${depName}` },
      })
      return
    }
  }

  // SCC (Strongly Connected Components) algorithm
  tarjan(depsGraph).forEach((scc) => {
    if (scc.size <= 1) return
    const route = [...scc, scc.values().next().value].slice(1).join('->')

    ctx.eslint.report({
      node,
      messageId: 'noCycleDeps',
      data: { route },
    })

    // Cleanup
    depsStore.forEach((depsData, filename) => {
      depsStore.set(filename, removePath(depsData, [...scc]))
    })
  })
}

/**
 * Check the set of all function names in Context
 *
 * (which is called in the nested scopes of ArtalkPlugin arrow-function scope)
 */
const checkNestedCtxRefs = (ctx: ArtalkPluginCheckerContext, m: Map<TSESTree.Node, string>) => {
  // console.debug('checkAllCtxFnCalls', m.values())
  // ...
  // TODO: Event Circular trigger Check

  const depsMap: DepsData = new Map()

  m.forEach((methodName, node) => {
    // Check dependency providers via `ctx.provide`
    if (methodName === 'provide') {
      const callExpr = node.parent
      if (!callExpr || callExpr.type !== 'CallExpression') return

      // Record
      const dm = getDepsMap(callExpr)
      if (!dm) return

      // Merge
      dm.forEach((providers, depName) => {
        if (!depsMap.has(depName)) depsMap.set(depName, new Set())
        providers.forEach((provider) => depsMap.get(depName)!.add(provider))
      })

      // Check
      const depsStoreShallowCopy = new Map(ctx.depsStore)
      depsStoreShallowCopy.set(ctx.eslint.filename, depsMap)
      checkCircularDependency(depsStoreShallowCopy, ctx, callExpr)
    }
  })

  // Overwrite the historical dependency data
  ctx.depsStore.set(ctx.eslint.filename, depsMap)
}

/**
 * Check the set of all function names in Context
 *
 * (which is called in the nested scopes of ArtalkPlugin arrow-function scope, excluding the top-level)
 */
const checkNestedCtxRefsNoTop = (
  ctx: ArtalkPluginCheckerContext,
  m: Map<TSESTree.Node, string>,
) => {
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
        ctx.eslint.report({
          node: parent,
          messageId: 'noLifeCycleEventInNestedBlocks',
          data: {
            eventName,
          },
        })
      }
    }

    // Disallow inject in nested blocks
    if (methodName === 'inject') {
      ctx.eslint.report({
        node: node.parent || node,
        messageId: 'noInjectInNestedBlocks',
      })
    }
  })
}

/**
 * Check the set of all function names in Context
 *
 * (which is called in the watchConf effect function scope)
 */
const checkWatchConfCalls = (ctx: ArtalkPluginCheckerContext, m: Map<TSESTree.Node, string>) => {
  const disallowedMethods = [...ctxEventFns]
  m.forEach((methodName, node) => {
    if (disallowedMethods.includes(methodName)) {
      ctx.eslint.report({
        node: node.parent || node,
        messageId: 'noEventInWatchConf',
        data: { functionName: `ctx.${methodName}` },
      })
    }
  })
}

/**
 * Check the ArtalkPlugin variable declaration
 */
export const checkPluginFunction = (
  ctx: ArtalkPluginCheckerContext,
  pluginFn: TSESTree.ArrowFunctionExpression,
) => {
  // Get the first parameter name as the Context reference
  if (pluginFn.params.length === 0) return // No ctx reference
  const ctxArg = pluginFn.params[0]
  if (ctxArg.type !== 'Identifier') return
  const ctxArgName = ctxArg.name

  // Visit the top-level scope of the ArtalkPlugin arrow-function
  const pluginFnScope = ctx.eslint.sourceCode.getScope(pluginFn.body)
  const topLevelCtxRefs = getCtxRefNamesInTopScope(ctxArgName, pluginFnScope)
  checkTopLevelCtxRefs(ctx, topLevelCtxRefs)

  // Visit all nested scopes (including the top-level) of the ArtalkPlugin arrow-function
  const nestedCtxRefsIncludeTop = getCtxRefNamesInNestedScope(ctxArgName, pluginFnScope, true)
  checkNestedCtxRefs(ctx, nestedCtxRefsIncludeTop)

  // Visit all nested scopes (excluding the top-level) of the ArtalkPlugin arrow-function
  const nestedCtxRefsExcludeTop = getCtxRefNamesInNestedScope(ctxArgName, pluginFnScope, false)
  checkNestedCtxRefsNoTop(ctx, nestedCtxRefsExcludeTop)

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

      // Get the references to Context in the watchConf effect function top scope
      const scope = ctx.eslint.sourceCode.getScope(watchConfEffectFn.body)
      getCtxRefNamesInTopScope(ctxArgName, scope).forEach((v, k) => watchConfCalls.set(k, v))
    }
  })
  checkWatchConfCalls(ctx, watchConfCalls)
}

export const checkInjectCallOutsideArtalkPlugin = (
  ctx: ArtalkPluginCheckerContext,
  node: TSESTree.Identifier,
) => {
  if (node.name !== 'inject') return
  const parent = node.parent
  if (parent.type !== 'MemberExpression') return
  if (parent.object.type !== 'Identifier') return
  if (!['ctx', 'context'].includes(parent.object.name)) return

  // traverse up to find the ArtalkPlugin arrow-function
  let curr: TSESTree.Node | undefined = parent
  let pluginFn: TSESTree.ArrowFunctionExpression | undefined
  const visited = new Set<TSESTree.Node>()
  while (curr) {
    if (visited.has(curr)) break
    visited.add(curr)
    if (curr.type === 'ArrowFunctionExpression') {
      pluginFn = curr
    }
    curr = curr.parent
  }

  const fail = () => {
    ctx.eslint.report({
      node,
      messageId: 'noInjectOutsidePlugin',
    })
  }

  // check if the ArtalkPlugin arrow-function is found
  if (!pluginFn) return fail()
  const varDecl = pluginFn.parent
  if (varDecl.type !== 'VariableDeclarator') return fail()
  const typeRef = varDecl.id.typeAnnotation?.typeAnnotation
  if (!typeRef || typeRef.type !== 'TSTypeReference') return fail()
  const typeNameId = typeRef.typeName
  if (typeNameId.type !== 'Identifier') return fail()
  if (typeNameId.name !== 'ArtalkPlugin') return fail()
}

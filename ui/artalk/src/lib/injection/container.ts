type PrimitiveKey = string | number | symbol
type Services = { [key: PrimitiveKey]: any }

export type Constructor<T, S = Services, D extends readonly (keyof S)[] = any> = (
  ...args: { [K in keyof D]: D[K] extends keyof S ? S[D[K]] : never }
) => T

export type Lifecycle = 'transient' | 'singleton'

export interface Provider<T, S = Services> {
  /** Implementation constructor */
  impl: Constructor<T>

  /** Dependency constructors */
  deps: readonly (keyof S)[]

  /** Lifecycle */
  lifecycle: Lifecycle
}

export interface ProvideFuncOptions {
  /** Lifecycle (default: 'singleton') */
  lifecycle?: Lifecycle
}

export interface DependencyContainer<S = Services> {
  provide<K extends keyof S, T extends S[K] = any, D extends readonly (keyof S)[] = any>(
    key: K,
    impl: Constructor<T, S, D>,
    deps?: D,
    opts?: ProvideFuncOptions,
  ): void

  inject<T = undefined, K extends keyof S = any>(key: K): T extends undefined ? S[K] : T
}

export function createInjectionContainer<S = Services>(): DependencyContainer<S> {
  const providers = new Map<PrimitiveKey, Provider<any, S>>()
  const initializedDeps = new Map<PrimitiveKey, any>()

  const provide: DependencyContainer<S>['provide'] = (key, impl, deps, opts = {}) => {
    providers.set(key, { impl, deps: deps || [], lifecycle: opts.lifecycle || 'singleton' })
  }

  const inject: DependencyContainer<S>['inject'] = (key) => {
    const provider = providers.get(key)
    if (!provider) {
      throw new Error(`No provide for ${String(key)}`)
    }

    if (provider.lifecycle === 'singleton' && initializedDeps.has(key)) {
      return initializedDeps.get(key)
    }

    const { impl, deps } = provider
    const params = deps.map((d) => inject(d))
    const resolved = impl(...params)

    initializedDeps.set(key, resolved)

    return resolved
  }

  return { provide, inject }
}

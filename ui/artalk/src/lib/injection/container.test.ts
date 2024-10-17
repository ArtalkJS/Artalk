import { describe, expect, it } from 'vitest'
import { createInjectionContainer } from './container'

describe('InjectionContainer', () => {
  it('should inject a singleton service', () => {
    const container = createInjectionContainer()

    // Mock Service
    const myService = {
      getValue: () => 'test value',
    }

    // Provide service as singleton
    container.provide('myService', () => myService)

    // Inject service and validate
    const injectedService = container.inject('myService')
    expect(injectedService).toBe(myService)
  })

  it('should inject a transient service', () => {
    const container = createInjectionContainer()

    // Provide service as transient
    container.provide('myService', () => ({ id: Math.random() }), [], { lifecycle: 'transient' })

    // Inject service twice and expect different instances
    const service1 = container.inject('myService')
    const service2 = container.inject('myService')
    expect(service1).not.toBe(service2)
  })

  it('should throw error when service is not provided', () => {
    const container = createInjectionContainer()

    // Try injecting a non-existent service
    expect(() => container.inject('nonExistentService')).toThrowError(
      'No provide for nonExistentService',
    )
  })

  it('should inject services with dependencies', () => {
    type Dep1 = { getName: () => string }
    type Dep2 = { getName: () => string }
    type Deps = { dep1: Dep1; dep2: Dep2; myService: MyService }

    const container = createInjectionContainer<Deps>()

    // Mock dependencies
    const dep1: Dep1 = { getName: () => 'Dependency 1' }
    const dep2: Dep2 = { getName: () => 'Dependency 2' }

    // Provide dependencies
    container.provide('dep1', () => dep1)
    container.provide('dep2', () => dep2)

    // Mock service depending on dep1 and dep2
    class MyService {
      constructor(
        public dep1: Dep1,
        public dep2: Dep2,
      ) {}

      getDeps() {
        return [this.dep1.getName(), this.dep2.getName()]
      }
    }

    container.provide('myService', (dep1, dep2) => new MyService(dep1, dep2), [
      'dep1',
      'dep2',
    ] as const)

    // Inject service and validate dependencies
    const myService = container.inject('myService')
    expect(myService.getDeps()).toEqual(['Dependency 1', 'Dependency 2'])
  })

  it('should inject singleton with dependencies only once', () => {
    const container = createInjectionContainer()

    // Mock dependencies
    const dep1 = { id: Math.random() }
    const dep2 = { id: Math.random() }

    // Provide dependencies as singletons
    container.provide('dep1', () => dep1)
    container.provide('dep2', () => dep2)

    // Provide service depending on dep1 and dep2
    container.provide('myService', (dep1, dep2) => ({ dep1, dep2 }), ['dep1', 'dep2'])

    // Inject service and dependencies
    const myService1 = container.inject('myService')
    const myService2 = container.inject('myService')

    // Expect same instance for singleton
    expect(myService1).toBe(myService2)
    expect(myService1.dep1).toBe(dep1)
    expect(myService1.dep2).toBe(dep2)
  })
})

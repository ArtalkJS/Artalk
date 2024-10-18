import { artalkPlugin } from './artalk-plugin'
import { setupTest } from './test-helper'

const { ruleTester } = setupTest()

const invalid = [
  {
    name: 'should not allow life-cycle event functions in nested blocks',
    code: `
      import type { ArtalkPlugin } from 'artalk'

      export const TestPlugin: ArtalkPlugin = (ctx) => {
        const foo = () => {
          const bar = () => {
            ctx.on('updated', () => {})
          }
        }
      }
    `,
    errorId: 'noLifeCycleEventInNestedBlocks',
  },
  {
    name: "should not allow 'off' call inside watchConf effect",
    code: `
      import type { ArtalkPlugin } from 'artalk'

      export const TestPlugin: ArtalkPlugin = (ctx) => {
        ctx.watchConf(['pageVote'], (conf) => {
          ctx.off('updated', () => {})
        })
      }
    `,
    errorId: 'noEventInWatchConf',
  },
  {
    name: "should not allow 'inject' call in nested blocks",
    code: `
      import type { ArtalkPlugin } from 'artalk'

      export const TestPlugin: ArtalkPlugin = (ctx) => {
        const fn = () => {
          const foo = ctx.inject('foo')
        }
      }
    `,
    errorId: 'noInjectInNestedBlocks',
  },
  {
    name: "should not allow 'inject' call outside ArtalkPlugin",
    code: `
      function fn(ctx) {
        const foo = ctx.inject('foo')
      }
    `,
    errorId: 'noInjectOutsidePlugin',
  },
  {
    name: 'should not allow circular dependency providing',
    code: `
      import type { ArtalkPlugin } from 'artalk'

      export const TestPlugin: ArtalkPlugin = (ctx) => {
        ctx.provide('foo', (foo) => {}, ['foo'])
      }
    `,
    errorId: 'noCycleDeps',
  },
  {
    name: 'should not allow multiple ArtalkPlugin in a single file',
    code: `
      import type { ArtalkPlugin } from 'artalk'

      export const TestPlugin: ArtalkPlugin = (ctx) => {}

      export const TestPlugin2: ArtalkPlugin = (ctx) => {}
    `,
    errorId: 'onePluginPerFile',
  },
]

for (const { name, code, errorId } of invalid) {
  ruleTester.run(name, artalkPlugin as any, {
    valid: [],
    invalid: [{ code, errors: [{ messageId: errorId }] }],
  })
}

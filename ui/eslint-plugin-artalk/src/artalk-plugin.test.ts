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
]

for (const { name, code, errorId } of invalid) {
  ruleTester.run(name, artalkPlugin as any, {
    valid: [],
    invalid: [{ code, errors: [{ messageId: errorId }] }],
  })
}

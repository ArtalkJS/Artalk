import * as vitest from 'vitest'
import { RuleTester } from '@typescript-eslint/rule-tester'

let ruleTester: RuleTester | undefined

export function setupTest() {
  if (ruleTester) return { ruleTester }

  RuleTester.afterAll = vitest.afterAll
  RuleTester.it = vitest.it
  RuleTester.itOnly = vitest.it.only
  RuleTester.describe = vitest.describe

  ruleTester = new RuleTester({
    languageOptions: {
      parserOptions: {
        projectService: {
          allowDefaultProject: ['*.ts*'],
        },
      },
    },
  })

  return { ruleTester }
}

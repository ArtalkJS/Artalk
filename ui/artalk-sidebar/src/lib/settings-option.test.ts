import { describe, expect, it } from 'vitest'
import YAML from 'yaml'
import {
  formatOptionValue,
  KEYWORD_FILE_SEPARATOR_PATH,
  patchOptionValue,
  type OptionNode,
} from './settings-option'

const separatorNode: OptionNode = {
  name: 'file_sep',
  path: KEYWORD_FILE_SEPARATOR_PATH,
  level: 3,
  type: 'string',
  title: 'FileSep',
}

describe('keyword file separator option', () => {
  it.each([
    ['line feed', '\n', '\\n'],
    ['CRLF', '\r\n', '\\r\\n'],
    ['tab', '\t', '\\t'],
    ['pipe', '|', '|'],
    ['literal escape sequence', '\\n', '\\\\n'],
  ])('round-trips %s', (_, storedValue, inputValue) => {
    expect(formatOptionValue(storedValue, separatorNode)).toBe(inputValue)
    expect(patchOptionValue(inputValue, separatorNode)).toBe(storedValue)
  })

  it('preserves custom separators and unknown escape sequences', () => {
    expect(patchOptionValue(' :: ', separatorNode)).toBe(' :: ')
    expect(patchOptionValue('\\x', separatorNode)).toBe('\\x')
  })

  it('preserves a line feed through the YAML read and save flow', () => {
    const document = YAML.parseDocument('moderator:\n  keywords:\n    file_sep: "\\n"\n')
    const path = KEYWORD_FILE_SEPARATOR_PATH.split('.')
    const storedValue = document.getIn(path)
    const inputValue = formatOptionValue(storedValue, separatorNode)

    document.setIn(path, patchOptionValue(inputValue, separatorNode))

    expect(document.getIn(path)).toBe('\n')
    expect(YAML.parse(document.toString()).moderator.keywords.file_sep).toBe('\n')
  })
})

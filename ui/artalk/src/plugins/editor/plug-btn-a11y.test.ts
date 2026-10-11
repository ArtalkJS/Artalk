import { describe, expect, it } from 'vitest'
import type PlugKit from './_kit'
import type EditorPlugin from './_plug'
import Emoticons from './emoticons'
import Preview from './preview'
import Upload from './upload'

function createKit() {
  const handlers: Record<string, (() => void)[]> = {}
  const events = {
    on: (name: string, fn: () => void) => {
      ;(handlers[name] ||= []).push(fn)
    },
  }
  const kit = {
    useEvents: () => events,
    useMounted: (fn: () => void) => events.on('mounted', fn),
    useUnmounted: (fn: () => void) => events.on('unmounted', fn),
    useConf: () => ({ imgUpload: true }),
    useUI: () => ({ $textarea: document.createElement('textarea') }),
  }
  const mount = () => handlers['mounted']?.forEach((fn) => fn())
  return { kit: kit as unknown as PlugKit, mount }
}

describe('editor plug buttons', () => {
  it.each([
    ['Emoticons', Emoticons],
    ['Preview', Preview],
    ['Upload', Upload],
  ] as [string, new (kit: PlugKit) => EditorPlugin][])(
    '%s icon has an image role for its aria-label',
    (_, Plug) => {
      const { kit, mount } = createKit()
      const plug = new Plug(kit)
      mount()

      const $icon = plug.$btn!.querySelector('[aria-label]')
      expect($icon).not.toBeNull()
      expect($icon!.getAttribute('role')).toBe('img')
    },
  )
})

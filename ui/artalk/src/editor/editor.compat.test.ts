import { describe, expect, it, vi } from 'vitest'
import { Editor } from './editor'
import type { CommentData } from '@/types'
import type { PluginManager } from '@/plugins/editor-kit'

describe('deprecated Editor aliases', () => {
  it('forwards to the current Editor API', () => {
    const editor = Object.create(Editor.prototype) as Editor
    const $el = document.createElement('div')
    const plugins = {} as PluginManager
    const comment = {} as CommentData

    vi.spyOn(editor, 'getEl').mockReturnValue($el)
    vi.spyOn(editor, 'getPlugins').mockReturnValue(plugins)
    const setReplyComment = vi.spyOn(editor, 'setReplyComment').mockImplementation(() => {})

    expect(editor.$el).toBe($el)
    expect(editor.getPlugs()).toBe(plugins)

    editor.setReply(comment, $el, true)
    expect(setReplyComment).toHaveBeenCalledWith(comment, $el)
  })
})

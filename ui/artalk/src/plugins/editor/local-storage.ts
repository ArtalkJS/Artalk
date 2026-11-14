import type PlugKit from './_kit'
import EditorPlugin from './_plug'
import $t from '@/i18n'
import type { Editor } from '@/types'
import { storageKey } from '@/lib/storage'

const legacyKey = 'ArtalkContent'
const loadedScopes = new WeakMap<Editor, string>()

export default class LocalStorage extends EditorPlugin {
  private key: string

  constructor(kit: PlugKit) {
    super(kit)
    this.key = this.getKey()
    const editor = kit.useEditor()
    const onContentUpdated = () => this.save()

    kit.useMounted(() => {
      const previous = loadedScopes.get(editor)
      // Replacing plugins must not reload or overwrite the current draft.
      if (previous !== this.key) {
        loadedScopes.set(editor, this.key)
        let content = ''
        try {
          const stored = window.localStorage.getItem(this.key)
          content = stored ?? window.localStorage.getItem(legacyKey) ?? ''
          // Claim the old unscoped draft once, without copying it to every widget.
          if (stored === null && content) {
            window.localStorage.setItem(this.key, content)
            window.localStorage.removeItem(legacyKey)
          }
        } catch {
          // The current editor still works without persistent storage.
        }
        if (previous || !editor.getContentRaw()) {
          editor.setContent(content)
          if (content.trim()) editor.showNotify($t('restoredMsg'), 'i')
        }
      }
      kit.useEvents().on('content-updated', onContentUpdated)
    })
    kit.useUnmounted(() => kit.useEvents().off('content-updated', onContentUpdated))
  }

  private getKey() {
    const { server, site, pageKey } = this.kit.useConf()
    return `${storageKey(legacyKey, server)}:${JSON.stringify([site, pageKey])}`
  }

  public save() {
    // Destination changes may reset edit mode before these plugins are replaced.
    if (this.key !== this.getKey()) return
    try {
      window.localStorage.setItem(this.key, this.kit.useEditor().getContentRaw().trim())
    } catch {
      // Keep the draft in the editor when storage is unavailable.
    }
  }
}

import Editor from '../editor'
import EditorPlug from '../editor-plug'

const LocalStorageKey = 'ArtalkContent'

export default class LocalStoragePlug extends EditorPlug {
  constructor(editor: Editor) {
    super(editor)

    this.kit.useMounted(() => {
      // load editor content from localStorage when init
      const localContent = window.localStorage.getItem(LocalStorageKey) || ''
      if (localContent.trim() !== '') {
        editor.showNotify(editor.$t('restoredMsg'), 'i')
        editor.setContent(localContent)
      }
    })

    this.kit.useUnmounted(() => {
    })

    this.kit.useContentUpdated(() => {
      this.save()
    })
  }

  // Save editor content to localStorage
  public save() {
    window.localStorage.setItem(LocalStorageKey, this.editor.getContentRaw().trim())
  }
}

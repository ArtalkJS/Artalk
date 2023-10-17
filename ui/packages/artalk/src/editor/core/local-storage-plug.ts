import $t from '@/i18n'
import EditorPlug from '../editor-plug'
import PlugKit from '../plug-kit'

const LocalStorageKey = 'ArtalkContent'

export default class LocalStoragePlug extends EditorPlug {
  constructor(kit: PlugKit) {
    super(kit)

    this.kit.useMounted(() => {
      // load editor content from localStorage when init
      const localContent = window.localStorage.getItem(LocalStorageKey) || ''
      if (localContent.trim() !== '') {
        this.kit.useEditor().showNotify($t('restoredMsg'), 'i')
        this.kit.useEditor().setContent(localContent)
      }
    })

    this.kit.useUnmounted(() => {
    })

    this.kit.useEvents().on('content-updated', () => {
      this.save()
    })
  }

  // Save editor content to localStorage
  public save() {
    window.localStorage.setItem(LocalStorageKey, this.kit.useEditor().getContentRaw().trim())
  }
}

import type PlugKit from './_kit'
import EditorPlugin from './_plug'
import $t from '@/i18n'

const LocalStorageKey = 'ArtalkContent'

export default class LocalStorage extends EditorPlugin {
  constructor(kit: PlugKit) {
    super(kit)

    const onContentUpdated = () => {
      this.save()
    }

    this.kit.useMounted(() => {
      // load editor content from localStorage when init
      const localContent = window.localStorage.getItem(LocalStorageKey) || ''
      if (localContent.trim() !== '') {
        this.kit.useEditor().showNotify($t('restoredMsg'), 'i')
        this.kit.useEditor().setContent(localContent)
      }

      // bind event
      this.kit.useEvents().on('content-updated', onContentUpdated)
    })

    this.kit.useUnmounted(() => {
      this.kit.useEvents().off('content-updated', onContentUpdated)
    })
  }

  // Save editor content to localStorage
  public save() {
    window.localStorage.setItem(LocalStorageKey, this.kit.useEditor().getContentRaw().trim())
  }
}

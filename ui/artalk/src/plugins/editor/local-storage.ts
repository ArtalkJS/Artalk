import $t from '@/i18n'
import type PlugKit from './_kit'
import EditorPlug from './_plug'

const LocalStorageKey = 'ArtalkContent'

export default class LocalStorage extends EditorPlug {
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

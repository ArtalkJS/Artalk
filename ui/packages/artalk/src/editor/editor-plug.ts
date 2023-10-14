import Editor from '@/editor/editor'
import PlugKit from './plug-kit'

/**
 * Editor 插件
 *
 * (使用 Interface x Abstract 合并声明：https://www.typescriptlang.org/docs/handbook/declaration-merging.html#merging-interfaces)
 */
interface EditorPlug {
  $btn?: HTMLElement
  $panel?: HTMLElement
  onMounted?(): void
  onUnmounted?(): void
  onPanelShow?(): void
  onPanelHide?(): void
  onHeaderInput?(key: string, $input: HTMLElement): void
  onContentUpdated?(content: string): void
  contentTransformer?(rawContent: string): string
}

class EditorPlug {
  protected kit: PlugKit

  constructor(protected editor: Editor) {
    this.kit = new PlugKit(this)
  }
}

export default EditorPlug

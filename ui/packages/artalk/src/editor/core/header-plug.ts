import User from '@/lib/user'
import $t from '@/i18n'
import Editor from '../editor'
import EditorPlug from '../editor-plug'
import EditPlug from './edit-plug'
import PlugKit from '../plug-kit'

export default class HeaderPlug extends EditorPlug {
  private get $inputs() {
    return this.kit.useEditor().getHeaderInputEls()
  }

  constructor(kit: PlugKit) {
    super(kit)

    const inputEventFns: {[name: string]: () => void} = {}

    // the input event
    const onInput = ($input: HTMLInputElement, key: string) => () => {
      if (this.kit.useDeps(EditPlug)?.getIsEditMode()) return // 评论编辑模式，不修改个人信息

      User.update({ [key]: $input.value.trim() })

      // trigger header input event
      this.kit.useEvents().trigger('header-input', { field: key, $input })
    }

    this.kit.useMounted(() => {
      // set placeholder and sync header input value
      Object.entries(this.$inputs).forEach(([key, $input]) => {
        $input.placeholder = `${$t(key as any)}`
        $input.value = User.data[key] || ''
      })

      // bind the event
      Object.entries(this.$inputs).forEach(([key, $input]) => {
        $input.addEventListener('input', inputEventFns[key] = onInput($input, key))
      })
    })

    this.kit.useUnmounted(() => {
      // unmount the event
      Object.entries(this.$inputs).forEach(([key, $input]) => {
        $input.removeEventListener('input', inputEventFns[key])
      })
    })
  }
}

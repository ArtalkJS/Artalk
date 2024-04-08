import EditorPlug from './_plug'
import type PlugKit from './_kit'

export default class HeaderEvent extends EditorPlug {
  private get $inputs() {
    return this.kit.useEditor().getHeaderInputEls()
  }

  constructor(kit: PlugKit) {
    super(kit)

    const inputEventFns: { [name: string]: () => void } = {}
    const changeEventFns: { [name: string]: () => void } = {}

    const trigger =
      (evt: 'header-input' | 'header-change', $input: HTMLInputElement, field: string) => () => {
        this.kit.useEvents().trigger(evt, { field, $input })
      }

    this.kit.useMounted(() => {
      // batch bind the events
      Object.entries(this.$inputs).forEach(([key, $input]) => {
        $input.addEventListener(
          'input',
          (inputEventFns[key] = trigger('header-input', $input, key)),
        )
        $input.addEventListener(
          'change',
          (changeEventFns[key] = trigger('header-change', $input, key)),
        )
      })
    })

    this.kit.useUnmounted(() => {
      // unmount the event
      Object.entries(this.$inputs).forEach(([key, $input]) => {
        $input.removeEventListener('input', inputEventFns[key])
        $input.removeEventListener('change', changeEventFns[key])
      })
    })
  }
}

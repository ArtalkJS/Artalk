import EditorPlugin from './_plug'
import type PlugKit from './_kit'

export default class SamplePlug extends EditorPlugin {
  constructor(kit: PlugKit) {
    super(kit)
  }
}

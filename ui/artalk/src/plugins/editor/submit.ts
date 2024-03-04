import type { CommentData } from '@/types'
import $t from '@/i18n'
import EditorPlug from './_plug'
import type PlugKit from './_kit'
import SubmitAddPreset from './submit-add'

interface CustomSubmit {
  activeCond: () => void
  pre?: () => void
  req?: () => Promise<CommentData>
  /** The `post` method is called after the comment is successfully submitted */
  post?: (nComment: CommentData) => void
}

export default class Submit extends EditorPlug {
  private customs: CustomSubmit[] = []
  private defaultPreset: SubmitAddPreset

  constructor(kit: PlugKit) {
    super(kit)

    this.defaultPreset = new SubmitAddPreset(this.kit)

    const onEditorSubmit = () => this.do()

    this.kit.useMounted(() => {
      // invoke `do()` when event `editor-submit` is triggered
      this.kit.useGlobalCtx().on('editor-submit', onEditorSubmit)
    })
    this.kit.useUnmounted(() => {
      this.kit.useGlobalCtx().off('editor-submit', onEditorSubmit)
    })
  }

  registerCustom(c: CustomSubmit) {
    this.customs.push(c)
  }

  private async do() {
    if (this.kit.useEditor().getContentFinal().trim() === '') {
      this.kit.useEditor().focus()
      return
    }

    const custom = this.customs.find(o => o.activeCond())

    this.kit.useEditor().showLoading()

    try {
      // pre submit
      if (custom?.pre) custom.pre()

      let nComment: CommentData

      // submit request
      if (custom?.req) nComment = await custom.req()
      else nComment = await this.defaultPreset.reqAdd()

      // post submit
      if (custom?.post) custom.post(nComment)
      else this.defaultPreset.postSubmitAdd(nComment)
    } catch (err: any) {
      // submit error
      console.error(err)
      this.kit.useEditor().showNotify(`${$t('commentFail')}: ${err.message || String(err)}`, 'e')
      return
    } finally {
      this.kit.useEditor().hideLoading()
    }

    this.kit.useEditor().reset() // 复原编辑器
    this.kit.useGlobalCtx().trigger('editor-submitted')
  }
}

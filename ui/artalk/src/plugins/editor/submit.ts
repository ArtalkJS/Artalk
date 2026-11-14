import EditorPlugin from './_plug'
import type PlugKit from './_kit'
import SubmitAddPreset from './submit-add'
import $t from '@/i18n'
import type { CommentData } from '@/types'

interface CustomSubmit {
  activeCond: () => void
  pre?: () => void
  req?: () => Promise<CommentData>
  /** The `post` method is called after the comment is successfully submitted */
  post?: (nComment: CommentData) => void
}

export default class Submit extends EditorPlugin {
  private customs: CustomSubmit[] = []
  private defaultPreset: SubmitAddPreset
  private submitting = false
  private mounted = false

  constructor(kit: PlugKit) {
    super(kit)

    this.defaultPreset = new SubmitAddPreset(this.kit)

    const onEditorSubmit = () => this.do()

    this.kit.useMounted(() => {
      this.mounted = true
      // invoke `do()` when event `editor-submit` is triggered
      this.kit.useEvents().on('editor-submit', onEditorSubmit)
    })
    this.kit.useUnmounted(() => {
      this.mounted = false
      if (this.submitting) this.kit.useEditor().hideLoading()
      this.kit.useEvents().off('editor-submit', onEditorSubmit)
    })
  }

  registerCustom(c: CustomSubmit) {
    this.customs.push(c)
  }

  private async do() {
    if (this.submitting || !this.mounted) return
    if (this.kit.useEditor().getContentFinal().trim() === '') {
      this.kit.useEditor().focus()
      return
    }

    const custom = this.customs.find((o) => o.activeCond())
    const scope = this.getScope()
    const active = () => this.mounted && scope === this.getScope()

    this.submitting = true
    this.kit.useEditor().showLoading()

    try {
      // pre submit
      if (custom?.pre) custom.pre()

      let nComment: CommentData

      // submit request
      if (custom?.req) nComment = await custom.req()
      else nComment = await this.defaultPreset.reqAdd()

      if (!active()) return
      // post submit
      if (custom?.post) custom.post(nComment)
      else this.defaultPreset.postSubmitAdd(nComment)
    } catch (err: any) {
      if (!active()) return
      // submit error
      console.error(err)
      this.kit.useEditor().showNotify(`${$t('commentFail')}: ${err.message || String(err)}`, 'e')
      return
    } finally {
      this.submitting = false
      if (this.mounted) this.kit.useEditor().hideLoading()
    }

    this.kit.useEditor().reset() // 复原编辑器
    this.kit.useEvents().trigger('editor-submitted')
  }

  private getScope() {
    const { server, site, pageKey } = this.kit.useConf()
    // A login challenge can refresh the token while retrying this submission.
    // The viewer's name and email, rather than that token, identify the draft.
    const { name, email } = this.kit.useUser().getData()
    return JSON.stringify([server, site, pageKey, name, email])
  }
}

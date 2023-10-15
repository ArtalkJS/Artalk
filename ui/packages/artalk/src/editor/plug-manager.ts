import type ArtalkConfig from '~/types/artalk-config'
import MoverPlug from './core/mover-plug'
import EmoticonsPlug from './plugs/emoticons-plug'
import UploadPlug from './plugs/upload-plug'
import PreviewPlug from './plugs/preview-plug'
import HeaderInputPlug from './core/header-input-plug'
import EditorPlug from './editor-plug'
import Editor from './editor'
import LocalStoragePlug from './core/local-storage-plug'
import HeaderPlug from './core/header-plug'
import TextareaPlug from './core/textarea-plug'
import SubmitBtnPlug from './core/submit-btn-plug'
import SubmitPlug from './core/submit-plug'
import ReplyPlug from './core/reply-plug'
import EditPlug from './core/edit-plug'
import ClosablePlug from './core/closable-plug'

/** The default enabled plugs */
const ENABLED_PLUGS: (typeof EditorPlug)[] = [
  // Core
  LocalStoragePlug,
  HeaderPlug, HeaderInputPlug, TextareaPlug,
  SubmitPlug, SubmitBtnPlug,
  MoverPlug, ReplyPlug, EditPlug,
  ClosablePlug,

  // Extensions
  EmoticonsPlug, UploadPlug, PreviewPlug
]

export default class PlugManager {
  plugs: EditorPlug[] = []
  openedPlug: EditorPlug|null = null

  constructor(
    public editor: Editor
  ) {
    // handle ui, clear and reset the plug btns and plug panels
    editor.getUI().$plugPanelWrap.innerHTML = ''
    editor.getUI().$plugPanelWrap.style.display = 'none'
    editor.getUI().$plugBtnWrap.innerHTML = ''

    // init the all enabled plugs
    const DISABLED = getDisabledPlugByConf(editor.conf)

    ENABLED_PLUGS
      .filter(p => !DISABLED.includes(p)) // 禁用的插件
      .forEach((Plug) => {
        // create the plug instance
        this.plugs.push(new Plug(this.editor))
      })

    // load the plug UI
    this.plugs.forEach((plug) => {
      this.loadPlugUI(plug)
    })
  }

  /** Load the plug btn and plug panel on editor ui */
  private loadPlugUI(plug: EditorPlug) {
    const $btn = plug.$btn
    if (!$btn) return
    this.editor.getUI().$plugBtnWrap.appendChild($btn)

    // bind the event when click plug btn
    $btn.onclick = $btn.onclick || (() => {
      // removing the active class from all the buttons
      this.editor.getUI().$plugBtnWrap
        .querySelectorAll('.active')
        .forEach(item => item.classList.remove('active'))

      // if the plug is not the same as the openedPlug,
      if (plug !== this.openedPlug) {
        // then open the plug current clicked plug panel
        this.openPlugPanel(plug)

        // add active class for current plug panel
        $btn.classList.add('active')
      } else {
        // then close the plug
        this.closePlugPanel()
      }
    })

    // initialization of plug panel
    const $panel = plug.$panel
    if ($panel) {
      $panel.style.display = 'none'
      this.editor.getUI().$plugPanelWrap.appendChild($panel)
    }
  }

  get<T extends typeof EditorPlug>(plug: T) {
    return this.plugs.find(p => p instanceof plug) as InstanceType<T> | undefined;
  }

  /** Open the editor plug panel */
  openPlugPanel(plug: EditorPlug) {
    this.plugs.forEach((aPlug) => {
      const plugPanel = aPlug.$panel
      if (!plugPanel) return

      if (aPlug === plug) {
        plugPanel.style.display = ''
        plug.onPanelShow && plug.onPanelShow()
      } else {
        plugPanel.style.display = 'none'
        plug.onPanelHide && plug.onPanelHide()
      }
    })

    this.editor.getUI().$plugPanelWrap.style.display = ''
    this.openedPlug = plug
  }

  /** Close the editor plug panel */
  closePlugPanel() {
    if (!this.openedPlug) return

    this.openedPlug.onPanelHide && this.openedPlug.onPanelHide()

    this.editor.getUI().$plugPanelWrap.style.display = 'none'
    this.openedPlug = null
  }

  /** Get the content which is transformed by plugs */
  getTransformedContent(rawContent: string) {
    let result = rawContent
    this.plugs.forEach((aPlug) => {
      if (!aPlug.contentTransformer) return
      result = aPlug.contentTransformer(result)
    })
    return result
  }

  // -------------------------------------------------------------------
  //  Events
  // -------------------------------------------------------------------

  /** Trigger event when mounted */
  triggerMounted() {
    this.plugs.forEach((aPlug) => aPlug.onMounted && aPlug.onMounted())
  }

  /** Trigger event when unmounted */
  triggerUnmounted() {
    this.plugs.forEach((aPlug) => aPlug.onUnmounted && aPlug.onUnmounted())
  }

  /** Trigger event when editor header input changed */
  triggerHeaderInputEvt(field: string, $input: HTMLInputElement) {
    this.plugs.forEach((aPlug) => aPlug.onHeaderInput && aPlug.onHeaderInput(field, $input))
  }

  /** Trigger event when editor content updated */
  triggerContentUpdatedEvt(content: string) {
    this.plugs.forEach((aPlug) => aPlug.onContentUpdated && aPlug.onContentUpdated(content))
  }
}

/** Get the name list of disabled plugs */
function getDisabledPlugByConf(conf: ArtalkConfig): (typeof EditorPlug)[] {
  return [
    {k: UploadPlug, v: conf.imgUpload},
    {k: EmoticonsPlug, v: conf.emoticons},
    {k: PreviewPlug, v: conf.preview},
    {k: MoverPlug, v: conf.editorTravel},
  ].filter(n => !n.v).flatMap(n => n.k)
}

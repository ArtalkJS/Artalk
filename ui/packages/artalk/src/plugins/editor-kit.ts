import type ArtalkPlug from '~/types/plug'
import type { EditorApi } from '~/types/editor'
import EventManager from '@/lib/event-manager'
import { ENABLED_PLUGS, getDisabledPlugByConf } from './editor'
import EditorPlug from './editor/_plug'
import PlugKit from './editor/_kit'

export interface EditorEventPayloadMap {
  'mounted': undefined
  'unmounted': undefined
  'header-input': { field: string, $input: HTMLInputElement }
  'content-updated': string
  'panel-show': EditorPlug
  'panel-hide': EditorPlug
  'panel-close': undefined

  'editor-close': undefined
  'editor-open': undefined
}

export const EditorKit: ArtalkPlug = (ctx) => {
  const editor = ctx.get('editor')

  // handle ui, clear and reset the plug btns and plug panels
  editor.getUI().$plugPanelWrap.innerHTML = ''
  editor.getUI().$plugPanelWrap.style.display = 'none'
  editor.getUI().$plugBtnWrap.innerHTML = ''

  const editorPlugs = new PlugManager(editor)
  ctx.inject('editorPlugs', editorPlugs)
}

export class PlugManager {
  private plugs: EditorPlug[] = []
  private openedPlug: EditorPlug|null = null
  private events = new EventManager<EditorEventPayloadMap>()

  getPlugs() { return this.plugs }
  getEvents() { return this.events }

  constructor(
    public editor: EditorApi
  ) {
    // init the all enabled plugs
    const DISABLED = getDisabledPlugByConf(editor.conf)

    ENABLED_PLUGS
      .filter(p => !DISABLED.includes(p)) // 禁用的插件
      .forEach((Plug) => {
        // create the plug instance
        const kit = new PlugKit(this)
        this.plugs.push(new Plug(kit))
      })

    // load the plug UI
    this.plugs.forEach((plug) => {
      this.loadPlugUI(plug)
    })

    // bind events
    this.events.on('panel-close', () => this.closePlugPanel())
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
        this.events.trigger('panel-show', plug)
      } else {
        plugPanel.style.display = 'none'
        this.events.trigger('panel-hide', plug)
      }
    })

    this.editor.getUI().$plugPanelWrap.style.display = ''
    this.openedPlug = plug
  }

  /** Close the editor plug panel */
  closePlugPanel() {
    if (!this.openedPlug) return

    this.editor.getUI().$plugPanelWrap.style.display = 'none'
    this.events.trigger('panel-hide', this.openedPlug)
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
}
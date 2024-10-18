import { getEnabledPlugs } from './editor'
import EditorPlugin from './editor/_plug'
import PlugKit from './editor/_kit'
import { EventManager } from '@/lib/event-manager'
import type {
  Editor,
  ArtalkPlugin,
  DataManager,
  ConfigManager,
  UserManager,
  EditorPluginManager as IEditorPluginManager,
  CheckerManager,
} from '@/types'
import type { Api } from '@/api'

export interface EditorEventPayloadMap {
  mounted: undefined
  unmounted: undefined
  'header-input': { field: string; $input: HTMLInputElement }
  'header-change': { field: string; $input: HTMLInputElement }
  'content-updated': string
  'panel-show': EditorPlugin
  'panel-hide': EditorPlugin
  'panel-close': undefined

  'editor-close': undefined
  'editor-open': undefined
  'editor-submit': undefined
  'editor-submitted': undefined
}

export const EditorKit: ArtalkPlugin = (ctx) => {
  ctx.provide(
    'editorPlugs',
    (editor, config, user, data, checkers, api) => {
      const pluginManager = new PluginManager({
        getArtalkRootEl: () => ctx.getEl(),
        getEditor: () => editor,
        getConf: () => config,
        getUser: () => user,
        getApi: () => api,
        getData: () => data,
        getCheckers: () => checkers,
        onSubmitted: () => ctx.trigger('editor-submitted'),
      })
      editor.setPlugins(pluginManager)
      return pluginManager
    },
    ['editor', 'config', 'user', 'data', 'checkers', 'api'] as const,
  )
}

export interface PluginManagerOptions {
  getArtalkRootEl: () => HTMLElement
  getEditor: () => Editor
  getConf: () => ConfigManager
  getUser: () => UserManager
  getApi: () => Api
  getData: () => DataManager
  getCheckers: () => CheckerManager
  onSubmitted: () => void
}

export class PluginManager implements IEditorPluginManager {
  private plugins: EditorPlugin[] = []
  private openedPlug: EditorPlugin | null = null
  private events = new EventManager<EditorEventPayloadMap>()

  constructor(public opts: PluginManagerOptions) {
    let confLoaded = false // config not loaded at first time
    this.opts
      .getConf()
      .watchConf(['imgUpload', 'emoticons', 'preview', 'editorTravel', 'locale'], (conf) => {
        // trigger unmount event will call all plugs' unmount function
        // (this will only be called while conf reloaded, not be called at first time)
        confLoaded && this.getEvents().trigger('unmounted')

        // reset the plugs
        this.clear()

        // init the all enabled plugs
        getEnabledPlugs(conf).forEach((Plug) => {
          // create the plug instance
          const kit = new PlugKit(this)
          this.plugins.push(new Plug(kit))
        })

        // trigger event for plug initialization
        this.getEvents().trigger('mounted')
        confLoaded = true

        // refresh the plug UI
        this.loadPluginUI()
      })

    this.events.on('panel-close', () => this.closePluginPanel())
    this.events.on('editor-submitted', () => opts.onSubmitted())
  }

  getPlugins() {
    return this.plugins
  }

  getEvents() {
    return this.events
  }

  getEditor() {
    return this.opts.getEditor()
  }

  getOptions() {
    return this.opts
  }

  private clear() {
    this.plugins = []
    this.events = new EventManager()
    if (this.openedPlug) this.closePluginPanel()
  }

  private loadPluginUI() {
    // handle ui, clear and reset the plug btns and plug panels
    this.getEditor().getUI().$plugPanelWrap.innerHTML = ''
    this.getEditor().getUI().$plugPanelWrap.style.display = 'none'
    this.getEditor().getUI().$plugBtnWrap.innerHTML = ''

    // load the plug UI
    this.plugins.forEach((plug) => this.loadPluginItem(plug))
  }

  /** Load the plug btn and plug panel on editor ui */
  private loadPluginItem(plug: EditorPlugin) {
    const $btn = plug.$btn
    if (!$btn) return
    this.getEditor().getUI().$plugBtnWrap.appendChild($btn)

    // bind the event when click plug btn
    !$btn.onclick &&
      ($btn.onclick = () => {
        // removing the active class from all the buttons
        this.opts
          .getEditor()
          .getUI()
          .$plugBtnWrap.querySelectorAll('.active')
          .forEach((item) => item.classList.remove('active'))

        // if the plug is not the same as the openedPlug,
        if (plug !== this.openedPlug) {
          // then open the plug current clicked plug panel
          this.openPluginPanel(plug)

          // add active class for current plug panel
          $btn.classList.add('active')
        } else {
          // then close the plug
          this.closePluginPanel()
        }
      })

    // initialization of plug panel
    const $panel = plug.$panel
    if ($panel) {
      $panel.style.display = 'none'
      this.getEditor().getUI().$plugPanelWrap.appendChild($panel)
    }
  }

  get<T extends typeof EditorPlugin>(plug: T) {
    return this.plugins.find((p) => p instanceof plug) as InstanceType<T> | undefined
  }

  /** Open the editor plug panel */
  openPluginPanel(plug: EditorPlugin) {
    this.plugins.forEach((aPlug) => {
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

    this.getEditor().getUI().$plugPanelWrap.style.display = ''
    this.openedPlug = plug
  }

  /** Close the editor plugin panel */
  closePluginPanel() {
    if (!this.openedPlug) return

    this.getEditor().getUI().$plugPanelWrap.style.display = 'none'
    this.events.trigger('panel-hide', this.openedPlug)
    this.openedPlug = null
  }

  /** Get the content which is transformed by plugs */
  getTransformedContent(rawContent: string) {
    let result = rawContent
    this.plugins.forEach((aPlug) => {
      if (!aPlug.contentTransformer) return
      result = aPlug.contentTransformer(result)
    })
    return result
  }
}

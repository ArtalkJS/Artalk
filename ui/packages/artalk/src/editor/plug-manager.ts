import type ArtalkConfig from '~/types/artalk-config'
import EmoticonsPlug from './plugs/emoticons-plug'
import UploadPlug from './plugs/upload-plug'
import PreviewPlug from './plugs/preview-plug'
import HeaderInputPlug from './plugs/header-input-plug'
import EditorPlug from './plugs/editor-plug'
import Editor from './editor'

/** The default enabled plugs */
const ENABLED_PLUGS = [ EmoticonsPlug, UploadPlug, PreviewPlug, HeaderInputPlug ]

/** Context of an editor plug manager */
export interface PlugManager {
  editor: Editor
  plugList: { [name: string]: EditorPlug }
  openedPlugName: string|null

  openPlugPanel(plugName: string): void
  closePlugPanel(): void

  triggerHeaderInputEvt(field: string, $input: HTMLInputElement): void
  triggerContentUpdatedEvt(content: string): void
  getTransformedContent(raw: string): string
}

/** Create an editor plug manager */
export function createPlugManager(editor: Editor): PlugManager {
  const ctx: PlugManager = {
    editor,
    plugList: {},
    openedPlugName: null,
    openPlugPanel: (p: string) => openPlugPanel(ctx, p),
    closePlugPanel: () => closePlugPanel(ctx),

    triggerHeaderInputEvt: (f, $) => triggerHeaderInputEvt(ctx, f, $),
    triggerContentUpdatedEvt: (c) => triggerContentUpdatedEvt(ctx, c),
    getTransformedContent: (r) => getTransformedContent(ctx, r),
  }

  // handle ui
  editor.getUI().$plugPanelWrap.innerHTML = ''
  editor.getUI().$plugPanelWrap.style.display = 'none'
  editor.getUI().$plugBtnWrap.innerHTML = ''

  // init the all enabled plugs
  const DISABLED = getDisabledPlugNames(editor.conf)

  ENABLED_PLUGS
    .filter(p => !DISABLED.includes(p.name)) // 禁用的插件
    .forEach((Plug) => {
      initPlugItem(ctx, Plug)
    })

  return ctx
}

/** Get the name list of disabled plugs */
function getDisabledPlugNames(conf: ArtalkConfig) {
  return [
    {k: 'upload', v: conf.imgUpload},
    {k: 'emoticons', v: conf.emoticons},
    {k: 'preview', v: conf.preview}
  ].filter(n => !n.v).flatMap(n => n.k)
}

/** Initialization of a plug item */
function initPlugItem(ctx: PlugManager, Plug: typeof ENABLED_PLUGS[number]) {
  const plugName = Plug.Name

  // create the plug instance
  const plug = new Plug(ctx.editor)
  ctx.plugList[plugName] = plug

  // gen plug ui
  genPlugBtn(ctx, plugName, plug)
}

/** Gen the plug btn (and plug panel) on editor ui */
function genPlugBtn(ctx: PlugManager, plugName: string, plug: EditorPlug) {
  const $btn = plug.getBtn()
  if (!$btn) return
  ctx.editor.getUI().$plugBtnWrap.appendChild($btn)
  $btn.onclick = $btn.onclick || (() => {
    // the event when click plug btn

    // removing the active class from all the buttons
    ctx.editor.getUI().$plugBtnWrap
      .querySelectorAll('.active')
      .forEach(item => item.classList.remove('active'))


    // if the plugName is the same as the openedPlugName, then close the plugPanel
    if (plugName === ctx.openedPlugName) {
      closePlugPanel(ctx)
      return
    }

    // open the plug current clicked plug panel
    openPlugPanel(ctx, plugName)

    // add active class for current plug panel
    $btn.classList.add('active')
  })

  // initialization of plug panel
  const $panel = plug.getPanel()
  if ($panel) {
    $panel.setAttribute('data-plug-name', plugName)
    $panel.style.display = 'none'
    ctx.editor.getUI().$plugPanelWrap.appendChild($panel)
  }
}

/** Open the editor plug panel */
function openPlugPanel(ctx: PlugManager, plugName: string) {
  Object.entries(ctx.plugList).forEach(([aPlugName, plug]) => {
    const plugPanel = plug.getPanel()
    if (!plugPanel) return

    if (aPlugName === plugName) {
      plugPanel.style.display = ''
      if (plug.onPanelShow) plug.onPanelShow()
    } else {
      plugPanel.style.display = 'none'
      if (plug.onPanelHide) plug.onPanelHide()
    }
  })

  ctx.editor.getUI().$plugPanelWrap.style.display = ''
  ctx.openedPlugName = plugName
}

/** Close the editor plug panel */
function closePlugPanel(ctx: PlugManager) {
  if (!ctx.openedPlugName) return

  const plug = ctx.plugList[ctx.openedPlugName]
  if (!plug) return

  if (plug.onPanelHide) plug.onPanelHide()

  ctx.editor.getUI().$plugPanelWrap.style.display = 'none'
  ctx.openedPlugName = null
}

// ============== Events ====================

/** Trigger event when editor header input changed */
function triggerHeaderInputEvt(ctx: PlugManager, field: string, $input: HTMLInputElement) {
  Object.entries(ctx.plugList).forEach(([plugName, plug]) => {
    if (!plug.onHeaderInput) return
    plug.onHeaderInput(field, $input)
  })
}

/** Trigger event when editor content updated */
function triggerContentUpdatedEvt(ctx: PlugManager, content: string) {
  Object.entries(ctx.plugList).forEach(([plugName, plug]) => {
    if (!plug.onContentUpdated) return
    plug.onContentUpdated(content)
  })
}

/** Get the content which is transformed by plugs */
function getTransformedContent(ctx: PlugManager, rawContent: string) {
  let result = rawContent
  Object.entries(ctx.plugList).forEach(([plugName, plug]) => {
    if (!plug.contentTransformer) return
    result = plug.contentTransformer(result)
  })
  return result
}

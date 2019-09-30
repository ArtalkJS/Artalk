import './PreviewPlug.scss'
import Editor from '../Editor'
import ArtalkContext from '~/src/ArtalkContext'
import Utils from '~/src/utils'

export default class PreviewPlug extends ArtalkContext {
  public elem: HTMLElement

  constructor (public editor: Editor) {
    super()

    this.initElem()
  }

  initElem () {
    this.elem = Utils.createElement('<div class="artalk-editor-plug-example"></div>')
  }

  getName () {
    return 'example'
  }

  getBtnHtml () {
    return '栗子'
  }

  getElem () {
    return this.elem
  }

  onShow () {}

  onHide () {}
}

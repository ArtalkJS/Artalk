import $ from 'jquery'
import './PreviewPlug.scss'

export default class PreviewPlug {
  constructor (editor) {
    this.editor = editor
    this.artalk = editor.artalk

    this.initElem()
  }

  initElem () {
    this.elem = $('<div class="artalk-editor-plug-example"></div>')
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

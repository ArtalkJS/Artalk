import $ from 'jquery'
import '../css/editor.scss'
import EmoticonsPlug from './editor-plugs/EmoticonsPlug.js'
import PreviewPlug from './editor-plugs/PreviewPlug.js'

class Editor {
  constructor (artalk) {
    this.artalk = artalk

    this.plugs = [new EmoticonsPlug(this), new PreviewPlug(this)]
    this.el = $(require('../tpl/editor.ejs')(this)).appendTo(this.artalk.el)

    this.initTextarea()
    this.initEditorPlug()
  }

  initTextarea () {
    this.textareaEl = this.el.find('.artalk-editor-textarea')
  }

  initEditorPlug () {
    this.plugWrapEl = this.el.find('.artalk-editor-plug-wrap')

    let openedPlugName = null
    this.el.find('.artalk-editor-plug-switcher').click((evt) => {
      let btnElem = $(evt.currentTarget)
      let plugIndex = btnElem.attr('data-plug-index')
      let plug = this.plugs[plugIndex]
      if (!plug) {
        throw Error(`Plug index=${plugIndex} was not found`)
      }

      this.el.find('.artalk-editor-plug-switcher.active').removeClass('active')

      if (openedPlugName === plug.getName()) {
        this.plugWrapEl.css('display', 'none')
        openedPlugName = null
        return
      }

      if (!this.plugWrapEl.find(`[data-plug-name="${plug.getName()}"]`).length) {
        // 需要初始化
        plug.getElem()
          .attr('data-plug-name', plug.getName())
          .css('display', 'none')
          .appendTo(this.plugWrapEl)
      }

      this.plugWrapEl.children().each((i, plugElem) => {
        plugElem = $(plugElem)
        if (plugElem.attr('data-plug-name') !== plug.getName()) {
          plugElem.css('display', 'none')
        } else {
          plugElem.css('display', '')
          plug.onShow()
        }
      })

      this.plugWrapEl.css('display', '')
      openedPlugName = plug.getName()

      btnElem.addClass('active')
    })
  }

  putContent (val) {
    var cursorPos = this.textareaEl.prop('selectionStart')
    var v = this.textareaEl.val()
    var textBefore = v.substring(0, cursorPos)
    var textAfter = v.substring(cursorPos, v.length)
    this.textareaEl.val(textBefore + val + textAfter)
    this.textareaEl.focus()
  }

  setContent (val) {
    this.textareaEl.text(val)
  }

  clearContent () {
    this.textareaEl.html()
  }

  getContent () {
    return this.textareaEl.val()
  }

  getContentMarked () {
    return this.artalk.marked(this.textareaEl.val())
  }
}

export default Editor

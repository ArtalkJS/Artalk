import $ from 'jquery'

export default class Ui {
  constructor (artalk) {
    this.artalk = artalk
    this.el = this.artalk.el
  }

  /**
   * 显示加载
   */
  showLoading (parentElem) {
    if (!parentElem) {
      parentElem = this.el
    }
    let loadingElem = $(parentElem).find('.artalk-loading')
    if (!loadingElem.length) {
      loadingElem = $(`<div class="artalk-loading" style="display: none;"><div class="artalk-loading-spinner"><svg viewBox="25 25 50 50"><circle cx="50" cy="50" r="20" fill="none" stroke-width="2" stroke-miterlimit="10"></circle></svg></div></div>`).appendTo(parentElem)
    }
    loadingElem.css('display', '')
  }

  /**
   * 隐藏加载
   */
  hideLoading (parentElem) {
    if (!parentElem) {
      parentElem = this.el
    }
    let loadingElem = $(parentElem).find('.artalk-loading')
    loadingElem.css('display', 'none')
  }

  /**
   * 设置全局错误
   */
  setGlobalError (html) {
    let elem = this.el.find('.artalk-error-layer')
    if (html === null) {
      elem.remove()
      return
    }
    if (!elem.length) {
      elem = $('<div class="artalk-error-layer"><span class="artalk-error-title">Artalk Error</span><span class="artalk-error-text"></span></div>').appendTo(this.el)
    }
    elem.find('.artalk-error-text').html('')
    $(html).appendTo(elem.find('.artalk-error-text'))
  }

  /**
   * 显示对话框
   */
  showDialog (parentElem, html, onConfirm, onCancel) {
    if (!parentElem) {
      throw Error('showDialog 未指定 parentElem')
    }

    let dialogElem = $(`<div class="artalk-layer-dialog-wrap"><div class="artalk-layer-dialog"><div class="artalk-layer-dialog-content"></div><div class="artalk-layer-dialog-action"></div>`).appendTo(parentElem)

    // 按钮
    let actionElem = dialogElem.find('.artalk-layer-dialog-action')
    let onclick = (func) => (evt) => {
      let returnVal = func(dialogElem, $(evt.currentTarget))
      if (returnVal === undefined || returnVal === true) {
        dialogElem.remove()
      }
    }
    if (typeof onConfirm === 'function') {
      $('<button data-action="confirm">确定</button>').click(onclick(onConfirm)).appendTo(actionElem)
    }
    if (typeof onCancel === 'function') {
      $('<button data-action="cancel">取消</button>').click(onclick(onCancel)).appendTo(actionElem)
    }

    // 内容
    $(html).appendTo(dialogElem.find('.artalk-layer-dialog-content'))
  }

  /**
   * 显示消息
   */
  showNotify (msg, type, wrapElem) {
    if (!wrapElem) {
      throw Error('wrapElem 未指定')
    }

    let colors = { s: '#57d59f', e: '#ff6f6c', w: '#ffc721', i: '#2ebcfc' }
    if (!colors[type]) {
      throw Error('showNotify 的 type 有问题！仅支持：' + Object.keys(colors).join(', '))
    }

    let timeout = 3000 // 持续显示时间 ms

    let notifyElem = $(`<div class="artalk-notify artalk-fade-in" style="background-color: ${colors[type]}"><span class="artalk-notify-content"></span></div>`)
    notifyElem.find('.artalk-notify-content').html($('<div/>').text(msg).html().replace('\n', '<br/>'))
    notifyElem.appendTo(wrapElem)

    let notifyRemove = () => {
      notifyElem.addClass('artalk-fade-out')
      setTimeout(() => {
        notifyElem.remove()
      }, 200)
    }

    let timeoutFn
    if (timeout > 0) {
      timeoutFn = setTimeout(() => {
        notifyRemove()
      }, timeout)
    }

    notifyElem.click(() => {
      notifyRemove()
      clearTimeout(timeoutFn)
    })
  }

  /**
   * 滚动到元素中心
   */
  scrollIntoView (elem) {
    if (!this.isVisible(elem)) {
      $('html,body').animate({
        scrollTop: $(elem).offset().top - $(window).height() / 2
      }, 500)
    }
  }

  /**
   * 元素是否用户可见
   */
  isVisible (elem) {
    let docViewTop = $(window).scrollTop()
    let docViewBottom = docViewTop + $(window).height()

    let elemTop = $(elem).offset().top
    let elemBottom = elemTop + $(elem).height()

    return ((elemBottom <= docViewBottom) && (elemTop >= docViewTop))
  }
}

import './_EmoticonsPlug/EmoticonsPlug.less'
import Editor from '../Editor'
import ArtalkContext from '~/src/ArtalkContext'
import Utils from '~/src/utils'

type EmoticonListType = {
  [grpName: string]: {
    inputType: 'emoticon'|'image'
    container: { [name: string]: string }
  }
}

export default class EmoticonsPlug extends ArtalkContext {
  public elem: HTMLElement
  public emoticons: EmoticonListType

  public listWrapElem: HTMLElement
  public typesElem: HTMLElement


  constructor (public editor: Editor) {
    super()

    this.emoticons = this.artalk.conf.emoticons
    this.initElem()
  }

  initElem () {
    this.elem = Utils.createElement(require('./_EmoticonsPlug/EmoticonsPlug.ejs')(this))

    this.listWrapElem = this.elem.querySelector('.artalk-emoticons-list-wrap')
    this.typesElem = this.elem.querySelector('.artalk-emoticons-types')

    // 绑定切换类型按钮
    this.typesElem.querySelectorAll('span').forEach((btn) => {
      btn.addEventListener('click', (evt) => {
        const btnElem = evt.currentTarget as HTMLElement
        const key = btnElem.getAttribute('data-key')
        if (key) this.openType(key)
      })
    })

    // 默认打开第一个类型
    if (Object.keys(this.emoticons).length > 0)
      this.openType(Object.keys(this.emoticons)[0])

    // 绑定点击表情
    this.listWrapElem.querySelectorAll('.artalk-emoticons-item').forEach((item: HTMLElement) => {
      item.onclick = (evt) => {
        const elem = evt.currentTarget as HTMLElement
        const inputType = elem.closest('.artalk-emoticons-list').getAttribute('data-input-type')

        const title = elem.getAttribute('title')
        const content = elem.getAttribute('data-content')
        if (inputType === 'image') {
          this.editor.insertContent(`:[${title}]`)
        } else {
          this.editor.insertContent(content)
        }
      }
    })
  }

  openType (key: string) {
    Array.from(this.listWrapElem.children).forEach((item: HTMLElement) => {
      if (item.getAttribute('data-key') !== key) {
        item.style.display = 'none'
      } else {
        item.style.display = ''
      }
    })

    this.typesElem.querySelectorAll('span.active').forEach(item => item.classList.remove('active'))
    this.typesElem.querySelector(`span[data-key="${key}"]`).classList.add('active')

    this.changeListHeight()
  }

  getName () {
    return 'emoticons'
  }

  getBtnHtml () {
    return '表情'
  }

  getElem () {
    return this.elem
  }

  changeListHeight () {
    /* const listWrapHeight = Utils.getHeight(this.listWrapElem)
    this.editor.plugWrapEl.style.height = `${listWrapHeight > 150 ? listWrapHeight : 150}px` */
  }

  onShow () {
    // 延迟执行，防止无法读取高度
    setTimeout(() => {
      this.changeListHeight()
    }, 30)
  }

  onHide () {
    this.elem.parentElement.style.height = ''
  }

  public transEmoticonImageText (text: string) {
    Object.entries(this.emoticons).forEach(([grpName, grp]) => {
      if (grp.inputType !== 'image') return
      Object.entries(grp.container).forEach(([name, imgSrc]) => {
        text = text.split(`:[${name}]`).join(`![${name}](${imgSrc}) `) // replaceAll(...)
      })
    })

    return text
  }
}

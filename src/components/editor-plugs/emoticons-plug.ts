import './emoticons-plug.less'

import Editor from '../editor'
import EditorPlug from './editor-plug'
import * as Utils from '~/src/lib/utils'

type EmoticonListType = {
  [grpName: string]: {
    inputType: 'emoticon'|'image'
    container: { [name: string]: string }
  }
}

export default class EmoticonsPlug extends EditorPlug {
  public $el!: HTMLElement
  public emoticons: EmoticonListType

  public listWrapEl!: HTMLElement
  public typesEl!: HTMLElement

  constructor (public editor: Editor) {
    super(editor)

    this.emoticons = this.ctx.conf.emoticons
    this.initEl()
  }

  initEl () {
    this.$el = Utils.createElement(`<div class="atk-editor-plug-emoticons"></div>`)

    // 表情列表
    this.listWrapEl = Utils.createElement(`<div class="atk-emoticons-list-wrap"></div>`)
    this.$el.append(this.listWrapEl)

    Object.entries(this.emoticons).forEach(([key, item]) => {
      const emoticonsEl = Utils.createElement(`<div class="atk-emoticons-list" style="display: none;"></div>`)
      this.listWrapEl.append(emoticonsEl)
      emoticonsEl.setAttribute('data-key', key)
      emoticonsEl.setAttribute('data-input-type', item.inputType)
      Object.entries(item.container).forEach(([name, content]) => {
        const itemEl = Utils.createElement(`<span class="atk-emoticons-item"></span>`)
        emoticonsEl.append(itemEl)
        itemEl.setAttribute('title', name)
        itemEl.setAttribute('data-content', content)
        if (item.inputType === 'image') {
          const imgEl = document.createElement('img')
          imgEl.src = content
          imgEl.alt = name
          itemEl.append(imgEl)
        } else {
          itemEl.innerText = content
        }
      })
    })

    // 表情分类
    this.typesEl = Utils.createElement(`<div class="atk-emoticons-types"></div>`)
    this.$el.append(this.typesEl)
    Object.entries(this.emoticons).forEach(([key, item]) => {
      const itemEl = Utils.createElement('<span />')
      this.typesEl.append(itemEl)
      itemEl.setAttribute('data-key', key)
      itemEl.innerText = key
    })

    // 绑定切换分类按钮
    this.typesEl.querySelectorAll('span').forEach((btn) => {
      btn.addEventListener('click', (evt) => {
        const btnEl = evt.currentTarget as HTMLElement
        const key = btnEl.getAttribute('data-key')
        if (key) this.openType(key)
      })
    })

    // 默认打开第一个分类
    if (Object.keys(this.emoticons).length > 0)
      this.openType(Object.keys(this.emoticons)[0])

    // 绑定点击表情
    this.listWrapEl.querySelectorAll<HTMLElement>('.atk-emoticons-item').forEach((item: HTMLElement) => {
      item.onclick = (evt) => {
        const elem = evt.currentTarget as HTMLElement
        const inputType = elem.closest('.atk-emoticons-list')!.getAttribute('data-input-type')

        const title = elem.getAttribute('title')
        const content = elem.getAttribute('data-content')
        if (inputType === 'image') {
          this.editor.insertContent(`:[${title}]`)
        } else {
          this.editor.insertContent(content || '')
        }
      }
    })
  }

  openType (key: string) {
    Array.from(this.listWrapEl.children).forEach((item) => {
      const el = item as HTMLElement
      if (el.getAttribute('data-key') !== key) {
        el.style.display = 'none'
      } else {
        el.style.display = ''
      }
    })

    this.typesEl.querySelectorAll('span.active').forEach(item => item.classList.remove('active'))
    this.typesEl.querySelector(`span[data-key="${key}"]`)?.classList.add('active')

    this.changeListHeight()
  }

  getName () {
    return 'emoticons'
  }

  getBtnHtml () {
    return '表情'
  }

  getEl () {
    return this.$el
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
    this.$el.parentElement!.style.height = ''
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

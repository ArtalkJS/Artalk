import './emoticons-plug.less'

import Editor from '../editor'
import EditorPlug from './editor-plug'
import * as Utils from '~/src/lib/utils'
import * as Ui from '~/src/lib/ui'
import { EmoticonListData, EmoticonGrpData } from '~/types/artalk-data'

type OwOFormatType = {
  [key: string] : {
    type: 'emoticon'|'emoji'|'image',
    container: {icon: string, text: string}[]
  }
}

export default class EmoticonsPlug extends EditorPlug {
  public $el!: HTMLElement
  public emoticons: EmoticonListData = []

  public $listWrap!: HTMLElement
  public $types!: HTMLElement

  constructor (public editor: Editor) {
    super(editor)

    this.$el = Utils.createElement(`<div class="atk-editor-plug-emoticons"></div>`)

    this.init()
  }

  async init() {
    // 数据处理
    Ui.showLoading(this.$el)

    if (typeof this.ctx.conf.emoticons === 'string') {
      this.emoticons = await this.remoteLoad(this.ctx.conf.emoticons)
    } else {
      this.emoticons = this.ctx.conf.emoticons
    }

    this.checkConvertOwO()

    if (!Array.isArray(this.emoticons)) {
      Ui.setError(this.$el, "表情包数据必须为 Array 类型")
      Ui.hideLoading(this.$el)
      return
    }

    // 加载子内容
    await Promise.all(this.emoticons.map(async (grp, index) => {
      if (typeof grp === "string") {
        const grpData = await this.remoteLoad(grp)
        if (grpData) {
          this.emoticons[index] = grpData
        }
      }
    }))

    Ui.hideLoading(this.$el)




    this.solveNullKey()
    this.solveSameKey()

    // 初始化元素
    this.initEmoticonsList()
  }

  async remoteLoad(url: string): Promise<any> {
    if (!url) return []

    try {
      const resp = await fetch(url)
      const json = await resp.json()
      return json
    } catch (err) {
      Ui.hideLoading(this.$el)
      Ui.setError(this.$el, `表情加载失败 ${String(err)}`)
      return []
    }
  }

  solveNullKey() {
    this.emoticons.forEach((grp) => {
      grp.items.forEach((item, index) => {
        if (!item.key) item.key = `${grp.name} ${index+1}`
      })
    })
  }

  solveSameKey() {
    const tmp: {[key: string]: number} = {}
    this.emoticons.forEach((grp) => {
      grp.items.forEach(item => {
        if (!item.key || String(item.key).trim() === "") return
        if (!tmp[item.key]) tmp[item.key] = 1
        else tmp[item.key]++

        if (tmp[item.key] > 1) item.key = `${item.key} ${tmp[item.key]}`
      })
    })
  }

  checkConvertOwO() {
    if (this.isOwOFormat(this.emoticons)) {
      this.emoticons = this.convertOwO(this.emoticons as any)
    }
  }

  isOwOFormat(data: any) {
    try {
      return (typeof data === 'object') && !!Object.values(data).length
        && Array.isArray(Object.keys(Object.values<any>(data)[0].container))
        && Object.keys(Object.values<any>(data)[0].container[0]).includes('icon')
    } catch { return false }
  }

  convertOwO(owoData: OwOFormatType): EmoticonListData {
    const dest: EmoticonListData = []
    Object.entries(owoData).forEach(([grpName, grp]) => {
      const nGrp: EmoticonGrpData = { name: grpName, type: grp.type, items: [] }
      grp.container.forEach((item, index) => {
        // 图片标签提取 src 属性值
        const iconStr = item.icon
        if (/<(img|IMG)/.test(iconStr)) {
          const find = /src=["'](.*?)["']/.exec(iconStr)
          if (find && find.length > 1) item.icon = find[1]
        }
        nGrp.items.push({ key: item.text || `${grpName} ${index+1}`, val: item.icon })
      })
      dest.push(nGrp)
    })
    return dest
  }

  initEmoticonsList () {
    // 表情列表
    this.$listWrap = Utils.createElement(`<div class="atk-emoticons-list-wrap"></div>`)
    this.$el.append(this.$listWrap)

    this.emoticons.forEach((grp, index) => {
      const emoticonsEl = Utils.createElement(`<div class="atk-emoticons-list" style="display: none;"></div>`)
      this.$listWrap.append(emoticonsEl)
      emoticonsEl.setAttribute('data-index', String(index))
      emoticonsEl.setAttribute('data-grp-name', grp.name)
      emoticonsEl.setAttribute('data-type', grp.type)
      grp.items.forEach((item) => {
        const $item = Utils.createElement(`<span class="atk-emoticons-item"></span>`)
        emoticonsEl.append($item)

        if (!!item.key && !(new RegExp(`^(${grp.name})?\\s?[0-9]+$`).test(item.key)))
          $item.setAttribute('title', item.key)

        if (grp.type === 'image') {
          const imgEl = document.createElement('img')
          imgEl.src = item.val
          imgEl.alt = item.key
          $item.append(imgEl)
        } else {
          $item.innerText = item.val
        }

        $item.onclick = () => {
          if (grp.type === 'image') {
            this.editor.insertContent(`:[${item.key}]`)
          } else {
            this.editor.insertContent(item.val || '')
          }
        }
      })
    })

    // 表情分类
    this.$types = Utils.createElement(`<div class="atk-emoticons-types"></div>`)
    this.$el.append(this.$types)
    this.emoticons.forEach((grp, index) => {
      const $item = Utils.createElement('<span />')
      $item.innerText = grp.name
      $item.setAttribute('data-index', String(index))
      $item.onclick = () => (this.openType(index))
      this.$types.append($item)
    })

    // 默认打开第一个分类
    if (this.emoticons.length > 0)
      this.openType(0)
  }

  openType (index: number) {
    Array.from(this.$listWrap.children).forEach((item) => {
      const el = item as HTMLElement
      if (el.getAttribute('data-index') !== String(index)) {
        el.style.display = 'none'
      } else {
        el.style.display = ''
      }
    })

    this.$types.querySelectorAll('span.active').forEach(item => item.classList.remove('active'))
    this.$types.querySelector(`span[data-index="${index}"]`)?.classList.add('active')

    this.changeListHeight()
  }

  static Name = 'emoticons'
  static BtnHTML = '表情'

  getEl () {
    return this.$el
  }

  changeListHeight () {
    /* const listWrapHeight = Utils.getHeight(this.$listWrapem)
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
    if (!this.emoticons || !Array.isArray(this.emoticons))
      return text

    this.emoticons.forEach((grp) => {
      if (grp.type !== 'image') return
      Object.entries(grp.items).forEach(([index, item]) => {
        text = text.split(`:[${item.key}]`).join(`<img src="${item.val}" atk-emoticon="${item.key}">`) // replaceAll(...)
      })
    })

    return text
  }
}

import './emoticons-plug.scss'

import * as Utils from '~/src/lib/utils'
import * as Ui from '~/src/lib/ui'
import { EmoticonListData, EmoticonGrpData } from '~/types/artalk-data'
import EditorPlug from './editor-plug'
import Editor from '../editor'

type OwOFormatType = {
  [key: string] : {
    type: 'emoticon'|'emoji'|'image',
    container: {icon: string, text: string}[]
  }
}

export default class EmoticonsPlug extends EditorPlug {
  public static Name = 'emoticons'

  declare protected $panel: HTMLElement

  public emoticons: EmoticonListData = []
  public loadingTask: Promise<void>|null = null

  public $grpWrap!: HTMLElement
  public $grpSwitcher!: HTMLElement

  public constructor(public editor: Editor) {
    super(editor)

    this.registerPanel(`<div class="atk-editor-plug-emoticons"></div>`)
    this.registerBtn(this.ctx.$t('emoticon'))
    this.registerContentTransformer((raw) => this.transEmoticonImageText(raw))

    // 表情包预加载
    window.setTimeout(() => {
      this.loadEmoticonsData()
    }, 1000) // 延迟 1s 加载
  }

  public onPanelShow() {
    ;(async () => {
      await this.loadEmoticonsData()

      // 初始化元素
      if (!this.isImgLoaded) {
        this.initEmoticonsList()
        this.isImgLoaded = true
      }

      // 延迟执行，防止无法读取高度
      setTimeout(() => {
        this.changeListHeight()
      }, 30)
    })()
  }

  public onPanelHide() {
    this.$panel.parentElement!.style.height = ''
  }

  isListLoaded = false
  isImgLoaded = false

  public async loadEmoticonsData() {
    if (this.isListLoaded) return
    if (this.loadingTask !== null) {
      await this.loadingTask
      return
    }

    // 数据处理
    this.loadingTask = (async () => {
      Ui.showLoading(this.$panel)
      this.emoticons = await this.handleData(this.ctx.conf.emoticons)
      Ui.hideLoading(this.$panel)
      this.loadingTask = null
      this.isListLoaded = true
    })()
    await this.loadingTask
  }

  private async handleData(data: any): Promise<any[]> {
    if (!Array.isArray(data) && ['object', 'string'].includes(typeof data)) {
      data = [data]
    }

    if (!Array.isArray(data)) {
      Ui.setError(this.$panel, "表情包数据必须为 Array/Object/String 类型")
      Ui.hideLoading(this.$panel)
      return []
    }

    const pushGrp = (grp: EmoticonGrpData) => {
      if (typeof grp !== 'object') return
      if (grp.name && data.find(o => o.name === grp.name)) return
      data.push(grp)
    }

    // 加载子内容
    const remoteLoad = async (d: any[]) => {
      await Promise.all(d.map(async (grp, index) => {
        if (typeof grp === 'object' && !Array.isArray(grp)) {
          pushGrp(grp)
        } else if (Array.isArray(grp)) {
          await remoteLoad(grp)
        } else if (typeof grp === "string") {
          const grpData = await this.remoteLoad(grp)

          if (Array.isArray(grpData)) await remoteLoad(grpData)
          else if (typeof grpData === 'object') pushGrp(grpData)
        }
      }))
    }
    await remoteLoad(data)

    // 处理数组数据
    data.forEach((item: any) => {
      if (this.isOwOFormat(item)) {
        const c = this.convertOwO(item)
        c.forEach((grp) => { pushGrp(grp) })
      } else if (Array.isArray(item)) {
        item.forEach((grp) => { pushGrp(grp) })
      }
    })

    // 剔除非法数据
    data = data.filter((item: any) => (typeof item === 'object' && !Array.isArray(item) && !!item && !!item.name))

    // console.log(data)

    this.solveNullKey(data)
    this.solveSameKey(data)

    return data
  }

  /** 远程加载 */
  private async remoteLoad(url: string): Promise<any> {
    if (!url) return []

    try {
      const resp = await fetch(url)
      const json = await resp.json()
      return json
    } catch (err) {
      Ui.hideLoading(this.$panel)
      Ui.setError(this.$panel, `表情加载失败 ${String(err)}`)
      return []
    }
  }

  /** 避免表情 item.key 为 null 的情况 */
  private solveNullKey(data: EmoticonGrpData[]) {
    data.forEach((grp) => {
      grp.items.forEach((item, index) => {
        if (!item.key) item.key = `${grp.name} ${index+1}`
      })
    })
  }

  /** 避免相同 item.key */
  private solveSameKey(data: EmoticonGrpData[]) {
    const tmp: {[key: string]: number} = {}
    data.forEach((grp) => {
      grp.items.forEach(item => {
        if (!item.key || String(item.key).trim() === "") return
        if (!tmp[item.key]) tmp[item.key] = 1
        else tmp[item.key]++

        if (tmp[item.key] > 1) item.key = `${item.key} ${tmp[item.key]}`
      })
    })
  }

  /** 判断是否为 OwO 格式 */
  private isOwOFormat(data: any) {
    try {
      return (typeof data === 'object') && !!Object.values(data).length
        && Array.isArray(Object.keys(Object.values<any>(data)[0].container))
        && Object.keys(Object.values<any>(data)[0].container[0]).includes('icon')
    } catch { return false }
  }

  /** 转换 OwO 格式数据 */
  private convertOwO(owoData: OwOFormatType): EmoticonListData {
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

  /** 初始化表情列表界面 */
  private initEmoticonsList() {
    // 表情列表
    this.$grpWrap = Utils.createElement(`<div class="atk-grp-wrap"></div>`)
    this.$panel.append(this.$grpWrap)

    this.emoticons.forEach((grp, index) => {
      const $grp = Utils.createElement(`<div class="atk-grp" style="display: none;"></div>`)
      this.$grpWrap.append($grp)
      $grp.setAttribute('data-index', String(index))
      $grp.setAttribute('data-grp-name', grp.name)
      $grp.setAttribute('data-type', grp.type)
      grp.items.forEach((item) => {
        const $item = Utils.createElement(`<span class="atk-item"></span>`)
        $grp.append($item)

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

    // 表情分类切换 bar
    if (this.emoticons.length > 1) {
      this.$grpSwitcher = Utils.createElement(`<div class="atk-grp-switcher"></div>`)
      this.$panel.append(this.$grpSwitcher)
      this.emoticons.forEach((grp, index) => {
        const $item = Utils.createElement('<span />')
        $item.innerText = grp.name
        $item.setAttribute('data-index', String(index))
        $item.onclick = () => (this.openGrp(index))
        this.$grpSwitcher.append($item)
      })
    }

    // 默认打开第一个分类
    if (this.emoticons.length > 0) this.openGrp(0)
  }

  /** 打开一个表情组 */
  public openGrp(index: number) {
    Array.from(this.$grpWrap.children).forEach((item) => {
      const el = item as HTMLElement
      if (el.getAttribute('data-index') !== String(index)) {
        el.style.display = 'none'
      } else {
        el.style.display = ''
      }
    })

    this.$grpSwitcher?.querySelectorAll('span.active').forEach(item => item.classList.remove('active'))
    this.$grpSwitcher?.querySelector(`span[data-index="${index}"]`)?.classList.add('active')

    this.changeListHeight()
  }

  private changeListHeight() {
    /* const listWrapHeight = Utils.getHeight(this.$grpWrapem)
    this.editor.plugWrapEl.style.height = `${listWrapHeight > 150 ? listWrapHeight : 150}px` */
  }

  /** 处理评论 content 中的表情内容 */
  public transEmoticonImageText(text: string) {
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

import '../style/editor.less'

import { CommentData } from '~/types/artalk-data'
import Context from '../context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import EditorHTML from './html/editor.html?raw'

import EmoticonsPlug from './editor-plugs/emoticons-plug'
import PreviewPlug from './editor-plugs/preview-plug'
import Api from '../api'

export default class Editor extends Component {
  private readonly LOADABLE_PLUG_LIST = [EmoticonsPlug, PreviewPlug]
  public plugList: { [name: string]: any } = {}

  public $header: HTMLElement
  public $textareaWrap: HTMLElement
  public $textarea: HTMLTextAreaElement
  public $plugWrap: HTMLElement
  public $bottom: HTMLElement
  public $plugBtnWrap: HTMLElement
  public $imgUploadBtn?: HTMLElement
  public $imgUploadInput?: HTMLInputElement
  public $submitBtn: HTMLButtonElement
  public $notifyWrap: HTMLElement

  private replyComment: CommentData|null = null
  private $sendReply: HTMLElement|null = null

  private isTraveling = false

  private get user () {
    return this.ctx.user
  }

  constructor (ctx: Context) {
    super(ctx)

    this.$el = Utils.createElement(EditorHTML)

    this.$header = this.$el.querySelector('.atk-header')!
    this.$textareaWrap = this.$el.querySelector('.atk-textarea-wrap')!
    this.$textarea = this.$el.querySelector('.atk-textarea')!
    this.$plugWrap = this.$el.querySelector('.atk-plug-wrap')!
    this.$bottom = this.$el.querySelector('.atk-bottom')!
    this.$plugBtnWrap = this.$el.querySelector('.atk-plug-btn-wrap')!
    this.$submitBtn = this.$el.querySelector('.atk-send-btn')!
    this.$notifyWrap = this.$el.querySelector('.atk-notify-wrap')!

    this.initLocalStorage()
    this.initHeader()
    this.initTextarea()
    this.initEditorPlug()
    this.initBottomPart()

    // 监听事件
    this.ctx.on('editor-open', () => (this.open()))
    this.ctx.on('editor-close', () => (this.close()))
    this.ctx.on('editor-reply', (p) => (this.setReply(p.data, p.$el, p.scroll)))
    this.ctx.on('editor-reply-cancel', () => (this.cancelReply()))
    this.ctx.on('editor-show-loading', () => (Ui.showLoading(this.$el)))
    this.ctx.on('editor-hide-loading', () => (Ui.hideLoading(this.$el)))
    this.ctx.on('editor-notify', (f) => (this.showNotify(f.msg, f.type)))
    this.ctx.on('editor-travel', ($el) => (this.travel($el)))
    this.ctx.on('editor-travel-back', () => (this.travelBack()))
    this.ctx.on('conf-updated', () => (this.refreshUploadBtn()))
  }

  initLocalStorage () {
    const localContent = window.localStorage.getItem('ArtalkContent') || ''
    if (localContent.trim() !== '') {
      this.showNotify(this.$t('restoredMsg'), 'i')
      this.setContent(localContent)
    }
    this.$textarea.addEventListener('input', () => {
      this.saveContent()
    })
  }

  initHeader () {
    Object.keys(this.user.data).forEach((field) => {
      const inputEl = this.getInputEl(field)
      if (inputEl && inputEl instanceof HTMLInputElement) {
        inputEl.value = this.user.data[field] || ''
        // 绑定事件
        inputEl.addEventListener('input', () => this.onHeaderInput(field, inputEl))
      }
    })

    // Link URL 自动补全协议
    const $linkInput = this.getInputEl('link')
    if ($linkInput) {
      $linkInput.addEventListener('change', () => {
        const link = $linkInput.value.trim()
        if (!!link && !/^(http|https):\/\//.test(link)) {
          $linkInput.value = `https://${link}`
          this.user.data.link = $linkInput.value
          this.saveUser()
        }
      })
    }

    // i18n patch
    [['nick', '昵称'], ['email', '邮箱'], ['link', '网址']].forEach((entry) => {
      const $input = this.getInputEl(entry[0])!
      $input.placeholder = $input.placeholder.replace(entry[1], this.$t(entry[0] as any))
    })
  }

  getInputEl (field: 'nick'|'email'|'link'|string) {
    const inputEl = this.$header.querySelector<HTMLInputElement>(`[name="${field}"]`)
    return inputEl
  }

  queryUserInfo = {
    timeout: <number|null>null,
    abortFunc: <(() => void)|null>null
  }

  /** header 输入框内容变化事件 */
  onHeaderInput(field: string, inputEl: HTMLInputElement) {
    this.user.data[field] = inputEl.value.trim()

    // 若修改的是 nick or email
    if (field === 'nick' || field === 'email') {
      this.user.data.token = '' // 清除 token 登陆状态
      this.user.data.isAdmin = false

      // 获取用户信息
      if (this.queryUserInfo.timeout !== null) window.clearTimeout(this.queryUserInfo.timeout) // 清除待发出的请求
      if (this.queryUserInfo.abortFunc !== null) this.queryUserInfo.abortFunc() // 之前发出未完成的请求立刻中止

      this.queryUserInfo.timeout = window.setTimeout(() => {
        this.queryUserInfo.timeout = null // 清理

        const {req, abort} = new Api(this.ctx).userGet(
          this.user.data.nick, this.user.data.email
        )
        this.queryUserInfo.abortFunc = abort
        req.then(data => {
          if (!data.is_login) {
            this.user.data.token = ''
            this.user.data.isAdmin = false
          }

          // 未读消息更新
          this.ctx.trigger('unread-update', { notifies: data.unread, })

          // 若用户为管理员，执行登陆操作
          if (this.user.checkHasBasicUserInfo() && !data.is_login && data.user && data.user.is_admin) {
            this.showLoginDialog()
          }

          // 自动填入 link
          if (data.user && data.user.link) {
            this.user.data.link = data.user.link
            this.getInputEl('link')!.value = data.user.link
          }
        })
        .catch(() => {})
        .finally(() => {
          this.queryUserInfo.abortFunc = null // 清理
        })
      }, 400) // 延迟执行，减少请求次数
    }

    this.saveUser()
  }

  showLoginDialog () {
    this.ctx.trigger('checker-admin', {
      onSuccess: () => {
      }
    })
  }

  saveUser () {
    this.user.save()
    this.ctx.trigger('user-changed', this.ctx.user.data)
  }

  saveContent () {
    window.localStorage.setItem('ArtalkContent', this.getContentOriginal().trim())
  }

  initTextarea () {
    // 占位符
    this.$textarea.placeholder = this.ctx.conf.placeholder || this.$t('placeholder')

    // 修复按下 Tab 输入的内容
    this.$textarea.addEventListener('keydown', (e) => {
      const keyCode = e.keyCode || e.which

      if (keyCode === 9) {
        e.preventDefault()
        this.insertContent('\t')
      }
    })

    // 输入框高度随内容而变化
    this.$textarea.addEventListener('input', (evt) => {
      this.adjustTextareaHeight()
    })
  }

  adjustTextareaHeight () {
    const diff = this.$textarea.offsetHeight - this.$textarea.clientHeight
    this.$textarea.style.height = '0px' // it's a magic. 若不加此行，内容减少，高度回不去
    this.$textarea.style.height = `${this.$textarea.scrollHeight + diff}px`
  }

  openedPlugName: string|null = null

  initEditorPlug () {
    this.plugList = {}
    this.$plugWrap.innerHTML = ''
    this.$plugWrap.style.display = 'none'
    this.openedPlugName = null
    this.$plugBtnWrap.innerHTML = ''

    // 初始化 plug 按钮
    this.LOADABLE_PLUG_LIST.forEach((PlugObj) => {
      if (PlugObj.Name === 'emoticons' && !this.conf.emoticons) return

      // 切换按钮
      let btnHtml = PlugObj.BtnHTML
      btnHtml = btnHtml.replace('表情', this.$t('emoticon')).replace('预览', this.$t('preview'))
      const btnElem = Utils.createElement(`<span class="atk-plug-btn" data-plug-name="${PlugObj.Name}">${btnHtml}</span>`)
      this.$plugBtnWrap.appendChild(btnElem)

      btnElem.addEventListener('click', () => {
        let plug = this.plugList[PlugObj.Name]
        if (!plug) {
          plug = new PlugObj(this)
          this.plugList[PlugObj.Name] = plug
        }

        this.$plugBtnWrap.querySelectorAll('.active').forEach(item => item.classList.remove('active'))

        // 若点击已打开的，则收起
        if (PlugObj.Name === this.openedPlugName) {
          plug.onHide()
          this.$plugWrap.style.display = 'none'
          this.openedPlugName = null
          return
        }

        if (this.$plugWrap.querySelector(`[data-plug-name="${PlugObj.Name}"]`) === null) {
          // 需要初始化
          const plugEl = plug.getEl()
          plugEl.setAttribute('data-plug-name', PlugObj.Name)
          plugEl.style.display = 'none'
          this.$plugWrap.appendChild(plugEl)
        }

        (Array.from(this.$plugWrap.children) as HTMLElement[]).forEach((plugItemEl: HTMLElement) => {
          const plugItemName = plugItemEl.getAttribute('data-plug-name')!
          if (plugItemName === PlugObj.Name) {
            plugItemEl.style.display = ''
            this.plugList[plugItemName].onShow()
          } else {
            plugItemEl.style.display = 'none'
            this.plugList[plugItemName].onHide()
          }
        })

        this.$plugWrap.style.display = ''
        this.openedPlugName = PlugObj.Name

        btnElem.classList.add('active')
      })
    })

    // 表情包插件预加载
    if (this.conf.emoticons) {
      const emoPlug = new EmoticonsPlug(this)
      this.plugList[EmoticonsPlug.Name] = emoPlug

      window.setTimeout(() => {
        emoPlug.loadEmoticonsData()
      }, 1000) // 延迟 1s 加载
    }

    this.initImgUpload()
  }

  /** 关闭编辑器插件 */
  closePlug () {
    this.$plugWrap.innerHTML = ''
    this.$plugWrap.style.display = 'none'
    this.openedPlugName = null
  }

  /** 允许的图片格式 */
  allowImgExts = ['png', 'jpg', 'jpeg', 'gif', 'bmp', 'svg', 'webp']

  /** 初始化图片上传功能 */
  initImgUpload() {
    this.$imgUploadBtn = Utils.createElement(`<span class="atk-plug-btn">${this.$t('image')}</span>`)
    this.$plugBtnWrap.querySelector('[data-plug-name="preview"]')!.before(this.$imgUploadBtn) // 显示在预览图标之前

    this.$imgUploadInput = document.createElement('input')
    this.$imgUploadInput.type = 'file'
    this.$imgUploadInput.style.display = 'none'
    this.$imgUploadInput.accept = this.allowImgExts.map(o => `.${o}`).join(',')
    this.$imgUploadBtn.after(this.$imgUploadInput)

    // 按钮点击
    this.$imgUploadBtn.onclick = () => {
      // 选择图片
      const $input = this.$imgUploadInput!
      $input.onchange = () => {
        (async () => { // 解决阻塞 UI 问题
          if (!$input.files || $input.files.length === 0) return
          const file = $input.files[0]
          this.uploadImg(file)
        })()
      }
      $input.click() // 显示选择图片对话框
    }

    // 统一从 FileList 获取文件并上传图片方法
    const uploadFromFileList = (files?: FileList) => {
      if (!files) return

      for (let i = 0; i < files.length; i++) {
        const file = files[i]
        this.uploadImg(file)
      }
    }

    // 拖拽图片
    // @link https://developer.mozilla.org/zh-CN/docs/Web/API/HTML_Drag_and_Drop_API/File_drag_and_drop
    // 阻止浏览器的默认释放行为
    this.$textarea.addEventListener('dragover', (evt) => {
      evt.stopPropagation()
      evt.preventDefault()
    })

    this.$textarea.addEventListener('drop', (evt) => {
      const files = evt.dataTransfer?.files
      if (files?.length) {
        evt.preventDefault()
        uploadFromFileList(files)
      }
    })

    // 粘贴图片
    this.$textarea.addEventListener('paste', (evt) => {
      const files = evt.clipboardData?.files
      if (files?.length) {
        evt.preventDefault()
        uploadFromFileList(files)
      }
    })
  }

  /** 刷新图片上传按钮 */
  refreshUploadBtn() {
    if (!this.$imgUploadBtn) return

    if (!this.ctx.conf.imgUpload) {
      this.$imgUploadBtn.setAttribute('atk-only-admin-show', '')
      this.ctx.trigger('check-admin-show-el')
    }
  }

  async uploadImg(file: File) {
    const fileExt = /[^.]+$/.exec(file.name)
    if (!fileExt || !this.allowImgExts.includes(fileExt[0])) return

    // 未登录提示
    if (!this.ctx.user.checkHasBasicUserInfo()) {
      this.showNotify('填入你的名字邮箱才能上传哦', 'w')
      return
    }

    // 插入图片前换一行
    let insertPrefix = '\n'
    if (this.$textarea.value.trim() === '') insertPrefix = ''

    // 插入占位加载文字
    const uploadPlaceholderTxt = `${insertPrefix}![](Uploading ${file.name}...)`
    this.insertContent(uploadPlaceholderTxt)

    // 上传图片
    let resp: any
    try {
      resp = await new Api(this.ctx).imgUpload(file)
    } catch (err: any) {
      console.error(err)
      this.showNotify(`${this.$t('uploadFail')}，${err.msg}`, 'e')
    }
    if (!!resp && resp.img_url) {
      let imgURL = resp.img_url as string

      // 若为相对路径，加上 artalk server
      if (!Utils.isValidURL(imgURL)) imgURL = Utils.getURLBasedOnApi(this.ctx, imgURL)

      // 上传成功插入图片
      this.setContent(this.$textarea.value.replace(uploadPlaceholderTxt, `${insertPrefix}![](${imgURL})`))
    } else {
      // 上传失败删除加载文字
      this.setContent(this.$textarea.value.replace(uploadPlaceholderTxt, ''))
    }
  }

  insertContent (val: string) {
    if ((document as any).selection) {
      this.$textarea.focus();
      (document as any).selection.createRange().text = val
      this.$textarea.focus()
    } else if (this.$textarea.selectionStart || this.$textarea.selectionStart === 0) {
      const sStart = this.$textarea.selectionStart
      const sEnd = this.$textarea.selectionEnd
      const sT = this.$textarea.scrollTop
      this.setContent(this.$textarea.value.substring(0, sStart) + val + this.$textarea.value.substring(sEnd, this.$textarea.value.length))
      this.$textarea.focus()
      this.$textarea.selectionStart = sStart + val.length
      this.$textarea.selectionEnd = sStart + val.length
      this.$textarea.scrollTop = sT
    } else {
      this.$textarea.focus()
      this.$textarea.value += val
    }
  }

  setContent (val: string) {
    this.$textarea.value = val
    this.saveContent()
    if (!!this.plugList && !!this.plugList.preview) {
      this.plugList.preview.updateContent()
    }

    // 延迟执行防止无效
    window.setTimeout(() => {
      this.adjustTextareaHeight()
    }, 80)
  }

  clearEditor () {
    this.setContent('')
    this.cancelReply()
  }

  /** 获取最终用于 submit 的数据 */
  getFinalContent () {
    let content = this.getContentOriginal()

    // 表情包处理
    if (this.plugList && this.plugList.emoticons) {
      const emoticonsPlug = this.plugList.emoticons as EmoticonsPlug
      content = emoticonsPlug.transEmoticonImageText(content)
    }

    return content
  }

  getContentOriginal () {
    return this.$textarea.value || '' // Tip: !!"0" === true
  }

  getContentMarked () {
    return Utils.marked(this.ctx, this.getFinalContent())
  }

  initBottomPart () {
    this.initReply()
    this.initSubmit()
  }

  initReply () {
    this.replyComment = null
    this.$sendReply = null
  }

  setReply (commentData: CommentData, $comment: HTMLElement, scroll = true) {
    if (this.replyComment !== null) {
      this.cancelReply()
    }

    if (this.$sendReply === null) {
      this.$sendReply = Utils.createElement(`<div class="atk-send-reply">${this.$t('reply')} <span class="atk-text"></span><span class="atk-cancel">×</span></div>`);
      this.$sendReply.querySelector<HTMLElement>('.atk-text')!.innerText = `@${commentData.nick}`
      this.$sendReply.addEventListener('click', () => {
        this.cancelReply()
      })
      this.$textareaWrap.append(this.$sendReply)
    }
    this.replyComment = commentData

    if (this.ctx.conf.editorTravel === true) {
      this.travel($comment)
    }

    if (scroll) Ui.scrollIntoView(this.$el)

    this.$textarea.focus()
  }

  cancelReply () {
    if (this.$sendReply !== null) {
      this.$sendReply.remove()
      this.$sendReply = null
    }
    this.replyComment = null

    if (this.ctx.conf.editorTravel === true) {
      this.travelBack()
    }
  }

  initSubmit () {
    this.$submitBtn.innerText = this.ctx.conf.sendBtn || this.$t('send')

    this.$submitBtn.addEventListener('click', (evt) => {
      const btnEl = evt.currentTarget
      this.submit()
    })
  }

  async submit () {
    if (this.getFinalContent().trim() === '') {
      this.$textarea.focus()
      return
    }

    this.ctx.trigger('editor-submit')

    Ui.showLoading(this.$el)

    try {
      const nComment = await new Api(this.ctx).add({
        content: this.getFinalContent(),
        nick: this.user.data.nick,
        email: this.user.data.email,
        link: this.user.data.link,
        rid: (this.replyComment === null) ? 0 : this.replyComment.id,
        page_key: (this.replyComment === null) ? this.ctx.conf.pageKey : this.replyComment.page_key,
        page_title: (this.replyComment === null) ? this.ctx.conf.pageTitle : undefined,
        site_name: (this.replyComment === null) ? this.ctx.conf.site : this.replyComment.site_name
      })

      // 回复不同页面的评论
      if (this.replyComment !== null && this.replyComment.page_key !== this.ctx.conf.pageKey) {
        window.open(`${this.replyComment.page_url}#atk-comment-${nComment.id}`)
      }

      this.ctx.trigger('list-insert', nComment)
      this.clearEditor() // 清空编辑器
      this.ctx.trigger('editor-submitted')
    } catch (err: any) {
      console.error(err)
      this.showNotify(`${this.$t('commentFail')}，${err.msg || String(err)}`, 'e')
      return
    } finally {
      Ui.hideLoading(this.$el)
    }
  }

  showNotify (msg: string, type) {
    Ui.showNotify(this.$notifyWrap, msg, type)
  }

  /** 关闭评论 */
  close () {
    if (!this.$textareaWrap.querySelector('.atk-comment-closed'))
      this.$textareaWrap.prepend(Utils.createElement('<div class="atk-comment-closed">仅管理员可评论</div>'))

    if (!this.user.data.isAdmin) {
      this.$textarea.style.display = 'none'
      this.closePlug()
      this.$bottom.style.display = 'none'
    } else {
      // 管理员一直打开评论
      this.$textarea.style.display = ''
      this.$bottom.style.display = ''
    }
  }

  /** 打开评论 */
  open () {
    this.$textareaWrap.querySelector('.atk-comment-closed')?.remove()
    this.$textarea.style.display = ''
    this.$bottom.style.display = ''
  }

  travel ($afterEl: HTMLElement) {
    if (this.isTraveling) return
    this.isTraveling = true
    this.$el.after(Utils.createElement('<div class="atk-editor-travel-placeholder"></div>'))

    const $travelPlace = Utils.createElement('<div></div>')
    $afterEl.after($travelPlace)
    $travelPlace.replaceWith(this.$el)

    this.$el.classList.add('atk-fade-in') // 添加渐入动画
  }

  travelBack () {
    if (!this.isTraveling) return
    this.isTraveling = false
    this.ctx.$root.querySelector('.atk-editor-travel-placeholder')?.replaceWith(this.$el)

    // 取消回复
    if (this.replyComment !== null) this.cancelReply()
  }

  initRemoteEmoticons () {

  }
}


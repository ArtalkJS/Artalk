import * as Utils from '@/lib/utils'
import $t from '@/i18n'
import type PlugKit from './_kit'
import EditorPlug from './_plug'

/** 允许的图片格式 */
const AllowImgExts = ['png', 'jpg', 'jpeg', 'gif', 'bmp', 'svg', 'webp']

export default class Upload extends EditorPlug {
  private $imgUploadInput?: HTMLInputElement

  constructor(kit: PlugKit) {
    super(kit)

    this.kit.useMounted(() => this.init())

    this.initDragImg()
  }

  private init() {
    this.$imgUploadInput = document.createElement('input')
    this.$imgUploadInput.type = 'file'
    this.$imgUploadInput.style.display = 'none'
    this.$imgUploadInput.accept = AllowImgExts.map(o => `.${o}`).join(',')

    // TODO: Use btn cannot refresh when mounted event is triggered
    const $btn = this.useBtn(`${$t('image')}`)
    $btn.after(this.$imgUploadInput)
    $btn.onclick = () => {
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

    if (!this.kit.useConf().imgUpload) {
      this.$btn!.setAttribute('atk-only-admin-show', '')
    }
  }

  private initDragImg() {
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
    const onDragover = (evt: Event) => {
      evt.stopPropagation()
      evt.preventDefault()
    }

    const onDrop = (evt: DragEvent) => {
      const files = evt.dataTransfer?.files
      if (files?.length) {
        evt.preventDefault()
        uploadFromFileList(files)
      }
    }

    // 粘贴图片
    const onPaste = (evt: ClipboardEvent) => {
      const files = evt.clipboardData?.files
      if (files?.length) {
        evt.preventDefault()
        uploadFromFileList(files)
      }
    }

    this.kit.useMounted(() => {
      this.kit.useUI().$textarea.addEventListener('dragover', onDragover)
      this.kit.useUI().$textarea.addEventListener('drop', onDrop)
      this.kit.useUI().$textarea.addEventListener('paste', onPaste)
    })
    this.kit.useUnmounted(() => {
      this.kit.useUI().$textarea.removeEventListener('dragover', onDragover)
      this.kit.useUI().$textarea.removeEventListener('drop', onDrop)
      this.kit.useUI().$textarea.removeEventListener('paste', onPaste)
    })
  }

  async uploadImg(file: File) {
    const fileExt = /[^.]+$/.exec(file.name)
    if (!fileExt || !AllowImgExts.includes(fileExt[0])) return

    // 未登录提示
    if (!this.kit.useUser().checkHasBasicUserInfo()) {
      this.kit.useEditor().showNotify($t('uploadLoginMsg'), 'w')
      return
    }

    // 插入图片前换一行
    let insertPrefix = '\n'
    if (this.kit.useUI().$textarea.value.trim() === '') insertPrefix = ''

    // 插入占位加载文字
    const uploadPlaceholderTxt = `${insertPrefix}![](Uploading ${file.name}...)`
    this.kit.useEditor().insertContent(uploadPlaceholderTxt)

    // 上传图片
    let resp: { public_url: string } | undefined
    try {
      const customUploaderFn = this.kit.useConf().imgUploader
      if (!customUploaderFn) {
        // 使用 Artalk 进行图片上传
        resp = (await this.kit.useApi().upload.upload({ file })).data
      } else {
        // 使用自定义的图片上传器
        resp = { public_url: await customUploaderFn(file) }
      }
    } catch (err: any) {
      console.error(err)
      this.kit.useEditor().showNotify(`${$t('uploadFail')}: ${err.msg}`, 'e')
    }
    if (!!resp && resp.public_url) {
      let imgURL = resp.public_url as string

      // 若为相对路径，加上 artalk server
      if (!Utils.isValidURL(imgURL)) imgURL = Utils.getURLBasedOnApi({
        base: this.kit.useConf().server,
        path: imgURL,
      })

      // 上传成功插入图片
      this.kit.useEditor().setContent(this.kit.useUI().$textarea.value.replace(uploadPlaceholderTxt, `${insertPrefix}![](${imgURL})`))
    } else {
      // 上传失败删除加载文字
      this.kit.useEditor().setContent(this.kit.useUI().$textarea.value.replace(uploadPlaceholderTxt, ''))
    }
  }
}

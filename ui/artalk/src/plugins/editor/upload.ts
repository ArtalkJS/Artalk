import type PlugKit from './_kit'
import EditorPlugin from './_plug'
import * as Utils from '@/lib/utils'
import $t from '@/i18n'

/** 允许的图片格式 */
const AllowImgExts = ['png', 'jpg', 'jpeg', 'gif', 'bmp', 'svg', 'webp']

export default class Upload extends EditorPlugin {
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
    this.$imgUploadInput.accept = AllowImgExts.map((o) => `.${o}`).join(',')

    // TODO: Use btn cannot refresh when mounted event is triggered
    const $btn = this.useBtn(
      `<i aria-label="${$t('uploadImage')}"><svg fill="currentColor" aria-hidden="true" height="14" viewBox="0 0 14 14" width="14"><path d="m0 1.94444c0-1.074107.870333-1.94444 1.94444-1.94444h10.11116c1.0741 0 1.9444.870333 1.9444 1.94444v10.11116c0 1.0741-.8703 1.9444-1.9444 1.9444h-10.11116c-1.074107 0-1.94444-.8703-1.94444-1.9444zm1.94444-.38888c-.21466 0-.38888.17422-.38888.38888v7.06689l2.33333-2.33333 2.33333 2.33333 3.88888-3.88889 2.3333 2.33334v-5.51134c0-.21466-.1742-.38888-.3888-.38888zm10.49996 8.09977-2.3333-2.33333-3.88888 3.8889-2.33333-2.33334-2.33333 2.33334v.8447c0 .2146.17422.3888.38888.3888h10.11116c.2146 0 .3888-.1742.3888-.3888zm-7.1944-6.54422c-.75133 0-1.36111.60978-1.36111 1.36111 0 .75134.60978 1.36111 1.36111 1.36111s1.36111-.60977 1.36111-1.36111c0-.75133-.60978-1.36111-1.36111-1.36111z"/></svg></i>`,
    )
    $btn.after(this.$imgUploadInput)
    $btn.onclick = () => {
      // 选择图片
      const $input = this.$imgUploadInput!
      $input.onchange = () => {
        ;(async () => {
          // 解决阻塞 UI 问题
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
    if (!fileExt || !AllowImgExts.includes(String(fileExt[0]).toLowerCase())) return

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
      this.kit.useEditor().showNotify(`${$t('uploadFail')}: ${err.message}`, 'e')
    }
    if (!!resp && resp.public_url) {
      let imgURL = resp.public_url as string

      // 若为相对路径，加上 artalk server
      if (!Utils.isValidURL(imgURL))
        imgURL = Utils.getURLBasedOnApi({
          base: this.kit.useConf().server,
          path: imgURL,
        })

      // 上传成功插入图片
      this.kit
        .useEditor()
        .setContent(
          this.kit
            .useUI()
            .$textarea.value.replace(uploadPlaceholderTxt, `${insertPrefix}![](${imgURL})`),
        )
    } else {
      // 上传失败删除加载文字
      this.kit
        .useEditor()
        .setContent(this.kit.useUI().$textarea.value.replace(uploadPlaceholderTxt, ''))
    }
  }
}

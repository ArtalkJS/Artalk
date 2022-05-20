import * as Utils from 'artalk/src/lib/utils'
import * as Ui from 'artalk/src/lib/ui'
import SidebarView from './sidebar-view'

export default class TransferView extends SidebarView {
  protected readonly viewName = 'transfer'
  public readonly viewAdminOnly = true
  viewTitle() { return '迁移' }

  protected tabs = {
    'import': '导入',
    'export': '导出',
  }
  protected activeTab = 'import'

  mount() {
    this.switchTab('import')
  }

  switchTab(tab: string) {
    if (tab === 'import') {
      this.initImport()
    } else if (tab === 'export') {
      // if (window.confirm(`开始导出？`))
      this.initExport()
      return false
    }
    return true
  }

  initImport() {
    this.$el.innerHTML =
    `<div class="atk-log-wrap" style="display: none;">
      <div class="atk-log-back-btn">返回</div>
      <div class="atk-log"></div>
    </div>
    <div class="atk-form">
    <div class="atk-label atk-data-file-label">Artrans 数据文件</div>
    <div class="atk-file-upload-group">
      <div class="atk-file-input-wrap atk-fade-in">
        <input type="file" name="AtkDataFile" accept=".artrans">
        <span class="atk-desc">使用「<a href="https://artalk.js.org/guide/transfer.html" target="_blank">转换工具</a>」将评论数据转为 Artrans 格式</span>
      </div>
      <div class="atk-uploading-wrap atk-fade-in">
        <div class="atk-progress">
          <div class="atk-bar"></div>
        </div>
        <div class="atk-status">上传中 <span class="atk-curt">0%</span>... <span class="atk-abort">取消</span></div>
      </div>
    </div>
    <div class="atk-label">目标站点名</div>
    <input type="text" name="AtkSiteName" placeholder="输入内容..." autocomplete="off">
    <div class="atk-label">目标站点 URL</div>
    <input type="text" name="AtkSiteURL" placeholder="输入内容..." autocomplete="off">
    <div class="atk-label">启动参数（可选）</div>
    <textarea name="AtkPayload"></textarea>
    <span class="atk-desc">参考「<a href="https://artalk.js.org/guide/transfer.html" target="_blank">文档 · 数据迁移</a>」</span>
    <button class="atk-btn" name="AtkSubmit">导入</button>
    </div>`

    const $form = this.$el.querySelector<HTMLSelectElement>('.atk-form')!

    const $fileWrap = $form.querySelector<HTMLElement>('.atk-file-input-wrap')!
    const $file = $fileWrap.querySelector<HTMLInputElement>('[name="AtkDataFile"]')!
    const $fileDesc = $fileWrap.querySelector<HTMLElement>('.atk-desc')!
    const fileDescOrgHTML = $fileDesc.innerHTML
    const restoreFileInput = () => {
      $fileDesc.innerHTML = fileDescOrgHTML
      $file.value = ''
    }

    // 文件上传
    const $uploadingWrap = $form.querySelector<HTMLElement>('.atk-uploading-wrap')!
    const $uploadProgress = $uploadingWrap.querySelector<HTMLElement>('.atk-progress')!
    const $uploadProgressBar = $uploadProgress.querySelector<HTMLElement>('.atk-bar')!
    const $uploadStatus = $uploadingWrap.querySelector<HTMLElement>('.atk-status')!
    const $uploadStatusCurt = $uploadStatus.querySelector<HTMLElement>('.atk-curt')!
    const $uploadAbortBtn = $uploadStatus.querySelector<HTMLElement>('.atk-abort')!

    const $siteName = $form.querySelector<HTMLInputElement>('[name="AtkSiteName"]')!
    const $siteURL = $form.querySelector<HTMLInputElement>('[name="AtkSiteURL"]')!
    const $payload = $form.querySelector<HTMLTextAreaElement>('[name="AtkPayload"]')!
    const $submitBtn = $form.querySelector<HTMLButtonElement>('[name="AtkSubmit"]')!
    const setError = (msg: string) => window.alert(msg)

    const showUploading = () => {
      $fileWrap.style.display = 'none'
      $uploadingWrap.style.display = ''
      setUploading(0)
    }

    const hideUploading = () => {
      $fileWrap.style.display = ''
      $uploadingWrap.style.display = 'none'
    }

    const setUploading = (progress: number) => {
      $uploadProgressBar.style.width = `${progress}%`
      $uploadStatusCurt.innerText = `${progress}%`
    }

    let UploadedFilename: string = ''

    hideUploading()

    const startUploadFile = async () => {
      UploadedFilename = ''

      const xhr = new XMLHttpRequest()

      // 进度条
      xhr.upload.addEventListener('progress', (evt) => {
        if (evt.loaded === evt.total) {
          // 上传完毕
          setUploading(100)
          return
        }

        const fileSize = $file.files![0].size
        if (evt.loaded <= fileSize) {
          // 正在上传
          const percent = Math.round(evt.loaded / fileSize * 100)
          setUploading(percent)
        }
      })

      // 创建上传参数
      const formData = new FormData()
      formData.append('file', $file.files![0])
      formData.append('token', this.ctx.user.data.token)

      // 开始上传
      xhr.open('post', `${this.ctx.conf.server}/api/admin/import-upload`)
      xhr.timeout = 2*60*1000 // 2min
      xhr.send(formData)

      // 上传成功事件
      xhr.onload = () => {
        const setErr = (msg: string): void => {
          restoreFileInput()
          hideUploading()
          alert(`文件上传失败，${msg}`)
        }

        if (xhr.status !== 200) {
          setErr(`响应状态码：${xhr.status}`)
          return
        }

        let json: any
        try {
          json = JSON.parse(xhr.response)
        } catch (err) {
          console.error(err)
          setErr(`JSON 解析失败：${err}`)
          return
        }

        if (!json.success || !json.data.filename) {
          setErr(`响应：${xhr.response}`)
          return
        }

        $fileDesc.innerHTML = '文件已成功上传，可以开始导入'
        UploadedFilename = json.data.filename
        hideUploading()
      }

      // 中止上传
      $uploadAbortBtn.onclick = () => {
        xhr.abort()
        restoreFileInput()
        hideUploading()
      }
    }

    // 文件上传操作
    $file.onchange = () => {
      if (!$file.files || $file.files.length === 0) return

      showUploading()
      setTimeout(async () => {
        await startUploadFile()
        hideUploading()
      }, 80)
    }

    // 开始导入按钮
    $submitBtn.onclick = () => {
      if (!UploadedFilename) {
        setError(`请先上传 Artrans 数据文件`)
        return
      }

      const siteName = $siteName.value.trim()
      const siteURL = $siteURL.value.trim()
      const payload = $payload.value.trim()

      // 请求 payload 参数制备
      let rData: any = {}
      if (payload) {
        // JSON 格式检验
        try {
          rData = JSON.parse(payload)
        } catch (err) {
          setError(`Payload JSON 格式有误：${String(err)}`)
          return
        }

        if (typeof rData !== 'object' || Array.isArray(rData)) {
          setError(`Payload 需为 JSON 对象`)
          return
        }
      }
      if (siteName) rData.t_name = siteName
      if (siteURL) rData.t_url = siteURL

      // 创建导入会话
      const createSession = (filename?: string) => {
        const $logWrap = this.$el.querySelector<HTMLElement>('.atk-log-wrap')!
        const $log = $logWrap.querySelector<HTMLElement>('.atk-log')!
        const $backBtn = this.$el.querySelector<HTMLElement>('.atk-log-back-btn')!

        $logWrap.style.display = ''
        $form.style.display = 'none'

        $backBtn.onclick = () => {
          $logWrap.style.display = 'none'
          $form.style.display = ''
          restoreFileInput()
          UploadedFilename = ''
          this.sidebar.reload()
        }

        if (filename) rData.json_file = filename

        // 创建 iframe
        const frameName = `f_${+new Date()}`
        const $frame = document.createElement('iframe')
        $frame.className = 'atk-iframe'
        $frame.name = frameName
        $log.innerHTML = ''
        $log.append($frame)

        const formParams: {[k: string]: string} = {
          payload: JSON.stringify(rData),
          token: this.ctx.user.data.token || '',
        }

        // 创建临时表单，初始化 iframe
        const $formTmp = document.createElement('form')
        $formTmp.style.display = 'none'
        $formTmp.setAttribute('method', 'post')
        $formTmp.setAttribute('action', `${this.ctx.conf.server}/api/admin/import`)
        $formTmp.setAttribute("target", frameName)

        Object.entries(formParams).forEach(([key, val]) => {
          const $inputTmp = document.createElement('input')
          $inputTmp.setAttribute('type', 'hidden')
          $inputTmp.setAttribute('name', key)
          $inputTmp.value = val
          $formTmp.appendChild($inputTmp)
        })

        $logWrap.append($formTmp)
        $formTmp.submit()
        $formTmp.remove()
      }

      // 直接开始会话
      createSession(UploadedFilename)
    }
  }

  async initExport() {
    Ui.showLoading(this.$el)

    try {
      const d = await this.ctx.getApi().export()
      this.download(`backup-${this.getYmdHisFilename()}.artrans`, d)
    } catch (err: any) {
      console.log(err)
      window.alert(`${String(err)}`)
      return
    } finally { Ui.hideLoading(this.$el) }
  }

  download(filename: string, text: string) {
    const el = document.createElement('a')
    el.setAttribute('href', `data:text/json;charset=utf-8,${encodeURIComponent(text)}`);
    el.setAttribute('download', filename)
    el.style.display = 'none'
    document.body.appendChild(el)
    el.click()
    document.body.removeChild(el)
  }

  getYmdHisFilename() {
    const date = new Date()

    const year = date.getFullYear()
    const month = date.getMonth() + 1
    const day = date.getDate()
    const hours = date.getHours()
    const minutes = date.getMinutes()
    const seconds = date.getSeconds()

    return `${year}${month}${day}-${hours}${Utils.padWithZeros(minutes, 2)}${Utils.padWithZeros(seconds, 2)}`
  }
}

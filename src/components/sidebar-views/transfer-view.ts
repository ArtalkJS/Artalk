import Api from '~/src/api'
import Context from '../../context'
import Component from '../../lib/component'
import * as Utils from '../../lib/utils'
import * as Ui from '../../lib/ui'
import SiteList from '../admin/site-list'
import Comment from '../comment'
import SidebarView from '../sidebar-view'

export default class TransferView extends SidebarView {ß
  static viewName = 'transfer'
  static viewTitle = '迁移'
  static viewAdminOnly = true

  viewTabs = {
    'import': '导入',
    'export': '导出',
  }
  viewActiveTab = 'import'

  mount(siteName: string) {
    this.switchTab('import', siteName)
  }

  switchTab(tab: string, siteName: string) {
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
    <div class="atk-label">数据类型</div>
    <select name="AtkDataType">
      <option value="artrans">Artrans (数据行囊)</option>
      <option value="artalk_v1">Artalk v1 (PHP 旧版)</option>
      <option value="typecho">Typecho</option>
      <option value="wordpress">WordPress</option>
      <option value="disqus">Disqus</option>
      <option value="commento">Commento</option>
      <option value="valine">Valine</option>
      <option value="twikoo">Twikoo</option>
    </select>
    <div class="atk-label atk-data-file-label">数据文件</div>
    <input type="file" name="AtkDataFile" accept="text/plain,.json">
    <div class="atk-label">目标站点名</div>
    <input type="text" name="AtkSiteName" placeholder="输入内容..." autocomplete="off">
    <div class="atk-label">目标站点 URL</div>
    <input type="text" name="AtkSiteURL" placeholder="输入内容..." autocomplete="off">
    <div class="atk-label">启动参数（可选）</div>
    <textarea name="AtkPayload"></textarea>
    <span class="atk-desc">启动参数查阅：“<a href="https://artalk.js.org/guide/transfer.html" target="_blank">文档 · 数据搬家</a>”</span>
    <button class="atk-btn" name="AtkSubmit">导入</button>
    </div>`

    const $form = this.$el.querySelector<HTMLSelectElement>('.atk-form')!
    const $dataType = $form.querySelector<HTMLSelectElement>('[name="AtkDataType"]')!
    const $dataFile = $form.querySelector<HTMLInputElement>('[name="AtkDataFile"]')!
    const $dataFileLabel = $form.querySelector<HTMLInputElement>('.atk-data-file-label')!
    const $siteName = $form.querySelector<HTMLInputElement>('[name="AtkSiteName"]')!
    const $siteURL = $form.querySelector<HTMLInputElement>('[name="AtkSiteURL"]')!
    const $payload = $form.querySelector<HTMLTextAreaElement>('[name="AtkPayload"]')!
    const $submitBtn = $form.querySelector<HTMLButtonElement>('[name="AtkSubmit"]')!
    const setError = (msg: string) => window.alert(msg)

    $dataType.onchange = () => {
      if (['typecho'].includes($dataType.value)) {
        $dataFile.style.display = 'none'
        $dataFileLabel.style.display = 'none'
      } else {
        $dataFile.style.display = ''
        $dataFileLabel.style.display = ''
      }
    }

    $submitBtn.onclick = () => {
      const dataType = $dataType.value.trim()
      const siteName = $siteName.value.trim()
      const siteURL = $siteURL.value.trim()
      const payload = $payload.value.trim()

      if (dataType === '') {
        setError('请选择数据类型')
        return
      }

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

        if (rData !instanceof Object) {
          setError(`Payload 需为 JSON 对象`)
          return
        }
      }
      if (siteName) rData.t_name = siteName
      if (siteURL) rData.t_url = siteURL

      // 创建导入会话
      const createSession = (dataStr?: string) => {
        const $logWrap = this.$el.querySelector<HTMLElement>('.atk-log-wrap')!
        const $log = $logWrap.querySelector<HTMLElement>('.atk-log')!
        const $backBtn = this.$el.querySelector<HTMLElement>('.atk-log-back-btn')!

        $logWrap.style.display = ''
        $form.style.display = 'none'

        $backBtn.onclick = () => {
          $logWrap.style.display = 'none'
          $form.style.display = ''
        }

        if (dataStr) rData.json_data = dataStr

        // 创建 iframe
        const frameName = `f_${+new Date()}`
        const $frame = document.createElement('iframe')
        $frame.className = 'atk-iframe'
        $frame.name = frameName
        $log.innerHTML = ''
        $log.append($frame)

        const formParams: {[k: string]: string} = {
          type: dataType,
          payload: JSON.stringify(rData),
          token: this.ctx.user.data.token || '',
        }

        // 创建临时表单，初始化 iframe
        const $formTmp = document.createElement('form')
        $formTmp.style.display = 'none'
        $formTmp.setAttribute('method', 'post')
        $formTmp.setAttribute('action', `${this.ctx.conf.server}/admin/import`)
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

      const reader = new FileReader()
      reader.onload = () => {
        const data = String(reader.result)
        createSession(data)
      }

      // 是否已选择文件
      if ($dataFile.files?.length) {
        // 先读取文件
        reader.readAsText($dataFile.files[0])
      } else {
        // 直接开始会话
        createSession()
      }
    }
  }

  async initExport() {
    Ui.showLoading(this.$el)

    try {
      const d = await new Api(this.ctx).export()
      this.download(`artrans-${this.getYmdHisFilename()}.json`, d)
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

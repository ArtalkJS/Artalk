<script setup lang="ts">
import { useNavStore } from '../stores/nav'
import { useUserStore } from '../stores/user'
import { artalk } from '../global'
import type { SiteData } from 'artalk/types/artalk-data'
import { storeToRefs } from 'pinia';

const nav = useNavStore()
const user = useUserStore()
const { curtTab } = storeToRefs(nav)

const importParams = ref({
  siteName: '',
  siteURL: '',
  payload: ''
})

const isLoading = ref(false)

const uploadApiURL = ref('')
const importTaskApiURL = ref('')

const uploadedFilename = ref('')

const importTaskStarted = ref(false)
const importTaskParams = ref<{[k:string]:string}>({})

const exportTaskStarted = ref(false)

onMounted(() => {
  nav.updateTabs({
    'import': '导入',
    'export': '导出',
  }, 'import')
  watch(curtTab, (tab) => {
    if (tab === 'export') {
      startExportTask()
      curtTab.value = 'import'
    }
  })

  uploadApiURL.value = `${artalk?.ctx.conf.server}/api/admin/import-upload`
  importTaskApiURL.value = `${artalk?.ctx.conf.server}/api/admin/import`
})

function setError(msg: string) {
  window.alert(msg)
}

function fileUploaded(filename: string) {
  uploadedFilename.value = filename
}

function startImportTask() {
  if (!uploadedFilename.value) {
    setError(`请先上传 Artrans 数据文件`)
    return
  }

  const p = importParams.value
  const siteName = p.siteName.trim()
  const siteURL = p.siteURL.trim()
  const payload = p.payload.trim()

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
  rData.json_file = uploadedFilename.value

  // 创建导入会话
  importTaskParams.value = {
    payload: JSON.stringify(rData),
    token: user.token,
  }
  importTaskStarted.value = true
}

function importTaskDone() {
  importTaskStarted.value = false
}

async function startExportTask() {
  if (exportTaskStarted.value) return
  exportTaskStarted.value = true
  isLoading.value = true
  try {
    const data = await artalk!.ctx.getApi().site.export()
    downloadFile(`backup-${getYmdHisFilename()}.artrans`, data)
  } catch (err: any) {
    console.log(err)
    window.alert(`${String(err)}`)
    return
  } finally {
    exportTaskStarted.value = false
    isLoading.value = false
  }
}

function downloadFile(filename: string, text: string) {
  const el = document.createElement('a')
  el.setAttribute('href', `data:text/json;charset=utf-8,${encodeURIComponent(text)}`);
  el.setAttribute('download', filename)
  el.style.display = 'none'
  document.body.appendChild(el)
  el.click()
  document.body.removeChild(el)
}

function getYmdHisFilename() {
  const date = new Date()

  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hours = date.getHours()
  const minutes = date.getMinutes()
  const seconds = date.getSeconds()

  return `${year}${month}${day}-${hours}${padWithZeros(minutes, 2)}${padWithZeros(seconds, 2)}`
}

function padWithZeros(vNumber: number, width: number) {
  let numAsString = vNumber.toString()
  while (numAsString.length < width) {
    numAsString = `0${numAsString}`
  }
  return numAsString
}
</script>

<template>
  <LoadingLayer v-if="isLoading" />
  <LogTerminal v-if="importTaskStarted" :api-url="importTaskApiURL" :req-params="importTaskParams" @back="importTaskDone()" />
  <div v-show="!importTaskStarted" class="atk-form">
    <div class="atk-label atk-data-file-label">Artrans 数据文件</div>
    <FileUploader :api-url="uploadApiURL" @done="fileUploaded">
      <template v-slot:tip>
        使用「<a href="https://artalk.js.org/guide/transfer.html" target="_blank">转换工具</a>」将评论数据转为 Artrans 格式
      </template>
      <template v-slot:done-msg>
        文件已成功上传，可以开始导入
      </template>
    </FileUploader>
    <div class="atk-label">目标站点名</div>
    <input
      type="text"
      name="AtkSiteName"
      placeholder="输入内容..."
      autocomplete="off"
      v-model="importParams.siteName"
    />
    <div class="atk-label">目标站点 URL</div>
    <input
      type="text"
      name="AtkSiteURL"
      placeholder="输入内容..."
      autocomplete="off"
      v-model="importParams.siteURL"
    />
    <div class="atk-label">启动参数（可选）</div>
    <textarea name="AtkPayload" v-model="importParams.payload"></textarea>
    <span class="atk-desc">
      参考「<a href="https://artalk.js.org/guide/transfer.html" target="_blank">文档 · 数据迁移</a>」
    </span>
    <button class="atk-btn" name="AtkSubmit" @click="startImportTask()">导入</button>
  </div>
</template>

<style scoped lang="scss">

</style>

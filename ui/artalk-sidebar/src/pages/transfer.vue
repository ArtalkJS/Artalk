<script setup lang="ts">
import { useNavStore } from '../stores/nav'
import { useUserStore } from '../stores/user'
import { artalk } from '../global'
import { storeToRefs } from 'pinia'

const nav = useNavStore()
const user = useUserStore()
const { curtTab } = storeToRefs(nav)
const { t } = useI18n()

const importParams = ref({
  siteName: '',
  siteURL: '',
  payload: '',
})

const isLoading = ref(false)

const uploadApiURL = ref('')
const importTaskApiURL = ref('')

const uploadedFilename = ref('')

const importTaskStarted = ref(false)
const importTaskParams = ref<Record<string, string>>({})

const exportTaskStarted = ref(false)

onMounted(() => {
  nav.updateTabs(
    {
      import: 'import',
      export: 'export',
    },
    'import',
  )
  watch(curtTab, (tab) => {
    if (tab === 'export') {
      startExportTask()
      curtTab.value = 'import'
    }
  })

  uploadApiURL.value = `${artalk?.ctx.conf.server}/api/v2/transfer/upload`
  importTaskApiURL.value = `${artalk?.ctx.conf.server}/api/v2/transfer/import`
})

function setError(msg: string) {
  window.alert(msg)
}

function fileUploaded(filename: string) {
  uploadedFilename.value = filename
}

function startImportTask() {
  if (!uploadedFilename.value) {
    setError(`Please upload a data file first`)
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
      setError(`Payload JSON invalid: ${err}`)
      return
    }

    if (typeof rData !== 'object' || Array.isArray(rData)) {
      setError(`Payload should be an object`)
      return
    }
  }
  if (siteName) rData.target_site_name = siteName
  if (siteURL) rData.target_site_url = siteURL
  rData.json_file = uploadedFilename.value

  // 创建导入会话
  importTaskParams.value = {
    ...rData,
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
    const res = await artalk!.ctx.getApi().transfer.exportArtrans()
    downloadFile(`backup-${getYmdHisFilename()}.artrans`, res.data.artrans)
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
  el.setAttribute('href', `data:text/json;charset=utf-8,${encodeURIComponent(text)}`)
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

const artransferToolHint = computed(() =>
  t('artransferToolHint', { link: '__LINK__' }).replace(
    '__LINK__',
    `<a href="https://artalk.js.org/guide/transfer.html" target="_blank">${t('artransfer')}</a>`,
  ),
)
</script>

<template>
  <LoadingLayer v-if="isLoading" />
  <LogTerminal
    v-if="importTaskStarted"
    :api-url="importTaskApiURL"
    :req-params="importTaskParams"
    @back="importTaskDone()"
  />
  <div v-show="!importTaskStarted" class="atk-form">
    <div class="atk-label atk-data-file-label">Artrans {{ t('dataFile') }}</div>
    <FileUploader :api-url="uploadApiURL" @done="fileUploaded">
      <template v-slot:tip>
        <span v-html="artransferToolHint" />
      </template>
      <template v-slot:done-msg>
        {{ t('uploadReadyToImport') }}
      </template>
    </FileUploader>
    <div class="atk-label">{{ t('targetSiteName') }}</div>
    <input
      type="text"
      name="AtkSiteName"
      :placeholder="t('inputHint')"
      autocomplete="off"
      v-model="importParams.siteName"
    />
    <div class="atk-label">{{ t('targetSiteURL') }}</div>
    <input
      type="text"
      name="AtkSiteURL"
      :placeholder="t('inputHint')"
      autocomplete="off"
      v-model="importParams.siteURL"
    />
    <div class="atk-label">{{ t('payload') }} ({{ t('optional') }})</div>
    <textarea name="AtkPayload" v-model="importParams.payload"></textarea>
    <span class="atk-desc">
      <a href="https://artalk.js.org/guide/transfer.html" target="_blank">
        {{ t('moreDetails') }}
      </a>
    </span>
    <button class="atk-btn" name="AtkSubmit" @click="startImportTask()">
      {{ t('import') }}
    </button>
  </div>
</template>

<style scoped lang="scss"></style>

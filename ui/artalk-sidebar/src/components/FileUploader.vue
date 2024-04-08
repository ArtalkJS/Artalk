<script setup lang="ts">
import { useUserStore } from '../stores/user'

const user = useUserStore()

const props = defineProps<{
  apiUrl: string
}>()

const emit = defineEmits<{
  /** 上传完毕 */
  (evt: 'done', filename: string): void
}>()

const { apiUrl } = toRefs(props)

let xhr: XMLHttpRequest | null = null
const fileInputEl = ref<HTMLInputElement | null>(null)
const remoteFilename = ref('')
const isUploading = ref(false)
const progress = ref(0)
const isDone = ref(false)

function reset() {
  remoteFilename.value = ''
  isUploading.value = false
  progress.value = 0
  isDone.value = false
  if (fileInputEl.value) fileInputEl.value.value = ''
}

async function startUploadFile(file: File) {
  remoteFilename.value = ''

  xhr = new XMLHttpRequest()

  // 进度条
  xhr.upload.addEventListener('progress', (evt) => {
    if (evt.loaded === evt.total) {
      // 上传完毕
      progress.value = 100
      return
    }

    const fileSize = file.size
    if (evt.loaded <= fileSize) {
      // 正在上传
      progress.value = Math.round((evt.loaded / fileSize) * 100)
    }
  })

  // 创建上传参数
  const formData = new FormData()
  formData.append('file', file)
  formData.append('token', user.token)

  // 开始上传
  xhr.open('post', apiUrl.value)
  xhr.timeout = 5 * 60 * 1000 // 5分钟超时
  xhr.send(formData)

  // 上传成功事件
  xhr.onload = () => {
    const setErr = (msg: string): void => {
      reset()
      isUploading.value = false
      alert(`File upload failed: ${msg}`)
    }

    if (!xhr) {
      setErr('xhr instance is null')
      return
    }

    const ok = xhr.status >= 200 && xhr.status <= 299
    if (!ok) {
      setErr(`Response HTTP Code: ${xhr.status}, Body: ${xhr.response}`)
      return
    }

    let json: any
    try {
      json = JSON.parse(xhr.response)
    } catch (err) {
      console.error(err)
      setErr(`JSON parse error: ${err}`)
      return
    }

    if (!json.filename) {
      setErr(`Response filename is empty: ${xhr.response}`)
      return
    }

    isDone.value = true
    remoteFilename.value = json.filename
    isUploading.value = false
    emit('done', remoteFilename.value)
  }
}

function onFileInputChange() {
  const files = fileInputEl.value?.files
  if (!files || files.length === 0) return

  isUploading.value = true
  setTimeout(async () => {
    await startUploadFile(files[0])
    isUploading.value = false
  }, 80)
}

function abortUpload() {
  xhr?.abort()
  reset()
  isUploading.value = false
}
</script>

<template>
  <div class="atk-file-upload-group">
    <div v-show="!isUploading" class="atk-file-input-wrap atk-fade-in">
      <input
        ref="fileInputEl"
        type="file"
        name="AtkDataFile"
        accept=".artrans"
        @change="onFileInputChange()"
      />
      <div class="atk-desc">
        <slot v-if="!isDone" name="tip"></slot>
        <slot v-if="isDone" name="done-msg"></slot>
      </div>
    </div>
    <div v-show="isUploading" class="atk-uploading-wrap atk-fade-in">
      <div class="atk-progress">
        <div class="atk-bar" :style="{ width: `${progress}%` }"></div>
      </div>
      <div class="atk-status">
        上传中
        <span class="atk-curt">{{ progress }}%</span>
        ...
        <span class="atk-abort" @click="abortUpload()">取消</span>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss"></style>

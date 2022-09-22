<script setup lang="ts">
const props = defineProps<{
  apiUrl: string
  reqParams: { [k: string]: string }
}>()

const emit = defineEmits<{
  (evt: 'back'): void
}>()

const logWrapEl = ref<HTMLElement|null>(null)

onMounted(() => {
  // 创建 iframe
  const frameName = `f_${+new Date()}`
  const $frame = document.createElement('iframe')
  $frame.className = 'atk-iframe'
  $frame.name = frameName
  logWrapEl.value!.append($frame)

  // 创建临时表单，初始化 iframe
  const $formTmp = document.createElement('form')
  $formTmp.style.display = 'none'
  $formTmp.setAttribute('method', 'post')
  $formTmp.setAttribute('action', props.apiUrl)
  $formTmp.setAttribute('target', frameName)

  Object.entries(props.reqParams).forEach(([key, val]) => {
    const $inputTmp = document.createElement('input')
    $inputTmp.setAttribute('type', 'hidden')
    $inputTmp.setAttribute('name', key)
    $inputTmp.value = val
    $formTmp.appendChild($inputTmp)
  })

  logWrapEl.value!.append($formTmp)
  $formTmp.submit()
  $formTmp.remove()
})

function back() {
  emit('back')
}
</script>

<template>
  <div class="atk-log-wrap">
    <div class="atk-log-back-btn" @click="back()">返回</div>
    <div ref="logWrapEl" class="atk-log"></div>
  </div>
</template>

<style scoped lang="scss">
.atk-log-wrap {
  .atk-log-back-btn {
    display: inline-block;
    padding: 5px 33px;
    cursor: pointer;
    user-select: none;
    border-right: 1px solid var(--at-color-border);
    border-left: 1px solid transparent;
    &:hover {
      background: #f4f4f4;
    }
  }

  .atk-log {}

  .atk-iframe {
    width: 100%;
    height: calc(100vh - 150px);
    border: 0;
    background: #f4f4f4;
    border: 3px solid #eee;
  }
}
</style>

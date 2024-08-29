<script setup lang="ts">
import { useNavStore } from '@/stores/nav'

const props = defineProps<{
  apiUrl: string
  reqParams: { [k: string]: string }
}>()

const emit = defineEmits<{
  (evt: 'back'): void
}>()

const { t } = useI18n()

const logWrapEl = ref<HTMLElement | null>(null)

onMounted(() => {
  // Create iframe element
  const frameName = `f_${+new Date()}`
  const $frame = document.createElement('iframe')
  $frame.className = 'atk-iframe'
  $frame.name = frameName
  logWrapEl.value!.append($frame)

  // on iframe done
  $frame.onload = () => {
    useNavStore().refreshSites()
  }

  // Crate temporary form for submitting and load iframe page
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
    <div class="atk-log-back-btn" @click="back()">{{ t('back') }}</div>
    <div ref="logWrapEl" class="atk-log"></div>
  </div>
</template>

<style scoped lang="scss">
.atk-log-wrap {
  margin-bottom: -40px;

  .atk-log-back-btn {
    display: inline-block;
    padding: 5px 33px;
    cursor: pointer;
    user-select: none;
    border-right: 1px solid var(--at-color-border);
    border-left: 1px solid transparent;
    &:hover {
      background: var(--at-color-bg-grey);
    }
  }

  .atk-log {
  }

  .atk-iframe {
    width: 100%;
    height: calc(100vh - 150px);
    border: 0;
    background: var(--at-color-bg-grey);
    border: 3px solid #eee;
  }
}
</style>

<script setup lang="ts">
const props = defineProps<{
  transparentBg?: boolean
  timeout?: number
}>()

const showSpinner = ref(false)
let timer: number | undefined

onMounted(() => {
  // spinner 延迟显示，若加载等待时间太短，没必要显示（闪一下即可）
  timer = window.setTimeout(() => {
    showSpinner.value = true
    timer = undefined
  }, props.timeout || 700)
})

onUnmounted(() => {
  window.clearTimeout(timer)
})
</script>

<template>
  <div
    class="atk-loading atk-fade-in"
    :style="{ background: props.transparentBg ? 'transparent' : undefined }"
  >
    <div v-if="showSpinner" class="atk-loading-spinner">
      <svg viewBox="25 25 50 50">
        <circle cx="50" cy="50" r="20" fill="none" stroke-width="2" stroke-miterlimit="10"></circle>
      </svg>
    </div>
  </div>
</template>

<style scoped lang="scss"></style>

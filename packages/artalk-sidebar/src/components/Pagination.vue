<script setup lang="ts">
const props = defineProps<{
  /** 总页数 */
  total: number

  /** 每页数量 */
  pageSize: number

  /** 加载中 */
  isLoading?: boolean
}>()

const { pageSize, total, isLoading } = toRefs(props)

const emit = defineEmits<{
  /** 翻页事件 */
  (evt: 'change', offset: number): void
}>()

const curtPage = ref(1)
const offset = computed(() => pageSize.value * (curtPage.value - 1))
const maxPage = computed(() => Math.ceil(total.value / pageSize.value))
const prevDisabled = computed(() => (curtPage.value - 1 < 1))
const nextDisabled = computed(() => (curtPage.value + 1 > maxPage.value))

const inputValue = ref(String(curtPage.value))
let inputTimer: number|undefined

function changePage(page: number) {
  curtPage.value = page
  emit('change', offset.value)
  fillInput(page)
}

/* 按钮点击操作 */
function prev() {
  const page = curtPage.value - 1
  if (page < 1) { return }
  changePage(page)
}

function next() {
  const page = curtPage.value + 1
  if (page > maxPage.value) { return }
  changePage(page)
}

function reset() {
  curtPage.value = 1
  fillInput(1)
}

/* 页码输入框操作 */
function fillInput(page: number) {
  inputValue.value = String(page)
}

function revokeInput() {
  fillInput(curtPage.value)
}

function triggerInput(now: boolean = false) {
  window.clearTimeout(inputTimer)

  const value = inputValue.value.trim()

  const modify = () => {
    if (value === '') { revokeInput();return }
    let page = Number(value)
    if (Number.isNaN(page)) { revokeInput();return }
    if (page < 1) { revokeInput();return }
    if (page > maxPage.value) { page = maxPage.value }
    changePage(page)
  }

  // 延迟 800ms 执行
  if (!now) inputTimer = window.setTimeout(() => modify(), 800)
  else modify()
}

function onInputKeydown(evt: KeyboardEvent) {
  const keyCode = evt.keyCode || evt.which

  if (keyCode === 38) {
    // 上键
    const page = Number(inputValue.value) + 1
    if (page > maxPage.value) { return }
    fillInput(page);
    triggerInput(false)
  } else if (keyCode === 40) {
    // 下键
    const page = Number(inputValue.value) - 1
    if (page < 1) { return }
    fillInput(page)
    triggerInput(false)
  } else if (keyCode === 13) {
    // 回车
    triggerInput(true)
  }
}

defineExpose({ prev, next, reset })
</script>

<template>
  <div class="atk-pagination-wrap">
    <div class="atk-pagination">
      <div
        class="atk-btn atk-btn-prev"
        :class="{ 'atk-disabled': prevDisabled }"
        @click="prev()"
      >Prev</div>
      <input
        type="text"
        class="atk-input"
        aria-label="Enter the number of page"
        v-model="inputValue"
        @input="triggerInput(false)"
        @keydown="onInputKeydown"
      />
      <div
        class="atk-btn atk-btn-next"
        :class="{ 'atk-disabled': nextDisabled }"
        @click="next()"
      >Next</div>
    </div>
    <LoadingLayer v-if="isLoading" />
  </div>
</template>

<style scoped lang="scss">
.atk-pagination {
  display: flex;
  flex-direction: row;
  justify-content: center;
  padding: 10px 0;
  position: relative;

  & > .atk-btn,
  & > .atk-input {
    font-size: 15px;
    height: 30px;
    border: 1px solid var(--at-color-border);
    border-radius: 3px;
    padding: 0 5px;
    text-align: center;
    background: var(--at-color-bg);
  }

  & > .atk-btn {
    user-select: none;
    width: 70px;
    cursor: pointer;
    display: flex;
    justify-content: center;
    align-items: center;

    &:hover {
      background: var(--at-color-bg-grey);
    }

    &.atk-disabled {
      color: var(--at-color-sub);

      &:hover {
        cursor: default;
        background: initial;
      }
    }
  }

  & > .atk-input {
    background: transparent;
    color: var(--at-color-font);
    font-size: 18px;
    width: 60px;
    outline: none;

    &:focus {
      border-color: var(--at-color-main);
    }
  }

  & > * {
    &:not(:last-child) {
      margin-right: 10px;
    }
  }
}
</style>

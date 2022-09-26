<script setup lang="ts">
const props = defineProps<{
  initValue?: string
  placeholder?: string
  validator?: (value: string) => boolean
}>()

const emit = defineEmits<{
  (evt: 'yes', value: string): boolean|void|Promise<boolean|void>
  (evt: 'no', value: string): boolean|void|Promise<boolean|void>
  (evt: 'close'): void
}>()

const inputEl = ref<HTMLInputElement|null>(null)
const inputVal = ref('')
const inputInvalid = ref(false)

onMounted(() => {
  inputVal.value = props.initValue || ''
  window.setTimeout(() => inputEl.value?.focus(), 80)
})

function onInput() {
  // 验证器
  if (props.validator) {
    inputInvalid.value = props.validator(inputVal.value)
  }
}

function onKeyUp(evt: KeyboardEvent) {
  if (evt.key === 'Enter' || evt.keyCode === 13) { // 按下回车键
    evt.preventDefault()
    submit('yes')
  }
}

async function submit(type: 'yes'|'no') {
  if (type == 'yes' && inputInvalid.value) return

  let isContinue: any = undefined

  const callback = emit(type as any, inputVal.value)
  if (callback instanceof (async () => {}).constructor) {
    isContinue = await callback
  } else {
    isContinue = callback
  }

  if (isContinue === undefined || isContinue === true) {
    emit('close')
  }
}
</script>

<template>
  <div class="atk-item-text-editor-layer">
    <div class="atk-edit-form">
      <input
        ref="inputEl"
        class="atk-main-input"
        type="text"
        :placeholder="props.placeholder || '输入内容...'"
        autocomplete="off"
        v-model="inputVal"
        @input="onInput()"
        @keyup="onKeyUp"
        :class="{ 'atk-invalid': inputInvalid }"
      >
    </div>
    <div class="atk-actions">
      <div
        class="atk-item atk-yes-btn"
        :class="{ 'atk-disabled': inputInvalid }"
        @click="submit('yes')"
      >
        <i class="atk-icon atk-icon-yes"></i>
      </div>
      <div
        class="atk-item atk-no-btn"
        @click="submit('no')"
      >
        <i class="atk-icon atk-icon-no"></i>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.atk-item-text-editor-layer {
  z-index: 999;
  background: var(--at-color-bg);
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: row;
  align-items: center;

  .atk-edit-form {
    flex: auto;
    padding-left: 20px;

    input {
      font-size: 17px;
      width: 100%;
      padding: 3px 5px;
      border: 0;
      border-bottom: 1px solid var(--at-color-border);
      outline: none;
      background: transparent;

      &.atk-invalid {}

      &:focus {
        border-bottom-color: var(--at-color-main);
      }
    }
  }

  .atk-actions {
    @extend .atk-list-btn-actions;

    .atk-yes-btn.atk-disabled {}
  }
}
</style>

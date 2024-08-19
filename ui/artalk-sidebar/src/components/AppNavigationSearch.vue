<script lang="ts" setup>
import type { SearchStateApi } from './AppNavigationMenu'
import { useNavStore } from '@/stores/nav'

const { t } = useI18n()
const nav = useNavStore()
const inputEl = ref<HTMLInputElement | null>(null)

const props = defineProps<{
  state: SearchStateApi
}>()

watch(
  () => props.state.show,
  (show) => {
    if (show) {
      nextTick(() => {
        inputEl.value?.focus()
      })
    } else {
      props.state.updateValue('')
      nav.searchResetEvent?.()
    }
  },
)

const onSearchSubmit = (evt: Event) => {
  const { value } = props.state
  if (value.trim() === '') {
    inputEl.value?.focus()
    return
  }

  nav.searchEvent?.(value.trim())
}
</script>

<template>
  <form v-if="props.state.show" class="search-layer atk-fade-in" @submit.prevent="onSearchSubmit">
    <div class="item back-btn" @click="props.state.hideSearch()">
      <div class="icon arrow"></div>
    </div>
    <input
      ref="inputEl"
      :value="props.state.value"
      type="text"
      :placeholder="t('searchHint')"
      required
      @input="props.state.updateValue(($event.target as HTMLInputElement).value)"
    />
    <button type="submit" class="item search-btn"></button>
  </form>
</template>

<style lang="scss" scoped>
.search-layer {
  z-index: 100;
  position: absolute;
  height: 100%;
  width: 100%;
  top: 0;
  background: var(--at-color-bg);
  display: flex;
  flex-direction: row;

  .back-btn {
    padding: 0 20px;
    border-right: 1px solid var(--at-color-border);
  }

  input {
    flex: 1;
    border: 0;
    outline: none;
    padding: 0 20px;
    font-size: 15px;
    background: transparent;
  }

  .search-btn {
    cursor: pointer;
    border: none;
    background: transparent;
  }

  .icon {
    width: 17px;
    height: 100%;
    background-color: var(--at-color-deep);
    background-size: contain;
    background-repeat: no-repeat;
    mask-repeat: no-repeat;
    mask-position: center;
    mask-size: 100%;

    &.arrow {
      $arrowImg: url("data:image/svg+xml,%3Csvg width='50' height='50' viewBox='0 0 50 50' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath fill-rule='evenodd' clip-rule='evenodd' d='M32.1433 7L15 24.4999L15.4676 24.9773L32.1433 42L34 40.1047L18.7133 24.5001L34 8.89529L32.1433 7Z' fill='black'/%3E%3C/svg%3E%0A");
      mask-image: $arrowImg;
    }
  }

  .item {
    padding: 0 20px;
    color: var(--at-color-sub);
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    cursor: pointer;
    white-space: nowrap;

    &.active {
      color: var(--at-color-deep);
    }

    &:hover {
      background: var(--at-color-bg-grey);
    }

    &.search-btn {
      margin-left: auto;
      padding: 0 20px;
      border-left: 1px solid var(--at-color-border);

      &::after {
        display: inline-block;
        content: '';
        height: 15px;
        width: 15px;
        background-repeat: no-repeat;
        background-position: center;
        background-image: url("data:image/svg+xml,%3Csvg width='15' height='15' viewBox='0 0 15 15' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath fill-rule='evenodd' clip-rule='evenodd' d='M9.89266 3.70572C8.18412 1.99718 5.41402 1.99718 3.70548 3.70572C1.99694 5.41427 1.99694 8.18436 3.70548 9.89291C5.41402 11.6015 8.18412 11.6015 9.89266 9.89291C11.6012 8.18436 11.6012 5.41427 9.89266 3.70572ZM2.8216 2.82184C5.0183 0.625142 8.57985 0.625142 10.7765 2.82184C12.8239 4.86916 12.9631 8.10202 11.1942 10.3106L13.4282 12.5446L12.5443 13.4284L10.3103 11.1945C8.10177 12.9633 4.86892 12.8241 2.8216 10.7768C0.624897 8.58009 0.624897 5.01854 2.8216 2.82184Z' fill='%23757575'/%3E%3C/svg%3E");
      }
    }
  }
}
</style>

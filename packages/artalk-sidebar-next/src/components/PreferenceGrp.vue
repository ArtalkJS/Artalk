<script setup lang="ts">
import { useNavStore } from '../stores/nav'

const nav = useNavStore()

const props = defineProps<{
  pfKey: string
  title: string
  subTitle: string
  level: number
  expand: boolean
}>()

const emits = defineEmits<{
  (evt: 'toggle', key?: string): void
}>()

function onHeadClick(evt: Event) {
  if (props.level !== 1) return
  emits('toggle', props.pfKey)
  nextTick(() => {
    nav.scrollToEl(evt.target as HTMLElement)
  })
}
</script>

<template>
  <div class="pf-grp" :class="[`level-${level}`, expand ? 'expand' : '']">
    <div class="pf-head" @click="onHeadClick">
      <div class="title">{{ props.title }}</div>
      <div v-if="!!props.subTitle" class="sub-title">{{ props.subTitle }}</div>
    </div>
    <div v-show="expand" class="pf-body">
      <slot />
    </div>
  </div>
</template>

<style scoped lang="scss">
.pf-grp.level-1 {
  & > .pf-head {
    margin-top: 30px;
    margin-bottom: 20px;
    cursor: pointer;

    .title {
      font-size: 1.4em;
      font-weight: bold;
      padding-left: 12px;

      &::before {
        position: absolute;
        top: 50%;
        transform: translateY(-50%);
        transition: height ease .2s;
        left: -10px;
        content: '';
        height: 10px;
        width: 10px;
        background: #B7DCFF;
      }
    }
  }

  & > .pf-body {
  }

  &.expand > .pf-head .title::before {
    height: 25px;
  }
}

.pf-grp.level-2 {
  & > .pf-head {
    margin-top: 15px;
    margin-bottom: 20px;

    .title {
      position: relative;
      font-weight: bold;
      font-size: 1.2em;
    }
  }
}

.pf-grp.level-1, .pf-grp.level-2 {
  & > .pf-head > .sub-title {
    padding: 0 10px 0 10px;
    margin-left: 4px;
    margin-top: 15px;
    border-left: 2px solid #eee;
  }
}


.pf-grp.level-3 {
  margin-left: 25px;

  & > .pf-head {
    margin-top: 15px;
    margin-bottom: 20px;

    & > .title {
      position: relative;
      font-weight: bold;
      font-size: 0.9em;
    }
  }
}

.pf-head {
  & > .title {
    position: relative;
  }

  & > .sub-title {
    font-size: 14px;
    margin-top: 5px;
  }
}
</style>

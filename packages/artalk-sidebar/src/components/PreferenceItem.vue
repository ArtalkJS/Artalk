<script setup lang="ts">
import settings from '../lib/settings'

const props = defineProps<{
  tplData: any
  path: (string|number)[]
}>()

const desc = computed(() => settings.extractItemDescFromComment(props.path))
const customValue = computed(() => settings.customs.value?.getIn(props.path) as any)

function onChange(value: boolean|string) {
  settings.customs.value?.setIn(props.path, value)
}
</script>

<template>
  <div class="pf-item">
    <div class="info">
      <div class="title">{{ desc.title }}</div>
        <div v-if="!!desc.subTitle" class="sub-title">{{ desc.subTitle }}</div>
      </div>

      <div class="value">
        <!-- 候选框 -->
        <template v-if="desc.opts !== null">
          <select :value="customValue" @change="onChange(($event.target as any).value)">
            <option
              v-for="item in desc.opts"
              :value="item"
            >{{ item }}</option>
          </select>
        </template>

        <!-- 开关 -->
        <template v-else-if="typeof tplData === 'boolean'">
          <input type="checkbox" :checked="customValue" @change="onChange(($event.target as any).checked)">
        </template>

        <!-- 文本框 -->
        <template v-else>
          <input
            type="text"
            :value="(typeof customValue === 'undefined' || customValue === null) ? '' : String(customValue)"
            @change="onChange(($event.target as any).value)"
          />
        </template>
      </div>
  </div>
</template>

<style scoped lang="scss">
.pf-item {
  display: flex;
  flex-direction: row;
  margin-bottom: 20px;

  & > .info {
    display: flex;
    justify-content: center;
    flex-direction: column;
    flex: 1;
    padding-right: 20px;

    .title {
    }

    .sub-title {
      font-size: 14px;
      margin-top: 4px;
      color: #697182;
    }
  }

  & > .value {
    flex: 1;
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    align-items: center;
    min-height: 35px;
  }
}
</style>

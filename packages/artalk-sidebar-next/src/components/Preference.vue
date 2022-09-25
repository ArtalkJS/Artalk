<script setup lang="ts">
import settings from '../lib/settings'

const props = defineProps<{
  tplPfItem: {
    key: string|number,
    valDefault: { [key: string]: any }|Array<any>|boolean|string|number,
    path: (string|number)[]
  },
  expandGrp?: string
}>()

const emits = defineEmits<{
  (evt: 'toggle', key?: string): void
}>()

const { tplPfItem: tplPf } = toRefs(props)

const path = computed(() => tplPf.value.path.join('.'))
const desc = computed(() => settings.extractItemDescFromComment(String(tplPf.value.key), path.value))
const isRoot = computed(() => tplPf.value.path.length === 1)
const level = computed(() => tplPf.value.path.length)

function onToggleGrp(key?: string) {
  emits('toggle', key)
}
</script>

<template>
  <div class="pf">
    <!-- 数组 -->
    <div v-if="Array.isArray(tplPfItem.valDefault)" class="arr-grp">
      <PreferenceGrp
        :title="desc.title"
        :sub-title="desc.subTitle"
        :level="level"
        :pf-key="String(tplPfItem.key)"
        :expand="level !== 1 || expandGrp === tplPfItem.key"
        @toggle="onToggleGrp(String(tplPfItem.key))"
      >
        <div v-for="(item, index) in tplPfItem.valDefault">
          <Preference :tpl-pf-item="{ key: index, valDefault: item, path: [...tplPfItem.path, index] }" />
        </div>
      </PreferenceGrp>
    </div>

    <!-- 对象 -->
    <div v-else-if="tplPfItem.valDefault !== null && typeof tplPfItem.valDefault === 'object'" class="obj-grp">
      <PreferenceGrp
        :title="desc.title"
        :sub-title="desc.subTitle"
        :level="level"
        :pf-key="String(tplPfItem.key)"
        :expand="level !== 1 || expandGrp === tplPfItem.key"
        @toggle="onToggleGrp(String(tplPfItem.key))"
      >
        <div v-for="key in Object.keys(tplPfItem.valDefault)">
          <Preference :tpl-pf-item="{ key, valDefault: tplPfItem.valDefault[key], path: [...tplPfItem.path, key] }" />
        </div>
      </PreferenceGrp>
    </div>

    <div class="form-grp" v-else>
      <div class="info">
        <div class="title">{{ desc.title }}</div>
        <div v-if="!!desc.subTitle" class="sub-title">{{ desc.subTitle }}</div>
      </div>

      <div class="value">
        <!-- 开关 -->
        <template v-if="typeof tplPfItem.valDefault === 'boolean'">
          <input type="checkbox" :checked="tplPfItem.valDefault">
        </template>

        <!-- 文本框 -->
        <template v-else>
          <input type="text" :value="String(tplPfItem.valDefault)" />
        </template>
      </div>
    </div>


  </div>
</template>

<style scoped lang="scss">
.pf {
  .form-grp {
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

      input[type="text"] {
        font-size: 17px;
        width: 100%;
        height: 35px;
        padding: 3px 5px;
        border: 0;
        border-bottom: 1px solid var(--at-color-border);
        outline: none;
        background: transparent;

        &:focus {
          border-bottom-color: var(--at-color-main);
        }
      }
    }
  }

  .arr-grp {
    background: #fff;
    margin-bottom: 10px;
    border-radius: 4px;
  }

  .obj-grp {
    background: #fff;
    margin-bottom: 10px;
    border-radius: 4px;
  }
}
</style>

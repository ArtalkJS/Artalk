<script setup lang="ts">
import YAML from 'yaml'
import { useNavStore } from '../stores/nav'
import { artalk } from '../global'
import settings from '../lib/settings'
import confTemplate from '../assets/artalk-go.example.yml?raw'
import { storeToRefs } from 'pinia'

const nav = useNavStore()
const { curtTab } = storeToRefs(nav)
const settingsTpl = YAML.parseDocument(confTemplate)

onMounted(() => {
  nav.updateTabs({
    save: '保存',
  })
  watch(curtTab, (tab) => {
    if (tab === 'save') {
      console.log(settings.customs.value?.toString())
    }
    nav.setTabActive('')
  })

  artalk!.ctx.getApi().system.getSettings().then((yamlStr) => {
    settings.customs.value = YAML.parseDocument(yamlStr)
  })
})
</script>

<template>
  <div class="settings">
    <PreferenceGrp
      :tpl-data="settingsTpl.toJS()"
      :path="[]"
    />
  </div>
</template>

<style scoped lang="scss">
.settings {
  padding: 20px 30px;

  :deep(input[type="text"]), :deep(select) {
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
</style>

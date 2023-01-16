<script setup lang="ts">
import YAML from 'yaml'
import { useNavStore } from '../stores/nav'
import { artalk } from '../global'
import settings from '../lib/settings'
import confTemplate from '../assets/artalk-go.example.yml?raw'
import { storeToRefs } from 'pinia'
import LoadingLayer from '../components/LoadingLayer.vue'

const nav = useNavStore()
const router = useRouter()
const { curtTab } = storeToRefs(nav)
const settingsTpl = YAML.parseDocument(confTemplate)
const isLoading = ref(false)

onBeforeMount(() => {
  settings.init()
})

onMounted(() => {
  nav.updateTabs({
    'sites': '站点',
    'transfer': '迁移',
  })

  watch(curtTab, (tab) => {
    if (tab === 'sites') router.replace('/sites')
    else if (tab === 'transfer') router.replace('/transfer')
  })

  artalk!.ctx.getApi().system.getSettings().then((yamlStr) => {
    settings.get().customs.value = YAML.parseDocument(yamlStr)
  })
})

function save() {
  let yamlStr = ''
  try {
    yamlStr = settings.get().customs.value?.toString() || ''
  } catch (err) {
    alert('配置文件生成失败：'+err)
    console.error(err)
    return
  }

  console.log(yamlStr)
  if (!yamlStr) {
    alert('配置文件生成失败：数据为空')
    return
  }

  if (isLoading.value) return
  isLoading.value = true
  artalk!.ctx.getApi().system.saveSettings(yamlStr).then(() => {
    alert('设置保存成功')
  }).catch((err) => {
    console.error(err)
    alert('设置保存失败：'+err)
  }).finally(() => {
    isLoading.value = false
  })
}
</script>

<template>
  <div class="settings">
    <div class="act-bar">
      <div class="status-text"></div>
      <button class="save-btn" @click="save()"><i class="atk-icon atk-icon-yes" /> 应用</button>
      <LoadingLayer v-if="isLoading" />
    </div>
    <div class="pfs">
      <PreferenceGrp
        :tpl-data="settingsTpl.toJS()"
        :path="[]"
      />
      <div class="notice">注：某些配置项可能需手动重启才能生效，详见 <a href="https://artalk.js.org/guide/backend/config.html" target="_blank">官方文档</a></div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.settings {
  .notice {
    font-size: 13px;
    background: #e5ecfa;
    color: #1967d2;
    text-align: center;
    padding: 8px 10px;
    margin-top: 10px;
    margin-bottom: 20px;
  }

  .act-bar {
    z-index: 999;
    position: fixed;
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-direction: row;
    height: 55px;
    width: 100%;
    bottom: 0;
    left: 0;
    background: rgba(255, 255, 255, 0.831);
    border-top: 1px solid var(--at-color-border);
    padding: 0 20px;

    .status-text {
      padding: 0 5px;
      flex: 1;
    }

    button {
      font-size: 16px;
      display: inline-flex;
      align-items: center;
      padding: 4px 16px;
      cursor: pointer;
      background: transparent;
      background: var(--at-color-main);
      color: #fff;
      border: 0;

      &:active {
        opacity: .9;
      }

      i {
        width: 16px;
        height: 16px;
        display: inline-block;
        margin-right: 8px;
        background-color: #fff;
      }
    }
  }

  .pfs {
    padding: 10px 30px;
  }

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

<script setup lang="ts">
import YAML from 'yaml'
import { shallowRef } from 'vue'
import { useNavStore } from '../stores/nav'
import { artalk } from '../global'
import settings, { type OptionNode } from '../lib/settings'
import { storeToRefs } from 'pinia'
import LoadingLayer from '../components/LoadingLayer.vue'

const nav = useNavStore()
const router = useRouter()
const { t } = useI18n()
const { curtTab } = storeToRefs(nav)
const isLoading = ref(false)
const tree = shallowRef<OptionNode>()

onMounted(() => {
  nav.updateTabs({
    'sites': 'site',
    'transfer': 'transfer',
  })

  watch(curtTab, (tab) => {
    if (tab === 'sites') router.replace('/sites')
    else if (tab === 'transfer') router.replace('/transfer')
  })

  Promise.all([
    artalk!.ctx.getApi().settings.getSettingsTemplate(''),
    artalk!.ctx.getApi().settings.getSettings(),
  ]).then(([template, custom]) => {
    const yamlObj = YAML.parseDocument(template.data.yaml)
    tree.value = settings.init(yamlObj).getTree()
    console.log(tree.value)
    settings.get().setCustoms(custom.data.yaml)
  })
})

function save() {
  let yamlStr = ''
  try {
    yamlStr = settings.get().getCustoms().value?.toString() || ''
  } catch (err) {
    alert('YAML export error: '+err)
    console.error(err)
    return
  }

  console.log(yamlStr)
  if (!yamlStr) {
    alert('YAML export error: data is empty')
    return
  }

  if (isLoading.value) return
  isLoading.value = true
  artalk!.ctx.getApi().settings.applySettings({
    yaml: yamlStr
  }).then(() => {
    alert(t('settingSaved'))
  }).catch((err) => {
    console.error(err)
    alert(t('settingSaveFailed')+': '+err)
  }).finally(() => {
    isLoading.value = false
  })
}
</script>

<template>
  <div class="settings">
    <div class="act-bar">
      <div class="status-text"></div>
      <button class="save-btn" @click="save()"><i class="atk-icon atk-icon-yes" /> {{ t('apply') }}</button>
      <LoadingLayer v-if="isLoading" />
    </div>
    <div v-if="tree" class="pfs">
      <PreferenceGrp :node="tree" />
      <div class="notice">{{ t('settingNotice') }}</div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.settings {
  .notice {
    font-size: 13px;
    background: var(--at-color-bg-light);
    color: var(--at-color-light);
    border-radius: 2px;
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
    background: var(--at-color-bg-transl);
    border-top: 1px solid var(--at-color-border);
    padding: 0 20px;

    .status-text {
      padding: 0 5px;
      flex: 1;
    }

    button {
      font-size: 14px;
      display: inline-flex;
      align-items: center;
      padding: 4px 16px;
      cursor: pointer;
      background: transparent;
      border-radius: 2px;
      background: var(--at-color-main);
      color: #fff;
      border: 0;

      &:active {
        opacity: .9;
      }

      i {
        margin-right: 8px;

        &::after {
          background-color: #fff;
        }
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

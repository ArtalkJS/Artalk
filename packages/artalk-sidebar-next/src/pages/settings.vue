<script setup lang="ts">
import YAML from 'yaml'
import { useNavStore } from '../stores/nav'
import { artalk } from '../global'
import settings from '../lib/settings'
import Preference from '../components/Preference.vue'
import confTemplate from '../assets/artalk-go.example.yml?raw'

const nav = useNavStore()
const yamlDocTpl = YAML.parseDocument(confTemplate)
const yamlDoc = ref<YAML.Document.Parsed<YAML.ParsedNode>>()

onMounted(() => {
  nav.updateTabs({})

  artalk!.ctx.getApi().system.getSettings().then((yamlStr) => {
    yamlDoc.value = YAML.parseDocument(yamlStr)
    // console.log(yamlDoc.value)
    // console.log(yamlDoc.value.toJS())

    yamlDoc.value.setIn(['email','enabled'], true)

    // console.log(yamlDoc.value.toString())
  })
})

const yamlTplObj = computed(() => yamlDocTpl.toJS() || {})
</script>

<template>
  <div class="settings">
    <div v-for="key in Object.keys(yamlTplObj)">
      <Preference :tpl-pf-item="{ key, valDefault: yamlTplObj[key], path: [key] }" />
    </div>
  </div>
</template>

<style scoped lang="scss">
.settings {
  padding: 20px 30px;
}
</style>

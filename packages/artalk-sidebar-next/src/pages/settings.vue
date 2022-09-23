<script setup lang="ts">
import YAML from 'yaml'
import { useNavStore } from '../stores/nav'
import { artalk } from '../global'
import settings from '../lib/settings'
import Preference from '../components/Preference.vue'

const nav = useNavStore()
const yamlDoc = ref<YAML.Document.Parsed<YAML.ParsedNode>>()

onMounted(() => {
  nav.updateTabs({})

  artalk!.ctx.getApi().system.getSettings().then((yamlStr) => {
    yamlDoc.value = YAML.parseDocument(yamlStr)
    console.log(yamlDoc.value.toJS())

    yamlDoc.value.setIn(['email','enabled'], true)

    console.log(yamlDoc.value.toString())
  })
})

const yamlObj = computed(() => yamlDoc.value?.toJS() || {})
</script>

<template>
  <div>
    <div v-for="key in Object.keys(yamlObj)">
      <Preference :pf="{ name: key, data: yamlObj[key], paths: [key] }" />
    </div>
  </div>
</template>

<style scoped lang="scss">

</style>

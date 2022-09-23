<script setup lang="ts">
const props = defineProps<{
  pf: {
    name: string|number,
    data: { [key: string]: any }|Array<any>|boolean|string|number,
    paths: (string|number)[]
  },
}>()

const { pf } = toRefs(props)

</script>

<template>
  <div class="preference" style="margin-left: 2em">
    {{ pf.name }} [{{ pf.paths.join('.') }}]:
    <template v-if="Array.isArray(pf.data)">
      <div v-for="(item, index) in pf.data">
        <Preference :pf="{ name: index, data: item, paths: [...pf.paths, index] }" />
      </div>
    </template>
    <template v-else-if="pf.data !== null && typeof pf.data === 'object'">
      <div v-for="objKey in Object.keys(pf.data)">
        <Preference :pf="{ name: objKey, data: pf.data[objKey], paths: [...pf.paths, objKey] }" />
      </div>
    </template>
    <template v-else-if="typeof pf.data === 'boolean'">
      <div>
        <input type="checkbox" :checked="pf.data">
        <label>{{ pf.name }}</label>
      </div>
    </template>
    <template v-else>
      <input type="text" :value="String(pf.data)" />
    </template>
  </div>
</template>

<style scoped lang="scss">

</style>

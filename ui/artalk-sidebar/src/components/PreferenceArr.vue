<script setup lang="ts">
import settings, { patchOptionValue, type OptionNode } from '../lib/settings'

const props = defineProps<{
  node: OptionNode
}>()

const customValue = ref<string[]>([])
const disabled = ref(false)

onMounted(() => {
  sync()
})

function sync() {
  const value = settings.get().getCustom(props.node.path)
  disabled.value = !!settings.get().getEnvByPath(props.node.path)
  if (typeof value === 'object' && 'toJSON' in value && typeof value.toJSON === 'function') {
    customValue.value = value.toJSON()
  } else if (typeof value === 'string') {
    customValue.value = value.split(' ')
  } else {
    customValue.value = []
  }
}

function save() {
  const v = patchOptionValue(customValue.value, props.node)
  settings.get().setCustom(props.node.path, v)
}

function onChange(index: number, val: string) {
  customValue.value[index] = val
  save()
}

function remove(index: number) {
  customValue.value.splice(index, 1)
  save()
}

function add() {
  customValue.value.push('')
  save()
}
</script>

<template>
  <div class="arr-grp">
    <div v-for="(item, index) in customValue" :key="index" class="arr-item">
      <input
        type="text"
        :value="String(item)"
        :disabled="disabled"
        @change="onChange(index, ($event.target as any).value)"
      />
      <button v-if="!disabled" class="act-btn" @click="remove(index)">-</button>
    </div>
    <div v-if="!disabled" class="act-grp">
      <button class="act-btn" @click="add()">+</button>
    </div>
  </div>
</template>

<style scoped lang="scss">
.arr-grp {
  width: 100%;
}

.arr-item {
  position: relative;
  margin-bottom: 20px;
  padding-right: 40px;

  .act-btn {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    right: 0;
  }
}

.act-grp {
  margin-left: 10px;

  .act-btn {
    padding: 2px 30px;
  }
}

.act-btn {
  display: inline-block;
  padding: 2px 10px;
  cursor: pointer;
  border: 0;
  color: var(--at-color-font);
  background: var(--at-color-bg-grey);
  border-radius: 2px;

  &:hover {
    color: var(--at-color-light);
    background: var(--at-color-bg-light);
  }
}
</style>

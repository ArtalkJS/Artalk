<script setup lang="ts">
import settings from '../lib/settings'

const props = defineProps<{
  tplData: Array<any>,
  path: (string|number)[]
}>()

const ci = getCurrentInstance()

const customValue = ref([])
const update = () => {
  customValue.value = (settings.customs.value?.getIn(props.path) as any)?.items || []
  ci?.proxy?.$forceUpdate()
}
watch(settings.customs, (customs) => {
  update()
})

function onChange(index: number, val: string) {
  settings.customs.value?.setIn([...props.path, index], val)
}

function remove(index: number) {
  settings.customs.value?.deleteIn([...props.path, index])
  update()
}

function add() {
  if (!customValue.value) settings.customs.value?.setIn([...props.path], [''])
  else settings.customs.value?.setIn([...props.path, customValue.value.length], '')
  update()
}
</script>

<template>
  <div class="arr-grp">
    <div v-for="(item, index) in customValue" class="arr-item">
      <input
        type="text"
        :value="String(item)"
        @change="onChange(index, ($event.target as any).value)"
      />
      <button class="act-btn" @click="remove(index)">-</button>
    </div>
    <div class="act-grp">
      <button class="act-btn" @click="add()">+</button>
    </div>
  </div>
</template>

<style scoped lang="scss">
.arr-item {
  position: relative;
  margin-bottom: 20px;
  margin-left: 10px;
  padding-right: 50px;

  .act-btn {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    right: 10px;
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
  color: #5f6369;
  background: #f4f4f4;
  border-radius: 2px;

  &:hover {
    color: #1967d2;
    background: #e5ecfa;
  }
}
</style>

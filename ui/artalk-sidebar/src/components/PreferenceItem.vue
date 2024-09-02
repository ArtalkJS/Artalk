<script setup lang="ts">
import settings, { patchOptionValue, type OptionNode } from '../lib/settings'
import { isSensitiveConfigPath } from '@/lib/settings-sensitive'

const props = defineProps<{
  node: OptionNode
}>()

const value = ref('')
const disabled = ref(false)
const sensitiveHidden = ref(true)

const { t } = useI18n()

onBeforeMount(() => {
  // initial value
  value.value = settings.get().getCustom(props.node.path)
  disabled.value = !!settings.get().getEnvByPath(props.node.path)
})

function onChange() {
  const v = patchOptionValue(value.value, props.node)
  settings.get().setCustom(props.node.path, v)
  // console.log('[SET]', props.node.path, v)
}

const envVariableName = computed(() => `ATK_${props.node.path.replace(/\./g, '_').toUpperCase()}`)
const isSensitive = computed(() => isSensitiveConfigPath(props.node.path))

function toggleSensitiveHidden() {
  sensitiveHidden.value = !sensitiveHidden.value
}
</script>

<template>
  <div class="pf-item">
    <div class="info">
      <div class="title" :title="envVariableName">{{ node.title }}</div>
      <div v-if="node.subTitle" class="sub-title">{{ node.subTitle }}</div>
    </div>

    <div class="value">
      <div v-if="disabled" class="disable-note">
        {{ t('envVarControlHint', { key: 'ATK_' + props.node.path.toUpperCase() }) }}
      </div>

      <!-- Array -->
      <template v-if="node.type === 'array'">
        <PreferenceArr :node="node" />
      </template>

      <!-- Dropdown -->
      <template v-else-if="node.selector">
        <select v-model="value" :disabled="disabled" @change="onChange">
          <option v-for="(item, i) in node.selector" :key="i" :value="item">
            {{ item }}
          </option>
        </select>
      </template>

      <!-- Toggle -->
      <template v-else-if="node.type === 'boolean'">
        <input v-model="value" type="checkbox" :disabled="disabled" @change="onChange" />
      </template>

      <!-- Text -->
      <template v-else>
        <input
          v-model="value"
          :type="!isSensitive || !sensitiveHidden ? 'text' : 'password'"
          :disabled="disabled"
          @change="onChange"
        />
        <div v-if="isSensitive" class="input-suffix">
          <div class="hidden-switch" @click="toggleSensitiveHidden()">
            <i :class="['atk-icon', `atk-icon-eye-${sensitiveHidden ? 'off' : 'on'}`]" />
          </div>
        </div>
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
    position: relative;
    flex: 1;
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    align-items: center;
    min-height: 35px;

    &:hover {
      .disable-note {
        opacity: 1;
      }
    }

    .disable-note {
      z-index: 999;
      padding: 3px 8px;
      background: rgba(105, 113, 130, 0.9);
      color: #fff;
      position: absolute;
      top: -30px;
      font-size: 13px;
      left: 0;
      opacity: 0;
      transition: opacity 0.2s;
    }

    .input-suffix {
      margin-left: 5px;
    }

    .hidden-switch {
      cursor: pointer;
      padding-left: 10px;

      .atk-icon {
        &::after {
          background-color: #697182;
        }

        &.atk-icon-eye-on::after {
          mask-image: url('@/assets/icon-eye-on.svg');
        }

        &.atk-icon-eye-off::after {
          mask-image: url('@/assets/icon-eye-off.svg');
        }
      }
    }
  }
}
</style>

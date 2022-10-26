<script setup lang="ts">
import settings from '../lib/settings'

const router = useRouter()

const props = defineProps<{
  tplData: Object|Array<any>
  path: (string|number)[]
}>()

const desc = computed(() => settings.get().extractItemDescFromComment(props.path))
const level = computed(() => props.path.length)

const expanded = ref(true)

onMounted(() => {
  if (level.value === 1) expanded.value = false
})

function onHeadClick(evt: Event) {
  if (level.value !== 1) return
  if (!expanded.value) {
    expanded.value = true
    // nextTick(() => {
    //   nav.scrollPageToEl(evt.target as HTMLElement)
    // })
  } else {
    expanded.value = false
  }
}
</script>

<template>
  <div class="pf-grp" :class="[`level-${level}`, expanded ? 'expand' : '']">
    <div v-if="level > 0" class="pf-head" @click="onHeadClick">
      <div class="title">{{ desc.title }}</div>
      <div v-if="!!desc.subTitle" class="sub-title">{{ desc.subTitle }}</div>
    </div>
    <div v-show="expanded" class="pf-body">
      <!-- Array -->
      <template v-if="Array.isArray(tplData)">
        <PreferenceArr :tpl-data="tplData" :path="path" />
      </template>
      <!-- Object -->
      <template v-else>
        <div v-for="[key, value] in Object.entries(tplData)">
          <!-- Admin Users -->
          <template v-if="key === 'admin_users'"></template>
          <!-- Grp -->
          <PreferenceGrp
            v-else-if="value !== null && typeof value === 'object'"
            :tpl-data="value"
            :path="[...path, key]"
          />
          <!-- Item Input -->
          <PreferenceItem
            v-else
            :tpl-data="value"
            :path="[...path, key]"
          />
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped lang="scss">
.pf-grp {
  background: #fff;
  margin-bottom: 10px;
  border-radius: 4px;
}

.pf-grp.level-1 {
  & > .pf-head {
    margin-top: 30px;
    margin-bottom: 20px;
    cursor: pointer;

    .title {
      font-size: 1.4em;
      font-weight: bold;
      padding-left: 12px;

      &::before {
        position: absolute;
        top: 50%;
        transform: translateY(-50%);
        transition: height ease .2s;
        left: -10px;
        content: '';
        height: 10px;
        width: 10px;
        background: #B7DCFF;
      }
    }
  }

  & > .pf-body {
  }

  &.expand > .pf-head .title::before {
    height: 25px;
  }
}

.pf-grp.level-2 {
  & > .pf-head {
    margin-top: 15px;
    margin-bottom: 20px;

    .title {
      position: relative;
      font-weight: bold;
      font-size: 1.2em;
    }
  }
}

.pf-grp.level-1, .pf-grp.level-2 {
  & > .pf-head > .sub-title {
    padding: 0 10px 0 10px;
    margin-left: 4px;
    margin-top: 15px;
    border-left: 2px solid #eee;
  }
}


.pf-grp.level-3 {
  margin-left: 10px;

  & > .pf-head {
    margin-top: 15px;
    margin-bottom: 20px;

    & > .title {
      position: relative;
      font-weight: bold;
      font-size: 0.9em;
    }
  }

  & > .pf-body {
    margin-left: 15px;
  }
}

.pf-head {
  & > .title {
    position: relative;
  }

  & > .sub-title {
    font-size: 14px;
    margin-top: 5px;
  }
}
</style>

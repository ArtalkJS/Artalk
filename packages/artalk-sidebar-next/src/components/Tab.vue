<script setup lang="ts">
const tabGrps = {
  '评论': ['全部', '待审', '个人'],
  '页面': [],
  '站点': [],
  '迁移': ['导入', '导出']
}

type View = keyof typeof tabGrps

let viewCurt = ref<View>('评论')
let tabCurt = ref('全部')

let displayType = ref<'views'|'tabs'>('tabs')

function toggleDisplay() {
  displayType.value = displayType.value !== 'tabs' ? 'tabs' : 'views'
}

function switchView(name: View) {
  viewCurt.value = name
  displayType.value = 'tabs'
}

function switchTab(name: string) {
  tabCurt.value = name
}
</script>

<template>
  <div class="tab">
    <div class="view" @click="toggleDisplay()">
      <div class="icon" :class="(displayType === 'tabs') ? 'menu' : 'arrow'"></div>
      <div class="text">{{ viewCurt }}</div>
    </div>

    <div class="tab-list">
      <!-- tabs -->
      <template v-if="displayType === 'tabs'">
        <div
          v-for="(tab) in tabGrps[viewCurt]"
          class="item"
          :class="{ active: tab === tabCurt }"
          @click="switchTab(tab)"
        >{{ tab }}</div>
      </template>

      <!-- views -->
      <template v-else>
        <div
          v-for="(view) in Object.keys(tabGrps)"
          class="item"
          :class="{ active: view === viewCurt }"
          @click="switchView(view as View)"
        >{{ view }}</div>
      </template>
    </div>
  </div>
</template>

<style scoped lang="scss">
.tab {
  display: flex;
  flex-direction: row;
  align-items: center;
  border-bottom: 1px solid #eceff2;
  height: 41px;

  .view {
    display: flex;
    align-items: center;
    height: 100%;
    padding: 0 20px;
    border-right: 1px solid #eceff2;
    cursor: pointer;
    user-select: none;

    .icon {
      width: 17px;
      height: 100%;
      background-repeat: no-repeat;
      background-position: 50%;
      margin-right: 10px;

      &.menu {
        background-image: url("data:image/svg+xml,%3Csvg width='14' height='11' viewBox='0 0 14 11' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cg clip-path='url(%23clip0_1_2)'%3E%3Crect width='14' height='11' fill='white'/%3E%3Cpath fill-rule='evenodd' clip-rule='evenodd' d='M0 0H14V1H0V0ZM0 5H14V6H0V5ZM14 10H0V11H14V10Z' fill='black'/%3E%3C/g%3E%3Cdefs%3E%3CclipPath id='clip0_1_2'%3E%3Crect width='14' height='11' fill='white'/%3E%3C/clipPath%3E%3C/defs%3E%3C/svg%3E");
      }

      &.arrow {
        background-image: url("data:image/svg+xml,%3Csvg width='9' height='14' viewBox='0 0 9 14' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Crect width='9' height='13' fill='white'/%3E%3Cpath fill-rule='evenodd' clip-rule='evenodd' d='M7.52898 0L1 6.52898L1.17809 6.70707L1.17805 6.70711L7.52897 13.058L8.23608 12.3509L2.41418 6.52902L8.23609 0.707107L7.52898 0Z' fill='black'/%3E%3C/svg%3E");
      }
    }

    .text {
      color: #2a2e2e;
    }
  }

  .tab-list {
    flex: auto;
    height: 100%;
    display: flex;
    flex-direction: row;
    overflow-y: auto;

    .item {
      padding: 0 20px;
      color: #757575;
      display: flex;
      align-items: center;
      justify-content: center;
      height: 100%;
      cursor: pointer;

      &.active {
        color: #2a2e2e;
      }

      &:hover {
        background: #f4f4f4;
      }
    }
  }
}
</style>

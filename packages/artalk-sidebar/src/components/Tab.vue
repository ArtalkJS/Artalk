<script setup lang="ts">
import { useNavStore } from '../stores/nav'
import { useUserStore } from '../stores/user'
import { storeToRefs } from 'pinia'

type PageItem = { label: string, link: string, hide?: boolean }

const router = useRouter()
const route = useRoute()
const { curtPage, curtTab, tabs } = storeToRefs(useNavStore())
const { isAdmin } = storeToRefs(useUserStore())
const indicator = ref<'pages'|'tabs'>('tabs')
const tabListEl = ref<HTMLElement|null>(null)

const pages = ref<{ [name: string]: PageItem }>({})

onMounted(() => {
  if (isAdmin.value) {
    pages.value = {
      comments: {
        label: '评论',
        link: '/comments',
      },
      pages: {
        label: '页面',
        link: '/pages',
      },
      sites: {
        label: '站点',
        link: '/sites',
      },
      transfer: {
        label: '迁移',
        link: '/transfer',
        hide: true,
      },
      settings: {
        label: '设置',
        link: '/settings'
      }
    }
  } else {
    pages.value = {
      comments: {
        label: '评论',
        link: '/comments',
      }
    }
  }

  tabListEl.value!.addEventListener('wheel', (evt) => {
    evt.preventDefault()
    tabListEl.value!.scrollLeft += evt.deltaY
  })

  curtPage.value = route.name.replace(/^\//, '')
})

function toggleIndicator() {
  indicator.value = (indicator.value !== 'tabs') ? 'tabs' : 'pages'
}

function switchPage(pageName: string) {
  indicator.value = 'tabs'

  router.replace(pages.value[pageName].link)
}

function switchTab(tabName: string) {
  curtTab.value = tabName
}

router.afterEach((to, from, failure) => {
  curtPage.value = to.name.replace(/^\//, '')
})
</script>

<template>
  <div class="tab">
    <div class="page" @click="toggleIndicator()">
      <div class="icon" :class="(indicator === 'tabs') ? 'menu' : 'arrow'"></div>
      <div class="text">{{ pages[curtPage]?.label || '' }}</div>
    </div>

    <div ref="tabListEl" class="tab-list">
      <!-- tabs -->
      <template v-if="indicator === 'tabs'">
        <div
          v-for="(tabLabel, tabName) in tabs"
          class="item"
          :class="{ active: curtTab === tabName }"
          @click="switchTab(tabName as string)"
        >{{ tabLabel }}</div>
      </template>

      <!-- pages -->
      <template v-else>
        <div
          v-for="(page, pageName) in pages"
          class="item"
          v-show="!page.hide"
          :class="{ active: pageName === curtPage }"
          @click="switchPage(pageName as string)"
        >{{ page.label }}</div>
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

  .page {
    display: flex;
    align-items: center;
    height: 100%;
    padding: 0 20px;
    border-right: 1px solid #eceff2;
    cursor: pointer;
    user-select: none;
    white-space: nowrap;

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

    &::-webkit-scrollbar {
      width: 0;
      height: 0;
      background: transparent;
    }

    .item {
      padding: 0 20px;
      color: #757575;
      display: flex;
      align-items: center;
      justify-content: center;
      height: 100%;
      cursor: pointer;
      white-space: nowrap;

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

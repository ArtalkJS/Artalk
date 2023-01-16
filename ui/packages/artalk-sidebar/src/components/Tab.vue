<script setup lang="ts">
import { useNavStore } from '../stores/nav'
import { useUserStore } from '../stores/user'
import { storeToRefs } from 'pinia'

type PageItem = { label: string, link: string, hide?: boolean }

const router = useRouter()
const route = useRoute()
const nav = useNavStore()

const { curtPage, curtTab, tabs, isSearchEnabled } = storeToRefs(nav)
const { isAdmin } = storeToRefs(useUserStore())
const indicator = ref<'pages'|'tabs'>('tabs')
const tabListEl = ref<HTMLElement|null>(null)

const pages = computed((): { [name: string]: PageItem } => {
  if (isAdmin.value) {
    return {
      comments: {
        label: '评论',
        link: '/comments',
      },
      pages: {
        label: '页面',
        link: '/pages',
      },
      users: {
        label: '用户',
        link: '/users',
      },
      sites: {
        label: '站点',
        link: '/sites',
        hide: true,
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
    return {
      comments: {
        label: '评论',
        link: '/comments',
      }
    }
  }
})

const isSearchShow = ref(false)
const searchValue = ref('')
const searchInputEl = ref<HTMLInputElement|null>(null)

onMounted(() => {
  tabListEl.value!.addEventListener('wheel', (evt) => {
    evt.preventDefault()
    tabListEl.value!.scrollLeft += evt.deltaY
  })
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

// @link https://router.vuejs.org/zh/guide/advanced/navigation-guards.html
router.beforeEach((to, from) => {
  isSearchEnabled.value = false
})

router.afterEach((to, from, failure) => {
  curtPage.value = to.name.replace(/^\//, '')
})

function showSearch() {
  isSearchShow.value = true
  nextTick(() => {
    searchInputEl.value?.focus()
  })
}

function hideSearch() {
  isSearchShow.value = false
  searchValue.value = ''
  if (nav.searchResetEvent) {
    nav.searchResetEvent()
  }
}

function onSearchSubmit(evt: Event) {
  if (searchValue.value.trim() === '') {
    searchInputEl.value?.focus()
    return
  }

  if (nav.searchEvent) {
    nav.searchEvent(searchValue.value.trim())
  }
}
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

        <div
          v-if="isSearchEnabled"
          class="item search-btn"
          @click="showSearch()"
        ></div>
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

    <form
      v-if="isSearchShow"
      class="search-layer atk-fade-in"
      @submit.prevent="onSearchSubmit"
    >
      <div class="item back-btn" @click="hideSearch()">
        <div class="icon arrow"></div>
      </div>
      <input ref="searchInputEl" type="text" placeholder="搜索关键字..." v-model="searchValue" required>
      <button type="submit" class="item search-btn"></button>
    </form>
  </div>
</template>

<style scoped lang="scss">
.tab {
  position: relative;
  display: flex;
  flex-direction: row;
  align-items: center;
  border-bottom: 1px solid #eceff2;
  height: 41px;

  .icon {
    width: 17px;
    height: 100%;
    background-repeat: no-repeat;
    background-position: 50%;

    &.menu {
      background-image: url("data:image/svg+xml,%3Csvg width='14' height='11' viewBox='0 0 14 11' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cg clip-path='url(%23clip0_1_2)'%3E%3Crect width='14' height='11' fill='white'/%3E%3Cpath fill-rule='evenodd' clip-rule='evenodd' d='M0 0H14V1H0V0ZM0 5H14V6H0V5ZM14 10H0V11H14V10Z' fill='black'/%3E%3C/g%3E%3Cdefs%3E%3CclipPath id='clip0_1_2'%3E%3Crect width='14' height='11' fill='none'/%3E%3C/clipPath%3E%3C/defs%3E%3C/svg%3E");
    }

    &.arrow {
      background-image: url("data:image/svg+xml,%3Csvg width='9' height='14' viewBox='0 0 9 14' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Crect width='9' height='13' fill='none'/%3E%3Cpath fill-rule='evenodd' clip-rule='evenodd' d='M7.52898 0L1 6.52898L1.17809 6.70707L1.17805 6.70711L7.52897 13.058L8.23608 12.3509L2.41418 6.52902L8.23609 0.707107L7.52898 0Z' fill='black'/%3E%3C/svg%3E");
    }
  }

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
      margin-right: 10px;
    }

    .text {
      color: #2a2e2e;
    }
  }

  .tab-list {
    position: relative;
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

    &.search-btn {
      margin-left: auto;
      padding: 0 20px;
      border-left: 1px solid #eceff2;

      &::after {
        display: inline-block;
        content: '';
        height: 15px;
        width: 15px;
        background-repeat: no-repeat;
        background-position: center;
        background-image: url("data:image/svg+xml,%3Csvg width='15' height='15' viewBox='0 0 15 15' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath fill-rule='evenodd' clip-rule='evenodd' d='M9.89266 3.70572C8.18412 1.99718 5.41402 1.99718 3.70548 3.70572C1.99694 5.41427 1.99694 8.18436 3.70548 9.89291C5.41402 11.6015 8.18412 11.6015 9.89266 9.89291C11.6012 8.18436 11.6012 5.41427 9.89266 3.70572ZM2.8216 2.82184C5.0183 0.625142 8.57985 0.625142 10.7765 2.82184C12.8239 4.86916 12.9631 8.10202 11.1942 10.3106L13.4282 12.5446L12.5443 13.4284L10.3103 11.1945C8.10177 12.9633 4.86892 12.8241 2.8216 10.7768C0.624897 8.58009 0.624897 5.01854 2.8216 2.82184Z' fill='%23757575'/%3E%3C/svg%3E");
      }
    }
  }

  .search-layer {
    z-index: 100;
    position: absolute;
    height: 100%;
    width: 100%;
    background: #fff;
    display: flex;
    flex-direction: row;

    .back-btn {
      padding: 0 20px;
      border-right: 1px solid #eceff2;
    }

    input {
      flex: 1;
      border: 0;
      outline: none;
      padding: 0 20px;
      font-size: 15px;
    }

    .search-btn {
      cursor: pointer;
      border: none;
      background: transparent;
    }
  }
}
</style>

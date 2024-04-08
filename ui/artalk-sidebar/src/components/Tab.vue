<script setup lang="ts">
import { useNavStore } from '../stores/nav'
import { useUserStore } from '../stores/user'
import { storeToRefs } from 'pinia'

type PageItem = { label: string; link: string; hide?: boolean }

const router = useRouter()
const route = useRoute()
const nav = useNavStore()
const { t } = useI18n()

const { curtPage, curtTab, tabs, isSearchEnabled } = storeToRefs(nav)
const { isAdmin } = storeToRefs(useUserStore())
const indicator = ref<'pages' | 'tabs'>('tabs')
const tabListEl = ref<HTMLElement | null>(null)

const pages = computed((): { [name: string]: PageItem } => {
  if (isAdmin.value) {
    return {
      comments: {
        label: 'comment',
        link: '/comments',
      },
      pages: {
        label: 'page',
        link: '/pages',
      },
      users: {
        label: 'user',
        link: '/users',
      },
      sites: {
        label: 'site',
        link: '/sites',
        hide: true,
      },
      transfer: {
        label: 'transfer',
        link: '/transfer',
        hide: true,
      },
      settings: {
        label: 'settings',
        link: '/settings',
      },
    }
  } else {
    return {
      comments: {
        label: 'comment',
        link: '/comments',
      },
    }
  }
})

const isSearchShow = ref(false)
const searchValue = ref('')
const searchInputEl = ref<HTMLInputElement | null>(null)

onMounted(() => {
  tabListEl.value!.addEventListener('wheel', (evt) => {
    evt.preventDefault()
    tabListEl.value!.scrollLeft += evt.deltaY
  })
})

function toggleIndicator() {
  indicator.value = indicator.value !== 'tabs' ? 'tabs' : 'pages'
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
      <div class="icon" :class="indicator === 'tabs' ? 'menu' : 'arrow'"></div>
      <div class="text" v-if="pages[curtPage]?.label">
        {{ t(pages[curtPage].label) }}
      </div>
    </div>

    <div ref="tabListEl" class="tab-list">
      <!-- tabs -->
      <template v-if="indicator === 'tabs'">
        <div
          v-for="(tabLabel, tabName) in tabs"
          :key="tabName"
          class="item"
          :class="{ active: curtTab === tabName }"
          @click="switchTab(tabName as string)"
        >
          {{ t(tabLabel) }}
        </div>

        <div v-if="isSearchEnabled" class="item search-btn" @click="showSearch()"></div>
      </template>

      <!-- pages -->
      <template v-else>
        <div
          v-for="(page, pageName) in pages"
          :key="pageName"
          class="item"
          v-show="!page.hide"
          :class="{ active: pageName === curtPage }"
          @click="switchPage(pageName as string)"
        >
          {{ t(page.label) }}
        </div>
      </template>
    </div>

    <form v-if="isSearchShow" class="search-layer atk-fade-in" @submit.prevent="onSearchSubmit">
      <div class="item back-btn" @click="hideSearch()">
        <div class="icon arrow"></div>
      </div>
      <input
        ref="searchInputEl"
        type="text"
        :placeholder="t('searchHint')"
        v-model="searchValue"
        required
      />
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
  border-bottom: 1px solid var(--at-color-border);
  height: 41px;

  .icon {
    width: 17px;
    height: 100%;
    background-color: var(--at-color-deep);
    background-size: contain;
    background-repeat: no-repeat;
    mask-repeat: no-repeat;
    -webkit-mask-repeat: no-repeat;
    mask-position: center;
    -webkit-mask-position: center;
    mask-size: 100%;
    -webkit-mask-size: 100%;

    &.menu {
      $menuImg: url("data:image/svg+xml,%3Csvg width='50' height='50' viewBox='0 0 50 50' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath fill-rule='evenodd' clip-rule='evenodd' d='M5 9H44V11.8182H5V9ZM5 23.0909H44V25.9091H5V23.0909ZM44 37.1818H5V40H44V37.1818Z' fill='black'/%3E%3C/svg%3E%0A");
      mask-image: $menuImg;
      -webkit-mask-image: $menuImg;
    }

    &.arrow {
      $arrowImg: url("data:image/svg+xml,%3Csvg width='50' height='50' viewBox='0 0 50 50' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath fill-rule='evenodd' clip-rule='evenodd' d='M32.1433 7L15 24.4999L15.4676 24.9773L32.1433 42L34 40.1047L18.7133 24.5001L34 8.89529L32.1433 7Z' fill='black'/%3E%3C/svg%3E%0A");
      mask-image: $arrowImg;
      -webkit-mask-image: $arrowImg;
    }
  }

  .page {
    display: flex;
    align-items: center;
    height: 100%;
    padding: 0 20px;
    border-right: 1px solid var(--at-color-border);
    cursor: pointer;
    user-select: none;
    white-space: nowrap;

    .icon {
      margin-right: 10px;
    }

    .text {
      color: var(--at-color-deep);
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
    color: var(--at-color-sub);
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    cursor: pointer;
    white-space: nowrap;

    &.active {
      color: var(--at-color-deep);
    }

    &:hover {
      background: var(--at-color-bg-grey);
    }

    &.search-btn {
      margin-left: auto;
      padding: 0 20px;
      border-left: 1px solid var(--at-color-border);

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
    background: var(--at-color-bg);
    display: flex;
    flex-direction: row;

    .back-btn {
      padding: 0 20px;
      border-right: 1px solid var(--at-color-border);
    }

    input {
      flex: 1;
      border: 0;
      outline: none;
      padding: 0 20px;
      font-size: 15px;
      background: transparent;
    }

    .search-btn {
      cursor: pointer;
      border: none;
      background: transparent;
    }
  }
}
</style>

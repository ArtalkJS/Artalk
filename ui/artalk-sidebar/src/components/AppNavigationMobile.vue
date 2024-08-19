<script lang="ts" setup>
/**
 * The part of the navigation bar that is displayed on mobile devices and the sidebar
 */
import { useNavigationMenu } from './AppNavigationMenu'

const { t } = useI18n()

const indicatorState = ref<'pages' | 'tabs'>('tabs')
const tabListEl = ref<HTMLElement | null>(null)
const nav = useNavigationMenu({
  onGoPage: () => {
    indicatorState.value = 'tabs'
  },
})

/**
 * For tab list for horizontal scroll
 */
const wheelHandler = (evt: WheelEvent) => {
  evt.preventDefault()
  tabListEl.value && (tabListEl.value!.scrollLeft += evt.deltaY)
}
onMounted(() => tabListEl.value!.addEventListener('wheel', wheelHandler))
onUnmounted(() => tabListEl.value!.removeEventListener('wheel', wheelHandler))

/**
 * Toggle indicator
 */
function toggleIndicator() {
  indicatorState.value = indicatorState.value !== 'tabs' ? 'tabs' : 'pages'
}
</script>

<template>
  <div class="top-navigation">
    <div class="page" @click="toggleIndicator()">
      <div class="icon" :class="indicatorState === 'tabs' ? 'menu' : 'arrow'"></div>
      <div v-if="nav.pages[nav.curtPage]?.label" class="text">
        {{ t(nav.pages[nav.curtPage].label) }}
      </div>
    </div>

    <div ref="tabListEl" class="tab-list">
      <!-- Tabs -->
      <template v-if="indicatorState === 'tabs'">
        <div
          v-for="(tabLabel, tabName) in nav.tabs"
          :key="tabName"
          class="item"
          :class="{ active: nav.curtTab === tabName }"
          @click="nav.goTab(tabName as string)"
        >
          {{ t(tabLabel) }}
        </div>

        <div v-if="nav.isSearchEnabled" class="item search-btn" @click="nav.showSearch()"></div>
      </template>

      <!-- Pages -->
      <template v-if="indicatorState === 'pages'">
        <div
          v-for="(page, pageName) in nav.pages"
          v-show="!page.hideOnMobile"
          :key="pageName"
          class="item"
          :class="{ active: pageName === nav.curtPage }"
          @click="nav.goPage(pageName as string)"
        >
          {{ t(page.label) }}
        </div>
      </template>
    </div>

    <AppNavigationSearch v-if="nav.isMobile" :state="nav.searchState" />
  </div>
</template>

<style lang="scss" scoped>
.top-navigation {
  position: relative;
  display: flex;
  flex-direction: row;
  align-items: center;
  border-bottom: 1px solid var(--at-color-border);
  height: 41px;

  @media (min-width: 1024px) {
    display: none;
  }

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
}
</style>

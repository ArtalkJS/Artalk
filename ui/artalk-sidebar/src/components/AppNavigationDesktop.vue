<script lang="ts" setup>
/**
 * The part of the navigation bar that is displayed on desktop
 */
import { useNavigationMenu } from './AppNavigationMenu'

const nav = useNavigationMenu()
const { t } = useI18n()

const getIconCSSAttribute = (icon: string) => `url('${icon.replace(/'/g, "\\'")}')`
</script>

<template>
  <div class="sidebar-navigation">
    <div
      v-for="(page, pageName) in nav.pages"
      :key="pageName"
      class="item"
      :class="{ active: pageName === nav.curtPage }"
      @click="nav.goPage(pageName as string)"
    >
      <span class="icon" :style="{ backgroundImage: getIconCSSAttribute(page.icon) }"></span>
      {{ t(page.label) }}
    </div>
  </div>

  <div v-if="!!Object.keys(nav.tabs).length" class="top-tabbar-wrap atk-sidebar-container">
    <div class="top-tabbar">
      <div
        v-for="(tabLabel, tabName) in nav.tabs"
        :key="tabName"
        class="item-wrap"
        :class="{ active: nav.curtTab === tabName }"
        @click="nav.goTab(tabName as string)"
      >
        <div class="item">{{ t(tabLabel) }}</div>
      </div>

      <div v-if="nav.isSearchEnabled" class="item-wrap search-btn" @click="nav.showSearch()">
        <div class="item">
          <div class="icon"></div>
        </div>
      </div>
    </div>

    <AppNavigationSearch v-if="!nav.isMobile" :state="nav.searchState" />
  </div>
</template>

<style lang="scss" scoped>
$colorMain: #8ecee2;

.sidebar-navigation {
  z-index: 10;
  width: 280px;
  position: fixed;
  left: 0;
  top: 61px;
  height: calc(100vh - 61px - 41px);
  padding: 20px;
  background: var(--at-color-bg);
  border-right: 1px solid var(--at-color-border);

  @media (max-width: 1023px) {
    display: none;
  }

  .item {
    position: relative;
    display: flex;
    align-items: center;
    padding: 6px 12px;
    font-size: 15px;
    border-radius: 6px;
    cursor: pointer;
    user-select: none;
    transition: background 0.2s;

    &:not(:last-child) {
      margin-bottom: 4px;
    }

    &.active {
      font-weight: bold;
      color: var(--at-color-deep);

      &::before {
        content: '';
        display: block;
        position: absolute;
        left: -8px;
        top: 5%;
        background: $colorMain;
        height: 90%;
        width: 4px;
        border-radius: 1000px;
      }
    }

    &:hover,
    &.active {
      background: var(--at-color-bg-grey-transl);
    }

    .icon {
      width: 20px;
      height: 20px;
      margin-right: 10px;
      background-size: contain;
      background-repeat: no-repeat;
    }
  }
}

.top-tabbar-wrap {
  @media (max-width: 1023px) {
    display: none;
  }
}

.top-tabbar {
  display: flex;
  align-items: center;
  height: 50px;
  margin: 0 auto;
  padding: 0 8px;
  border-bottom: 1px solid var(--at-color-border);

  .item-wrap {
    margin-bottom: -2px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    height: 100%;
    padding: 0;
    cursor: pointer;
    user-select: none;
    border-bottom: 2px solid transparent;

    &.active {
      font-weight: bold;
      border-color: $colorMain;
    }

    &:not(:last-child) {
      margin-right: 7px;
    }

    &:hover .item {
      background: var(--at-color-bg-grey-transl);
    }

    &.search-btn {
      margin-left: auto;

      .icon {
        background-image: url('@/assets/nav-icon-search.svg');
      }
    }

    .item {
      border-radius: 4px;
      padding: 5px 16px;
      transition: background 0.2s;

      .icon {
        width: 20px;
        height: 20px;
        background-size: contain;
        background-repeat: no-repeat;
      }
    }
  }
}
</style>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useUserStore } from '../stores/user'
import { useNavStore } from '../stores/nav'
import { isOpenFromSidebar } from '../global'

const nav = useNavStore()
const router = useRouter()
const user = useUserStore()
const { t } = useI18n()
const { site: curtSite, is_admin: isAdmin, avatar } = storeToRefs(user)
const { darkMode } = storeToRefs(nav)

const avatarClickHandler = () => {
  if (!isOpenFromSidebar()) logout()
}

const logout = () => {
  if (!window.confirm(t('logoutConfirm'))) return

  useUserStore().logout()
  nextTick(() => {
    router.replace('/login')
  })
}
</script>

<template>
  <div class="header">
    <template v-if="isAdmin">
      <div
        class="avatar clickable"
        :class="{ active: nav.siteSwitcherShow }"
        @click="nav.showSiteSwitcher()"
      >
        <div class="site">{{ curtSite.substring(0, 1) || '_' }}</div>
      </div>
    </template>
    <template v-else>
      <div class="avatar" @click="avatarClickHandler">
        <img :src="avatar" />
      </div>
    </template>

    <div class="title">
      <div class="text show-mobile">{{ isAdmin ? t('ctrlCenter') : t('msgCenter') }}</div>
      <div class="text show-desktop">
        <template v-if="!isAdmin">{{ t('msgCenter') }}</template>
        <template v-else-if="!!curtSite">{{ curtSite }}</template>
        <template v-else>
          <img src="../assets/favicon.png" class="artalk-logo" draggable="false" />
        </template>
      </div>
    </div>

    <div class="close-btn"></div>
    <div
      class="dark-mode-toggle"
      :data-darkmode="darkMode ? 'on' : 'off'"
      @click="nav.toggleDarkMode()"
    ></div>
  </div>
  <SiteSwitcher v-if="isAdmin" />
</template>

<style scoped lang="scss">
.header {
  display: flex;
  flex-direction: row;
  align-items: center;
  border-bottom: 1px solid var(--at-color-border);

  @media (min-width: 1024px) {
    background: var(--at-sidebar-header-bg);
  }

  .avatar {
    position: relative;
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-right: 1px solid transparent;
    padding: 0 10px 0 15px;

    .site {
      user-select: none;
      color: #fff;
      background: #697182;
      height: 30px;
      width: 30px;
      border-radius: 2px;
      text-align: center;
      line-height: 30px;
      font-size: 13px;
    }

    &.clickable {
      cursor: pointer;

      &:hover,
      &.active {
        background: var(--at-color-bg-grey);
        border-right: 1px solid var(--at-color-border);
      }

      &::after {
        content: '';
        margin-left: 10px;
        vertical-align: middle;
        border-top: 5px solid #747474;
        border-left: 3px solid transparent;
        border-right: 3px solid transparent;
        margin-top: -1px;
        display: inline-block;
      }
    }

    img {
      height: 30px;
      width: 30px;
      border-radius: 2px;
      background: #697182;
      text-align: center;
      line-height: 30px;
      font-size: 13px;
      color: #fff;
    }
  }

  .title {
    flex: auto;
    text-align: center;
    user-select: none;

    .text {
      display: inline-block;
      position: relative;
      color: var(--at-color-deep);
      font-size: 20px;

      &.show-mobile {
        display: none;
      }

      &.show-desktop {
        display: inline-block;

        &::after {
          display: none;
        }
      }

      @media (max-width: 1023px) {
        &.show-mobile {
          display: inline-block;
        }

        &.show-desktop {
          display: none;
        }
      }

      &::after {
        content: '';
        position: absolute;
        bottom: 0;
        left: 0;
        width: 100%;
        height: 6px;
        background: #0083ff;
        opacity: 0.4;
      }
    }

    .artalk-logo {
      border-radius: 2px;
      height: 40px;
      width: 40px;
      background: #697182;
      text-align: center;
      line-height: 30px;
      font-size: 13px;
      color: #fff;
    }
  }

  .close-btn {
    width: 60px;
    height: 60px;
    margin-left: 10px;
  }

  .dark-mode-toggle {
    display: none;
    width: 60px;
    height: 60px;
    margin-left: 10px;
    justify-content: center;
    align-items: center;
    user-select: none;
    cursor: pointer;

    @media (min-width: 1024px) {
      display: flex;
    }

    &:hover {
      &::after {
        background-color: var(--at-color-bg-grey);
      }
    }

    &::after {
      content: '';
      display: block;
      transition: background-color 0.2s;
      width: 35px;
      height: 35px;
      background-repeat: no-repeat;
      background-position: center;
      background-size: 60%;
      border-radius: 50%;
    }

    &[data-darkmode='on']::after {
      background-image: url('@/assets/icon-darkmode-on.svg');
    }

    &[data-darkmode='off']::after {
      background-image: url('@/assets/icon-darkmode-off.svg');
    }
  }
}
</style>

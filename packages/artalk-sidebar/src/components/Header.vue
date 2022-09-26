<script setup lang="ts">
import MD5 from '../lib/md5'
import { storeToRefs } from 'pinia'
import { useUserStore } from '../stores/user'
import { useNavStore } from '../stores/nav'
import { artalk } from '../global'

const nav = useNavStore()
const user = useUserStore()
const { site: curtSite, isAdmin, email } = storeToRefs(user)

const userAvatarImgURL = computed(() =>
  `${(artalk?.ctx.conf.gravatar.mirror || '').replace(/\/$/, '')}/${MD5(email.value)}`
  + `?d=${encodeURIComponent(artalk?.ctx.conf.gravatar.default || 'mp')}&s=80`)
</script>

<template>
  <div class="header">
    <template v-if="isAdmin">
      <div
        class="avatar clickable"
        :class="{ 'active': nav.siteSwitcherShow }"
        @click="nav.showSiteSwitcher()"
      >
        <div class="site">{{ curtSite.substring(0, 1) }}</div>
      </div>
    </template>
    <template v-else>
      <div class="avatar">
        <img :src="userAvatarImgURL">
      </div>
    </template>

    <div class="title">
      <div class="text">{{ isAdmin ? '控制中心' : '通知中心' }}</div>
    </div>

    <div class="close-btn"></div>
  </div>
  <SiteSwitcher v-if="isAdmin" />
</template>

<style scoped lang="scss">
.header {
  display: flex;
  flex-direction: row;
  align-items: center;
  border-bottom: 1px solid #eceff2;

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

      &:hover, &.active {
        background: #f4f4f4;
        border-right: 1px solid #eceff2;
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
      color: #FFF;
    }
  }

  .title {
    flex: auto;
    text-align: center;
    user-select: none;

    .text {
      display: inline-block;
      position: relative;
      color: #2a2e2e;
      font-size: 20px;

      &::after {
        content: "";
        position: absolute;
        bottom: 0;
        left: 0;
        width: 100%;
        height: 6px;
        background: #0083ff;
        opacity: .4;
      }
    }
  }

  .close-btn {
    width: 60px;
    height: 60px;
    margin-left: 10px;
  }
}
</style>

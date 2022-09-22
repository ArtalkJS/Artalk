<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useUserStore } from '../stores/user'
import { useNavStore } from '../stores/nav'

const nav = useNavStore()
const user = useUserStore()
const { site: curtSite } = storeToRefs(user)
</script>

<template>
  <div class="header">
    <div class="avatar clickable" :class="{ 'active': nav.siteSwitcherShow }" @click="nav.showSiteSwitcher()">
      <div class="site">{{ curtSite.substring(0, 1) }}</div>
    </div>
    <div class="title">
      <div class="text">控制中心</div>
    </div>
    <div class="close-btn"></div>
  </div>
  <SiteSwitcher />
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

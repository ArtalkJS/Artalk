<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { artalk } from '../global'
import type { PageData } from 'artalk/types/artalk-data'
import { useNavStore } from '../stores/nav'
import { useUserStore } from '../stores/user'

const nav = useNavStore()
const { site } = storeToRefs(useUserStore())
const pages = ref<PageData[]>([])
const curtEditPageID = ref<number|null>(null)

onMounted(() => {
  nav.updateTabs({

  }, '')

  artalk?.ctx.getApi().page.pageGet(site.value).then(data => {
    pages.value = data.pages
  })
})

function editPage(page: PageData) {
  curtEditPageID.value = page.id
}
</script>

<template>
  <div class="atk-page-list-wrap">
    <div class="atk-header-action-bar">
    <span class="atk-update-all-title-btn"><i class="atk-icon atk-icon-sync"></i> <span class="atk-text">更新标题</span></span>
    <span class="atk-cache-flush-all-btn"><span class="atk-text">缓存清除</span></span>
    <span class="atk-cache-warm-up-btn"><span class="atk-text">缓存预热</span></span>
    </div>
    <div class="atk-page-list">
      <div v-for="(page) in pages" class="atk-page-item">
        <div class="atk-page-main">
          <div class="atk-title">{{ page.title }}</div>
          <div class="atk-sub">{{ page.url }}</div>
        </div>
        <div class="atk-page-actions">
          <div class="atk-item atk-edit-btn" @click="editPage(page)">
            <i class="atk-icon atk-icon-edit"></i>
          </div>
        </div>
        <PageEditor v-if="curtEditPageID === page.id" :page="page" @close="curtEditPageID = null" />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.atk-page-list-wrap {
  .atk-header-action-bar {
    display: block;
    overflow: hidden;
    padding: 10px 15px 0 15px;

    & > span {
      display: inline-block;
      padding: 2px 10px;
      cursor: pointer;
      font-size: 13px;

      i {
        display: inline-block;
        width: 16px;
        height: 16px;
        vertical-align: middle;
        background-color: var(--at-color-meta);
        margin-right: 4px;
      }

      &:hover {
        background: var(--at-color-bg-grey);
      }
    }
  }
}

.atk-page-list {
  .atk-page-item {
    display: flex;
    flex-direction: row;
    position: relative;
    min-height: 120px;
    align-items: center;

    &:not(:last-child) {
      border-bottom: 1px solid var(--at-color-border);
    }
  }

  .atk-page-main {
    display: flex;
    flex-direction: column;
    flex: auto;
    padding: 20px 30px;

    .atk-title {
      color: var(--at-color-font);
      font-size: 21px;
      margin-bottom: 10px;
      cursor: pointer;
    }

    .atk-sub {
      color: var(--at-color-sub);
      font-size: 14px;
      cursor: pointer;
    }
  }

  :deep(.atk-page-actions) {
    @extend .atk-list-btn-actions;
  }
}
</style>

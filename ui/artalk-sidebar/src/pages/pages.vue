<script setup lang="ts">
import { storeToRefs } from 'pinia'
import type { ArtalkType } from 'artalk'
import { artalk } from '../global'
import { useNavStore } from '../stores/nav'
import { useUserStore } from '../stores/user'
import Pagination from '../components/Pagination.vue'

const nav = useNavStore()
const { site: curtSite } = storeToRefs(useUserStore())
const pages = ref<ArtalkType.PageData[]>([])
const curtEditPageID = ref<number | null>(null)
const { t } = useI18n()

const pageSize = ref(20)
const pageTotal = ref(0)
const search = ref('')
const pagination = ref<InstanceType<typeof Pagination>>()
const showActBarBorder = ref(false)
const refreshBtn = ref({
  isRun: false,
  statusText: '',
})
let pageRequest: AbortController | undefined
let refreshRequest: AbortController | undefined
let refreshTimer: number | undefined

onBeforeUnmount(() => {
  pageRequest?.abort()
  cancelRefreshTask()
  nav.setPageLoading(false)
})

onMounted(() => {
  nav.updateTabs(
    {
      all: 'all',
    },
    'all',
  )

  // Users search
  nav.enableSearch(
    (value: string) => {
      search.value = value
      fetchPages(0)
    },
    () => {
      if (search.value === '') return
      search.value = ''
      fetchPages(0)
    },
  )

  fetchPages(0)

  // Refresh task status recovery
  watchRefreshTask()
})

watch(curtSite, () => {
  pagination.value?.reset()
  fetchPages(0)
})

onMounted(() => nav.scrollableArea?.addEventListener('scroll', scrollHandler))
onUnmounted(() => nav.scrollableArea?.removeEventListener('scroll', scrollHandler))

function scrollHandler() {
  showActBarBorder.value = nav.scrollableArea!.scrollTop > 10
}

function editPage(page: ArtalkType.PageData) {
  curtEditPageID.value = page.id
}

function fetchPages(offset: number) {
  pageRequest?.abort()
  const request = (pageRequest = new AbortController())
  if (offset === 0) pagination.value?.reset()
  nav.setPageLoading(true)
  artalk?.ctx
    .getApi()
    .pages.getPages(
      {
        site_name: curtSite.value,
        offset: offset,
        limit: pageSize.value,
        search: search.value,
        has_comments: true,
      },
      { signal: request.signal },
    )
    .then((res) => {
      if (request.signal.aborted) return
      pageTotal.value = res.data.count
      pages.value = res.data.pages
      nav.scrollPageToTop()
    })
    .catch((err) => {
      if (!request.signal.aborted) alert(err.message || t('opFailed'))
    })
    .finally(() => {
      if (!request.signal.aborted) nav.setPageLoading(false)
    })
}

function onChangePage(offset: number) {
  fetchPages(offset)
}

function onPageItemUpdate(page: ArtalkType.PageData) {
  const index = pages.value.findIndex((p) => p.id === page.id)
  if (index != -1) {
    const orgPage = pages.value[index]
    Object.keys(page).forEach((key) => {
      ;(orgPage as any)[key] = (page as any)[key]
    })
  }
}

function onPageItemRemove(id: number) {
  const index = pages.value.findIndex((p) => p.id === id)
  pages.value.splice(index, 1)
}

function cancelRefreshTask() {
  refreshRequest?.abort()
  window.clearTimeout(refreshTimer)
}

async function watchRefreshTask() {
  cancelRefreshTask()
  const request = (refreshRequest = new AbortController())
  try {
    const { data } = await artalk!.ctx.getApi().pages.getPageFetchStatus({ signal: request.signal })
    if (request.signal.aborted) return
    if (data.is_progress) {
      refreshBtn.value = { isRun: true, statusText: data.msg }
      refreshTimer = window.setTimeout(watchRefreshTask, 1000)
    } else if (refreshBtn.value.isRun) {
      setRefreshTaskDone()
    }
  } catch (err) {
    if (request.signal.aborted) return
    console.error(err)
    setRefreshTaskDone(t('opFailed'))
  }
}

function setRefreshTaskDone(text = t('updateComplete')) {
  refreshBtn.value.statusText = text
  refreshTimer = window.setTimeout(() => {
    refreshBtn.value.isRun = false
  }, 1500)
}

async function refreshAllPages() {
  if (refreshBtn.value.isRun) return
  cancelRefreshTask()
  const request = (refreshRequest = new AbortController())
  refreshBtn.value.isRun = true
  refreshBtn.value.statusText = t('updateReady')

  try {
    await artalk!.ctx
      .getApi()
      .pages.fetchAllPages({ site_name: curtSite.value }, { signal: request.signal })
  } catch (err: any) {
    if (request.signal.aborted) return
    alert(err.msg)
    setRefreshTaskDone(t('opFailed'))
    return
  }

  if (!request.signal.aborted) watchRefreshTask()
}

function cacheFlush() {
  artalk!.ctx
    .getApi()
    .cache.flushCache()
    .then((res) => alert(res.data.msg))
    .catch(() => alert(t('opFailed')))
}

function cacheWarm() {
  artalk!.ctx
    .getApi()
    .cache.warmUpCache()
    .then((res) => alert(res.data.msg))
    .catch(() => alert(t('opFailed')))
}

function openPage(url: string) {
  window.open(url)
}
</script>

<template>
  <div class="atk-page-list-wrap">
    <div class="atk-header-action-bar" :class="{ bordered: showActBarBorder }">
      <span class="atk-update-all-title-btn" @click="refreshAllPages()">
        <i class="atk-icon atk-icon-sync" :class="{ 'atk-rotate': refreshBtn.isRun }"></i>
        <span class="atk-text">
          {{ refreshBtn.isRun ? refreshBtn.statusText : t('updateTitle') }}
        </span>
      </span>
      <span class="atk-cache-flush-all-btn" @click="cacheFlush()">
        <span class="atk-text">{{ t('cacheClear') }}</span>
      </span>
      <span class="atk-cache-warm-up-btn" @click="cacheWarm()">
        <span class="atk-text">{{ t('cacheWarm') }}</span>
      </span>
    </div>
    <div class="atk-page-list">
      <div v-for="page in pages" :key="page.id" class="atk-page-item">
        <div class="atk-page-main">
          <div class="atk-title" @click="openPage(page.url)">
            {{ page.title }}
          </div>
          <div class="atk-sub" @click="openPage(page.url)">{{ page.url }}</div>
        </div>
        <div class="atk-page-actions">
          <div class="atk-item atk-edit-btn" @click="editPage(page)">
            <i class="atk-icon atk-icon-edit"></i>
          </div>
        </div>
        <PageEditor
          v-if="curtEditPageID === page.id"
          :page="page"
          @close="curtEditPageID = null"
          @update="onPageItemUpdate"
          @remove="onPageItemRemove"
        />
      </div>
    </div>
    <Pagination
      ref="pagination"
      :page-size="pageSize"
      :total="pageTotal"
      :disabled="nav.isPageLoading"
      @change="onChangePage"
    />
  </div>
</template>

<style scoped lang="scss">
.atk-page-list-wrap {
  .atk-header-action-bar {
    position: sticky;
    top: 0;
    display: flex;
    align-items: center;
    overflow: hidden;
    padding: 10px 15px 0 15px;
    background: var(--at-color-bg);
    z-index: 10;
    border-bottom: 1px solid transparent;
    transition: 0.3s ease-out padding;

    &.bordered {
      padding-bottom: 10px;
      border-color: var(--at-color-border);
    }

    & > span {
      display: inline-flex;
      align-items: center;
      flex-direction: row;
      padding: 2px 10px;
      cursor: pointer;
      font-size: 13px;

      i {
        display: inline-block;
        width: 14px;
        height: 14px;
        margin-right: 5px;

        &::after {
          background-color: var(--at-color-meta);
        }
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

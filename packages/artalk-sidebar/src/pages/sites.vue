<script setup lang="ts">
import { useNavStore } from '../stores/nav'
import { artalk, bootParams } from '../global'
import type { SiteData } from 'artalk/types/artalk-data'

const nav = useNavStore()
const sites = ref<SiteData[]>([])
const curtEditSite = ref<SiteData|null>(null)
const showSiteCreate = ref(false)
const siteCreateInitVal = ref()

onMounted(() => {
  nav.updateTabs({

  }, '站点')

  nav.setPageLoading(true)
  artalk?.ctx.getApi().site.siteGet().then(gotSites => {
    sites.value = gotSites
  }).finally(() => {
    nav.setPageLoading(false)
  })

  // 通过启动参数打开站点创建
  const vp = bootParams.viewParams
  if (vp && vp.create_name && vp.create_urls) {
    siteCreateInitVal.value = { name: vp.create_name, urls: vp.create_urls }
    showSiteCreate.value = true
    nextTick(() => {
      siteCreateInitVal.value = null
      bootParams.viewParams = null
    })
  }
})

function create() {
  curtEditSite.value = null
  showSiteCreate.value = true
}

const sitesGrouped = computed(() => {
  if (sites.value.length === 0) return []

  const grp: SiteData[][] = []
  let j = -1
  for (let i = 0; i < sites.value.length; i++) {
    const item = sites.value[i]
    if (i % 4 === 0) { // 每 4 个一组
      grp.push([])
      j++
    }
    grp[j].push(item)
  }
  return grp
})

function edit(site: SiteData) {
  showSiteCreate.value = false
  curtEditSite.value = site
}

function onNewSiteCreated(siteNew: SiteData) {
  sites.value.push(siteNew)
  showSiteCreate.value = false
  nav.refreshSites()
}

function onSiteItemUpdate(site: SiteData) {
  const index = sites.value.findIndex(s => s.id === site.id)
  if (index != -1) {
    const orgSite = sites.value[index]
    Object.keys(site).forEach(key => {
      ;(orgSite as any)[key] = (site as any)[key]
    })
  }
  nav.refreshSites()
}

function onSiteItemRemove(id: number) {
  const index = sites.value.findIndex(p => p.id === id)
  sites.value.splice(index, 1)
  nav.refreshSites()
}
</script>

<template>
  <div class="atk-site-list">
    <div class="atk-header">
      <div class="atk-title">共 {{ sites.length }} 个站点</div>
      <div class="atk-actions">
        <div class="atk-item atk-site-add-btn" @click="create()"><i class="atk-icon atk-icon-plus"></i></div>
      </div>
    </div>
    <SiteCreate
      v-if="showSiteCreate"
      :init-val="siteCreateInitVal"
      @close="showSiteCreate = false"
      @done="onNewSiteCreated"
    />
    <div class="atk-site-rows-wrap">
      <template v-for="(sites) in sitesGrouped">
        <template v-if="curtEditSite !== null">
          <SiteEditor
            v-if="!!sites.includes(curtEditSite)"
            :site="curtEditSite"
            @close="curtEditSite = null"
            @update="onSiteItemUpdate"
            @remove="onSiteItemRemove"
          />
        </template>
        <div class="atk-site-row">
          <div
            v-for="(site) in sites"
            class="atk-site-item"
            :class="{ 'atk-active': curtEditSite === site }"
            @click="edit(site)"
          >
            <div class="atk-site-logo">{{ site.name.substring(0, 1) }}</div>
            <div class="atk-site-name">{{ site.name }}</div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped lang="scss">
.atk-site-list {
  & > .atk-header {
    display: flex;
    flex-direction: row;
    padding: 15px 30px;
    align-items: center;

    .atk-title {
      flex: auto;
      padding-right: 10px;
    }

    .atk-actions {
      display: flex;
      flex-direction: row;

      .atk-item {
        display: flex;
        height: 30px;
        width: 30px;
        justify-content: center;
        align-items: center;
        user-select: none;
        cursor: pointer;
        border-radius: 2px;

        &:hover {
          background: var(--at-color-bg-grey);
        }
      }
    }
  }

  .atk-site-rows-wrap {
    position: relative;

    .atk-site-row {
      display: flex;
      flex-direction: row;
      padding: 10px 20px;
      padding-bottom: 0;
    }

    .atk-site-item {
      display: flex;
      flex-basis: 25%;
      flex-direction: column;
      align-items: center;
      padding-bottom: 5px;
      user-select: none;
      cursor: pointer;
      border: 1px solid transparent;

      .atk-site-logo {
        margin: 15px;
        text-align: center;
        font-size: 20px;
        height: 65px;
        width: 65px;
        line-height: 65px;
        background: #687a86;
        color: #fff;
        border-radius: 4px;
      }

      .atk-site-name {
        text-align: center;
        font-size: 15px;
        color: var(--at-color-sub);
        padding: 0 17px;
        word-break: break-word;
      }

      &.atk-active {
        background-color: var(--at-color-bg-grey);
        border: 1px solid var(--at-color-border);
        margin-top: -1px;
        border-radius: 0 0 4px 4px;

        .atk-site-name {
          color: var(--at-color-deep);
        }
      }

      &:hover {
        .atk-site-name {
          color: var(--at-color-font);
        }
      }
    }
  }

  :deep(.atk-site-edit), :deep(.atk-site-add) {
    position: relative;
    min-height: 120px;
    width: 100%;
    border-top: 1px solid var(--at-color-border);
    border-bottom: 1px solid var(--at-color-border);
    margin-bottom: -10px;

    .atk-header {
      display: flex;
      flex-direction: row;
      align-items: center;
      padding: 10px 30px 0 35px;
      justify-content: space-between;

      .atk-site-info {
        .atk-site-name {
          cursor: pointer;
          display: inline-block;
          font-size: 23px;
          position: relative;
          line-height: 1.6em;

          &:after {
            content: ' ';
            position: absolute;
            width: 100%;
            height: 6px;
            background: var(--at-color-main);
            opacity: .4;
            left: 0;
            bottom: 6px;
          }
        }

        .atk-site-urls {
          display: flex;
          width: 100%;
          margin-top: 6px;
          flex-wrap: wrap;
          min-height: 23px;
          margin-bottom: 15px;

          .atk-url-item {
            background: var(--at-color-bg-grey);
            color: var(--at-color-font);
            border-radius: 2px;
            padding: 0 8px;
            font-size: 13px;
            margin-bottom: 3px;
            margin-right: 3px;
            cursor: pointer;

            &:hover {
            }
          }
        }
      }

      .atk-close-btn {
        width: 50px;
        height: 50px;
        display: flex;
        justify-content: center;
        align-items: center;
        cursor: pointer;

        &:hover i {
          background-color: var(--at-color-red);
        }
      }
    }

    .atk-main {
      position: relative;
      display: flex;
      flex-direction: row;
      padding: 0 30px 6px 35px;
      padding-bottom: 10px;

      .atk-site-text-actions {
        @extend .atk-list-text-actions;
        height: 90px;
        padding: 0;
        padding-left: 10px;

        .atk-item {
          margin-bottom: 17px;
          margin-right: 25px;
        }
      }

      .atk-site-btn-actions {
        @extend .atk-list-btn-actions;

        padding-right: 9px;
      }

      .atk-item-text-editor-layer {
        padding: 10px 20px;
      }
    }
  }

  :deep(.atk-site-add) {
    position: relative;

    .atk-header {
      .atk-title {
        font-size: 20px;
      }
    }

    .atk-form {
      padding: 20px 40px;
    }
  }
}
</style>

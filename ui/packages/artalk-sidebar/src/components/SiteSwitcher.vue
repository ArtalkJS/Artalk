<script setup lang="ts">
import type { SiteData } from 'artalk/types/artalk-data'
import { storeToRefs } from 'pinia';
import { useNavStore } from '../stores/nav'
import { useUserStore } from '../stores/user'

const el = ref<HTMLElement|null>(null)

const nav = useNavStore()
const { siteSwitcherShow: curtShow, sites } = storeToRefs(nav)
const { site: curtSite } = storeToRefs(useUserStore())

const router = useRouter()

onMounted(() => {
  nav.refreshSites()
})

interface IDisplaySite {
  label: string
  name: string
  logoText: string
}

function switchSite(siteName: string) {
  if (siteName === '__SITE_MANAGEMENT__') {
    router.replace('/sites')
  } else {
    curtSite.value = siteName
  }

  curtShow.value = false
}

const displaySites = computed(() => {
  const displays: IDisplaySite[] = []
  displays.push({ label: '所有站点', name: '__ATK_SITE_ALL', logoText: '_' })
  sites.value.forEach((site) => {
    displays.push({
      label: site.name, name: site.name,
      logoText: site.name.substring(0, 1)
    })
  })
  displays.push({ label: '站点管理', name: '__SITE_MANAGEMENT__', logoText: '+' })
  return displays
})

function outsideChecker(evt: MouseEvent) {
  const isClickInside = el.value?.contains(evt.target as any)
  if (!isClickInside) {
    curtShow.value = false
  }
}

watch(curtShow, (value) => {
  if (value) {
    setTimeout(() => {
      document.addEventListener('click', outsideChecker)
    }, 80)
  } else {
    document.removeEventListener('click', outsideChecker)
  }
})
</script>

<template>
  <Transition>
    <div ref="el" v-show="curtShow" class="atk-site-list-floater">
      <div class="atk-sites">
        <div
          v-for="(site) in displaySites"
          class="atk-site-item"
          :class="{ 'atk-active': curtSite === site.name }"
          @click="switchSite(site.name)"
        >
          <div class="atk-site-logo">{{ site.logoText }}</div>
          <div class="atk-site-name">{{ site.label }}</div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped lang="scss">
.atk-site-list-floater {
  @extend .atk-slim-scrollbar;

  position: absolute;
  z-index: 9999;
  left: calc((100% - 80%) / 2);
  top: 70px;
  pointer-events: all;
  background: #FFF;
  width: 80%;
  max-height: 40%;
  margin: 0 auto;
  border: 1px solid #eceff2;
  border-radius: 4px;
  overflow-y: auto;
  transition: all 0.2s ease;

  .atk-sites {
    display: flex;
    flex-direction: column;
    padding: 7px 0;

    .atk-site-item {
      display: flex;
      flex-direction: row;
      cursor: pointer;
      align-items: center;
      padding: 0 10px;
      margin-bottom: 1px;

      .atk-site-logo {
        width: 20px;
        height: 20px;
        line-height: 20px;
        background: #697182;
        margin: 10px;
        border-radius: 3px;
        text-align: center;
        color: #FFF;
        font-size: 12px;
      }

      .atk-site-name {
        color: #2a2e2e;
        font-size: 17px;
        margin-left: 7px;
      }

      &.atk-active, &:hover {
        background: #F4F4F4;
      }
    }
  }
}

.v-enter-from, .v-leave-to {
  opacity: 0;
  transform: scale3d(1.05, 1.05, 1.05);
}
</style>

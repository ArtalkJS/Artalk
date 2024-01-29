<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { useNavStore } from '../stores/nav'
import { useUserStore } from '../stores/user'
import { bootParams } from '@/global'

const el = ref<HTMLElement|null>(null)

const nav = useNavStore()
const { siteSwitcherShow: curtShow, sites } = storeToRefs(nav)
const { site: curtSite } = storeToRefs(useUserStore())

const router = useRouter()
const { t } = useI18n()

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
  displays.push({ label: t('allSites'), name: '', logoText: '_' })
  sites.value.forEach((site) => {
    displays.push({
      label: site.name, name: site.name,
      logoText: site.name.substring(0, 1)
    })
  })
  displays.push({ label: t('siteManage'), name: '__SITE_MANAGEMENT__', logoText: '+' })
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

function logout() {
  useUserStore().logout()
  nextTick(() => {
    router.replace('/login')
  })
}
</script>

<template>
  <Transition>
    <div ref="el" v-show="curtShow" class="atk-site-list-floater">
      <div class="atk-sites">
        <div
          v-for="(site, i) in displaySites"
          :key="i"
          class="atk-site-item"
          :class="{ 'atk-active': curtSite === site.name }"
          @click="switchSite(site.name)"
        >
          <div class="atk-site-logo">{{ site.logoText }}</div>
          <div class="atk-site-name">{{ site.label }}</div>
        </div>

        <!-- Logout Button -->
        <div v-if="!bootParams.user?.email" class="atk-site-item" @click="logout()">
          <svg class="atk-site-logo" stroke="currentColor" fill="currentColor" stroke-width="0" viewBox="-6 -6 36 36" xmlns="http://www.w3.org/2000/svg"><path d="M5 22C4.44772 22 4 21.5523 4 21V3C4 2.44772 4.44772 2 5 2H19C19.5523 2 20 2.44772 20 3V6H18V4H6V20H18V18H20V21C20 21.5523 19.5523 22 19 22H5ZM18 16V13H11V11H18V8L23 12L18 16Z"></path></svg>
          <div class="atk-site-name">{{ $t('logout') }}</div>
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
  background: var(--at-color-bg);
  width: 80%;
  max-height: 40%;
  margin: 0 auto;
  border: 1px solid var(--at-color-border);
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
        color: var(--at-color-deep);
        font-size: 17px;
        margin-left: 7px;
      }

      &.atk-active, &:hover {
        background: var(--at-color-bg-grey);
      }
    }
  }
}

.v-enter-from, .v-leave-to {
  opacity: 0;
  transform: scale3d(1.05, 1.05, 1.05);
}
</style>

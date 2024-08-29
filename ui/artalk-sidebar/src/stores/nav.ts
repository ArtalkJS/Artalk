import { defineStore } from 'pinia'
import type { ArtalkType } from 'artalk'
import { artalk, bootParams, getArtalk } from '@/global'
import { useMobileWidth } from '@/hooks/MobileWidth'

type TabsObj = { [name: string]: string }

export const useNavStore = defineStore('nav', () => {
  const curtTab = ref('')
  const tabs = ref<TabsObj>({})

  const curtPage = ref('comments')
  const sites = ref<ArtalkType.SiteData[]>([])
  const siteSwitcherShow = ref(false)

  const isSearchEnabled = ref(false)
  const searchEvent = ref<((val: string) => void) | null>(null)
  const searchResetEvent = ref<(() => void) | null>(null)

  const isPageLoading = ref(false)
  const scrollableArea = ref<HTMLElement | null>(null)

  const darkMode = ref(bootParams.darkMode)
  watch(darkMode, (val) => {
    getArtalk()?.setDarkMode(val)
    if (val != window.matchMedia('(prefers-color-scheme: dark)').matches)
      localStorage.setItem('ATK_SIDEBAR_DARK_MODE', val ? '1' : '0')
    else localStorage.removeItem('ATK_SIDEBAR_DARK_MODE') // enable auto switch
  })

  const updateTabs = (aTabs: TabsObj, activeTab?: string) => {
    tabs.value = aTabs
    if (activeTab) curtTab.value = activeTab
  }

  const setTabActive = (tabName: string) => {
    curtTab.value = tabName
  }

  const showSiteSwitcher = () => {
    siteSwitcherShow.value = true
  }

  const hideSiteSwitcher = () => {
    siteSwitcherShow.value = false
  }

  const toggleSiteSwitcher = () => {
    siteSwitcherShow.value = !siteSwitcherShow.value
  }

  const scrollPageToTop = () => {
    scrollableArea.value?.scrollTo(0, 0)
  }

  const scrollPageToEl = (el: HTMLElement) => {
    scrollableArea.value?.scrollTo(0, el.offsetTop)
  }

  const setPageLoading = (pageLoading: boolean) => {
    isPageLoading.value = pageLoading
  }

  const refreshSites = () => {
    artalk?.ctx
      .getApi()
      .sites.getSites()
      .then((res) => {
        sites.value = res.data.sites
      })
  }

  const enableSearch = (searchEvt: (val: string) => void, searchResetEvt: () => void) => {
    isSearchEnabled.value = true
    searchEvent.value = searchEvt
    searchResetEvent.value = searchResetEvt
  }

  const toggleDarkMode = () => {
    darkMode.value = !darkMode.value
  }

  useRouter().beforeEach((to, from) => {
    isSearchEnabled.value = false
  })

  const isMobile = useMobileWidth()

  return {
    sites,
    curtPage,
    curtTab,
    tabs,
    siteSwitcherShow,
    scrollableArea,
    isPageLoading,
    updateTabs,
    setTabActive,
    showSiteSwitcher,
    hideSiteSwitcher,
    toggleSiteSwitcher,
    scrollPageToTop,
    scrollPageToEl,
    setPageLoading,
    refreshSites,
    isSearchEnabled,
    searchEvent,
    searchResetEvent,
    enableSearch,
    isMobile,
    darkMode,
    toggleDarkMode,
  }
})

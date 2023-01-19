import { defineStore } from 'pinia'
import { artalk } from '../global'
import type { SiteData } from 'artalk/types/artalk-data'

type TabsObj = {[name: string]: string}

export const useNavStore = defineStore('nav', () => {
  const curtTab = ref('')
  const tabs = ref<TabsObj>({})

  const curtPage = ref('comments')
  const sites = ref<SiteData[]>([])
  const siteSwitcherShow = ref(false)

  const isSearchEnabled = ref(false)
  const searchEvent = ref<((val: string) => void)|null>(null)
  const searchResetEvent = ref<(() => void)|null>(null)

  const isPageLoading = ref(false)
  const scrollableArea = ref<HTMLElement|null>(null)

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
    artalk?.ctx.getApi().site.siteGet().then((respSites) => {
      sites.value = respSites
    })
  }

  const enableSearch = (searchEvt: ((val: string) => void), searchResetEvt: () => void) => {
    isSearchEnabled.value = true
    searchEvent.value = searchEvt
    searchResetEvent.value = searchResetEvt
  }

  return {
    sites, curtPage, curtTab, tabs, siteSwitcherShow, scrollableArea, isPageLoading,
    updateTabs, setTabActive,
    showSiteSwitcher, hideSiteSwitcher, toggleSiteSwitcher,
    scrollPageToTop, scrollPageToEl, setPageLoading,
    refreshSites,
    isSearchEnabled, searchEvent, searchResetEvent, enableSearch,
  }
})

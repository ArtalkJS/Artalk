import { defineStore } from 'pinia'

type TabsObj = {[name: string]: string}

export const useNavStore = defineStore('nav', () => {
  const curtPage = ref('comments')
  const curtTab = ref('')
  const tabs = ref<TabsObj>({})
  const siteSwitcherShow = ref(false)
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

  const scrollToTop = () => {
    scrollableArea.value?.scrollTo(0, 0)
  }

  return {
    curtPage, curtTab, tabs, siteSwitcherShow, scrollableArea,
    updateTabs, setTabActive,
    showSiteSwitcher, hideSiteSwitcher, toggleSiteSwitcher,
    scrollToTop,
  }
})

import { defineStore } from 'pinia'

type TabsObj = {[name: string]: string}

export const useNavStore = defineStore('nav', () => {
  const curtPage = ref('comments')
  const curtTab = ref('')
  const tabs = ref<TabsObj>({})

  const updateTabs = (aTabs: TabsObj, activeTab?: string) => {
    tabs.value = aTabs
    if (activeTab) curtTab.value = activeTab
  }

  const setTabActive = (tabName: string) => {
    curtTab.value = tabName
  }

  return { curtPage, curtTab, tabs, updateTabs, setTabActive }
})

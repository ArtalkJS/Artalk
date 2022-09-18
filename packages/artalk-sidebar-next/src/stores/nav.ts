import { defineStore } from 'pinia'

type TabsObj = {[name: string]: string}

export const useNavStore = defineStore('nav', () => {
  const curtPage = ref('comments')
  const curtTab = ref('')
  const tabs = ref<TabsObj>({})

  const updateTabs = (aTabs: TabsObj) => {
    tabs.value = aTabs
  }

  const setTabActive = (tabName: string) => {
    curtTab.value = tabName
  }

  return { curtPage, curtTab, tabs, updateTabs, setTabActive }
})

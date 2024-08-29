import { storeToRefs } from 'pinia'
import { RouteLocation } from 'vue-router'
import CommentsIcon from '@/assets/nav-icon-comments.svg'
import PagesIcon from '@/assets/nav-icon-pages.svg'
import UsersIcon from '@/assets/nav-icon-users.svg'
import SitesIcon from '@/assets/nav-icon-sites.svg'
import TransferIcon from '@/assets/nav-icon-transfer.svg'
import SettingsIcon from '@/assets/nav-icon-settings.svg'
import { useUserStore } from '@/stores/user'
import { useNavStore } from '@/stores/nav'

export type PageItem = { label: string; link: string; hideOnMobile?: boolean; icon: string }

/**
 * Pages for admin
 */
export const AdminPages: Record<string, PageItem> = {
  comments: {
    label: 'comment',
    link: '/comments',
    icon: CommentsIcon,
  },
  pages: {
    label: 'page',
    link: '/pages',
    icon: PagesIcon,
  },
  users: {
    label: 'user',
    link: '/users',
    icon: UsersIcon,
  },
  sites: {
    label: 'site',
    link: '/sites',
    hideOnMobile: true,
    icon: SitesIcon,
  },
  transfer: {
    label: 'transfer',
    link: '/transfer',
    hideOnMobile: true,
    icon: TransferIcon,
  },
  settings: {
    label: 'settings',
    link: '/settings',
    icon: SettingsIcon,
  },
}

/**
 * Pages for user
 */
export const UserPages: Record<string, PageItem> = {
  comments: {
    // Only show comments page for user
    ...AdminPages.comments,
  },
}

export interface NavigationStoreProps {
  onGoPage?: (pageName: string) => void
  onGoTab?: (tabName: string) => void
}

export interface SearchStateApi {
  show: boolean
  value: string
  updateValue: (val: string) => void
  showSearch: () => void
  hideSearch: () => void
}

/**
 * Navigation menu hook only use for AppNavigation component
 *
 * (do not expose to other components, others should call `useNavStore`)
 */
export const useNavigationMenu = (props: NavigationStoreProps = {}) => {
  const router = useRouter()
  const { is_admin: isAdmin } = storeToRefs(useUserStore())
  const { tabs, curtPage, curtTab, isMobile, isSearchEnabled } = storeToRefs(useNavStore())

  /**
   * Pages for current user
   */
  const pages = computed((): Record<string, PageItem> => (isAdmin.value ? AdminPages : UserPages))

  /**
   * State for search input
   */
  const searchState = reactive<SearchStateApi>({
    show: false,
    value: '',
    updateValue(val: string) {
      searchState.value = val
    },
    showSearch() {
      searchState.show = true
    },
    hideSearch() {
      searchState.show = false
    },
  })

  /**
   * Navigate to page
   */
  const goPage = (pageName: string) => {
    props.onGoPage?.(pageName)
    router.replace(pages.value[pageName].link)
  }

  /**
   * Navigate to tab
   */
  const goTab = (tabName: string) => {
    props.onGoTab?.(tabName)
    curtTab.value = tabName
  }

  /**
   * Sync current page name from route
   */
  const syncCurtPage = (to?: RouteLocation) =>
    (curtPage.value = String((to || useRoute()).name).replace(/^\//, ''))

  onMounted(() => {
    syncCurtPage()
  })

  router.afterEach((to, from, failure) => {
    syncCurtPage(to)
    searchState.value = ''
    searchState.hideSearch()
  })

  return reactive({
    pages,
    tabs,
    curtPage,
    curtTab,
    goPage,
    goTab,
    isMobile,
    isSearchEnabled,
    searchState,
    showSearch: searchState.showSearch,
    hideSearch: searchState.hideSearch,
  })
}

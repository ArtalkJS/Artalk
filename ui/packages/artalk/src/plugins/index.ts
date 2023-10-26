import ArtalkPlugin from '~/types/plugin'
import { EditorKit } from './editor-kit'
import * as Stat from './stat'
import { ListWithEditor } from './list-with-editor'
import { VersionCheck } from './version-check'
import { Unread } from './unread'
import { ListCount } from './list-count'
import { ListSidebarBtn } from './list-sidebar-btn'
import { ListUnreadBadge } from './list-unread-badge'
import { ListGoto } from './list-goto'
import { ListCopyright } from './list-copyright'
import { ListNoComment } from './list-no-comment'
import { ListDropdown } from './list-dropdown'
import { ListTimeTicking } from './list-time-ticking'
import { ListErrorDialog } from './list-error-dialog'
import { ListLoading } from './list-loading'
import { ListFetch } from './list-fetch'
import { ListReachBottom } from './list-reach-bottom'
import { ListGotoFirst } from './list-goto-first'

const ListPlugins: ArtalkPlugin[] = [
  ListFetch, ListLoading,
  ListWithEditor, ListCount, ListSidebarBtn,
  ListUnreadBadge, ListDropdown, ListGoto, ListNoComment, ListCopyright,
  ListTimeTicking, ListErrorDialog, ListReachBottom, ListGotoFirst,
]

export const DefaultPlugins: ArtalkPlugin[] = [
  EditorKit, Stat.PvCountWidget, VersionCheck, Unread,

  ...ListPlugins
]

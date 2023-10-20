import ArtalkPlugin from '~/types/plugin'
import { EditorKit } from './editor-kit'
import * as Stat from './stat'
import { ListCloseEditor } from './list-close-editor'
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

const ListPlugins: ArtalkPlugin[] = [
  ListCloseEditor, ListCount, ListSidebarBtn,
  ListUnreadBadge, ListDropdown, ListGoto, ListNoComment, ListCopyright,
  ListTimeTicking, ListErrorDialog
]

export const DefaultPlugins: ArtalkPlugin[] = [
  EditorKit, Stat.PvCountWidget, VersionCheck, Unread,

  ...ListPlugins
]

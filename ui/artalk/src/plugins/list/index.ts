import type { ArtalkPlugin } from '@/types'
import { WithEditor } from './with-editor'
import { Unread } from './unread'
import { Count } from './count'
import { SidebarBtn } from './sidebar-btn'
import { UnreadBadge } from './unread-badge'
import { GotoDispatcher } from './goto-dispatcher'
import { GotoFocus } from './goto-focus'
import { Copyright } from './copyright'
import { NoComment } from './no-comment'
import { Dropdown } from './dropdown'
import { TimeTicking } from './time-ticking'
import { ErrorDialog } from './error-dialog'
import { Loading } from './loading'
import { Fetch } from './fetch'
import { ReachBottom } from './reach-bottom'
import { GotoFirst } from './goto-first'

const ListPlugins: ArtalkPlugin[] = [
  Fetch,
  Loading,
  Unread,
  WithEditor,
  Count,
  SidebarBtn,
  UnreadBadge,
  Dropdown,
  GotoDispatcher,
  GotoFocus,
  NoComment,
  Copyright,
  TimeTicking,
  ErrorDialog,
  ReachBottom,
  GotoFirst,
]

export { ListPlugins }

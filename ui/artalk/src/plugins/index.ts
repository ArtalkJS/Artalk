import { Markdown } from './markdown'
import { EditorKit } from './editor-kit'
import { ListPlugins } from './list'
import { Notifies } from './notifies'
import { PvCountWidget } from './stat'
import { VersionCheck } from './version-check'
import { AdminOnlyElem } from './admin-only-elem'
import { DarkMode } from './dark-mode'
import { PageVoteWidget } from './page-vote'
import { Services } from '@/services'
import type { ArtalkPlugin } from '@/types'

export const DefaultPlugins: ArtalkPlugin[] = [
  ...Services,
  Markdown,
  EditorKit,
  AdminOnlyElem,
  ...ListPlugins,
  Notifies,
  PvCountWidget,
  VersionCheck,
  DarkMode,
  PageVoteWidget,
]

import ArtalkPlug from '~/types/plug'
import { EditorKit } from './editor-kit'
import * as Stat from './stat'
import { ListCloseEditor } from './list-close-editor'
import { VersionCheck } from './version-check'
import { Unread } from './unread'

export const DefaultPlugins: ArtalkPlug[] = [
  EditorKit, Stat.PvCountWidget, ListCloseEditor, VersionCheck, Unread
]

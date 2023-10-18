import ArtalkPlug from '~/types/plug'
import { EditorKit } from './editor-kit'
import * as Stat from './stat'
import { ListCloseEditor } from './list-close-editor'

export const DefaultPlugins: ArtalkPlug[] = [
  EditorKit, Stat.PvCountWidget, ListCloseEditor
]

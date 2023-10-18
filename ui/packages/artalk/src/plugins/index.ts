import ArtalkPlug from '~/types/plug'
import { EditorKit } from './editor-kit'
import * as Stat from './stat'

export const DefaultPlugins: ArtalkPlug[] = [
  EditorKit, Stat.PvCountWidget
]

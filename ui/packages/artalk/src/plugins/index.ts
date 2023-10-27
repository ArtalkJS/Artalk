import type ArtalkPlugin from '~/types/plugin'
import { EditorKit } from './editor-kit'
import { ListPlugins } from './list'
import { PvCountWidget } from './stat'
import { VersionCheck } from './version-check'

export const DefaultPlugins: ArtalkPlugin[] = [
  EditorKit, ...ListPlugins,
  PvCountWidget, VersionCheck
]

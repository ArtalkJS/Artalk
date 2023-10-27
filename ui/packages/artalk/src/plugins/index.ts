import type ArtalkPlugin from '~/types/plugin'
import { ConfRemoter } from './conf-remoter'
import { Markdown } from './markdown'
import { EditorKit } from './editor-kit'
import { ListPlugins } from './list'
import { PvCountWidget } from './stat'
import { VersionCheck } from './version-check'

export const DefaultPlugins: ArtalkPlugin[] = [
  ConfRemoter,
  Markdown, EditorKit, ...ListPlugins,
  PvCountWidget, VersionCheck
]

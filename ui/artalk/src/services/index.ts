import { ApiService } from './api'
import { UserService } from './user'
import { CheckersService } from './checkers'
import { EditorService } from './editor'
import { I18nService } from './i18n'
import { DataService } from './data'
import { LayersService } from './layer'
import { ListService } from './list'
import { SidebarService } from './sidebar'
import { ArtalkPlugin } from '@/types'

export const Services: Set<ArtalkPlugin> = new Set([
  DataService,
  ApiService,
  I18nService,
  UserService,
  EditorService,
  ListService,
  CheckersService,
  LayersService,
  SidebarService,
])

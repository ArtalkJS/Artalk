import type {
  DataManager,
  Editor,
  EventManager,
  Api,
  ApiHandlers,
  ConfigManager,
  UserManager,
  LayerManager,
  List,
  SidebarLayer,
  EditorPluginManager,
  CheckerManager,
} from '@/types'

import type { CommentService } from '@/services/comments'

export interface Services {
  config: ConfigManager
  events: EventManager
  data: DataManager
  api: Api
  apiHandlers: ApiHandlers
  /** @internal Comment request lifecycle. */
  comments: CommentService
  editor: Editor
  editorPlugs: EditorPluginManager | undefined
  list: List
  sidebar: SidebarLayer
  checkers: CheckerManager
  layers: LayerManager
  user: UserManager
}

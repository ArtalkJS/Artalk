import type {
  Api,
  ApiHandlers,
  ConfigManager,
  Editor,
  EventManager,
  I18n,
  LayerManager,
  UserManager,
} from 'artalk'

export interface AuthContext {
  getEditor: () => Editor
  getApi: () => Api
  getApiHandlers: () => ApiHandlers
  getUser: () => UserManager
  getEvents: () => EventManager
  getConf: () => ConfigManager
  getLayers: () => LayerManager
  $t(key: keyof I18n, args?: { [key: string]: string }): string
}

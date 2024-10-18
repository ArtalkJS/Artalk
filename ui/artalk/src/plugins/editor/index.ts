import type EditorPlugin from './_plug'
import LocalStorage from './local-storage'
import Textarea from './textarea'
import SubmitBtn from './submit-btn'
import Submit from './submit'
import StateReply from './state-reply'
import StateEdit from './state-edit'
import Closable from './closable'
import HeaderEvent from './header-event'
import HeaderUser from './header-user'
import HeaderLink from './header-link'
import Mover from './mover'
import Emoticons from './emoticons'
import Upload from './upload'
import Preview from './preview'
import type { Config } from '@/types'

/** The default enabled plugs */
const EDITOR_PLUGS: (typeof EditorPlugin)[] = [
  // Core
  LocalStorage,
  HeaderEvent,
  HeaderUser,
  HeaderLink,
  Textarea,
  Submit,
  SubmitBtn,
  Mover,
  StateReply,
  StateEdit,
  Closable,

  // Extensions
  Emoticons,
  Upload,
  Preview,
]

/**
 * Get the enabled plugs by config
 */
export function getEnabledPlugs(
  conf: Pick<Config, 'imgUpload' | 'emoticons' | 'preview' | 'editorTravel'>,
): (typeof EditorPlugin)[] {
  // The reference map of config and plugs
  // (for check if the plug is enabled)
  const confRefs = new Map<typeof EditorPlugin, any>()
  confRefs.set(Upload, conf.imgUpload)
  confRefs.set(Emoticons, conf.emoticons)
  confRefs.set(Preview, conf.preview)
  confRefs.set(Mover, conf.editorTravel)

  return EDITOR_PLUGS.filter((p) => !confRefs.has(p) || !!confRefs.get(p))
}

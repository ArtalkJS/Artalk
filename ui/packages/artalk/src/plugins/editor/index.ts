import type { ArtalkConfig } from '~/types'
import EditorPlug from './_plug'
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

/** The default enabled plugs */
export const ENABLED_PLUGS: (typeof EditorPlug)[] = [
  // Core
  LocalStorage,
  HeaderEvent, HeaderUser, HeaderLink,
  Textarea,
  Submit, SubmitBtn,
  Mover, StateReply, StateEdit,
  Closable,

  // Extensions
  Emoticons, Upload, Preview
]

/** Get the name list of disabled plugs */
export function getDisabledPlugByConf(conf: ArtalkConfig): (typeof EditorPlug)[] {
  return [
    {k: Upload, v: conf.imgUpload},
    {k: Emoticons, v: conf.emoticons},
    {k: Preview, v: conf.preview},
    {k: Mover, v: conf.editorTravel},
  ].filter(n => !n.v).flatMap(n => n.k)
}

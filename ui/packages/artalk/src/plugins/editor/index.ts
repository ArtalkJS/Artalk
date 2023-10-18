import type ArtalkConfig from '~/types/artalk-config'
import EditorPlug from './_plug'
import LocalStoragePlug from './local-storage'
import TextareaPlug from './textarea'
import SubmitBtnPlug from './submit-btn'
import SubmitPlug from './submit'
import ReplyPlug from './state-reply'
import EditPlug from './state-edit'
import ClosablePlug from './closable'
import HeaderEvent from './header-event'
import HeaderUser from './header-user'
import HeaderLink from './header-link'
import MoverPlug from './mover'
import EmoticonsPlug from './emoticons'
import UploadPlug from './upload'
import PreviewPlug from './preview'

/** The default enabled plugs */
export const ENABLED_PLUGS: (typeof EditorPlug)[] = [
  // Core
  LocalStoragePlug,
  HeaderEvent, HeaderUser, HeaderLink,
  TextareaPlug,
  SubmitPlug, SubmitBtnPlug,
  MoverPlug, ReplyPlug, EditPlug,
  ClosablePlug,

  // Extensions
  EmoticonsPlug, UploadPlug, PreviewPlug
]

/** Get the name list of disabled plugs */
export function getDisabledPlugByConf(conf: ArtalkConfig): (typeof EditorPlug)[] {
  return [
    {k: UploadPlug, v: conf.imgUpload},
    {k: EmoticonsPlug, v: conf.emoticons},
    {k: PreviewPlug, v: conf.preview},
    {k: MoverPlug, v: conf.editorTravel},
  ].filter(n => !n.v).flatMap(n => n.k)
}

import type ArtalkConfig from '~/types/artalk-config'
import EditorPlug from './_plug'
import LocalStoragePlug from './local-storage-plug'
import HeaderPlug from './header-plug'
import TextareaPlug from './textarea-plug'
import SubmitBtnPlug from './submit-btn-plug'
import SubmitPlug from './submit-plug'
import ReplyPlug from './reply-plug'
import EditPlug from './edit-plug'
import ClosablePlug from './closable-plug'
import HeaderInputPlug from './header-input-plug'
import MoverPlug from './mover-plug'
import EmoticonsPlug from './emoticons-plug'
import UploadPlug from './upload-plug'
import PreviewPlug from './preview-plug'

/** The default enabled plugs */
export const ENABLED_PLUGS: (typeof EditorPlug)[] = [
  // Core
  LocalStoragePlug,
  HeaderPlug, HeaderInputPlug, TextareaPlug,
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

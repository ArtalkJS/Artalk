import type ArtalkConfig from '~/types/artalk-config'
import EditorPlug from './editor-plug'
import LocalStoragePlug from './core/local-storage-plug'
import HeaderPlug from './core/header-plug'
import TextareaPlug from './core/textarea-plug'
import SubmitBtnPlug from './core/submit-btn-plug'
import SubmitPlug from './core/submit-plug'
import ReplyPlug from './core/reply-plug'
import EditPlug from './core/edit-plug'
import ClosablePlug from './core/closable-plug'
import HeaderInputPlug from './core/header-input-plug'
import MoverPlug from './core/mover-plug'
import EmoticonsPlug from './plugs/emoticons-plug'
import UploadPlug from './plugs/upload-plug'
import PreviewPlug from './plugs/preview-plug'

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

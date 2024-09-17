import {
  TbLayoutSidebarRightExpandFilled,
  TbMailFilled,
  TbEyeFilled,
  TbTransformFilled,
  TbLocationFilled,
  TbCardsFilled,
  TbPhotoSearch,
  TbMath,
  TbPlug,
  TbLanguage,
  TbTerminal,
  TbApi,
  TbSocial,
} from 'react-icons/tb'
import { BiSolidNotification, BiSolidBadgeCheck } from 'react-icons/bi'
import { RiLoader4Fill, RiPushpinLine, RiRobot2Fill, RiUpload2Fill } from 'react-icons/ri'
import { BsFillShieldLockFill } from 'react-icons/bs'
import { PiMoonFill, PiSmileyWinkBold } from 'react-icons/pi'
import { GrUpgrade } from 'react-icons/gr'
import { LuListTree, LuNewspaper } from 'react-icons/lu'
import { FaArrowTrendUp } from 'react-icons/fa6'
import { FaMarkdown, FaRegSave, FaSortAmountUpAlt } from 'react-icons/fa'
import { HiOutlineDocumentSearch } from 'react-icons/hi'
import { IoSearch, IoSend } from 'react-icons/io5'
import { IoMdLocate } from 'react-icons/io'
import { TFunction } from 'i18next'

export interface FeatureItem {
  icon: React.FC
  name: string
  desc: string
  link: string
}

export const getFeatures = (t: TFunction): FeatureItem[] => [
  {
    icon: TbLayoutSidebarRightExpandFilled,
    name: t('feature_sidebar_name'),
    desc: t('feature_sidebar_desc'),
    link: 'https://artalk.js.org/guide/frontend/sidebar.html',
  },
  {
    icon: TbSocial,
    name: t('feature_social_login_name'),
    desc: t('feature_social_login_desc'),
    link: 'https://artalk.js.org/guide/frontend/auth.html',
  },
  {
    icon: TbMailFilled,
    name: t('feature_email_notification_name'),
    desc: t('feature_email_notification_desc'),
    link: 'https://artalk.js.org/guide/backend/email.html',
  },
  {
    icon: IoSend,
    name: t('feature_diverse_push_name'),
    desc: t('feature_diverse_push_desc'),
    link: 'https://artalk.js.org/guide/backend/admin_notify.html',
  },
  {
    icon: BiSolidNotification,
    name: t('feature_site_notification_name'),
    desc: t('feature_site_notification_desc'),
    link: 'https://artalk.js.org/guide/frontend/sidebar.html',
  },
  {
    icon: RiRobot2Fill,
    name: t('feature_captcha_name'),
    desc: t('feature_captcha_desc'),
    link: 'https://artalk.js.org/guide/backend/captcha.html',
  },
  {
    icon: BsFillShieldLockFill,
    name: t('feature_comment_moderation_name'),
    desc: t('feature_comment_moderation_desc'),
    link: 'https://artalk.js.org/guide/backend/moderator.html',
  },
  {
    icon: RiUpload2Fill,
    name: t('feature_image_upload_name'),
    desc: t('feature_image_upload_desc'),
    link: 'https://artalk.js.org/guide/backend/img-upload.html',
  },
  {
    icon: FaMarkdown,
    name: t('feature_markdown_name'),
    desc: t('feature_markdown_desc'),
    link: 'https://artalk.js.org/guide/intro.html',
  },
  {
    icon: PiSmileyWinkBold,
    name: t('feature_emoji_pack_name'),
    desc: t('feature_emoji_pack_desc'),
    link: 'https://artalk.js.org/guide/frontend/emoticons.html',
  },
  {
    icon: TbCardsFilled,
    name: t('feature_multi_site_name'),
    desc: t('feature_multi_site_desc'),
    link: 'https://artalk.js.org/guide/backend/multi-site.html',
  },
  {
    icon: BiSolidBadgeCheck,
    name: t('feature_admin_name'),
    desc: t('feature_admin_desc'),
    link: 'https://artalk.js.org/guide/backend/multi-site.html',
  },
  {
    icon: LuNewspaper,
    name: t('feature_page_management_name'),
    desc: t('feature_page_management_desc'),
    link: 'https://artalk.js.org/guide/frontend/sidebar.html',
  },
  {
    icon: TbEyeFilled,
    name: t('feature_page_view_statistics_name'),
    desc: t('feature_page_view_statistics_desc'),
    link: 'https://artalk.js.org/guide/frontend/pv.html',
  },
  {
    icon: LuListTree,
    name: t('feature_hierarchical_structure_name'),
    desc: t('feature_hierarchical_structure_desc'),
    link: 'https://artalk.js.org/guide/frontend/config.html#nestmax',
  },
  {
    icon: FaArrowTrendUp,
    name: t('feature_comment_voting_name'),
    desc: t('feature_comment_voting_desc'),
    link: 'https://artalk.js.org/guide/frontend/config.html#vote',
  },
  {
    icon: FaSortAmountUpAlt,
    name: t('feature_comment_sorting_name'),
    desc: t('feature_comment_sorting_desc'),
    link: 'https://artalk.js.org/guide/frontend/config.html#listsort',
  },
  {
    icon: IoSearch,
    name: t('feature_comment_search_name'),
    desc: t('feature_comment_search_desc'),
    link: 'https://artalk.js.org/guide/frontend/sidebar.html',
  },
  {
    icon: RiPushpinLine,
    name: t('feature_comment_pinning_name'),
    desc: t('feature_comment_pinning_desc'),
    link: 'https://artalk.js.org/guide/frontend/sidebar.html',
  },
  {
    icon: HiOutlineDocumentSearch,
    name: t('feature_view_author_only_name'),
    desc: t('feature_view_author_only_desc'),
    link: 'https://artalk.js.org/guide/frontend/config.html',
  },
  {
    icon: IoMdLocate,
    name: t('feature_comment_jump_name'),
    desc: t('feature_comment_jump_desc'),
    link: 'https://artalk.js.org/guide/intro.html',
  },
  {
    icon: FaRegSave,
    name: t('feature_auto_save_name'),
    desc: t('feature_auto_save_desc'),
    link: 'https://artalk.js.org/guide/frontend/config.html',
  },
  {
    icon: TbLocationFilled,
    name: t('feature_ip_region_name'),
    desc: t('feature_ip_region_desc'),
    link: 'https://artalk.js.org/guide/frontend/ip-region.html',
  },
  {
    icon: TbTransformFilled,
    name: t('feature_data_migration_name'),
    desc: t('feature_data_migration_desc'),
    link: 'https://artalk.js.org/guide/transfer.html',
  },
  {
    icon: TbPhotoSearch,
    name: t('feature_image_lightbox_name'),
    desc: t('feature_image_lightbox_desc'),
    link: 'https://artalk.js.org/guide/frontend/lightbox.html',
  },
  {
    icon: RiLoader4Fill,
    name: t('feature_image_lazy_load_name'),
    desc: t('feature_image_lazy_load_desc'),
    link: 'https://artalk.js.org/guide/frontend/img-lazy-load.html',
  },
  {
    icon: TbMath,
    name: t('feature_latex_name'),
    desc: t('feature_latex_desc'),
    link: 'https://artalk.js.org/guide/frontend/latex.html',
  },
  {
    icon: PiMoonFill,
    name: t('feature_night_mode_name'),
    desc: t('feature_night_mode_desc'),
    link: 'https://artalk.js.org/guide/frontend/config.html#darkmode',
  },
  {
    icon: TbPlug,
    name: t('feature_extension_plugin_name'),
    desc: t('feature_extension_plugin_desc'),
    link: 'https://artalk.js.org/develop/plugin.html',
  },
  {
    icon: TbLanguage,
    name: t('feature_multi_language_name'),
    desc: t('feature_multi_language_desc'),
    link: 'https://artalk.js.org/guide/frontend/i18n.html',
  },
  {
    icon: TbTerminal,
    name: t('feature_command_line_name'),
    desc: t('feature_command_line_desc'),
    link: 'https://artalk.js.org/guide/backend/config.html',
  },
  {
    icon: TbApi,
    name: t('feature_api_documentation_name'),
    desc: t('feature_api_documentation_desc'),
    link: 'https://artalk.js.org/http-api.html',
  },
  {
    icon: GrUpgrade,
    name: t('feature_program_upgrade_name'),
    desc: t('feature_program_upgrade_desc'),
    link: 'https://artalk.js.org/guide/backend/update.html',
  },
]

import React from 'react'
import { TbLayoutSidebarRightExpandFilled, TbMailFilled, TbEyeFilled, TbTransformFilled, TbLocationFilled, TbCardsFilled, TbPhotoSearch, TbMath, TbPlug, TbLanguage, TbTerminal, TbApi, TbSocial } from 'react-icons/tb'
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

interface FuncItem {
  icon: React.ReactNode
  name: string
  desc: string
  link: string
}

export const FuncList: FuncItem[] = [
  {
    icon: <TbLayoutSidebarRightExpandFilled />,
    name: '侧边栏',
    desc: '快速管理、直观浏览',
    link: 'https://artalk.js.org/guide/frontend/sidebar.html'
  },
  {
    icon: <TbSocial />,
    name: '社交登录',
    desc: '通过社交账号快速登录',
    link: 'https://artalk.js.org/guide/frontend/auth.html'
  },
  {
    icon: <TbMailFilled />,
    name: '邮件通知',
    desc: '多种发送方式、邮件模板',
    link: 'https://artalk.js.org/guide/backend/email.html'
  },
  {
    icon: <IoSend />,
    name: '多元推送',
    desc: '多种推送方式、通知模版',
    link: 'https://artalk.js.org/guide/backend/admin_notify.html'
  },
  {
    icon: <BiSolidNotification />,
    name: '站内通知',
    desc: '红点标记、提及列表',
    link: 'https://artalk.js.org/guide/frontend/sidebar.html'
  },
  {
    icon: <RiRobot2Fill />,
    name: '验证码',
    desc: '多种验证类型、频率限制',
    link: 'https://artalk.js.org/guide/backend/captcha.html'
  },
  {
    icon: <BsFillShieldLockFill />,
    name: '评论审核',
    desc: '内容检测、垃圾拦截',
    link: 'https://artalk.js.org/guide/backend/moderator.html'
  },
  {
    icon: <RiUpload2Fill />,
    name: '图片上传',
    desc: '自定义上传、支持图床',
    link: 'https://artalk.js.org/guide/backend/img-upload.html'
  },
  {
    icon: <FaMarkdown />,
    name: 'Markdown',
    desc: '支持 Markdown 语法',
    link: 'https://artalk.js.org/guide/intro.html'
  },
  {
    icon: <PiSmileyWinkBold />,
    name: '表情包',
    desc: '兼容 OwO，快速集成',
    link: 'https://artalk.js.org/guide/frontend/emoticons.html'
  },
  {
    icon: <TbCardsFilled />,
    name: '多站点',
    desc: '站点隔离、集中管理',
    link: 'https://artalk.js.org/guide/backend/multi-site.html'
  },
  {
    icon: <BiSolidBadgeCheck />,
    name: '管理员',
    desc: '密码验证、徽章标识',
    link: 'https://artalk.js.org/guide/backend/multi-site.html'
  },
  {
    icon: <LuNewspaper />,
    name: '页面管理',
    desc: '快速查看、标题一键跳转',
    link: 'https://artalk.js.org/guide/frontend/sidebar.html',
  },
  {
    icon: <TbEyeFilled />,
    name: '浏览量统计',
    desc: '轻松统计网页浏览量',
    link: 'https://artalk.js.org/guide/frontend/pv.html'
  },
  {
    icon: <LuListTree />,
    name: '层级结构',
    desc: '嵌套分页列表、滚动加载',
    link: 'https://artalk.js.org/guide/frontend/config.html#nestmax',
  },
  {
    icon: <FaArrowTrendUp />,
    name: '评论投票',
    desc: '赞同或反对评论',
    link: 'https://artalk.js.org/guide/frontend/config.html#vote',
  },
  {
    icon: <FaSortAmountUpAlt />,
    name: '评论排序',
    desc: '多种排序方式，自由选择',
    link: 'https://artalk.js.org/guide/frontend/config.html#listsort',
  },
  {
    icon: <IoSearch />,
    name: '评论搜索',
    desc: '快速搜索评论内容',
    link: 'https://artalk.js.org/guide/frontend/sidebar.html',
  },
  {
    icon: <RiPushpinLine />,
    name: '评论置顶',
    desc: '重要消息置顶显示',
    link: 'https://artalk.js.org/guide/frontend/sidebar.html',
  },
  {
    icon: <HiOutlineDocumentSearch />,
    name: '仅看作者',
    desc: '仅显示作者的评论',
    link: 'https://artalk.js.org/guide/frontend/config.html',
  },
  {
    icon: <IoMdLocate />,
    name: '评论跳转',
    desc: '快速跳转到引用的评论',
    link: 'https://artalk.js.org/guide/intro.html',
  },
  {
    icon: <FaRegSave />,
    name: '自动保存',
    desc: '输入内容防丢功能',
    link: 'https://artalk.js.org/guide/frontend/config.html',
  },
  {
    icon: <TbLocationFilled />,
    name: 'IP 属地',
    desc: '用户 IP 属地展示',
    link: 'https://artalk.js.org/guide/frontend/ip-region.html'
  },
  {
    icon: <TbTransformFilled />,
    name: '数据迁移',
    desc: '自由迁移、快速备份',
    link: 'https://artalk.js.org/guide/transfer.html'
  },
  {
    icon: <TbPhotoSearch />,
    name: '图片灯箱',
    desc: '图片灯箱快速集成',
    link: 'https://artalk.js.org/guide/frontend/lightbox.html'
  },
  {
    icon: <RiLoader4Fill />,
    name: '图片懒加载',
    desc: '延迟加载图片，优化体验',
    link: 'https://artalk.js.org/guide/frontend/img-lazy-load.html'
  },
  {
    icon: <TbMath />,
    name: 'Latex',
    desc: 'Latex 公式解析集成',
    link: 'https://artalk.js.org/guide/frontend/latex.html'
  },
  {
    icon: <PiMoonFill />,
    name: '夜间模式',
    desc: '夜间模式切换',
    link: 'https://artalk.js.org/guide/frontend/config.html#darkmode'
  },
  {
    icon: <TbPlug />,
    name: '扩展插件',
    desc: '创造更多可能性',
    link: 'https://artalk.js.org/develop/'
  },
  {
    icon: <TbLanguage />,
    name: '多语言',
    desc: '多国语言切换',
    link: 'https://artalk.js.org/guide/frontend/i18n.html'
  },
  {
    icon: <TbTerminal />,
    name: '命令行',
    desc: '命令行操作管理能力',
    link: 'https://artalk.js.org/guide/backend/config.html'
  },
  {
    icon: <TbApi />,
    name: 'API 文档',
    desc: '提供 OpenAPI 格式文档',
    link: 'https://artalk.js.org/develop/'
  },
  {
    icon: <GrUpgrade />,
    name: '程序升级',
    desc: '版本检测，一键升级',
    link: 'https://artalk.js.org/guide/backend/update.html'
  }
]

export const FullFeatureList: React.FC = () => {
  // let func list item group by every two items
  const FuncListGrouped = FuncList.reduce<FuncItem[][]>((result, current, index) => {
    if (index % 3 === 0) result.push([current])
    else result[result.length - 1].push(current)
    return result
  }, [])

  return (
    <div className='func-list'>
      {FuncListGrouped.map((row) => (
        <div className='row'>
          {row.map((item) => (
            <a className='item' href={item.link} target='_blank'>
              <div className="header">
                {item.icon}
                <span className="text">{item.name}</span>
              </div>
              <div className="body">
                <div className="desc">{item.desc}</div>
              </div>
            </a>
          ))}
        </div>
      ))}
    </div>
  )
}

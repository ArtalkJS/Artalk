import React from 'react'
import { TbLayoutSidebarRightExpandFilled, TbMailFilled, TbEyeFilled, TbTransformFilled, TbLocationFilled, TbCardsFilled, TbPhotoSearch, TbMath, TbPlug, TbLanguage, TbTerminal, TbApi } from 'react-icons/tb'
import { BiSolidNotification, BiSolidBadgeCheck } from 'react-icons/bi'
import { RiRobot2Fill, RiUpload2Fill } from 'react-icons/ri'
import { BsFillShieldLockFill } from 'react-icons/bs'
import { PiSmileyWinkBold } from 'react-icons/pi'

interface FuncItem {
  icon: React.ReactNode
  name: string
  desc: string
  link: string
}

const FuncList: FuncItem[] = [
  {
    icon: <TbLayoutSidebarRightExpandFilled />,
    name: '侧边栏',
    desc: '快速管理、直观浏览',
    link: 'https://artalk.js.org/guide/frontend/sidebar.html'
  },
  {
    icon: <TbMailFilled />,
    name: '邮件通知',
    desc: '多方式、邮件模板',
    link: 'https://artalk.js.org/guide/backend/email.html'
  },
  {
    icon: <BiSolidNotification />,
    name: '多元推送',
    desc: '支持 Telegram 等消息推送',
    link: 'https://artalk.js.org/guide/backend/admin_notify.html'
  },
  {
    icon: <RiRobot2Fill />,
    name: '验证码',
    desc: '多方式、频率限制',
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
    icon: <TbCardsFilled />,
    name: '多站点',
    desc: '站点隔离、集中管理',
    link: 'https://artalk.js.org/guide/backend/multi-site.html'
  },
  {
    icon: <PiSmileyWinkBold />,
    name: '表情包',
    desc: '动态加载、兼容 OwO',
    link: 'https://artalk.js.org/guide/frontend/emoticons.html'
  },
  {
    icon: <BiSolidBadgeCheck />,
    name: '管理员',
    desc: '密码验证、徽章标识',
    link: 'https://artalk.js.org/guide/backend/multi-site.html'
  },
  {
    icon: <TbEyeFilled />,
    name: '浏览量统计',
    desc: '轻松统计网页浏览量',
    link: 'https://artalk.js.org/guide/frontend/pv.html'
  },
  {
    icon: <TbTransformFilled />,
    name: '数据迁移',
    desc: '自由迁移、快速备份',
    link: 'https://artalk.js.org/guide/transfer.html'
  },
  {
    icon: <TbLocationFilled />,
    name: 'IP 属地',
    desc: '用户 IP 属地展示',
    link: 'https://artalk.js.org/guide/frontend/ip-region.html'
  },
  {
    icon: <TbPhotoSearch />,
    name: '图片灯箱',
    desc: '快速集成图片灯箱',
    link: 'https://artalk.js.org/guide/frontend/lightbox.html'
  },
  {
    icon: <TbMath />,
    name: 'Latex',
    desc: '一键集成 Latex 公式解析',
    link: 'https://artalk.js.org/guide/frontend/latex.html'
  },
  {
    icon: <TbPlug />,
    name: '自定义插件',
    desc: '创造更多可能性',
    link: 'https://artalk.js.org/develop/'
  },
  {
    icon: <TbLanguage />,
    name: '多语言',
    desc: '支持多国语言切换',
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

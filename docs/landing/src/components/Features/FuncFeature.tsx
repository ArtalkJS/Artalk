import React from 'react'
import './FuncFeature.scss'
import { FeatureBase } from './FeatureBase'
import { FeatureTitle } from '../FeatureTitle'
import { FeatureDesc } from '../FeatureDesc'
import { Reveal } from '../Reveal'
import { FaArrowRight } from 'react-icons/fa'

const FuncGrps: {name: string, items: string[], link?: string}[] = [
  { name: '社交登录', items: ['Github', 'GitLab', 'Twitter', 'Facebook', 'Mastodon', 'Google', 'Microsoft', 'Apple', 'Discord', 'Slack', 'Tiktok', 'Steam'], link: 'https://artalk.js.org/guide/frontend/auth.html' },
  { name: '邮箱发送', items: ['SMTP', '阿里云邮件', 'sendmail'], link: 'https://artalk.js.org/guide/backend/email.html' },
  { name: '验证码', items: ['Turnstile', 'reCAPTCHA', 'hCaptcha', '极验'], link: 'https://artalk.js.org/guide/backend/captcha.html' },
  { name: '消息推送', items: ['Telegram', '飞书', '钉钉', 'Bark', 'WebHook', 'Slack', 'LINE'], link: 'https://artalk.js.org/guide/backend/admin_notify.html' },
  { name: '评论审核', items: ['Akismet', '腾讯云', '阿里云', '离线词库'], link: 'https://artalk.js.org/guide/backend/moderator.html' },
  { name: '表情包', items: ['标准格式', 'OwO 格式'], link: 'https://artalk.js.org/guide/frontend/emoticons.html' },
  { name: '图片上传', items: ['本地保存', 'UPGIT', '自定义函数'], link: 'https://artalk.js.org/guide/backend/img-upload.html' },
  { name: '富文本', items: ['Markdown', 'Latex'], link: 'https://artalk.js.org/guide/frontend/latex.html' },
  { name: '代码高亮', items: ['hanabi', 'highlight.js'] },
  { name: '用户头像', items: ['Gravatar', '自定义函数'], link: 'https://artalk.js.org/guide/frontend/config.html' },
  { name: '图片灯箱', items: ['LightGallery', 'FancyBox', 'lightbox2'], link: 'https://artalk.js.org/guide/frontend/lightbox.html' },
  { name: '数据库', items: ['SQLite', 'MySQL', 'PostgreSQL', 'SQL Server'], link: 'https://artalk.js.org/guide/backend/config.html' },
  { name: '高速缓存', items: ['内建缓存', 'Redis', 'Memcache'], link: 'https://artalk.js.org/guide/backend/config.html' },
  { name: '程序部署', items: ['二进制文件', 'Docker 镜像'], link: 'https://artalk.js.org/guide/deploy.html' },
  { name: '操作系统', items: ['Linux', 'Windows', 'macOS', 'FreeBSD'] },
  { name: '平台架构', items: ['x86', 'ARM'] },
  { name: '评论导入', items: ['Typecho', 'WordPress', 'Valine', 'Waline', 'Disqus', 'Commento', 'Twikoo', 'Artrans'], link: 'https://artalk.js.org/guide/transfer.html' },
]

export const FuncFeature: React.FC = () => {
  return (
    <FeatureBase className='func'>
      <FeatureTitle text='可选集成' color='#292F33' />

      <div className="content">
        <main>
          <Reveal delay={250}>
            <FeatureDesc>
            Artalk 提供丰富第三方接入能力。
            </FeatureDesc>
          </Reveal>
        </main>
        <aside>
            <svg className="ship-icon" width="177" height="177" viewBox="0 0 177 177" fill="none" xmlns="http://www.w3.org/2000/svg">
              <g clip-path="url(#clip0_174_446)"><path d="M78.6667 0H98.3333V24.5833H78.6667V0Z" fill="#292F33"/><path d="M78.6667 9.8335H98.3333V29.5002H78.6667V9.8335Z" fill="#D1D3D4"/><path d="M88.5 9.8335H98.3333V29.5002H88.5V9.8335Z" fill="#A7A9AC"/><path d="M88.5 93.4165H24.5833C24.5833 93.4165 31.2454 132.75 48.5521 162.25C65.8538 191.75 88.5 162.25 88.5 162.25C88.5 162.25 111.146 191.75 128.448 162.25C145.755 132.75 152.417 93.4165 152.417 93.4165H88.5Z" fill="#66757F"/><path d="M88.5 93.4165H24.5833C24.5833 93.4165 31.2454 132.75 48.5521 162.25C65.8538 191.75 88.5 162.25 88.5 162.25V93.4165Z" fill="#99AAB5"/>
              <g fill="#00A6ED" className="wave">
                <path d="m10.7 134.3c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14z"/>
              </g>
              <path d="M0 147.5H177V177H0V147.5Z" fill="#00A6ED"/>
              <path d="M147.5 44.2502H135.454L134.264 24.5835H42.7357L41.5508 44.2502H29.5V73.7502H39.766L38.5762 93.4168H138.424L137.234 73.7502H147.5V44.2502Z" fill="#E6E7E8"/><path d="M49.1667 73.75H127.833V93.4167H49.1667V73.75Z" fill="#D1D3D4"/><path d="M39.3333 54.0835H137.667V63.9168H39.3333V54.0835Z" fill="#6D6E71"/><path d="M49.1667 34.4165H127.833V44.2498H49.1667V34.4165Z" fill="#BCBEC0"/><path d="M29.0723 113.083H147.928C148.916 109.411 149.737 106.067 150.386 103.25H26.6139C27.2629 106.067 28.084 109.411 29.0723 113.083Z" fill="#BE1931"/><path d="M88.5 113.083H147.928C148.916 109.411 149.737 106.067 150.386 103.25H88.5V113.083Z" fill="#A0041E"/><path d="M59 83.5835H118V93.4168H59V83.5835Z" fill="#BCBEC0"/><path d="M78.6667 0H88.5V9.83333H78.6667V0Z" fill="#58595B"/></g><defs><clipPath id="clip0_174_446"><rect width="177" height="177" fill="white"/></clipPath></defs>
            </svg>
        </aside>
      </div>

      <div className="func-grps">
        {FuncGrps.map((grp) => (
          <div className="grp">
            {grp.link ? (
              <a className="name" href={grp.link} target='_blank'>{grp.name}<FaArrowRight /></a>
            ) : (
              <div className='name'>{grp.name}</div>
            )}
            <div className="items">
              {grp.items.map((item) => (
                <div className='item'>{item}</div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </FeatureBase>
  )
}

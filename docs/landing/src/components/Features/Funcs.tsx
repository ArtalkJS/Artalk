import React from 'react'
import { useTranslation } from 'react-i18next'
import './Funcs.scss'
import { FaArrowRight } from 'react-icons/fa'
import { Title } from '../FeatureComps/Title'
import { Desc } from '../FeatureComps/Desc'
import { Reveal } from '../Reveal'
import { Base } from '../FeatureComps/Base'

export const FuncsFeature: React.FC = () => {
  const { t } = useTranslation()

  const FuncGrps: { name: string; items: string[]; link?: string }[] = [
    {
      name: t('func_social_login'),
      items: [
        'Github',
        'GitLab',
        'Twitter',
        'Facebook',
        'Mastodon',
        'Google',
        'Microsoft',
        'Apple',
        'Discord',
        'Slack',
        'Tiktok',
        'Steam',
      ],
      link: 'https://artalk.js.org/guide/frontend/auth.html',
    },
    {
      name: t('func_email'),
      items: ['SMTP', t('func_email_aliyun'), 'sendmail'],
      link: 'https://artalk.js.org/guide/backend/email.html',
    },
    {
      name: t('func_captcha'),
      items: ['Turnstile', 'reCAPTCHA', 'hCaptcha', t('func_captcha_geetest')],
      link: 'https://artalk.js.org/guide/backend/captcha.html',
    },
    {
      name: t('func_message_pusher'),
      items: [
        'Telegram',
        t('func_message_pusher_lark'),
        t('func_message_pusher_dingtalk'),
        'Bark',
        'WebHook',
        'Slack',
        'LINE',
      ],
      link: 'https://artalk.js.org/guide/backend/admin_notify.html',
    },
    {
      name: t('func_moderator'),
      items: [
        'Akismet',
        t('func_moderator_tencent'),
        t('func_moderator_aliyun'),
        t('func_moderator_offline'),
      ],
      link: 'https://artalk.js.org/guide/backend/moderator.html',
    },
    {
      name: t('func_emoji'),
      items: [t('func_emoji_standard_format'), t('func_emoji_owo_format')],
      link: 'https://artalk.js.org/guide/frontend/emoticons.html',
    },
    {
      name: t('func_img_upload'),
      items: [t('func_img_upload_local'), 'UPGIT', t('func_img_upload_function')],
      link: 'https://artalk.js.org/guide/backend/img-upload.html',
    },
    {
      name: t('func_rich_text'),
      items: ['Markdown', 'Latex'],
      link: 'https://artalk.js.org/guide/frontend/latex.html',
    },
    { name: t('func_code_highlight'), items: ['hanabi', 'highlight.js'] },
    {
      name: t('func_avatar'),
      items: ['Gravatar', t('func_avatar_function')],
      link: 'https://artalk.js.org/guide/frontend/config.html',
    },
    {
      name: t('func_img_lightbox'),
      items: ['LightGallery', 'FancyBox', 'lightbox2'],
      link: 'https://artalk.js.org/guide/frontend/lightbox.html',
    },
    {
      name: t('func_database'),
      items: ['SQLite', 'MySQL', 'PostgreSQL', 'SQL Server'],
      link: 'https://artalk.js.org/guide/backend/config.html',
    },
    {
      name: t('func_cache'),
      items: [t('func_cache_internal'), 'Redis', 'Memcache'],
      link: 'https://artalk.js.org/guide/backend/config.html',
    },
    {
      name: t('func_deploy'),
      items: [t('func_deploy_bin'), t('func_deploy_docker')],
      link: 'https://artalk.js.org/guide/deploy.html',
    },
    { name: t('func_os'), items: ['Linux', 'Windows', 'macOS', 'FreeBSD'] },
    { name: t('func_platform'), items: ['x86', 'ARM'] },
    {
      name: t('func_transfer_import'),
      items: [
        'Typecho',
        'WordPress',
        'Valine',
        'Waline',
        'Disqus',
        'Commento',
        'Twikoo',
        'Artrans',
      ],
      link: 'https://artalk.js.org/guide/transfer.html',
    },
  ]

  return (
    <Base className="funcs-feature">
      <Title text={t('feature_func_title')} color="var(--color-font)" />

      <div className="content">
        <main>
          <Reveal delay={250}>
            <Desc>{t('feature_func_desc_line_1')}</Desc>
          </Reveal>
        </main>
        <aside>
          <svg
            className="ship-icon"
            width="177"
            height="177"
            viewBox="0 0 177 177"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <g clipPath="url(#clip0_174_446)">
              <path d="M78.6667 0H98.3333V24.5833H78.6667V0Z" fill="#292F33" />
              <path d="M78.6667 9.8335H98.3333V29.5002H78.6667V9.8335Z" fill="#D1D3D4" />
              <path d="M88.5 9.8335H98.3333V29.5002H88.5V9.8335Z" fill="#A7A9AC" />
              <path
                d="M88.5 93.4165H24.5833C24.5833 93.4165 31.2454 132.75 48.5521 162.25C65.8538 191.75 88.5 162.25 88.5 162.25C88.5 162.25 111.146 191.75 128.448 162.25C145.755 132.75 152.417 93.4165 152.417 93.4165H88.5Z"
                fill="#66757F"
              />
              <path
                d="M88.5 93.4165H24.5833C24.5833 93.4165 31.2454 132.75 48.5521 162.25C65.8538 191.75 88.5 162.25 88.5 162.25V93.4165Z"
                fill="#99AAB5"
              />
              <g fill="#00A6ED" className="wave">
                <path d="m10.7 134.3c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14zm24.4 0c.2-9.5 12.7 14 12.7 14l-25.4 0c0 0 12.5-4.5 12.7-14z" />
              </g>
              <path d="M0 147.5H177V177H0V147.5Z" fill="#00A6ED" />
              <path
                d="M147.5 44.2502H135.454L134.264 24.5835H42.7357L41.5508 44.2502H29.5V73.7502H39.766L38.5762 93.4168H138.424L137.234 73.7502H147.5V44.2502Z"
                fill="#E6E7E8"
              />
              <path d="M49.1667 73.75H127.833V93.4167H49.1667V73.75Z" fill="#D1D3D4" />
              <path d="M39.3333 54.0835H137.667V63.9168H39.3333V54.0835Z" fill="#6D6E71" />
              <path d="M49.1667 34.4165H127.833V44.2498H49.1667V34.4165Z" fill="#BCBEC0" />
              <path
                d="M29.0723 113.083H147.928C148.916 109.411 149.737 106.067 150.386 103.25H26.6139C27.2629 106.067 28.084 109.411 29.0723 113.083Z"
                fill="#BE1931"
              />
              <path
                d="M88.5 113.083H147.928C148.916 109.411 149.737 106.067 150.386 103.25H88.5V113.083Z"
                fill="#A0041E"
              />
              <path d="M59 83.5835H118V93.4168H59V83.5835Z" fill="#BCBEC0" />
              <path d="M78.6667 0H88.5V9.83333H78.6667V0Z" fill="#58595B" />
            </g>
            <defs>
              <clipPath id="clip0_174_446">
                <rect width="177" height="177" fill="white" />
              </clipPath>
            </defs>
          </svg>
        </aside>
      </div>

      <div className="func-grps">
        {FuncGrps.map((grp, i) => (
          <div key={i} className="grp">
            {grp.link ? (
              <a className="name" href={grp.link} target="_blank" rel="noreferrer">
                {grp.name}
                <FaArrowRight />
              </a>
            ) : (
              <div className="name">{grp.name}</div>
            )}
            <div className="items">
              {grp.items.map((item, j) => (
                <div key={j} className="item">
                  {item}
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </Base>
  )
}

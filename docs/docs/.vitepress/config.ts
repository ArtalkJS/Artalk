import { defineConfig } from 'vitepress'
import iterator from 'markdown-it-for-inline'
import * as Version from '../code/ArtalkVersion.json'

export default defineConfig({
  title: 'Artalk',
  description: '一款简洁的自托管评论系统',
  lang: 'zh-CN',

  head: [
    ['link', { rel: 'icon', href: '/favicon.png' }],
    [
      'meta',
      {
        name: 'viewport',
        content:
          'width=device-width, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0, user-scalable=no, target-densitydpi=device-dpi',
      },
    ],
    // light gallery
    [
      'link',
      {
        href: 'https://cdnjs.cloudflare.com/ajax/libs/lightgallery/2.3.0/css/lightgallery.css',
        rel: 'stylesheet',
      },
    ],
    [
      'script',
      {
        src: 'https://cdnjs.cloudflare.com/ajax/libs/lightgallery/2.3.0/lightgallery.min.js',
      },
    ],
    // katex
    [
      'link',
      {
        href: 'https://unpkg.com/katex@0.15.3/dist/katex.min.css',
        rel: 'stylesheet',
      },
    ],
    ['script', { src: 'https://unpkg.com/katex@0.15.3/dist/katex.min.js' }],
    // ['script', { src: 'https://unpkg.com/@artalk/plugin-katex/dist/artalk-plugin-katex.js' }],
  ],

  lastUpdated: true,

  markdown: {
    // @link https://github.com/shikijs/shiki
    theme: {
      light: 'github-light',
      dark: 'github-dark',
    },
    config: (md) => {
      md.use(iterator, 'artalk_version', 'text', function (tokens, idx) {
        tokens[idx].content = tokens[idx].content.replace(
          /:ArtalkVersion:/g,
          Version.latest,
        )
      })
      md.use(iterator, 'artalk_version_link', 'link_open', (tokens, idx) => {
        const href = tokens[idx].attrGet('href')
        tokens[idx].attrSet(
          'href',
          href.replace(/:ArtalkVersion:/g, Version.latest),
        )
      })
    },
  },

  locales: {
    root: {
      label: '简体中文',
      lang: 'zh',
      link: 'https://artalk.js.org/',
    },
    en: {
      label: 'English',
      lang: 'en',
      link: 'https://artalk-js-org.translate.goog/guide/intro?_x_tr_sl=auto&_x_tr_tl=en-US&_x_tr_hl=en-US&_x_tr_pto=wapp',
    },
  },

  themeConfig: {
    sidebar: {
      '/guide/': [
        {
          text: '快速开始',
          collapsed: false,
          items: [
            { text: '项目介绍', link: '/guide/intro.md' },
            { text: '程序部署', link: '/guide/deploy.md' },
            { text: '数据迁移', link: '/guide/transfer.md' },
          ],
        },
        {
          text: '核心指南',
          collapsed: false,
          items: [
            { text: '侧边栏', link: '/guide/frontend/sidebar.md' },
            { text: '邮件通知', link: '/guide/backend/email.md' },
            { text: '多元推送', link: '/guide/backend/admin_notify.md' },
            { text: '评论审核', link: '/guide/backend/moderator.md' },
            { text: '验证码', link: '/guide/backend/captcha.md' },
            { text: '图片上传', link: '/guide/backend/img-upload.md' },
            { text: '账户与多站点', link: '/guide/backend/multi-site.md' },
            { text: '解析相对路径', link: '/guide/backend/relative-path.md' },
            { text: '如何配置界面', link: '/guide/backend/fe-control.md' },
            { text: 'English (i18n)', link: '/guide/frontend/i18n.md' },
          ],
        },
        {
          text: '进阶指南',
          collapsed: false,
          items: [
            { text: '表情包', link: '/guide/frontend/emoticons.md' },
            { text: '浏览量统计', link: '/guide/frontend/pv.md' },
            { text: 'Latex', link: '/guide/frontend/latex.md' },
            { text: '图片灯箱', link: '/guide/frontend/lightbox.md' },
            { text: 'IP 属地', link: '/guide/frontend/ip-region.md' },
            { text: '精简版本', link: '/guide/frontend/artalk-lite.md' },
            { text: '置入博客', link: '/guide/frontend/import-blog.md' },
            { text: '置入框架', link: '/guide/frontend/import-framework.md' },
            { text: '前端 API', link: '/develop/fe-api.md' },
            { text: '开发文档', link: '/develop/index.md' },
          ],
        },
        {
          text: '配置文档',
          collapsed: true,
          items: [
            { text: '后端配置', link: '/guide/backend/config.md' },
            { text: '前端配置', link: '/guide/frontend/config.md' },
          ],
        },
        {
          text: '部署说明',
          collapsed: true,
          items: [
            { text: '守护进程', link: '/guide/backend/daemon.md' },
            { text: '反向代理', link: '/guide/backend/reverse-proxy.md' },
            { text: '编译构建', link: '/guide/backend/build.md' },
            { text: '程序升级', link: '/guide/backend/update.md' },
            { text: 'Docker', link: '/guide/backend/docker.md' },
          ],
        },
        {
          text: '更多内容',
          collapsed: true,
          items: [
            { text: '安全防范', link: '/guide/security.md' },
            { text: '扩展阅读', link: '/guide/extras.md' },
            { text: '案例展示', link: '/guide/cases.md' },
            { text: '关于我们', link: '/guide/about.md' },
          ],
        },
      ],
      '/develop/': [
        {
          text: '开发文档',
          items: [
            { text: '开发说明', link: '/develop/index.md' },
            { text: '前端 API', link: '/develop/fe-api.md' },
            { text: '前端 Event', link: '/develop/event.md' },
            { text: '插件开发', link: '/develop/plugs.md' },
            {
              text: '🔬 HTTP API',
              link: 'https://artalk.js.org/http-api.html',
            },
            {
              text: '📖 Contributing Guide',
              link: 'https://github.com/ArtalkJS/Artalk/blob/master/CONTRIBUTING.md',
            },
            {
              text: '🔖 Code of Conduct',
              link: 'https://github.com/ArtalkJS/Artalk/blob/master/CODE_OF_CONDUCT.md',
            },
          ],
        },
      ],
    },

    nav: [
      // NavbarItem
      {
        text: '介绍',
        link: '/guide/intro',
      },
      {
        text: '部署',
        link: '/guide/deploy',
      },
      {
        text: '配置',
        link: '/guide/backend/config',
      },
      {
        text: '迁移',
        link: '/guide/transfer',
      },
      {
        text: '案例',
        link: '/guide/cases',
      },
      {
        text: '开发',
        link: '/develop/',
      },
      // NavbarGroup
      {
        text: '传送',
        items: [
          {
            text: '代码仓库',
            link: 'https://github.com/ArtalkJS/Artalk',
          },
          // {
          //   text: "文档镜像 (国内)",
          //   link: "https://artalk-docs.qwqaq.com",
          // },
        ],
      },
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/ArtalkJS/Artalk' },
    ],

    search: {
      provider: 'algolia',
      options: {
        appId: '2WNJ32WVTY',
        apiKey: '6c6ebc345a87b738264f19095b78c91c',
        indexName: 'artalk-js',
        searchParameters: {
          facetFilters: ['lang:zh-CN'],
        },
      },
    },

    editLink: {
      pattern: 'https://github.com/ArtalkJS/Artalk/edit/master/docs/:path',
      text: '完善文档',
    },
  },

  vite: {
    server: {
      open: '/guide/intro.html',
    },
  },
})

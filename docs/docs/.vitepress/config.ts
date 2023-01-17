import { defineConfig } from 'vitepress'
import iterator from 'markdown-it-for-inline'
import * as ArtalkCDN from '../code/ArtalkCDN.json'
import * as Versions from '../code/ArtalkVersion.json'

export default defineConfig({
  title: 'Artalk',
  description: '一款简洁的自托管评论系统',
  lang: 'zh-CN',

  head: [
    ['link', { rel: 'icon', href: '/favicon.png' }],
    ['meta', { name: 'viewport', content: 'width=device-width, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0, user-scalable=no, target-densitydpi=device-dpi' }],
    // artalk
    ['link', { href: ArtalkCDN.CSS, rel: 'stylesheet' }],
    // ['script', { src: ArtalkCDN.JS }],
    // light gallery
    ['link', { href: 'https://npm.elemecdn.com/lightgallery@2.3.0/css/lightgallery.css', rel: 'stylesheet' }],
    ['script', { src: 'https://npm.elemecdn.com/lightgallery@2.3.0/lightgallery.min.js' }],
    // katex
    // ['link', { href: "https://npm.elemecdn.com/katex@0.15.3/dist/katex.min.css", rel: 'stylesheet' }],
    // ['script', { src: 'https://npm.elemecdn.com/katex@0.15.3/dist/katex.min.js' }],
    // ['script', { src: 'https://npm.elemecdn.com/@artalkjs/plugin-katex/dist/artalk-plugin-katex.js' }],
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
        tokens[idx].content = tokens[idx].content.replace(/:ArtalkVersion:/g, Versions.Artalk.replace(/^v/, ''));
      });
    },
  },

  themeConfig: {
    sidebar: {
      "/guide/": [
        {
          text: "快速开始",
          collapsible: true,
          items: [
            { text: '项目介绍', link: '/guide/intro.md' },
            { text: '程序部署', link: '/guide/deploy.md' },
            { text: '数据迁移', link: '/guide/transfer.md' },
          ]
        },
        {
          text: "前端",
          collapsible: true,
          items: [
            { text: '前端配置', link: '/guide/frontend/config.md' },
            { text: '侧边栏', link: '/guide/frontend/sidebar.md' },
            { text: '表情包', link: '/guide/frontend/emoticons.md' },
            { text: '浏览量统计', link: '/guide/frontend/pv.md' },
            { text: 'Latex', link: '/guide/frontend/latex.md' },
            { text: '图片灯箱', link: '/guide/frontend/lightbox.md' },
            { text: '多语言 (i18n)', link: '/guide/frontend/i18n.md' },
            { text: '置入博客', link: '/guide/frontend/import-blog.md' },
            { text: '置入框架', link: '/guide/frontend/import-framework.md' },
            { text: '精简版本', link: '/guide/frontend/artalk-lite.md' },
            { text: '扩展插件', link: '/guide/frontend/plugs.md' },
          ],
        },
        {
          text: '后端',
          collapsible: true,
          items: [
            { text: '后端配置', link: '/guide/backend/config.md' },
            { text: 'Docker', link: '/guide/backend/docker.md' },
            { text: '管理员 × 多站点', link: '/guide/backend/multi-site.md' },
            { text: '邮件通知', link: '/guide/backend/email.md' },
            { text: '多元推送', link: '/guide/backend/admin_notify.md' },
            { text: '图片上传', link: '/guide/backend/img-upload.md' },
            { text: '评论审核', link: '/guide/backend/moderator.md' },
            { text: '验证码', link: '/guide/backend/captcha.md' },
            { text: '在后端控制前端', link: '/guide/backend/fe-control.md' },
            { text: '相对 / 绝对路径', link: '/guide/backend/relative-path.md' },
            { text: '守护进程', link: '/guide/backend/daemon.md' },
            { text: '反向代理', link: '/guide/backend/reverse-proxy.md' },
            { text: '编译构建', link: '/guide/backend/build.md' },
            { text: '程序升级', link: '/guide/backend/update.md' },
          ]
        },
        {
          text: '更多内容',
          collapsible: true,
          items: [
            { text: '安全防范', link: '/guide/security.md' },
            { text: '扩展阅读', link: '/guide/extras.md' },
            { text: '案例展示', link: '/guide/cases.md' },
            { text: '关于我们', link: '/guide/about.md' },
          ]
        }
      ],
      "/develop/": [
        {
          text: '开发文档',
          items: [
            { text: '开发说明', link: '/develop/index.md', },
            { text: 'HTTP API', link: '/develop/api.md', },
            { text: 'Frontend Event', link: '/develop/event.md', },
          ]
        }
      ]
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
        link: '/guide/frontend/config',
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
            text: '前端仓库',
            link: 'https://github.com/ArtalkJS/Artalk',
          },
          {
            text: '后端仓库',
            link: 'https://github.com/ArtalkJS/ArtalkGo',
          },
          {
            text: '文档仓库',
            link: 'https://github.com/ArtalkJS/Docs',
          },
          {
            text: '文档镜像 (国内)',
            link: 'https://artalk-docs.qwqaq.com'
          }
        ],
      },
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/ArtalkJS/Artalk' }
    ],

    algolia: {
      appId: 'BH4D9OD16A',
      apiKey: '37ab96e3f5a774cbbf0571b035b42adb',
      indexName: 'artalk-js',
      searchParameters: {
        facetFilters: ['lang:zh-CN']
      }
    },

    editLink: {
      repo: 'ArtalkJS/Docs',
      branch: 'master',
      dir: 'docs',
      text: '完善文档'
    },
  }
})

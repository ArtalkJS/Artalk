import { defineConfig } from 'vitepress'

export const zh = defineConfig({
  lang: 'zh-CN',
  description: '由 Golang 驱动的自托管评论系统',

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
            { text: '社交登录', link: '/guide/frontend/auth.md' },
            { text: '评论审核', link: '/guide/backend/moderator.md' },
            { text: '验证码', link: '/guide/backend/captcha.md' },
            { text: '图片上传', link: '/guide/backend/img-upload.md' },
            { text: '账户与多站点', link: '/guide/backend/multi-site.md' },
            { text: '解析相对路径', link: '/guide/backend/relative-path.md' },
          ],
        },
        {
          text: '进阶指南',
          collapsed: false,
          items: [
            { text: '表情包', link: '/guide/frontend/emoticons.md' },
            { text: '浏览量统计', link: '/guide/frontend/pv.md' },
            { text: 'LaTeX', link: '/guide/frontend/latex.md' },
            { text: '图片灯箱', link: '/guide/frontend/lightbox.md' },
            { text: '图片懒加载', link: '/guide/frontend/img-lazy-load.md' },
            { text: 'IP 属地', link: '/guide/frontend/ip-region.md' },
            { text: '多语言', link: '/guide/frontend/i18n.md' },
            { text: '开发文档', link: '/develop/index.md' },
          ],
        },
        {
          text: '配置文档',
          collapsed: false,
          items: [
            { text: '环境变量', link: '/guide/env.md' },
            { text: '配置文件', link: '/guide/backend/config.md' },
            { text: '界面配置', link: '/guide/frontend/config.md' },
          ],
        },
        {
          text: '部署说明',
          collapsed: true,
          items: [
            { text: '守护进程', link: '/guide/backend/daemon.md' },
            { text: '反向代理', link: '/guide/backend/reverse-proxy.md' },
            { text: '编译构建', link: '/develop/contributing.md' },
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
            { text: '贡献指南', link: '/develop/contributing.md' },
            { text: '置入博客', link: '/develop/import-blog.md' },
            { text: '置入框架', link: '/develop/import-framework.md' },
            { text: '前端 API', link: '/develop/fe-api.md' },
            { text: '前端 Event', link: '/develop/event.md' },
            { text: '插件开发', link: '/develop/plugin.md' },
            { text: '兼容性', link: '/develop/compatibility.md' },
            {
              text: 'HTTP API',
              link: 'https://artalk.js.org/http-api.html',
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
    ],

    editLink: {
      pattern: 'https://github.com/ArtalkJS/Artalk/edit/master/docs/docs/:path',
      text: '完善文档',
    },
  },
})

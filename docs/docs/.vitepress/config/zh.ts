import { defineConfig } from 'vitepress'

export const zh = defineConfig({
  lang: 'zh-CN',
  description: '由 Golang 驱动的自托管评论系统',

  themeConfig: {
    sidebar: {
      '/zh/guide/': [
        {
          text: '快速开始',
          collapsed: false,
          items: [
            { text: '项目介绍', link: '/zh/guide/intro.md' },
            { text: '程序部署', link: '/zh/guide/deploy.md' },
            { text: '数据迁移', link: '/zh/guide/transfer.md' },
          ],
        },
        {
          text: '核心指南',
          collapsed: false,
          items: [
            { text: '侧边栏', link: '/zh/guide/frontend/sidebar.md' },
            { text: '邮件通知', link: '/zh/guide/backend/email.md' },
            { text: '多元推送', link: '/zh/guide/backend/admin_notify.md' },
            { text: '社交登录', link: '/zh/guide/frontend/auth.md' },
            { text: '评论审核', link: '/zh/guide/backend/moderator.md' },
            { text: '验证码', link: '/zh/guide/backend/captcha.md' },
            { text: '图片上传', link: '/zh/guide/backend/img-upload.md' },
            { text: '账户与多站点', link: '/zh/guide/backend/multi-site.md' },
            { text: '解析相对路径', link: '/zh/guide/backend/relative-path.md' },
          ],
        },
        {
          text: '进阶指南',
          collapsed: false,
          items: [
            { text: '表情包', link: '/zh/guide/frontend/emoticons.md' },
            { text: '浏览量统计', link: '/zh/guide/frontend/pv.md' },
            { text: 'LaTeX', link: '/zh/guide/frontend/latex.md' },
            { text: '图片灯箱', link: '/zh/guide/frontend/lightbox.md' },
            { text: '图片懒加载', link: '/zh/guide/frontend/img-lazy-load.md' },
            { text: '投票功能', link: '/zh/guide/frontend/voting.md' },
            { text: 'IP 属地', link: '/zh/guide/frontend/ip-region.md' },
            { text: '多语言', link: '/zh/guide/frontend/i18n.md' },
            { text: '开发文档', link: '/zh/develop/index.md' },
          ],
        },
        {
          text: '配置文档',
          collapsed: false,
          items: [
            { text: '环境变量', link: '/zh/guide/env.md' },
            { text: '配置文件', link: '/zh/guide/backend/config.md' },
            { text: '界面配置', link: '/zh/guide/frontend/config.md' },
          ],
        },
        {
          text: '部署说明',
          collapsed: true,
          items: [
            { text: '守护进程', link: '/zh/guide/backend/daemon.md' },
            { text: '反向代理', link: '/zh/guide/backend/reverse-proxy.md' },
            { text: '编译构建', link: '/zh/develop/contributing.md' },
            { text: '程序升级', link: '/zh/guide/backend/update.md' },
            { text: 'Docker', link: '/zh/guide/backend/docker.md' },
          ],
        },
        {
          text: '更多内容',
          collapsed: true,
          items: [
            { text: '安全防范', link: '/zh/guide/security.md' },
            { text: '扩展阅读', link: '/zh/guide/extras.md' },
            { text: '案例展示', link: '/zh/guide/cases.md' },
            { text: '关于我们', link: '/zh/guide/about.md' },
          ],
        },
      ],
      '/zh/develop/': [
        {
          text: '开发文档',
          items: [
            { text: '开发说明', link: '/zh/develop/index.md' },
            { text: '贡献指南', link: '/zh/develop/contributing.md' },
            { text: '置入博客', link: '/zh/develop/import-blog.md' },
            { text: '置入框架', link: '/zh/develop/import-framework.md' },
            { text: '前端 API', link: '/zh/develop/fe-api.md' },
            { text: '前端 Events', link: '/zh/develop/event.md' },
            { text: '前端 Types', link: 'https://artalk.js.org/typedoc/' },
            { text: '插件开发', link: '/zh/develop/plugin.md' },
            { text: '兼容性', link: '/zh/develop/compatibility.md' },
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
        link: '/zh/guide/intro',
      },
      {
        text: '部署',
        link: '/zh/guide/deploy',
      },
      {
        text: '配置',
        link: '/zh/guide/backend/config',
      },
      {
        text: '迁移',
        link: '/zh/guide/transfer',
      },
      {
        text: '案例',
        link: '/zh/guide/cases',
      },
      {
        text: '开发',
        link: '/zh/develop/',
      },
    ],

    editLink: {
      pattern: 'https://github.com/ArtalkJS/Artalk/edit/master/docs/docs/:path',
      text: '完善文档',
    },
  },
})

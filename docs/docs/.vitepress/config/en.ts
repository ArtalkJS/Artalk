import { defineConfig } from 'vitepress'

export const en = defineConfig({
  lang: 'en-US',
  description: 'A Self-hosted Comment System.',

  themeConfig: {
    sidebar: {
      '/en/guide/': [
        {
          text: 'Quick Start',
          collapsed: false,
          items: [
            { text: 'Introduction', link: '/en/guide/intro.md' },
            { text: 'Getting Started', link: '/en/guide/deploy.md' },
            { text: 'Data Migration', link: '/en/guide/transfer.md' },
          ],
        },
        {
          text: 'Basic Usage',
          collapsed: false,
          items: [
            { text: 'Sidebar', link: '/en/guide/frontend/sidebar.md' },
            { text: 'Email Notification', link: '/en/guide/backend/email.md' },
            { text: 'Multi-channel Notification', link: '/en/guide/backend/admin_notify.md' },
            { text: 'Social Login', link: '/en/guide/frontend/auth.md' },
            { text: 'Comment Moderation', link: '/en/guide/backend/moderator.md' },
            { text: 'Captcha', link: '/en/guide/backend/captcha.md' },
            { text: 'Image Upload', link: '/en/guide/backend/img-upload.md' },
            { text: 'Admins and Multi-Site', link: '/en/guide/backend/multi-site.md' },
            { text: 'Resolve Relative Path', link: '/en/guide/backend/relative-path.md' },
          ],
        },
        {
          text: 'Extensions',
          collapsed: false,
          items: [
            { text: 'Emoticons', link: '/en/guide/frontend/emoticons.md' },
            { text: 'Page View Statistics', link: '/en/guide/frontend/pv.md' },
            { text: 'LaTeX', link: '/en/guide/frontend/latex.md' },
            { text: 'Image Lightbox', link: '/en/guide/frontend/lightbox.md' },
            { text: 'Image Lazy Load', link: '/en/guide/frontend/img-lazy-load.md' },
            { text: 'Voting', link: '/zh/guide/frontend/voting.md' },
            { text: 'IP Region', link: '/en/guide/frontend/ip-region.md' },
            { text: 'Localization', link: '/en/guide/frontend/i18n.md' },
            { text: 'Development Documentation', link: '/en/develop/index.md' },
          ],
        },
        {
          text: 'Configurations',
          collapsed: false,
          items: [
            { text: 'Environment Variables', link: '/en/guide/env.md' },
            { text: 'Configuration File', link: '/en/guide/backend/config.md' },
            { text: 'UI Configuration', link: '/en/guide/frontend/config.md' },
          ],
        },
        {
          text: 'Deployment details',
          collapsed: false,
          items: [
            { text: 'Daemon Process', link: '/en/guide/backend/daemon.md' },
            { text: 'Reverse Proxy', link: '/en/guide/backend/reverse-proxy.md' },
            {
              text: 'Compile Source',
              link: 'https://github.com/ArtalkJS/Artalk/blob/master/CONTRIBUTING.md',
            },
            { text: 'Program Upgrade', link: '/en/guide/backend/update.md' },
            { text: 'Docker', link: '/en/guide/backend/docker.md' },
          ],
        },
      ],
      '/en/develop/': [
        {
          text: 'Development Documentation',
          items: [
            { text: 'Development Instructions', link: '/en/develop/index.md' },
            {
              text: 'Contribution Guide',
              link: 'https://github.com/ArtalkJS/Artalk/blob/master/CONTRIBUTING.md',
              target: '_blank',
            },
            { text: 'Import to Blog', link: '/en/develop/import-blog.md' },
            { text: 'Import to Framework', link: '/en/develop/import-framework.md' },
            { text: 'Frontend API', link: '/en/develop/fe-api.md' },
            { text: 'Frontend Events', link: '/en/develop/event.md' },
            { text: 'Frontend Types', link: 'https://artalk.js.org/typedoc/' },
            { text: 'Plugin Development', link: '/en/develop/plugin.md' },
            { text: 'Compatibility', link: '/en/develop/compatibility.md' },
            {
              text: 'HTTP API',
              link: 'https://artalk.js.org/http-api.html',
            },
          ],
        },
      ],
    },

    nav: [
      {
        text: 'Introduction',
        link: '/en/guide/intro',
      },
      {
        text: 'Installation',
        link: '/en/guide/deploy',
      },
      {
        text: 'Configuration',
        link: '/en/guide/backend/config',
      },
      {
        text: 'Migration',
        link: '/en/guide/transfer',
      },
      {
        text: 'Development',
        link: '/en/develop/',
      },
    ],

    editLink: {
      pattern: 'https://github.com/ArtalkJS/Artalk/edit/master/docs/docs/:path',
      text: 'Improve this document',
    },
  },
})

import { defineConfig } from 'vitepress'
import {
  GitChangelog,
  GitChangelogMarkdownSection,
} from '@nolebase/vitepress-plugin-git-changelog/vite'
import * as Version from '../../code/ArtalkVersion.json'

export const shared = defineConfig({
  title: 'Artalk',

  /* prettier-ignore */
  head: [
    ['link', { rel: 'icon', type: 'image/png', href: '/favicon.png' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:locale', content: 'en' }],
    ['meta', { property: 'og:title', content: 'Artalk - A Self-hosted Comment System'}],
    ['meta', { property: 'og:site_name', content: 'Artalk' }],
    ['meta', { property: 'og:image', content: 'https://artalk.js.org/assets/images/artalk-banner.png' }],
    ['meta', { property: 'og:url', content: 'https://artalk.js.org/' }],
    ['meta', { name: 'viewport', content: 'width=device-width, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0, user-scalable=no, target-densitydpi=device-dpi' }],
  ],

  markdown: {
    // @link https://github.com/shikijs/shiki
    theme: {
      light: 'github-light',
      dark: 'github-dark',
    },
    config: (md) => {
      const renderVersion = (c: string) => String(c).replace(/:ArtalkVersion:/g, Version.latest)
      md.core.ruler.push('artalk_version', (state) => {
        state.tokens?.forEach((token) => {
          if (token.type === 'inline') {
            token.children?.forEach((child) => {
              if (['text', 'link_open', 'code_inline'].includes(child.type))
                child.content = renderVersion(child.content)
            })
          } else if (token.type === 'fence') {
            token.content = renderVersion(token.content)
          }
        })
      })
    },
  },

  sitemap: {
    hostname: 'https://artalk.js.org',
    transformItems(items) {
      return items.filter((item) => !item.url.includes('migration'))
    },
  },

  lastUpdated: true,
  cleanUrls: false,
  metaChunk: true,

  themeConfig: {
    socialLinks: [{ icon: 'github', link: 'https://github.com/ArtalkJS/Artalk' }],

    search: {
      provider: 'algolia',
      options: {
        appId: '2WNJ32WVTY',
        apiKey: '6c6ebc345a87b738264f19095b78c91c',
        indexName: 'artalk-js',
      },
    },
  },

  vite: {
    server: {
      open: '/zh/guide/intro.html',
    },
    plugins: [
      GitChangelog({
        repoURL: () => 'https://github.com/ArtalkJS/Artalk',
        mapAuthors: [
          {
            name: 'qwqcode',
            username: 'qwqcode',
            mapByEmailAliases: ['qwqcode@gmail.com', '22412567+qwqcode@users.noreply.github.com'],
          },
          {
            name: 'pluveto',
            username: 'pluveto',
            mapByEmailAliases: ['i@pluvet.com', '50045289+pluveto@users.noreply.github.com'],
          },
          {
            name: 'Mr.Hope',
            username: 'Mister-Hope',
            mapByEmailAliases: ['mister-hope@outlook.com'],
          },
        ],
      }),
      GitChangelogMarkdownSection(),
    ],
  },
})

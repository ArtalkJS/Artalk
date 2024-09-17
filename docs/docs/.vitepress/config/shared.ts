import { defineConfig } from 'vitepress'
import iterator from 'markdown-it-for-inline'
import * as Version from '../../code/ArtalkVersion.json'

export const shared = defineConfig({
  title: 'Artalk',

  rewrites: {
    'zh/:rest*': ':rest*',
  },

  /* prettier-ignore */
  head: [
    ['link', { rel: 'icon', type: 'image/png', href: '/favicon.png' }],
    ['meta', { name: 'theme-color', content: '#007bff' }],
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
      md.use(iterator, 'artalk_version', 'text', function (tokens, idx) {
        tokens[idx].content = tokens[idx].content.replace(/:ArtalkVersion:/g, Version.latest)
      })
      md.use(iterator, 'artalk_version_link', 'link_open', (tokens, idx) => {
        const href = tokens[idx].attrGet('href')
        tokens[idx].attrSet('href', href.replace(/:ArtalkVersion:/g, Version.latest))
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
  cleanUrls: true,
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
      open: '/guide/intro.html',
    },
  },
})

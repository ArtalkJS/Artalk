<p align="center">
<img src="https://user-images.githubusercontent.com/22412567/137740516-d9e97af0-fb3b-4dab-b331-671a9a2a3a63.png" alt="Artalk" width="100%">
</p>

# Artalk

[![npm 版本](https://img.shields.io/npm/v/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![npm 下载量](https://img.shields.io/npm/dt/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![尺寸](https://badgen.net/bundlephobia/minzip/artalk?style=flat-square)](https://bundlephobia.com/package/artalk)
[![CircleCI](https://circleci.com/gh/ArtalkJS/Artalk/tree/master.svg?style=svg)](https://circleci.com/gh/ArtalkJS/Artalk/tree/master)

> 🌌 A self-hosted comment system

[简体中文](./README.md) / [Documentation](https://artalk.js.org) / [Releases](https://github.com/ArtalkJS/Artalk/releases) / [Artalk](https://github.com/ArtalkJS/Artalk)

---

- 🍃 Lightweight (~30kB)
- 👨‍👧‍👦 Secure (Self-hosted)
- 🐳 Easy to use (Very Easy)
- 🍱 Golang backend (Fast and Cross Platform)
- 🌊 TypeScript × Vanilla × Vite (No Dependencies)

## Features

- Sidebar: Multi-site centralized management
- Notification Center: Red badge alert / Mark as read
- Verification: User Verification Badge / Password access
- Moderation: Anti-spam detection / Captcha frequency limit
- Emoticons: Insert emoji / Quickly import presets
- Email Notify: Template customization / Send to multi-admin
- Site Isolation: Multi-site management / Admin assignment
- Page Management: Provide page title / Easy to look up
- Image Upload: Upload to localhost / various remote image host
- Private Space Mode: Only visible to yourself / Message board
- Multiple Msg Pushers: Support Telegram / Slack / LINE
- Nested Comments: Switchable to Flat Mode
- Comment Vote: For or against comments
- Comment Sort: Sort by popularity or time
- Comment PIN: Pin important comments
- Track only author: Show only author's comments
- Asynchronous: Send comments without waiting
- Infinite Scrolling: Provide various comment content pagination custom
- AutoSave: A anti-lost editor
- AutoFill: Autofill user profile
- Preview: Real-time preview of markdown
- Dark Mode: Prevents eye disease damage
- Collapse Comment: I do not want you to see this comment
- Data Backup: Prevent comment data loss
- Data Migration: Switch back and forth between different commenting systems
- Multiple comments on one page: Multiple comment areas on one page (seems useless
- Markdown: Markdown syntax was supported by default
- LaTex support: Import KaTex plugin to integrate LaTex parser for artalk
- Using [Vite](https://github.com/vitejs/vite): The ultimate developer experience

## Getting Started

Reference to：[**Documentation**](https://artalk.js.org/guide/deploy.html)

```sh
$ pnpm add artalk
```

<!-- prettier-ignore-start -->


```ts
import Artalk from 'artalk'

Artalk.init({
  el:        '#Comments',
  pageKey:   'http://your_domain/post/1', // Page Link
  pageTitle: 'The title of your page',    // Page Title
  server:    'http://localhost:8080/api', // Server URL
  site:      'Site Name',
  locale:    'en'
})
```

<!-- prettier-ignore-end -->

### Docker

```sh
docker run -d \
  --name artalk \
  -p 8080:23366 \
  -v $(pwd)/data:/data \
  artalk/artalk-go
```

### Docker Compose

```sh
mkdir Artalk
cd Artalk

vim docker-compose.yaml
```

```yaml
version: '3.5'
services:
  artalk:
    container_name: artalk
    image: artalk/artalk-go
    ports:
      - 8080:23366
    volumes:
      - ./data:/data
```

```sh
docker-compose up -d
```

## Development

see [CONTRIBUTING.md](./CONTRIBUTING.md)

## Contributors

[![](https://contrib.rocks/image?repo=ArtalkJS/Artalk)](https://github.com/ArtalkJS/Artalk/graphs/contributors)

## Supporters

[![Stargazers repo roster for @ArtalkJS/Artalk](https://reporoster.com/stars/ArtalkJS/Artalk)](https://github.com/ArtalkJS/Artalk/stargazers)

## Feedback

Thanks for the help and feedback provided by the community, if you have good suggestions or comments, please go to [issues](https://github.com/ArtalkJS/Artalk/issues) to let us know at any time.

## Stargazers over time

[![Stargazers over time](https://starchart.cc/ArtalkJS/Artalk.svg)](https://starchart.cc/ArtalkJS/Artalk)

## License

[MIT](./LICENSE)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_shield)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_large)

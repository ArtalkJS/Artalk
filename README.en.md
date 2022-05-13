<p align="center">
<img src="https://user-images.githubusercontent.com/22412567/137740516-d9e97af0-fb3b-4dab-b331-671a9a2a3a63.png" alt="Artalk" width="100%">
</p>

# Artalk

[![](https://img.shields.io/npm/v/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![](https://img.shields.io/npm/dt/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![](https://badgen.net/bundlephobia/minzip/artalk?style=flat-square)](https://bundlephobia.com/package/artalk)
[![CircleCI](https://circleci.com/gh/ArtalkJS/Artalk/tree/master.svg?style=svg)](https://circleci.com/gh/ArtalkJS/Artalk/tree/master)

> 🌌 A Self-hosted comment system

[简体中文](./README.md) / [Documentation](https://artalk.js.org) / [Releases](https://github.com/ArtalkJS/ArtalkGo/releases) / [ArtalkGo](https://github.com/ArtalkJS/ArtalkGo)

---

- 🍃 Lightweight (~30kB)
- 👨‍👧‍👦 Secure (Self-hosted)
- 🐳 Easy to use (Very Easy)
- 🍱 Golang backend (Fast and Cross Platform)
- 🌊 TypeScript × Vanilla × Vite (No Dependencies)

## Features

- Sidebar: Multi-site centralized management
- Notification Center: Red badge alert / Read mark
- Authentication: User badge customization / Password verification
- Moderation: Anti-spam detection / captcha frequency limit
- Emoji: Insert emoji / Quickly import emoji
- Email reminder: Template customization / multi-admin notification
- Site isolation: Multi-site management / admin assignment
- Page management: Title can be displayed / easy to look up
- Image upload: Upload to local / various remote image host
- Private Space Mode: Only visible to yourself / Message board
- Multiple push: Support Telegram / Slack / LINE
- Nesting comments: Switchable to flat mode
- Comment voting: For or against comments
- Comment sorting: Sort by popularity or time
- Comment PIN: Pin the important comments on top
- Look at author only: Show only author's comments
- Asynchronous processing: Send comments without waiting
- Scroll loading: Various comment content pagination customization
- Auto save: A anti-lost editor
- Autofill: client-autofill
- Real-time preview: Real-time preview of comment content
- Dark Mode: Prevents eye disease damage
- Comment Folding: I do not want you to see this comment
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

```ts
import Artalk from 'artalk'

new Artalk({
  el:        '#Comments',
  pageKey:   'http://your_domain/post/1', // Page Link
  pageTitle: 'The title of your page',    // Page Title
  server:    'http://localhost:8080/api', // Server URL
  site:      'Site Name',
  i18n:      'en-US'
})
```

### Docker

```sh
# Create a directory for Artalk
mkdir Artalk
cd Artalk

# Download config template
curl -L https://raw.githubusercontent.com/ArtalkJS/ArtalkGo/master/artalk-go.example.yml > conf.yml

docker run -d \
  --name artalk \
  -p 0.0.0.0:8080:23366 \
  -v $(pwd)/conf.yml:/conf.yml \
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
version: "3.5"
services:
  artalk:
    container_name: artalk
    image: artalk/artalk-go
    ports:
      - 8080:23366
    volumes:
      - ./conf.yml:/conf.yml
      - ./data:/data
```

```sh
docker-compose up -d
```

## Contributors

[![](https://contrib.rocks/image?repo=ArtalkJS/Artalk)](https://github.com/ArtalkJS/Artalk/graphs/contributors)

## Supporters

[![Stargazers repo roster for @ArtalkJS/Artalk](https://reporoster.com/stars/ArtalkJS/Artalk)](https://github.com/ArtalkJS/Artalk/stargazers)

## Feedback

Thanks for the help and feedback provided by the community, if you have good suggestions or comments, please go to [issues](https://github.com/ArtalkJS/Artalk/issues) to let us know at any time.

## Stargazers over time

[![Stargazers over time](https://starchart.cc/ArtalkJS/Artalk.svg)](https://starchart.cc/ArtalkJS/Artalk)

## License

[LGPL-3.0](./LICENSE)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_shield)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_large)

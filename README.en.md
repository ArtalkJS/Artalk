<p align="center">
<img src="https://user-images.githubusercontent.com/22412567/171680920-6e74b77c-c565-487b-bff1-4f94976ecbe7.png" alt="Artalk" width="100%">
</p>

# Artalk

[![npm version](https://img.shields.io/npm/v/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![npm downloads](https://img.shields.io/npm/dt/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![Docker Pulls](https://img.shields.io/docker/pulls/artalk/artalk-go?style=flat-square)](https://hub.docker.com/r/artalk/artalk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/ArtalkJS/Artalk?style=flat-square)](https://goreportcard.com/report/github.com/ArtalkJS/Artalk)
[![CircleCI](https://img.shields.io/circleci/build/gh/ArtalkJS/Artalk?style=flat-square)](https://circleci.com/gh/ArtalkJS/Artalk/tree/master)
[![Codecov](https://img.shields.io/codecov/c/gh/ArtalkJS/Artalk?style=flat-square)](https://codecov.io/gh/ArtalkJS/Artalk)
[![npm bundle size](https://img.shields.io/bundlephobia/minzip/artalk?style=flat-square)](https://bundlephobia.com/package/artalk)

[中文](./README.md) • [Official Website](https://artalk.js.org) • [Latest Release](https://github.com/ArtalkJS/Artalk/releases) • [Changelog](https://github.com/ArtalkJS/Artalk/blob/master/CHANGELOG.md)

🌌 Artalk is a simple yet feature-rich commenting system that you can deploy out of the box and integrate into any blog, website, or web application.

- 🍃 Frontend ~40KB, pure Vanilla JS
- 🍱 Backend Golang, efficient lightweight cross-platform
- 🐳 One-click deployment via Docker, convenient and fast
- 🌈 Open-source program, self-hosted, privacy-first

## Features

<!-- prettier-ignore-start -->

<!-- features -->
* [Sidebar](https://artalk.js.org/guide/frontend/sidebar.html): Quick management, intuitive browsing
* [Social Login](https://artalk.js.org/guide/frontend/auth.html): Quick login via social accounts
* [Email Notification](https://artalk.js.org/guide/backend/email.html): Multiple sending methods, email templates
* [Multi-channel Push](https://artalk.js.org/guide/backend/admin_notify.html): Multiple push methods, notification templates
* [In-site Notifications](https://artalk.js.org/guide/frontend/sidebar.html): Red dot marking, mention list
* [CAPTCHA](https://artalk.js.org/guide/backend/captcha.html): Multiple verification types, frequency limitation
* [Comment Moderation](https://artalk.js.org/guide/backend/moderator.html): Content detection, spam interception
* [Image Upload](https://artalk.js.org/guide/backend/img-upload.html): Custom upload, support for image hosting
* [Markdown](https://artalk.js.org/guide/intro.html): Support Markdown syntax
* [Emoticons](https://artalk.js.org/guide/frontend/emoticons.html): Compatible with OwO, quick integration
* [Multi-site](https://artalk.js.org/guide/backend/multi-site.html): Site isolation, centralized management
* [Admin](https://artalk.js.org/guide/backend/multi-site.html): Password verification, badge identification
* [Page Management](https://artalk.js.org/guide/frontend/sidebar.html): Quick view, one-click title jump
* [Page Views Statistics](https://artalk.js.org/guide/frontend/pv.html): Easily track page views
* [Hierarchy](https://artalk.js.org/guide/frontend/config.html#nestmax): Nested paginated lists, scroll loading
* [Comment Voting](https://artalk.js.org/guide/frontend/config.html#vote): Upvote or downvote comments
* [Comment Sorting](https://artalk.js.org/guide/frontend/config.html#listsort): Multiple sorting options, freedom to choose
* [Comment Search](https://artalk.js.org/guide/frontend/sidebar.html): Quickly search comment content
* [Comment Pinning](https://artalk.js.org/guide/frontend/sidebar.html): Pin important messages
* [Author-only View](https://artalk.js.org/guide/frontend/config.html): Show only comments from the author
* [Comment Jumping](https://artalk.js.org/guide/intro.html): Quickly jump to referenced comments
* [Auto Save](https://artalk.js.org/guide/frontend/config.html): Content input auto-saving
* [IP Geolocation](https://artalk.js.org/guide/frontend/ip-region.html): User IP location display
* [Data Migration](https://artalk.js.org/guide/transfer.html): Free migration, quick backup
* [Image Lightbox](https://artalk.js.org/guide/frontend/lightbox.html): Quickly integrate image lightbox
* [Image Lazy Loading](https://artalk.js.org/guide/frontend/img-lazy-load.html): Delay loading images, optimize experience
* [Latex](https://artalk.js.org/guide/frontend/latex.html): Latex formula parsing integration
* [Dark Mode](https://artalk.js.org/guide/frontend/config.html#darkmode): Dark mode switching
* [Extension Plugins](https://artalk.js.org/develop/): Create more possibilities
* [Multi-language](https://artalk.js.org/guide/frontend/i18n.html): Multi-language switching
* [Command Line](https://artalk.js.org/guide/backend/config.html): Command line operation management capability
* [API Documentation](https://artalk.js.org/develop/): Provides OpenAPI format documentation
* [Program Upgrade](https://artalk.js.org/guide/backend/update.html): Version detection, one-click upgrade
<!-- /features -->

<!-- prettier-ignore-end -->

## Installation

Deploy via Docker with one click:

```bash
docker run -d --name artalk -p 8080:23366 -v $(pwd)/data:/data artalk/artalk-go
```

Integrate Artalk into your webpage:

<!-- prettier-ignore-start -->

```ts
Artalk.init({
  el:      '#Comments',
  site:    'Artalk Blog',
  server:  'https://artalk.example.com',
  pageKey: '/2018/10/02/hello-world.html'
})
```

<!-- prettier-ignore-end -->

[**Learn More →**](https://artalk.js.org/guide/deploy.html)

## For Developers

Pull requests are welcome!

See [Development](https://artalk.js.org/develop/) and [Contributing](./CONTRIBUTING.md) for information on working with the codebase, getting a local development setup, and contributing changes.

## Contributors

Your contributions enrich the open-source community, fostering learning, inspiration, and innovation. We deeply value your involvement. Thank you for being a vital part of our community! 🥰

[![](https://contrib.rocks/image?repo=ArtalkJS/Artalk)](https://github.com/ArtalkJS/Artalk/graphs/contributors)

## Supporters

[![Stargazers repo roster for @ArtalkJS/Artalk](https://reporoster.com/stars/ArtalkJS/Artalk)](https://github.com/ArtalkJS/Artalk/stargazers)

## Repobeats Analytics

![Alt](https://repobeats.axiom.co/api/embed/a9fc9191ac561bc5a8ee2cddc81e635ecaebafb6.svg 'Repobeats analytics image')

## Stargazers over time

[![Stargazers over time](https://starchart.cc/ArtalkJS/Artalk.svg)](https://starchart.cc/ArtalkJS/Artalk)

## License

[MIT](./LICENSE)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_shield)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_large)

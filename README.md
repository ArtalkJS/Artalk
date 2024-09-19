<p align="center">
<img src="https://user-images.githubusercontent.com/22412567/171680920-6e74b77c-c565-487b-bff1-4f94976ecbe7.png" alt="Artalk" width="100%">
</p>

# Artalk

[![npm version](https://img.shields.io/npm/v/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![npm downloads](https://img.shields.io/npm/dt/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![Docker Pulls](https://img.shields.io/docker/pulls/artalk/artalk-go?style=flat-square)](https://hub.docker.com/r/artalk/artalk-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/artalkjs/artalk/v2.svg)](https://pkg.go.dev/github.com/artalkjs/artalk/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/ArtalkJS/Artalk?style=flat-square)](https://goreportcard.com/report/github.com/ArtalkJS/Artalk)
[![CircleCI](https://img.shields.io/circleci/build/gh/ArtalkJS/Artalk?style=flat-square)](https://circleci.com/gh/ArtalkJS/Artalk/tree/master)
[![Codecov](https://img.shields.io/codecov/c/gh/ArtalkJS/Artalk?style=flat-square)](https://codecov.io/gh/ArtalkJS/Artalk)
[![npm bundle size](https://img.shields.io/bundlephobia/minzip/artalk?style=flat-square)](https://bundlephobia.com/package/artalk)

[Homepage](https://artalk.js.org) • [Documentation](https://artalk.js.org/en/guide/deploy.html) • [Latest Release](https://github.com/ArtalkJS/Artalk/releases) • [Changelog](https://github.com/ArtalkJS/Artalk/blob/master/CHANGELOG.md) • [简体中文](./README.zh.md)

Artalk is an intuitive yet feature-rich comment system, ready for immediate deployment into any blog, website, or web application.

- 🍃 Client ~40KB, crafted with pure Vanilla JS, framework-agnostic
- 🍱 Server powered by Golang, offering efficient and lightweight cross-platform performance
- 🐳 One-click deployment via Docker, ensuring ease and speed
- 🌈 Open-source software, self-hosted with privacy as a priority

## Features

<!-- prettier-ignore-start -->

<!-- features -->
* [Sidebar](https://artalk.js.org/guide/frontend/sidebar.html): Quick management, intuitive browsing
* [Social Login](https://artalk.js.org/guide/frontend/auth.html): Fast login via social accounts
* [Email Notification](https://artalk.js.org/guide/backend/email.html): Various sending methods, email templates
* [Diverse Push](https://artalk.js.org/guide/backend/admin_notify.html): Multiple push methods, notification templates
* [Site Notification](https://artalk.js.org/guide/frontend/sidebar.html): Red dot marks, mention list
* [Captcha](https://artalk.js.org/guide/backend/captcha.html): Various verification types, frequency limits
* [Comment Moderation](https://artalk.js.org/guide/backend/moderator.html): Content detection, spam interception
* [Image Upload](https://artalk.js.org/guide/backend/img-upload.html): Custom upload, supports image hosting
* [Markdown](https://artalk.js.org/guide/intro.html): Supports Markdown syntax
* [Emoji Pack](https://artalk.js.org/guide/frontend/emoticons.html): Compatible with OwO, quick integration
* [Multi-Site](https://artalk.js.org/guide/backend/multi-site.html): Site isolation, centralized management
* [Admin](https://artalk.js.org/guide/backend/multi-site.html): Password verification, badge identification
* [Page Management](https://artalk.js.org/guide/frontend/sidebar.html): Quick view, one-click title navigation
* [Page View Statistics](https://artalk.js.org/guide/frontend/pv.html): Easily track page views
* [Hierarchical Structure](https://artalk.js.org/guide/frontend/config.html#nestmax): Nested paginated list, infinite scroll
* [Comment Voting](https://artalk.js.org/guide/frontend/config.html#vote): Upvote or downvote comments
* [Comment Sorting](https://artalk.js.org/guide/frontend/config.html#listsort): Various sorting options, freely selectable
* [Comment Search](https://artalk.js.org/guide/frontend/sidebar.html): Quick comment content search
* [Comment Pinning](https://artalk.js.org/guide/frontend/sidebar.html): Pin important messages
* [View Author Only](https://artalk.js.org/guide/frontend/config.html): Show only the author's comments
* [Comment Jump](https://artalk.js.org/guide/intro.html): Quickly jump to quoted comment
* [Auto Save](https://artalk.js.org/guide/frontend/config.html): Content loss prevention
* [IP Region](https://artalk.js.org/guide/frontend/ip-region.html): Display user's IP region
* [Data Migration](https://artalk.js.org/guide/transfer.html): Free migration, quick backup
* [Image Lightbox](https://artalk.js.org/guide/frontend/lightbox.html): Quick integration of image lightbox
* [Image Lazy Load](https://artalk.js.org/guide/frontend/img-lazy-load.html): Lazy load images, optimize experience
* [Latex](https://artalk.js.org/guide/frontend/latex.html): Integrate Latex formula parsing
* [Night Mode](https://artalk.js.org/guide/frontend/config.html#darkmode): Switch to night mode
* [Extension Plugin](https://artalk.js.org/develop/plugin.html): Create more possibilities
* [Multi-Language](https://artalk.js.org/guide/frontend/i18n.html): Switch between multiple languages
* [Command Line](https://artalk.js.org/guide/backend/config.html): Command line operation management
* [API Documentation](https://artalk.js.org/http-api.html): Provides OpenAPI format documentation
* [Program Upgrade](https://artalk.js.org/guide/backend/update.html): Version check, one-click upgrade
<!-- /features -->

<!-- prettier-ignore-end -->

## Installation

Deploy Artalk Server with Docker in one step:

```bash
docker run -d \
    --name artalk \
    -p 8080:23366 \
    -v $(pwd)/data:/data \
    -e "TZ=America/New_York" \
    -e "ATK_LOCALE=en" \
    -e "ATK_SITE_DEFAULT=Artalk Blog" \
    -e "ATK_SITE_URL=https://example.com" \
    artalk/artalk-go
```

Integrate Artalk Client into your webpage:

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

We offer various installation methods, including binary files, go install, and package managers for Linux distributions.

[**Learn More →**](https://artalk.js.org/en/guide/deploy.html)

## For Developers

Pull requests are welcome!

See [Development](https://artalk.js.org/en/develop/) and [Contributing](./CONTRIBUTING.md) for information on working with the codebase, getting a local development setup, and contributing changes.

## Contributors

Your contributions enrich the open-source community, fostering learning, inspiration, and innovation. We deeply value your involvement. Thank you for being a vital part of our community! 🥰

[![](https://contrib.rocks/image?repo=ArtalkJS/Artalk)](https://github.com/ArtalkJS/Artalk/graphs/contributors)

## Supporters

[![Stargazers repo roster for @ArtalkJS/Artalk](https://reporoster.com/stars/ArtalkJS/Artalk)](https://github.com/ArtalkJS/Artalk/stargazers)

## Repobeats Analytics

![Alt](https://repobeats.axiom.co/api/embed/a9fc9191ac561bc5a8ee2cddc81e635ecaebafb6.svg 'Repobeats analytics image')

## Stargazers over time

<a href="https://trendshift.io/repositories/6290" target="_blank"><img src="https://trendshift.io/api/badge/repositories/6290" alt="ArtalkJS%2FArtalk | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>

[![Stargazers over time](https://starchart.cc/ArtalkJS/Artalk.svg)](https://starchart.cc/ArtalkJS/Artalk)

## License

[MIT](./LICENSE)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_shield)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_large)

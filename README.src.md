<!--@nrg.languages=en,zh-->
<!--@nrg.defaultLanguage=en-->
<!--@nrg.fileNamePattern.zh=README.zh.md-->

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

[Homepage](https://artalk.js.org) • [Documentation](https://artalk.js.org/en/guide/deploy.html) • [Latest Release](https://github.com/ArtalkJS/Artalk/releases) • [Changelog](https://github.com/ArtalkJS/Artalk/blob/master/CHANGELOG.md) • [简体中文](./README.zh.md)<!--en-->
[官方网站](https://artalk.js.org) • [最新版本](https://github.com/ArtalkJS/Artalk/releases) • [更新日志](https://github.com/ArtalkJS/Artalk/blob/master/CHANGELOG.md) • [English](./README.md)<!--zh-->

Artalk is an intuitive yet feature-rich comment system, ready for immediate deployment into any blog, website, or web application.<!--en-->
Artalk 是一款简单易用但功能丰富的评论系统，你可以开箱即用地部署并置入任何博客、网站、Web 应用。<!--zh-->

- 🍃 Client ~40KB, crafted with pure Vanilla JS, framework-agnostic<!--en-->
- 🍃 前端 ~40KB，纯天然 Vanilla JS<!--zh-->
- 🍱 Server powered by Golang, offering efficient and lightweight cross-platform performance<!--en-->
- 🍱 后端 Golang，高效轻量跨平台<!--zh-->
- 🐳 One-click deployment via Docker, ensuring ease and speed<!--en-->
- 🐳 通过 Docker 一键部署，方便快捷<!--zh-->
- 🌈 Open-source software, self-hosted with privacy as a priority<!--en-->
- 🌈 开源程序，自托管，隐私至上<!--zh-->

## Features<!--en-->
## 特性<!--zh-->

<!-- prettier-ignore-start -->

<!-- features -->
* [Sidebar](https://artalk.js.org/guide/frontend/sidebar.html): Quick management, intuitive browsing<!--en-->
* [侧边栏](https://artalk.js.org/guide/frontend/sidebar.html): 快速管理、直观浏览<!--zh-->
* [Social Login](https://artalk.js.org/guide/frontend/auth.html): Fast login via social accounts<!--en-->
* [社交登录](https://artalk.js.org/guide/frontend/auth.html): 通过社交账号快速登录<!--zh-->
* [Email Notification](https://artalk.js.org/guide/backend/email.html): Various sending methods, email templates<!--en-->
* [邮件通知](https://artalk.js.org/guide/backend/email.html): 多种发送方式、邮件模板<!--zh-->
* [Diverse Push](https://artalk.js.org/guide/backend/admin_notify.html): Multiple push methods, notification templates<!--en-->
* [多元推送](https://artalk.js.org/guide/backend/admin_notify.html): 多种推送方式、通知模版<!--zh-->
* [Site Notification](https://artalk.js.org/guide/frontend/sidebar.html): Red dot marks, mention list<!--en-->
* [站内通知](https://artalk.js.org/guide/frontend/sidebar.html): 红点标记、提及列表<!--zh-->
* [Captcha](https://artalk.js.org/guide/backend/captcha.html): Various verification types, frequency limits<!--en-->
* [验证码](https://artalk.js.org/guide/backend/captcha.html): 多种验证类型、频率限制<!--zh-->
* [Comment Moderation](https://artalk.js.org/guide/backend/moderator.html): Content detection, spam interception<!--en-->
* [评论审核](https://artalk.js.org/guide/backend/moderator.html): 内容检测、垃圾拦截<!--zh-->
* [Image Upload](https://artalk.js.org/guide/backend/img-upload.html): Custom upload, supports image hosting<!--en-->
* [图片上传](https://artalk.js.org/guide/backend/img-upload.html): 自定义上传、支持图床<!--zh-->
* [Markdown](https://artalk.js.org/guide/intro.html): Supports Markdown syntax<!--en-->
* [Markdown](https://artalk.js.org/guide/intro.html): 支持 Markdown 语法<!--zh-->
* [Emoji Pack](https://artalk.js.org/guide/frontend/emoticons.html): Compatible with OwO, quick integration<!--en-->
* [表情包](https://artalk.js.org/guide/frontend/emoticons.html): 兼容 OwO，快速集成<!--zh-->
* [Multi-Site](https://artalk.js.org/guide/backend/multi-site.html): Site isolation, centralized management<!--en-->
* [多站点](https://artalk.js.org/guide/backend/multi-site.html): 站点隔离、集中管理<!--zh-->
* [Admin](https://artalk.js.org/guide/backend/multi-site.html): Password verification, badge identification<!--en-->
* [管理员](https://artalk.js.org/guide/backend/multi-site.html): 密码验证、徽章标识<!--zh-->
* [Page Management](https://artalk.js.org/guide/frontend/sidebar.html): Quick view, one-click title navigation<!--en-->
* [页面管理](https://artalk.js.org/guide/frontend/sidebar.html): 快速查看、标题一键跳转<!--zh-->
* [Page View Statistics](https://artalk.js.org/guide/frontend/pv.html): Easily track page views<!--en-->
* [浏览量统计](https://artalk.js.org/guide/frontend/pv.html): 轻松统计网页浏览量<!--zh-->
* [Hierarchical Structure](https://artalk.js.org/guide/frontend/config.html#nestmax): Nested paginated list, infinite scroll<!--en-->
* [层级结构](https://artalk.js.org/guide/frontend/config.html#nestmax): 嵌套分页列表、滚动加载<!--zh-->
* [Comment Voting](https://artalk.js.org/guide/frontend/config.html#vote): Upvote or downvote comments<!--en-->
* [评论投票](https://artalk.js.org/guide/frontend/config.html#vote): 赞同或反对评论<!--zh-->
* [Comment Sorting](https://artalk.js.org/guide/frontend/config.html#listsort): Various sorting options, freely selectable<!--en-->
* [评论排序](https://artalk.js.org/guide/frontend/config.html#listsort): 多种排序方式，自由选择<!--zh-->
* [Comment Search](https://artalk.js.org/guide/frontend/sidebar.html): Quick comment content search<!--en-->
* [评论搜索](https://artalk.js.org/guide/frontend/sidebar.html): 快速搜索评论内容<!--zh-->
* [Comment Pinning](https://artalk.js.org/guide/frontend/sidebar.html): Pin important messages<!--en-->
* [评论置顶](https://artalk.js.org/guide/frontend/sidebar.html): 重要消息置顶显示<!--zh-->
* [View Author Only](https://artalk.js.org/guide/frontend/config.html): Show only the author's comments<!--en-->
* [仅看作者](https://artalk.js.org/guide/frontend/config.html): 仅显示作者的评论<!--zh-->
* [Comment Jump](https://artalk.js.org/guide/intro.html): Quickly jump to quoted comment<!--en-->
* [评论跳转](https://artalk.js.org/guide/intro.html): 快速跳转到引用的评论<!--zh-->
* [Auto Save](https://artalk.js.org/guide/frontend/config.html): Content loss prevention<!--en-->
* [自动保存](https://artalk.js.org/guide/frontend/config.html): 输入内容防丢功能<!--zh-->
* [IP Region](https://artalk.js.org/guide/frontend/ip-region.html): Display user's IP region<!--en-->
* [IP 属地](https://artalk.js.org/guide/frontend/ip-region.html): 用户 IP 属地展示<!--zh-->
* [Data Migration](https://artalk.js.org/guide/transfer.html): Free migration, quick backup<!--en-->
* [数据迁移](https://artalk.js.org/guide/transfer.html): 自由迁移、快速备份<!--zh-->
* [Image Lightbox](https://artalk.js.org/guide/frontend/lightbox.html): Quick integration of image lightbox<!--en-->
* [图片灯箱](https://artalk.js.org/guide/frontend/lightbox.html): 图片灯箱快速集成<!--zh-->
* [Image Lazy Load](https://artalk.js.org/guide/frontend/img-lazy-load.html): Lazy load images, optimize experience<!--en-->
* [图片懒加载](https://artalk.js.org/guide/frontend/img-lazy-load.html): 延迟加载图片，优化体验<!--zh-->
* [Latex](https://artalk.js.org/guide/frontend/latex.html): Integrate Latex formula parsing<!--en-->
* [Latex](https://artalk.js.org/guide/frontend/latex.html): Latex 公式解析集成<!--zh-->
* [Night Mode](https://artalk.js.org/guide/frontend/config.html#darkmode): Switch to night mode<!--en-->
* [夜间模式](https://artalk.js.org/guide/frontend/config.html#darkmode): 夜间模式切换<!--zh-->
* [Extension Plugin](https://artalk.js.org/develop/plugin.html): Create more possibilities<!--en-->
* [扩展插件](https://artalk.js.org/develop/plugin.html): 创造更多可能性<!--zh-->
* [Multi-Language](https://artalk.js.org/guide/frontend/i18n.html): Switch between multiple languages<!--en-->
* [多语言](https://artalk.js.org/guide/frontend/i18n.html): 多国语言切换<!--zh-->
* [Command Line](https://artalk.js.org/guide/backend/config.html): Command line operation management<!--en-->
* [命令行](https://artalk.js.org/guide/backend/config.html): 命令行操作管理能力<!--zh-->
* [API Documentation](https://artalk.js.org/http-api.html): Provides OpenAPI format documentation<!--en-->
* [API 文档](https://artalk.js.org/http-api.html): 提供 OpenAPI 格式文档<!--zh-->
* [Program Upgrade](https://artalk.js.org/guide/backend/update.html): Version check, one-click upgrade<!--en-->
* [程序升级](https://artalk.js.org/guide/backend/update.html): 版本检测，一键升级<!--zh-->
<!-- /features -->

<!-- prettier-ignore-end -->

## Installation<!--en-->
## 安装<!--zh-->

Deploy Artalk Server with Docker in one step:<!--en-->
通过 Docker 一键部署 Artalk 服务器：<!--zh-->

```bash
docker run -d \
    --name artalk \
    -p 8080:23366 \
    -v $(pwd)/data:/data \
    -e "TZ=America/New_York" \<!--en-->
    -e "TZ=Asia/Shanghai" \<!--zh-->
    -e "ATK_LOCALE=en" \<!--en-->
    -e "ATK_LOCALE=zh-CN" \<!--zh-->
    -e "ATK_SITE_DEFAULT=Artalk Blog" \<!--en-->
    -e "ATK_SITE_DEFAULT=Artalk 的博客" \<!--zh-->
    -e "ATK_SITE_URL=https://example.com" \
    artalk/artalk-go
```

Integrate Artalk Client into your webpage:<!--en-->
在你的网页中引入 Artalk 客户端:<!--zh-->

<!-- prettier-ignore-start -->

```ts
Artalk.init({
  el:      '#Comments',
  site:    'Artalk Blog',<!--en-->
  site:    'Artalk 的博客',<!--zh-->
  server:  'https://artalk.example.com',
  pageKey: '/2018/10/02/hello-world.html'
})
```

<!-- prettier-ignore-end -->

We offer various installation methods, including binary files, go install, and package managers for Linux distributions.<!--en-->
我们提供多种安装方法，包括二进制文件、`go install` 和通过 Linux 发行版的包管理器安装。<!--zh-->

[**Learn More →**](https://artalk.js.org/en/guide/deploy.html)<!--en-->
[**了解更多 →**](https://artalk.js.org/zh/guide/deploy.html)<!--zh-->

## For Developers<!--en-->
## 参与开发<!--zh-->

Pull requests are welcome!<!--en-->
我们欢迎你的 Pull Request！<!--zh-->

See [Development](https://artalk.js.org/en/develop/) and [Contributing](./CONTRIBUTING.md) for information on working with the codebase, getting a local development setup, and contributing changes.<!--en-->
有关如何使用代码库、设置本地开发环境和贡献更改的信息，请参阅 [开发文档](https://artalk.js.org/zh/develop/) 和 [贡献指南](./CONTRIBUTING.md)。<!--zh-->

## Contributors<!--en-->
## 贡献者们<!--zh-->

Your contributions enrich the open-source community, fostering learning, inspiration, and innovation. We deeply value your involvement. Thank you for being a vital part of our community! 🥰<!--en-->
你的贡献丰富了开源社区，促进了学习、灵感和创新。我们非常重视你的参与。感谢你成为我们社区的重要一员！🥰<!--zh-->

[![](https://contrib.rocks/image?repo=ArtalkJS/Artalk)](https://github.com/ArtalkJS/Artalk/graphs/contributors)

## Supporters<!--en-->
## 支持者们<!--zh-->

[![Stargazers repo roster for @ArtalkJS/Artalk](https://reporoster.com/stars/ArtalkJS/Artalk)](https://github.com/ArtalkJS/Artalk/stargazers)

## Repobeats Analytics<!--en-->
## Repobeats 分析<!--zh-->

![Alt](https://repobeats.axiom.co/api/embed/a9fc9191ac561bc5a8ee2cddc81e635ecaebafb6.svg 'Repobeats analytics image')

## Stargazers over time<!--en-->
## Star 趋势<!--zh-->

<a href="https://trendshift.io/repositories/6290" target="_blank"><img src="https://trendshift.io/api/badge/repositories/6290" alt="ArtalkJS%2FArtalk | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>

[![Stargazers over time](https://starchart.cc/ArtalkJS/Artalk.svg)](https://starchart.cc/ArtalkJS/Artalk)

## License<!--en-->
## 开源许可协议<!--zh-->

[MIT](./LICENSE)<!--en-->
[MIT](./LICENSE) (麻省理工学院许可证)<!--zh-->

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_shield)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_large)

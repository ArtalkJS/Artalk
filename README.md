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

[English](./README.en.md) • [官方网站](https://artalk.js.org) • [最新版本](https://github.com/ArtalkJS/Artalk/releases) • [更新日志](https://github.com/ArtalkJS/Artalk/blob/master/CHANGELOG.md)

Artalk 是一款简单易用但功能丰富的评论系统，你可以开箱即用地部署并置入任何博客、网站、Web 应用。

- 🍃 前端 ~40KB，纯天然 Vanilla JS
- 🍱 后端 Golang，高效轻量跨平台
- 🐳 通过 Docker 一键部署，方便快捷
- 🌈 开源程序，自托管，隐私至上

## 特性

<!-- prettier-ignore-start -->

<!-- features -->
* [侧边栏](https://artalk.js.org/guide/frontend/sidebar.html): 快速管理、直观浏览
* [社交登录](https://artalk.js.org/guide/frontend/auth.html): 通过社交账号快速登录
* [邮件通知](https://artalk.js.org/guide/backend/email.html): 多种发送方式、邮件模板
* [多元推送](https://artalk.js.org/guide/backend/admin_notify.html): 多种推送方式、通知模版
* [站内通知](https://artalk.js.org/guide/frontend/sidebar.html): 红点标记、提及列表
* [验证码](https://artalk.js.org/guide/backend/captcha.html): 多种验证类型、频率限制
* [评论审核](https://artalk.js.org/guide/backend/moderator.html): 内容检测、垃圾拦截
* [图片上传](https://artalk.js.org/guide/backend/img-upload.html): 自定义上传、支持图床
* [Markdown](https://artalk.js.org/guide/intro.html): 支持 Markdown 语法
* [表情包](https://artalk.js.org/guide/frontend/emoticons.html): 兼容 OwO，快速集成
* [多站点](https://artalk.js.org/guide/backend/multi-site.html): 站点隔离、集中管理
* [管理员](https://artalk.js.org/guide/backend/multi-site.html): 密码验证、徽章标识
* [页面管理](https://artalk.js.org/guide/frontend/sidebar.html): 快速查看、标题一键跳转
* [浏览量统计](https://artalk.js.org/guide/frontend/pv.html): 轻松统计网页浏览量
* [层级结构](https://artalk.js.org/guide/frontend/config.html#nestmax): 嵌套分页列表、滚动加载
* [评论投票](https://artalk.js.org/guide/frontend/config.html#vote): 赞同或反对评论
* [评论排序](https://artalk.js.org/guide/frontend/config.html#listsort): 多种排序方式，自由选择
* [评论搜索](https://artalk.js.org/guide/frontend/sidebar.html): 快速搜索评论内容
* [评论置顶](https://artalk.js.org/guide/frontend/sidebar.html): 重要消息置顶显示
* [仅看作者](https://artalk.js.org/guide/frontend/config.html): 仅显示作者的评论
* [评论跳转](https://artalk.js.org/guide/intro.html): 快速跳转到引用的评论
* [自动保存](https://artalk.js.org/guide/frontend/config.html): 输入内容防丢功能
* [IP 属地](https://artalk.js.org/guide/frontend/ip-region.html): 用户 IP 属地展示
* [数据迁移](https://artalk.js.org/guide/transfer.html): 自由迁移、快速备份
* [图片灯箱](https://artalk.js.org/guide/frontend/lightbox.html): 图片灯箱快速集成
* [图片懒加载](https://artalk.js.org/guide/frontend/img-lazy-load.html): 延迟加载图片，优化体验
* [Latex](https://artalk.js.org/guide/frontend/latex.html): Latex 公式解析集成
* [夜间模式](https://artalk.js.org/guide/frontend/config.html#darkmode): 夜间模式切换
* [扩展插件](https://artalk.js.org/develop/): 创造更多可能性
* [多语言](https://artalk.js.org/guide/frontend/i18n.html): 多国语言切换
* [命令行](https://artalk.js.org/guide/backend/config.html): 命令行操作管理能力
* [API 文档](https://artalk.js.org/develop/): 提供 OpenAPI 格式文档
* [程序升级](https://artalk.js.org/guide/backend/update.html): 版本检测，一键升级
<!-- /features -->

<!-- prettier-ignore-end -->

## 安装

通过 Docker 一键部署：

```bash
docker run -d --name artalk -p 8080:23366 -v $(pwd)/data:/data artalk/artalk-go
```

在网页中引入 Artalk:

<!-- prettier-ignore-start -->

```ts
Artalk.init({
  el:      '#Comments',
  site:    'Artalk 的博客',
  server:  'https://artalk.example.com',
  pageKey: '/2018/10/02/hello-world.html'
})
```

<!-- prettier-ignore-end -->

[**了解更多 →**](https://artalk.js.org/guide/deploy.html)

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

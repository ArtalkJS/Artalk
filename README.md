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

> 🌌 Golang 自托管评论系统

[English](./README.en.md) / [官方文档](https://artalk.js.org) / [最新版本](https://github.com/ArtalkJS/Artalk/releases)

---

- 🍃 轻量 (~30kB)
- 👨‍👧‍👦 安全 (自托管)
- 🐳 易上手 (防脱发)
- 🍱 Golang 后端 (快速 / 跨平台)
- 🌊 Vanilla × TypeScript × Vite (纯天然 / 无依赖)

## 特性

- 侧 边 栏 ：支持多站点集中化管理
- 通知中心：红点的标记 / 提及列表
- 身份验证：徽标自定义 / 密码验证
- 评论审核：反垃圾检测 / 频率限制
- 表情符号：插入表情包 / 快速导入
- 邮件提醒：模版自定义 / 多管理员
- 站点隔离：管理员分配 / 多个站点
- 页面管理：标题可显示 / 快速跳转
- 图片上传：上传到本地 / 多种图床
- 多元推送：支持 钉钉 飞书 TG
- 无限层级：可切换为平铺模式
- 评论投票：赞同还是反对评论
- 评论排序：按热度或时间排序
- 评论置顶：重要消息置顶显示
- 只看作者：仅显示作者的评论
- 说说模式：仅自己可发布评论
- 异步处理：发送评论无需等待
- 滚动加载：评论内容分页处置
- 自动保存：用户输入防丢功能
- 自动填充：用户链接自动填充
- 实时预览：评论内容实时预览
- 夜间模式：防止眼部疾病伤害
- 评论折叠：这个不打算给你康
- 数据备份：防止评论数据丢失
- 数据迁移：快速切换评论系统
- Markdown：默认支持 MD 语法
- 支持 Latex：提供集成 Katex 插件
- 使用 [Vite](https://github.com/vitejs/vite)：属于开发者的极致体验

更多内容请查看：[**介绍文档**](https://artalk.js.org/guide/intro.html)

## 安装

参考：[**部署文档**](https://artalk.js.org/guide/deploy.html)

```bash
npm install artalk
```

```ts
import Artalk from 'artalk'

Artalk.init({
  el:      '#Comments',
  site:    'Artalk 的博客',
  server:  'https://artalk.example.com',
  pageKey: '/2018/10/02/hello-world.html'
})
```

### Docker

```bash
docker run -d \
  --name artalk \
  -p 8080:23366 \
  -v $(pwd)/data:/data \
  --restart=always \
  artalk/artalk-go
```

### Docker Compose

```bash
vim docker-compose.yaml
```

```yaml
version: "3.5"
services:
  artalk:
    container_name: artalk
    image: artalk/artalk-go
    restart: always
    ports:
      - 8080:23366
    volumes:
      - ./data:/data
```

```bash
docker-compose up -d
```

### 在 sealos 上一键部署

[![](https://raw.githubusercontent.com/labring-actions/templates/main/Deploy-on-Sealos.svg)](https://cloud.sealos.io/?openapp=system-template%3FtemplateName%3Dartalk)

## Development

see [CONTRIBUTING.md](./CONTRIBUTING.md)

## Contributors

[![](https://contrib.rocks/image?repo=ArtalkJS/Artalk)](https://github.com/ArtalkJS/Artalk/graphs/contributors)

## Supporters

[![Stargazers repo roster for @ArtalkJS/Artalk](https://reporoster.com/stars/ArtalkJS/Artalk)](https://github.com/ArtalkJS/Artalk/stargazers)

## Repobeats Analytics

![Alt](https://repobeats.axiom.co/api/embed/a9fc9191ac561bc5a8ee2cddc81e635ecaebafb6.svg "Repobeats analytics image")

## 特别致谢

感谢社区提供的帮助与反馈，若有好的建议与意见欢迎前往 [ISSUES](https://github.com/ArtalkJS/Artalk/issues) 随时告知。

## TODOs

- [x] Golang 后端
- [x] 多数据库支持
  - [x] SQLite
  - [x] MySQL
  - [x] Postgres
  - [x] SQLServer
- [x] 多缓存数据库支持
  - [x] In-memory (内建缓存)
  - [x] Redis
  - [x] Memcache
- [x] 多站点支持
- [x] 评论获取分页
- [x] 评论点赞投票
- [x] 浏览量统计
- [x] 评论分页加载
- [x] 评论置顶 / 精华
- [x] 评论排序 (热度 / 时间)
- [x] 通知中心
  - [x] 提及
  - [x] 全部
  - [x] 我的
  - [x] 待审
- [x] 管理员控制台
  - [x] 评论
  - [x] 页面
  - [x] 站点
  - [x] 用户
  - [x] 设置 (GUI)
- [x] 数据导入 ([Artransfer](https://artalk.js.org/guide/transfer.html))
  - [x] Artrans
  - [x] WordPress
  - [x] Typecho ([插件](https://github.com/ArtalkJS/Artrans-Typecho) / [Artransfer-CLI](https://github.com/ArtalkJS/Artransfer-CLI))
  - [x] Valine
  - [x] Disqus
  - [x] Commento
  - [x] Waline ([Artransfer-CLI](https://github.com/ArtalkJS/Artransfer-CLI))
  - [x] Twikoo
  - [x] Artalk v1 (PHP)
- [x] 数据导出
- [x] 邮件多种发送方式
  - [x] SMTP
  - [x] 阿里云邮件
  - [x] 系统调用 sendmail
- [x] 邮件多模板自定义
- [x] 邮件异步队列发送
  - [ ] 队列持久化
- [x] 用户已读标记
- [x] 验证码
  - [x] 图片验证码
  - [x] [Turnstile](https://www.cloudflare.com/products/turnstile/)
  - [x] [ReCAPTCHA](https://www.google.com/recaptcha/about/)
  - [x] [hCaptcha](https://www.hcaptcha.com/)
  - [x] [极验](https://www.geetest.com/)滑动验证码
- [x] 反垃圾
  - [x] [Akismet](https://akismet.com/)
  - [x] [阿里云内容安全](https://help.aliyun.com/document_detail/28417.html)
  - [x] [腾讯云内容安全](https://cloud.tencent.com/document/product/1124/64508)
  - [x] 关键字过滤
- [x] 评论通知管理员 ([notify](https://github.com/nikoksr/notify))
  - [x] Telegram Bot
  - [x] 飞书 WebHook Bot
  - [x] 钉钉
  - [x] Bark
  - [x] Slack
  - [x] LINE
  - [x] 自定义 HTTP 回调
- [ ] 命令行管理
- [ ] 博客邮件订阅
- [x] 用户鉴权机制
- [x] 跨域非法请求阻止
- [x] 全局验证码操作次数限制
- [x] JWT 登录状态验证
- [x] 时区自定义
- [x] 只看作者功能
- [x] 图片上传
- [x] 图片上传到图床 ([upgit](https://github.com/pluveto/upgit))
- [ ] 图片管理
- [ ] 附件上传 / 管理
- [ ] 表情包统一管理
  - [ ] 导入表情包
  - [ ] 表情包图片地址控制
- [ ] AT 提及 (@)
- [ ] 评论话题 (#)
- [ ] 评论标签分类系统
- [ ] 主题样式更换
- [ ] 规范化 API
- [ ] 扩展中心
- [ ] 开放用户注册
- [ ] 第三方登录接入
- [x] 多语言 / 国际化 (i18n)
- [x] 一键升级

## Stargazers over time

[![Stargazers over time](https://starchart.cc/ArtalkJS/Artalk.svg)](https://starchart.cc/ArtalkJS/Artalk)

## License

[MIT](./LICENSE)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_shield)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_large)

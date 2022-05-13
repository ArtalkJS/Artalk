<p align="center">
<img src="https://user-images.githubusercontent.com/22412567/137740516-d9e97af0-fb3b-4dab-b331-671a9a2a3a63.png" alt="Artalk" width="100%">
</p>

# Artalk

[![](https://img.shields.io/npm/v/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![](https://img.shields.io/npm/dt/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![](https://badgen.net/bundlephobia/minzip/artalk?style=flat-square)](https://bundlephobia.com/package/artalk)
[![CircleCI](https://circleci.com/gh/ArtalkJS/Artalk/tree/master.svg?style=svg)](https://circleci.com/gh/ArtalkJS/Artalk/tree/master)

> 🌌 Golang 自托管评论系统

[English](./README.en.md) / [官方文档](https://artalk.js.org) / [最新后端](https://github.com/ArtalkJS/ArtalkGo/releases) / [ArtalkGo](https://github.com/ArtalkJS/ArtalkGo)

---

- 🍃 轻量 (~30kB)
- 👨‍👧‍👦 安全 (自托管)
- 🐳 易上手 (防秃顶)
- 🍱 Golang 后端 (快速 / 跨平台)
- 🌊 TypeScript × Vanilla × Vite (纯天然 / 无依赖)

## 特性

- 侧 边 栏 ：多站点集中化管理
- 通知中心：红点的标记 / 已读记录
- 身份验证：徽标自定义 / 密码验证
- 评论审核：反垃圾检测 / 验证码频率限制
- 表情符号：插入表情包 / 快速导入表情包
- 邮件提醒：模版自定义 / 多管理员通知
- 站点隔离：多站点管理 / 管理员分配
- 页面管理：标题可显示 / 快速跳转
- 图片上传：上传到本地 / 多种图床
- 树洞模式：仅自己可见 / 说说功能
- 多元推送：支持钉钉飞书 TG
- 无限层级：可切换为平铺模式
- 评论投票：赞同还是反对评论
- 评论排序：按热度或时间排序
- 评论置顶：重要消息置顶显示
- 只看作者：仅显示作者的评论
- 异步处理：发送评论无需等待
- 滚动加载：评论内容分页处置
- 自动保存：用户输入防丢功能
- 自动填充：用户链接自动填充
- 实时预览：评论内容实时预览
- 暗黑模式：防止眼部疾病伤害
- 评论折叠：这个不打算给你康
- 数据备份：防止评论数据丢失
- 数据迁移：在不同评论系统之间来回切换
- 一页多评：一页多个评论区（似乎没啥用
- Markdown：语法默认支持
- 支持 Latex：引入 Artalk 的 Katex 插件
- 使用 [Vite](https://github.com/vitejs/vite)：开发者的极致体验

## 食用方针

前往：[“**文档 · 部署**”](https://artalk.js.org/guide/deploy.html)

```sh
$ pnpm add artalk
```

```ts
import Artalk from 'artalk'

new Artalk({
  el:        '#Comments',
  pageKey:   'http://your_domain/post/1', // 页面链接
  pageTitle: '关于如何引入 Artalk 这档子事', // 页面标题
  server:    'http://localhost:8080/api', // 后端地址
  site:      'Artalk 的博客 (你的站点名)',
})
```

### Docker

```sh
# 为 Artalk 创建一个目录
mkdir Artalk
cd Artalk

# 下载配置文件模版
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

## TODOs 

- [x] [Golang 后端](https://github.com/ArtalkJS/ArtalkGo)
- [x] 多数据库支持
  - [x] SQLite
  - [x] MySQL
  - [x] Postgres
  - [x] SQLServer
- [x] 多缓存支持
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
  - [x] 评论 (增/删/改)
  - [x] 页面 (增/删/改)
  - [x] 站点 (增/删/改)
  - [x] 数据分页
  - [ ] 配置 (GUI)
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
- [ ] 数据备份同步
- [x] 邮件异步队列发送
- [ ] 邮件队列持久化
- [x] 邮件多种发送方式
  - [x] SMTP
  - [x] 阿里云邮件
  - [x] 系统调用 sendmail
- [x] 邮件多模板自定义
- [x] 用户已读标记
- [x] 验证码
  - [x] 图片验证码
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
  - [ ] 自定义 HTTP 回调
- [ ] 命令行管理
- [ ] 博客邮件订阅
- [x] 用户鉴权机制
- [x] 跨域非法请求阻止
- [x] 全局验证码操作次数限制
- [x] JWT 登陆状态验证
- [x] 时区自定义
- [x] 只看作者功能
- [ ] AT 提及 (@)
- [ ] 评论提及 (#)
- [ ] 表情包统一管理
  - [ ] 导入表情包
  - [ ] 表情包图片地址控制
- [x] 图片上传
- [x] 图片上传到图床 ([upgit](https://github.com/pluveto/upgit))
- [ ] 图片管理
- [ ] 附件上传 / 管理
- [ ] 评论标签分类系统
- [ ] 主题样式更换
- [ ] 规范化 API
- [ ] 扩展中心
- [ ] 开放用户注册
- [ ] 接入第三方登录
- [ ] 国际化 (i18n)
- [x] 一键升级

## 特别致谢

感谢社区提供的帮助与反馈，若有好的建议与意见欢迎前往 [ISSUES](https://github.com/ArtalkJS/Artalk/issues) 随时告知。

## Stargazers over time

[![Stargazers over time](https://starchart.cc/ArtalkJS/Artalk.svg)](https://starchart.cc/ArtalkJS/Artalk)

## License

[LGPL-3.0](./LICENSE)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_shield)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_large)

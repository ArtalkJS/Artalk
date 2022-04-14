<p align="center">
<img src="https://user-images.githubusercontent.com/22412567/137740516-d9e97af0-fb3b-4dab-b331-671a9a2a3a63.png" alt="Artalk" width="100%">
</p>

# Artalk

[![](https://img.shields.io/npm/v/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![](https://img.shields.io/npm/dt/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![](https://img.shields.io/github/last-commit/ArtalkJS/Artalk/master.svg?style=flat-square)](https://github.com/ArtalkJS/Artalk/commits/master)
[![](https://img.shields.io/github/issues-raw/ArtalkJS/Artalk.svg?style=flat-square)](https://github.com/ArtalkJS/Artalk/issues)
[![](https://img.shields.io/github/issues-pr-raw/ArtalkJS/Artalk.svg?style=flat-square)](https://github.com/ArtalkJS/Artalk/pulls)
[![](https://img.shields.io/travis/com/ArtalkJS/Artalk?style=flat-square)](https://travis-ci.com/ArtalkJS/Artalk)
[![](https://img.shields.io/github/license/ArtalkJS/Artalk.svg?style=flat-square)](./LICENSE)

> 🌌 一款简洁的自托管评论系统 | A Selfhosted Comment System.

前往：[“**官方文档**”](https://artalk.js.org)

---

- 轻量级 (~30kB gzipped)
- 自托管 (数据安全)
- 易上手 (防秃顶)
- Golang 后端 (易部署 / 跨平台)
- TypeScript & Vanilla (纯天然无添加 / 无依赖)

## 特性

- 侧 边 栏 ：所见即所得的管理方式
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
- Markdown 语法默认支持
- 支持 Latex：需引入 Artalk 的 Katex 插件
- [Vite](https://github.com/vitejs/vite)：开发者的极致体验

## 食用方针

前往：[“**文档 · 部署**”](https://artalk.js.org/guide/deploy.html)

```sh
$ pnpm add artalk
```

```ts
import Artalk from 'artalk'

new Artalk({
  el:        '#Comments',
  pageKey:   '<页面链接>',
  pageTitle: '<页面标题>',
  server:    '<后端地址>',
  site:      '<站点名称>',
})
```

## 特别致谢

感谢社区提供的帮助与反馈，若有好的建议与意见欢迎前往 [ISSUES](https://github.com/ArtalkJS/Artalk/issues) 随时告知。

## Stargazers over time

[![Stargazers over time](https://starchart.cc/ArtalkJS/Artalk.svg)](https://starchart.cc/ArtalkJS/Artalk)

## License

[GPL-3.0](./LICENSE)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_shield)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalk?ref=badge_large)

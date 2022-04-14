<p align="center">
<img src="./docs/artalk-go.png" alt="Artalk" width="100%">
</p>

# ArtalkGo

[![GitHub issues](https://img.shields.io/github/issues/ArtalkJS/ArtalkGo)](https://github.com/ArtalkJS/ArtalkGo/issues)
[![](https://img.shields.io/github/issues-pr/ArtalkJS/ArtalkGo)](https://github.com/ArtalkJS/ArtalkGo/pulls)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalkGo.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalkGo?ref=badge_shield)
[![CI to Docker Hub](https://github.com/ArtalkJS/ArtalkGo/actions/workflows/dockerhub.yml/badge.svg)](https://github.com/ArtalkJS/ArtalkGo/actions/workflows/dockerhub.yml)

> ArtalkGo: Golang backend of Artalk.

前往：[“**官方文档 · 后端部分**”](https://artalk.js.org/guide/backend/)

---

- 高效快速
- 异步执行
- 跨平台兼容
- 轻量级部署

## Supports

- 运行环境：支持 Linux, Windows, Darwin (x64 + ARM)
- 数据存储：支持 SQLite, MySQL, PostgreSQL, SQL Server
- 邮件发送：支持 SMTP, 阿里云邮件, 调用 sendmail 发送邮件
- 高效缓存：支持 Redis, Memory cache

## Build

### 编译二进制文件

```sh
$ make all
```

编译后二进制文件将输出到 `bin/` 目录下

### Docker 镜像制作

```sh
# 制作镜像
$ make docker-docker

# 发布镜像
$ make docker-push
```

## TODOs 

- [ ] 命令行管理
- [x] 多站点支持
- [x] 多数据库支持
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
- [x] 数据导入 ([Artransfer](https://github.com/ArtalkJS/Artransfer))
  - [x] Artrans
  - [x] Artalk v1 (PHP)
  - [x] WordPress
  - [x] Typecho
  - [x] Valine
  - [x] Disqus
  - [x] Commento
  - [x] Twikoo
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
- [x] 反垃圾
  - [x] Akismet
  - [ ] 阿里云服务
  - [ ] 腾讯云服务
  - [ ] 关键字 / 正则表达式过滤
- [x] 评论通知管理员 ([notify](https://github.com/nikoksr/notify))
  - [x] Telegram Bot
  - [x] 飞书 WebHook Bot
  - [x] 钉钉
  - [x] Bark
  - [x] Slack
  - [x] LINE
  - [ ] 自定义 HTTP 回调
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

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalkGo.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalkGo?ref=badge_large)

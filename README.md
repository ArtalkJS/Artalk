<p align="center">
<img src="./docs/artalk-go.png" alt="Artalk" width="100%">
</p>

# ArtalkGo

> ArtalkGo: Golang backend of Artalk.

## QuickStart

1. 前往 Release 页下载已编译二进制文件
2. 编辑 `artalk-go.yml` 配置程序
3. 执行 `./artalk-go serve` 运行程序
4. 反代设定的端口到 80 并套上 CDN (Nginx)
5. 持久化运行 artalk-go (tmux, sysctl)

(目前部署较为繁琐，将推出 docker 镜像，更新待续...)

## Features

- 高效
- 跨平台
- 轻松部署
- 多站点支持
- 多数据库类型支持

## TODOs 

- [ ] 命令行管理
- [x] 多站点支持
- [x] 多数据库支持
- [x] 评论获取分页
- [x] 评论点赞投票
- [ ] 实时浏览量统计
- [x] 评论分页加载
- [x] 数据
- [x] 通知中心
  - [x] 提及
  - [x] 全部
  - [x] 我的
  - [x] 待审
- [x] 管理员控制台
  - [x] 评论 (增/删/改)
  - [x] 页面 (增/删/改)
  - [x] 站点 (增/删/改)
  - [ ] 配置 (GUI)
  - [ ] 数据分页
- [x] 数据导入
  - [x] Artalk v1 (PHP)
  - [ ] WordPress
  - [ ] Typecho
- [ ] 数据导出
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
- [ ] 评论通知 WebHook
  - [ ] Telegram Bot
  - [ ] 自定义 HTTP 回调
- [ ] 博客邮件订阅
- [x] 用户鉴权机制
- [x] 跨域非法请求阻止
- [x] 全局验证码操作次数限制
- [x] JWT 登陆状态验证
- [x] 时区自定义
- [ ] 评论置顶 / 精华
- [ ] AT 提及 (@)
- [ ] 评论提及 (#)
- [ ] 表情包统一管理
  - [ ] 导入表情包
  - [ ] 表情包图片地址控制
- [ ] 图片上传 / 管理
- [ ] 附件上传 / 管理
- [ ] 评论标签分类系统
- [ ] 主题样式更换
- [ ] 规范化 API
- [ ] 扩展中心
- [ ] 开放用户注册
- [ ] 接入第三方登录
- [ ] 国际化 (i18n)
- [ ] 在线升级

## Supports

- 跨平台：支持 Linux, Win, Darwin
- 数据存储：支持 SQLite, MySQL, PostgreSQL, SQL Server...
- 高效缓存：支持 Redis, Memory...
- 邮件发送：支持 SMTP, 阿里云邮件, 系统调用 sendmail 等发送邮件

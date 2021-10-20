<p align="center">
<img src="./docs/artalk-go.png" alt="Artalk" width="100%">
</p>

# ArtalkGo

> ArtalkGo: Golang backend of Artalk.

- 高效
- 跨平台
- 轻松部署
- 多站点支持
- 多数据库类型支持

# 部署方针

## 普通方式

1. 前往 [Release](https://github.com/ArtalkJS/ArtalkGo/releases) 下载已编译二进制文件
2. 编辑配置文件
   ```sh
   $ curl -L https://raw.githubusercontent.com/ArtalkJS/ArtalkGo/main/artalk-go.example.yml > artalk-go.yml
   $ vim artalk-go.yml
   ```
3. 运行程序 `./artalk-go serve`
4. 反代设定的端口到 80 并套上 CDN (Nginx, Apache)
5. 持久化运作 artalk-go 程序 (tmux, sysctl)

## Docker（推荐）

```sh
# 为 ArtalkGo 创建一个目录
$ mkdir ArtalkGo
$ cd ArtalkGo

# 下载配置文件模版
$ curl -L https://raw.githubusercontent.com/ArtalkJS/ArtalkGo/main/artalk-go.example.yml > conf.yml

# 编译配置文件
$ vim conf.yml

# 拉取 docker 镜像
$ docker pull artalk/artalk-go

# 新建 docker 容器
$ docker run -d \
   --name artalk-go \
   -p 23366:23366 \
   -v $(pwd)/conf.yml:/conf.yml \
   -v $(pwd)/data:/data \
   artalk/artalk-go
```

- 默认监听 `localhost:23366`
- 配置文件 `./conf.yml`
- 数据目录 `./data/`

# 编译

## 编译二进制文件

```sh
$ make all
```

编译后二进制文件将输出到 `bin/` 目录下

## Docker 镜像制作

```sh
## 制作镜像
$ make docker-docker

# 发布镜像
$ make docker-push
```

## Supports

- 跨平台：支持 Linux, Win, Darwin
- 数据存储：支持 SQLite, MySQL, PostgreSQL, SQL Server...
- 高效缓存：支持 Redis, Memory...
- 邮件发送：支持 SMTP, 阿里云邮件, 系统调用 sendmail 等发送邮件

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

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalkGo.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalkGo?ref=badge_shield)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalkGo.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalkGo?ref=badge_large)
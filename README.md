<p align="center">
<img src="./docs/artalk-go.png" alt="Artalk" width="100%">
</p>

# ArtalkGo

[![CircleCI](https://circleci.com/gh/ArtalkJS/ArtalkGo/tree/master.svg?style=svg)](https://circleci.com/gh/ArtalkJS/ArtalkGo/tree/master)
[![CI to Docker Hub](https://github.com/ArtalkJS/ArtalkGo/actions/workflows/dockerhub.yml/badge.svg)](https://github.com/ArtalkJS/ArtalkGo/actions/workflows/dockerhub.yml) [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalkGo.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalkGo?ref=badge_shield)
[![Docker Pulls](https://img.shields.io/docker/pulls/artalk/artalk-go?style=flat-square)](https://hub.docker.com/r/artalk/artalk-go)

> ArtalkGo: Golang backend of Artalk.

前往：[“**官方文档 · 后端部分**”](https://artalk.js.org/guide/backend/config.html)

---

- 高效快速
- 异步执行
- 跨平台兼容
- 轻量级部署

## Supports

- 运行环境：支持 Linux, Windows, Darwin (x64 + ARM)
- 数据存储：支持 SQLite, MySQL, PostgreSQL, SQL Server
- 邮件发送：支持 SMTP, 阿里云邮件, 调用 sendmail 发送邮件
- 高效缓存：支持 Redis, Memcache, In-Memory (BigCache)

## Build

> 注：绝大多数情况下，你无需进行以下手动编译操作，见：[后端部署文档](https://artalk.js.org/guide/backend/install.html)

### 编译二进制文件

```sh
$ make all
```

编译后二进制文件将输出到 `bin/` 目录下

### Docker Compose 编译运行

```sh
# 克隆项目
$ git clone https://github.com/ArtalkJS/ArtalkGo
$ cd ArtalkGo

# 构建镜像
$ docker compose build

# 运行
$ docker compose up -d
```

### Docker 镜像构建

```sh
# 克隆项目
$ git clone https://github.com/ArtalkJS/ArtalkGo
$ cd ArtalkGo

# 构建镜像
$ make docker-docker

# 发布镜像
$ make docker-push
```

## TODOs

Reference to https://github.com/ArtalkJS/Artalk#todos

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FArtalkJS%2FArtalkGo.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FArtalkJS%2FArtalkGo?ref=badge_large)

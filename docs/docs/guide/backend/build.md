# 后端构建

## 编译运行

```bash
# 拉取代码
git clone https://github.com/ArtalkJS/Artalk.git Artalk

# 编译程序
cd Artalk && make all

# 配置文件
cp conf/artalk.example.yml artalk.yml
vim artalk.yml

# 运行程序
./bin/artalk help
./bin/artalk -c artalk.yml server
```

## 构建二进制文件

```bash
# 拉取代码
git clone https://github.com/ArtalkJS/Artalk.git

# 执行编译
make all
```

编译二进制文件将会输出到 `bin` 目录中

## Docker Compose 编译运行

```bash
# 拉取代码
git clone https://github.com/ArtalkJS/Artalk
cd Artalk

# 构建镜像
docker-compose build

# 运行
docker-compose up -d
```

## Docker 镜像构建

```bash
# 拉取代码
git clone https://github.com/ArtalkJS/Artalk
cd Artalk

# 构建镜像
make docker-docker

# 发布镜像
make docker-push
```

## DevOps

程序由 [GitHub Actions](https://github.com/ArtalkJS/Artalk/actions) 自动化构建

| 版本发布                                                                                                                                                                 | Docker 镜像构建发布                                                                                                                                                       | 二进制构建发布                                                                                                                                                               | 前端构建发布                                                                                                                                                                           |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [![Release](https://github.com/ArtalkJS/Artalk/actions/workflows/release.yml/badge.svg?branch=master)](https://github.com/ArtalkJS/Artalk/actions/workflows/release.yml) | [![Build Docker](https://github.com/ArtalkJS/Artalk/actions/workflows/build-docker.yml/badge.svg)](https://github.com/ArtalkJS/Artalk/actions/workflows/build-docker.yml) | [![Release](https://github.com/ArtalkJS/Artalk/actions/workflows/build-app.yml/badge.svg?branch=master)](https://github.com/ArtalkJS/Artalk/actions/workflows/build-app.yml) | [![Release](https://github.com/ArtalkJS/Artalk/actions/workflows/build-frontend.yml/badge.svg?branch=master)](https://github.com/ArtalkJS/Artalk/actions/workflows/build-frontend.yml) |

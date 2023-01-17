# 后端构建

## 编译运行

```bash
# 拉取代码
git clone https://github.com/ArtalkJS/ArtalkGo.git ArtalkGo

# 编译程序
cd ArtalkGo && make all

# 配置文件
cp artalk-go.example.yml artalk-go.yml
vim artalk-go.yml

# 运行程序
./bin/artalk-go help
./bin/artalk-go -c artalk-go.yml server
```

## 构建二进制文件

```bash
# 拉取代码
git clone https://github.com/ArtalkJS/ArtalkGo.git

# 执行编译
make all
```

编译二进制文件将会输出到 `bin` 目录中

## Docker Compose 编译运行

```bash
# 拉取代码
git clone https://github.com/ArtalkJS/ArtalkGo
cd ArtalkGo

# 构建镜像
docker-compose build

# 运行
docker-compose up -d
```

## Docker 镜像构建

```bash
# 拉取代码
git clone https://github.com/ArtalkJS/ArtalkGo
cd ArtalkGo

# 构建镜像
make docker-docker

# 发布镜像
make docker-push
```

## DevOps

后端构建目前已交给 [GitHub Actions](https://github.com/ArtalkJS/ArtalkGo/actions) 自动完成

|Docker 镜像构建|Release 编译|
|-|-|
|[![CI to Docker Hub](https://github.com/ArtalkJS/ArtalkGo/actions/workflows/dockerhub.yml/badge.svg)](https://github.com/ArtalkJS/ArtalkGo/actions/workflows/dockerhub.yml)|[![Release Build](https://github.com/ArtalkJS/ArtalkGo/actions/workflows/release.yml/badge.svg?branch=master)](https://github.com/ArtalkJS/ArtalkGo/actions/workflows/release.yml)|

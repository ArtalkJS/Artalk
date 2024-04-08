# Docker

Artalk 提供后端程序的 Docker 镜像，以便加速部署流程，提供一个良好的部署体验。

[Docker Hub](https://hub.docker.com/r/artalk/artalk-go) 镜像版本随代码仓库的 [Releases](https://github.com/ArtalkJS/Artalk/releases) 保持同步。

## 镜像拉取

`docker pull artalk/artalk-go`

## 容器创建

:::tip

推荐使用 Docker Compose，可参考 [程序部署](/guide/deploy) 页面的步骤。

:::

Docker 容器创建：

```bash
# 为 Artalk 创建一个目录
mkdir Artalk
cd Artalk

# 拉取 docker 镜像
docker pull artalk/artalk-go

# 生成配置文件
docker run -it -v $(pwd)/data:/data --rm artalk/artalk-go gen config data/artalk.yml

# 编辑配置文件
vim data/artalk.yml

# 运行 docker 容器
docker run -d \
  --name artalk \
  -p 0.0.0.0:8080:23366 \
  -v $(pwd)/data:/data \
  artalk/artalk-go
```

然后，在前端配置填入后端地址：

```js
Artalk.init({ server: 'http://your_domain:8080' })
```

## 重启

修改配置文件后，需要重启才能生效。

```bash
# Docker Compose
docker-compose restart

# Docker
docker restart artalk
```

## 停止

```bash
# Docker Compose
docker-compose stop

# Docker
docker stop artalk
```

## 升级

删除现有容器，拉取最新镜像，然后重新创建容器即可。

### Docker Compose

```bash
docker-compose down
docker-compose pull
docker-compose up -d
```

### Docker

```bash
docker stop artalk
docker rm artalk
docker pull artalk/artalk-go
```

::: tip
突破性变动请注意查看版本更新日志中的提示：[CHANGELOG.md](https://github.com/ArtalkJS/Artalk/blob/master/CHANGELOG.md)
:::

## 拉取历史镜像

镜像会随代码仓库 tags 自动构建发布，您可拉取不同版本号的镜像。

```bash
docker pull artalk/artalk-go@版本号
```

## 进入容器

```bash
# Docker Compose
docker-compose exec artalk bash

# Docker
docker exec -it artalk bash
```

## 多平台兼容性

Docker 镜像暂仅提供 x86、arm64 的镜像构建，如需更多平台架构版本，请下载 [二进制构建部署](/guide/deploy.md#普通方式部署)。

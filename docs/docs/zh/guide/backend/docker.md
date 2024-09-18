# Docker

Artalk 提供后端程序的 Docker 镜像，以便加速部署流程，提供一个良好的部署体验。

[Docker Hub](https://hub.docker.com/r/artalk/artalk-go) 镜像版本随代码仓库的 [Releases](https://github.com/ArtalkJS/Artalk/releases) 保持同步。

## 镜像拉取

```bash
docker pull artalk/artalk-go
```

## 容器创建

:::tip

推荐使用 Docker Compose，可参考 [程序部署](../deploy) 页面的步骤。

:::

## 生成配置文件

```bash
docker run -it -v $(pwd)/data:/data --rm artalk/artalk-go gen config data/artalk.yml
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

## 测试版本

Nightly 镜像为测试版本，每日更新，通过代码仓库最新代码自动构建。

```bash
docker pull artalk/artalk-go:nightly
```

## 历史版本

镜像会随代码仓库 tags 自动构建发布，您可拉取不同版本号的镜像。

```bash
docker pull artalk/artalk-go:版本号
```

## 进入容器

```bash
# Docker Compose
docker-compose exec artalk bash

# Docker
docker exec -it artalk bash
```

## 多平台兼容性

Docker 镜像暂仅提供 x86、arm64 的镜像构建，如需更多平台架构版本，请下载 [二进制构建部署](../deploy.md#二进制文件)。

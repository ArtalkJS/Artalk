# 📦 程序部署

该指南将帮助你在服务器上部署 Artalk。之后，你可以将 Artalk 客户端集成到你的网站或博客中，让用户能够在你的网站上畅所欲言。

## Docker

以下是一个简单的 Artalk **服务器** 和 **客户端** 部署示例。

### 启动服务器

推荐使用 Docker 部署，预先安装 [Docker 引擎](https://docs.docker.com/engine/install/) 并创建一个工作目录，然后执行命令在后台启动容器：

```bash
docker run -d \
    --name artalk \
    -p 8080:23366 \
    -v $(pwd)/data:/data \
    -e "TZ=Asia/Shanghai" \
    -e "ATK_LOCALE=zh-CN" \
    -e "ATK_SITE_DEFAULT=Artalk 的博客" \
    -e "ATK_SITE_URL=https://example.com" \
    artalk/artalk-go
```

（注意：我们也提供了 [Docker Compose](#docker-compose) 的配置文件）。

执行命令创建管理员账户：

```bash
docker exec -it artalk artalk admin
```

浏览器输入 `http://artalk.example.com:8080` 进入 Artalk 后台登录界面。

### 集成客户端

在网页中引入 Artalk 程序内嵌的的前端 JS 和 CSS 资源并初始化 Artalk：

<!-- prettier-ignore-start -->

```html
<!-- CSS -->
<link href="http://artalk.example.com:8080/dist/Artalk.css" rel="stylesheet" />

<!-- JS -->
<script src="http://artalk.example.com:8080/dist/Artalk.js"></script>

<!-- Artalk -->
<div id="Comments"></div>
<script>
Artalk.init({
  el:        '#Comments',                // 绑定元素的 Selector
  pageKey:   '/post/1',                  // 固定链接
  pageTitle: '关于引入 Artalk 的这档子事',  // 页面标题 (留空自动获取)
  server:    'http://artalk.example.com:8080',  // 后端地址
  site:      'Artalk 的博客',             // 你的站点名
})
</script>
```
<!-- prettier-ignore-end -->

评论框输入管理员用户名和邮箱，「控制台」按钮将出现在评论框右下角。

在控制台，你可以根据喜好[配置评论系统](./backend/config.md)、[将评论迁移到 Artalk](./transfer.md)。

🥳 你已成功完成 Artalk 部署！

## 二进制文件

1. [GitHub Release](https://github.com/ArtalkJS/Artalk/releases) 下载程序压缩包
2. 解压 `tar -zxvf artalk_版本号_系统_架构.tar.gz`
3. 运行 `./artalk server`
4. 在你的网页中配置和初始化 Artalk 客户端：

   ```js
   Artalk.init({ server: 'http://artalk.example.com:23366' })
   ```

进阶操作：

- [守护进程 (Systemd, Supervisor)](./backend/daemon.md)
- [反向代理 (Caddy, Nginx, Apache)](./backend/reverse-proxy.md)
- [自编译 (通过最新代码构建)](../develop/contributing.md)

## Go 模块

如果你已经安装了 Golang 工具链，可以运行以下命令来编译和安装最新版本的 Artalk：

```bash
go install github.com/artalkjs/artalk/v2@latest
```

然后运行服务器：

```bash
artalk server
```

客户端集成步骤详见[此处](#集成客户端)。

## Linux 发行版

**Arch Linux (AUR)**：

```bash
paru -S artalk
```

**NixOS**：

```bash
nix-env -iA nixpkgs.artalk
```

**Termux**：

```bash
pkg install artalk
```

[![Packaging status](https://repology.org/badge/vertical-allrepos/artalk.svg)](https://repology.org/project/artalk/versions)

## Docker Compose

创建一个工作目录，并编辑 `docker-compose.yml` 文件：

```yaml
version: '3.8'
services:
  artalk:
    container_name: artalk
    image: artalk/artalk-go
    restart: unless-stopped
    ports:
      - 8080:23366
    volumes:
      - ./data:/data
    environment:
      - TZ=Asia/Shanghai
      - ATK_LOCALE=zh-CN
      - ATK_SITE_DEFAULT=Artalk 的博客
      - ATK_SITE_URL=https://your_domain
```

创建容器运行 Artalk 服务器：

```bash
docker-compose up -d
```

客户端集成步骤详见[此处](#集成客户端)。

::: details Compose 常用命令

```bash
docker-compose restart  # 重启容器
docker-compose stop     # 暂停容器
docker-compose down     # 删除容器
docker-compose pull     # 更新镜像
docker-compose exec artalk bash # 进入容器
```

:::

更多信息：[Docker](./backend/docker.md) / [环境变量](./env.md)

## 前端项目 (Node.js)

通过 NPM 安装 Artalk：

```bash
npm install artalk
```

在你的 Web 项目中引入 Artalk：

```js
import 'artalk/dist/Artalk.css'
import Artalk from 'artalk'

Artalk.init({
  // ...
})
```

更多参考：

- [置入博客文档](../develop/import-blog.md)
- [置入框架文档](../develop/import-framework.md)
- [前端 API](../develop/fe-api.md)
- [前端配置](./frontend/config.md)

## 前端 CDN 资源

::: tip Artalk 最新版本

当前 Artalk 前端最新版本号为： **:ArtalkVersion:**

若需升级前端将 URL 中的版本号数字部分替换即可。
:::

Artalk 后端程序内嵌了前端 JS、CSS 文件，使用公共 CDN 资源请注意版本兼容性。

**CDNJS**

> <https://cdnjs.cloudflare.com/ajax/libs/artalk/:ArtalkVersion:/Artalk.js>
>
> <https://cdnjs.cloudflare.com/ajax/libs/artalk/:ArtalkVersion:/Artalk.css>

::: details 查看更多
**SUSTech Mirrors (境内)**

> <https://mirrors.sustech.edu.cn/cdnjs/ajax/libs/artalk/:ArtalkVersion:/Artalk.js>
>
> <https://mirrors.sustech.edu.cn/cdnjs/ajax/libs/artalk/:ArtalkVersion:/Artalk.css>

**UNPKG**

> <https://unpkg.com/artalk@:ArtalkVersion:/dist/Artalk.js>
>
> <https://unpkg.com/artalk@:ArtalkVersion:/dist/Artalk.css>

**JS DELIVR**

> <https://cdn.jsdelivr.net/npm/artalk@:ArtalkVersion:/dist/Artalk.js>
>
> <https://cdn.jsdelivr.net/npm/artalk@:ArtalkVersion:/dist/Artalk.css>

:::

## 数据导入

从其他评论系统导入数据：[数据迁移](./transfer.md)。

## ArtalkLite

ArtalkLite 是一个轻量级的精简 Artalk 客户端，体积更小、更简约。查看：[ArtalkLite](./frontend/artalk-lite.md)。

## 开发环境

请参考：[开发者指南](https://github.com/ArtalkJS/Artalk/blob/master/CONTRIBUTING.md)。
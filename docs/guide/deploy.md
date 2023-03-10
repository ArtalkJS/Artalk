# 📦 程序部署

## Docker 部署

推荐使用 Docker 部署，需预先安装 [Docker 引擎](https://docs.docker.com/engine/install/)，服务器执行命令创建容器：

```bash
docker run -d \
    --name artalk \
    -p 8080:23366 \
    -v $(pwd)/data:/data \
    --restart=always \
    artalk/artalk-go
```

> 假设域名 `http://your_domain` 已正确添加 DNS 记录并指向你的服务器 IP

浏览器打开 `http://your_domain:8080` 将出现 Artalk 后台登陆界面。

执行命令创建管理员账户：

```bash
docker exec -it artalk artalk admin
```

在你的网站引入 Artalk 程序内嵌的的前端 CSS、JS 资源并初始化：

> 注：将 `http://your_domain:8080` 改为你的服务器域名，或使用 [公共 CDN 资源](#cdn-资源)。

```html
<!-- CSS -->
<link href="http://your_domain:8080/dist/Artalk.css" rel="stylesheet">

<!-- JS -->
<script src="http://your_domain:8080/dist/Artalk.js"></script>

<!-- Artalk -->
<div id="Comments"></div>
<script>
Artalk.init({
  el:        '#Comments',                // 绑定元素的 Selector
  pageKey:   '/post/1',                  // 固定链接 (留空自动获取)
  pageTitle: '关于引入 Artalk 的这档子事',  // 页面标题 (留空自动获取)
  server:    'http://your_domain:8080',  // 后端地址
  site:      'Artalk 的博客',             // 你的站点名
})
</script>
```

在评论框输入管理员的用户名和邮箱，控制台入口按钮将出现在评论框右下角位置。

在控制台，你可以根据喜好配置评论系统、[将评论迁移到 Artalk](./transfer.md)。

祝贺！你已成功完成 Artalk 部署 🥳

## 普通方式部署

1. 前往 [GitHub Release](https://github.com/ArtalkJS/Artalk/releases) 下载程序压缩包

  > 注：我们已将 ArtalkGo 与主仓库合并，历史旧版可在[此页面](https://github.com/ArtalkJS/ArtalkGo/releases)查看，日后新版将在[新页面](https://github.com/ArtalkJS/Artalk/releases)发布。

2. 提取压缩包：`tar -zxvf artalk_版本号_系统_架构.tar.gz`
3. 运行程序 `./artalk server`
4. 前端配置

    ```js
    Artalk.init({ server: "http://your_domain:23366" })
    ```

**其它可选操作：**

- [“反向代理端口到 80 / 443 (Nginx, Apache)”](./backend/reverse-proxy.md)
- ["持久化运作 (tmux, systemd, supervisor)"](./backend/daemon.md)

**附表：文件名释义表**

|文件名|操作系统|CPU 架构|
|:-|:-:|:-:|
|artalk_linux_amd64.tar.gz|Linux|x86_64|
|artalk_linux_arm64.tar.gz|Linux|ARM64|
|artalk_linux_arm7.tar.gz|Linux|ARMv7|
|artalk_windows_amd64.zip|Windows|x86_64|
|artalk_darwin_arm64.tar.gz|macOS|Apple Silicon|
|artalk_darwin_amd64.tar.gz|macOS|Intel Chip|

## Docker Compose 部署

提供 docker-compose.yaml 文件可供参考：

```yaml
version: "3.5"
services:
  artalk:
    container_name: artalk
    image: artalk/artalk-go
    restart: always
    ports:
      - 8080:23366
    volumes:
      - ./data:/data
```

在与配置文件相同的目录执行命令创建容器：

```bash
docker-compose up -d
```

::: details 一些 Docker Compose 常用命令

```bash
docker-compose restart  # 重启容器
docker-compose stop     # 暂停容器
docker-compose down     # 删除容器
docker-compose pull     # 更新镜像
docker-compose exec artalk bash # 进入容器
```

:::

> 详细可见：[“后端 · Docker”](./backend/docker.md)

## 自行编译并运行

可参考：[“后端构建”](./backend/build.md)

## CDN 资源

::: tip Artalk 最新版本

当前 Artalk 前端最新版本号为： **:ArtalkVersion:**

若需升级前端，请将 URL 中的版本号数字部分替换即可。
:::

Artalk 后端程序内嵌了前端 JS、CSS 文件，使用公共 CDN 资源请注意前后端版本的兼容性。

Artalk 静态资源通过上游 [CDNJS](https://cdnjs.com/) 分发，国内有许多镜像可供选择：

**BootCDN (国内)**

> https://cdn.bootcdn.net/ajax/libs/artalk/:ArtalkVersion:/Artalk.js
>
> https://cdn.bootcdn.net/ajax/libs/artalk/:ArtalkVersion:/Artalk.css


**ElemeCDN (国内)**

> https://npm.elemecdn.com/artalk@:ArtalkVersion:/dist/Artalk.js
>
> https://npm.elemecdn.com/artalk@:ArtalkVersion:/dist/Artalk.css

**CDNJS**

> https://cdnjs.cloudflare.com/ajax/libs/artalk/:ArtalkVersion:/Artalk.js
>
> https://cdnjs.cloudflare.com/ajax/libs/artalk/:ArtalkVersion:/Artalk.css

**UNPKG**

> https://unpkg.com/artalk@:ArtalkVersion:/dist/Artalk.js
> 
> https://unpkg.com/artalk@:ArtalkVersion:/dist/Artalk.css

**JS DELIVR**

> https://cdn.jsdelivr.net/npm/artalk@:ArtalkVersion:/dist/Artalk.js
> 
> https://cdn.jsdelivr.net/npm/artalk@:ArtalkVersion:/dist/Artalk.css

## ArtalkLite

可选择精简版 [ArtalkLite](./frontend/artalk-lite.md)：体积更小、更简约。

## Node 环境

```bash
pnpm add artalk
```

引入到你的项目：

```js
import 'artalk/dist/Artalk.css'
import Artalk from 'artalk'

Artalk.init({
  // ...
})
```

## 何时引入、何时 init？

- 可以在任意位置引入 JS 和 CSS 资源，但需确保 JS 引入在执行 `Artalk.init({})` 前。
- 执行 `Artalk.init({ el: '#x' })` 时，需要确保 `<div id="x"></div>` 存在于页面当中。

可参考：[“前端框架引入”](./frontend/import-framework.md) / [“博客引入”](./frontend/import-blog.md)

## 数据导入

从其他评论系统导入数据：[“数据迁移”](./transfer.md)

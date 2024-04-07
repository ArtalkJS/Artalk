# 图片上传

Artalk 提供图片上传功能，支持限制图片大小、上传频率等，你还能结合 UpGit 将图片上传到图床。

你可以在[控制中心](/guide/frontend/sidebar.md#控制中心)找到「设置」界面修改此配置。

## 配置文件

完整的 `img-upload` 配置如下：

```yaml
# 图片上传
img_upload:
  enabled: true              # 总开关
  path: "./data/artalk-img/" # 图片存放路径
  max_size: 5                # 图片大小限制 (单位：MB)
  public_path: null          # 指定图片链接基础路径 (默认为 "/static/images/")
  # 使用 upgit 将图片上传到 GitHub 或图床
  upgit:
    enabled: false  # 启用 upgit
    exec: "upgit -c <upgit配置文件路径> -t /artalk-img"
    del_local: true # 上传后删除本地的图片
```

## 使用 Upgit 上传到图床

[UpGit](https://github.com/pluveto/upgit) 支持将图片上传到 Github、Gitee、腾讯云 COS、七牛云、又拍云、SM.MS 等图床或代码仓库。

首先，根据 [README.md](https://github.com/pluveto/upgit) 的说明，下载 UpGit 并完成你需要上传的目标图床的配置。

然后，将 UpGit 加入系统的环境变量中，在 `~/.bashrc` 加入：

```bash
export PATH=$PATH:/path/to/upgit
```

（或者直接移入 `/usr/bin`）

最后，在 Artalk 的 `img_upload.upgit` 字段填入 UpGit 启动参数：

```yaml
  upgit:
    enabled: true  # 启用 upgit
    exec: "upgit -c <upgit配置文件路径> -t /artalk-img"
    del_local: true # 上传后删除本地的图片
```

::: warning 更新注意
从 `v2.8.4` 版本开始，为了增强安全性，Artalk 不再允许指定 UpGit 的可执行文件路径，请将其加入系统的环境变量中。:)
:::

### Docker 挂载 UpGit

如果你使用 Docker 部署 Artalk，可以将 UpGit 可执行文件挂载到容器中：

```bash
docker run -d --name artalk -v /path/to/upgit:/usr/bin/upgit -v /path/to/artalk:/app/data -p 8080:23366 artalk
```

## 上传频率限制

频率限制跟随 `captcha` 验证码配置，当超出限制将弹出验证码。

可参考：[后端 · 验证码](/guide/backend/captcha.md)

## path

`img_upload.path` 为上传的图片文件「本地存放目录」路径，该目录会被 Artalk 映射到可访问的：

```
http://<后端地址>/static/images/
```

## public_path

`img.public_path` 为空的默认值为：`/static/images/`

当该项为「相对路径」时，例如：`/static/images/` 前端上传图片得到的 HTML 标签将为：

```html
<img src="http://<后端地址>/static/images/1.png">
```

注：这里的 `<后端地址>` 是前端 `conf.server` 配置。

当该项为「完整 URL 路径」时，例如：`https://cdn.github.com/img/` 时，图片标签将为：

```html
<img src="https://cdn.github.com/img/1.png">
```

提示：这个配置可以结合负载均衡等场景使用。

## 在前端自定义上传 API

前端提供了配置项 `imgUploader`，你可以自定义前端图片上传时请求的 API，例如：

```js
Artalk.init({
  imgUploader: async (file) => {
    const form = new FormData()
    form.set('file', file)

    const imgUrl = await fetch("https://api.example.org/upload", {
      method: 'POST',
      body: form
    })

    return imgUrl
  }
})
```

参考：[前端配置文档](../frontend/config.md#imguploader)

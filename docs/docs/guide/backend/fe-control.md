# 在后端控制前端

你可以在后端完全控制前端的行为，我们推荐使用这种方式部署 Artalk。

## 在前端引入内置的资源

后端 ArtalkGo 服务器中内置了可供前端引入的 JS 和 CSS 资源文件：

```html
<!-- CSS -->
<link href="https://<artalk_go_server>/dist/Artalk.css" rel="stylesheet">

<!-- JS -->
<script src="https://<artalk_go_server>/dist/Artalk.js"></script>
```

> 提示：将 `<artalk_go_server>` 替换为你的 ArtalkGo 服务器地址。

这样如果升级后端 ArtalkGo 程序，前端无需更换新版 ArtalkJS 的引入地址，来使之与后端程序兼容。

注：内置的前端 JS 和 CSS 始终与后端版本兼容，但不保证是最新的版本。

## 在后端控制前端的配置

你能够在后端控制 [前端的配置](/guide/frontend/config)。

这个功能处于「默认关闭」状态，首先你需要在前端启用它：

```diff
new Artalk({
+  useBackendConf: true,
})
```

然后，在后端 ArtalkGo 的配置文件中添加 `frontend` 字段内容，例如 `artalk-go.yml` 添加：

```yaml
frontend:
  placeholder: "键入内容..."
  noComment: "「此时无声胜有声」"
  sendBtn: "发送评论"
  emoticons: "https://raw.githubusercontent.com/ArtalkJS/Emoticons/master/grps/default.json"
  # ----- 此处省略 -------
  # 与前端配置项名称保持一致
```

这样就无需在前端改动配置，前端的配置始终跟随后端。

一份完整的后端 `frontend` 字段配置文件可供参考：[artalk-go.frontend.example.yml](https://github.com/ArtalkJS/ArtalkGo/blob/master/artalk-go.frontend.example.yml)

::: tip

如果你的表情包配置项 [emoticons](/guide/frontend/emoticons) 需传递 Object 而非 URL，可以将其转为 JSON 字符串，例如：

```yaml
frontend:
  emoticons: '{"表情": { "test": "tttt..." }}'
```

:::

# 在后端控制前端

你可以在后端完全控制前端界面的行为（包括设定几乎所有[前端配置项](/guide/frontend/config)），我们推荐使用这种方式配置 Artalk 界面，并且这个特性「默认开启」。

## 如何修改前端界面配置

在侧边栏的[控制中心](/guide/frontend/sidebar.md#控制中心)提供了图形设置界面，你能快速轻松地对前端界面配置进行修改。

<img src="../../images/sidebar/fe-setting.png" width="400px">

## 在前端引入内置资源

后端 Artalk 服务器中内置了可供前端引入的 JS 和 CSS 资源文件：

```html
<!-- CSS -->
<link href="https://<artalk_go_server>/dist/Artalk.css" rel="stylesheet" />

<!-- JS -->
<script src="https://<artalk_go_server>/dist/Artalk.js"></script>
```

> 提示：将 `<artalk_go_server>` 替换为你的 Artalk 服务器地址。

这样能让前后端始终保持兼容性，而无需在程序升级后手动更换 Artalk 前端资源的引入地址。

## 进阶内容

### 手动编辑配置文件

Artalk 提供图形界面简化配置，一般情况无需手动修改。

在 Artalk 的配置文件 `artalk.yml` 中配置 `frontend` 字段来控制前端界面，例如：

```yaml
frontend:
  placeholder: 键入内容...
  noComment: 「此时无声胜有声」
  sendBtn: 发送评论
  emoticons: 'https://raw.githubusercontent.com/ArtalkJS/Emoticons/master/grps/default.json'
  # ----- 此处省略 -------
  # 与前端配置项名称保持一致
```

一份完整的后端 `frontend` 字段配置文件可供参考：[artalk.example.zh-CN.yml](https://github.com/ArtalkJS/Artalk/blob/master/conf/artalk.example.zh-CN.yml)

### 关闭 “后端控制前端” 功能

后端控制前端功能默认开启，我们不建议关闭该功能。

在前端 `Artalk.init({ ... })` 时定义的配置会被配置文件中的 `frontend` 字段配置所覆盖，但如果有需要也可以在前端禁用这个特性：

```diff
Artalk.init({
+  useBackendConf: false,
})
```

### 表情包配置

建议 `frontend.emoticons` 填入表情包的 URL，而非 Object。

参考文档：[表情包](/guide/frontend/emoticons)

如需传递 Object，可以将其转为 JSON 字符串格式，例如：

```yaml
frontend:
  emoticons: '{"表情": { "test": "tttt..." }}'
```

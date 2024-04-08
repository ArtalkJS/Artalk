# 前端配置

```js
const artalk = Artalk.init({ ... })

artalk.update({ ... })
```

- 默认配置：[defaults.ts](https://github.com/ArtalkJS/Artalk/blob/master/ui/artalk/src/defaults.ts)
- 声明文件：[config.ts](https://github.com/ArtalkJS/Artalk/blob/master/ui/artalk/src/types/config.ts)

## 轻松配置

推荐在侧边栏 “[控制中心](/guide/frontend/sidebar.md)” 通过图形界面修改前端的配置，而无需在代码中对界面进行设定。

注：前端的设定可能会被后端覆盖，更多内容参考：[在后端控制前端](/guide/backend/fe-control.md)

## 基本配置（必填项）

### el

**装载元素**（填入需要绑定的元素 Selector）

- 类型：`String|HTMLElement`
- 默认值：`undefined`

> 例如：`#Comments` 对应元素 `<div id="Comments"></div>`

### pageKey

**页面地址**（相对路径 / 完整 URL）

- 类型：`String`
- 默认值：`location.pathname`

可留空自动获取页面的相对路径。

可以填写由博客系统生成的 `固定链接`，但建议使用相对路径以便日后切换域名。

参考：[关于使用相对 / 绝对路径](/guide/backend/relative-path.md)

### pageTitle

**页面标题**（用于管理列表显示，邮件通知等）

- 类型：`String`
- 默认值：`document.title`

可留空自动获取页面的 `<title>` 标签值。

### server

**后端程序 API 地址**

- 类型：`String`
- 默认值：`undefined`

部署后端，确保后端地址前端可访问

> 例如：http://yourdomain.xxx

::: warning 更新注意

从 v2.2.6 版本开始，`server` 无需在结尾带上 `/api/` 路径。

:::

### site

**站点名称**

- 类型：`String`
- 默认值：`undefined`

可留空使用后端配置的 “默认站点”，

Artalk 支持多站点统一管理，此项用于站点隔离。

### useBackendConf

**跟随后端的配置**

- 类型：`Boolean`
- 默认值：`true`（默认启用）

可以在后端的配置文件中定义前端的配置，让前端配置始终跟随后端。

详情可参考：[在后端控制前端](/guide/backend/fe-control)

## 国际化 (i18n)

### locale

**语言**

- 类型：`String|Object|"auto"`
- 默认值：`"zh-CN"`

遵循 Unicode BCP 47 规范，该项默认为 "zh-CN" (简体中文)。

欢迎提交 PR 帮助翻译 Artalk 的多语言，为社区贡献一份力量！

详情参考：[多语言](./i18n.md)

## 请求

### reqTimeout

**请求超时**

- 类型：`Number`
- 默认值：`15000`

当请求时间大于该值，自动断开请求并报错。（单位：毫秒）

## 表情包

### emoticons

**表情包**

- 类型：`Object|Array|String|Boolean`
- 默认值："[https://cdn.jsdelivr.net/gh/ArtalkJS/Emoticons/grps/default.json](https://cdn.jsdelivr.net/gh/ArtalkJS/Emoticons/grps/default.json)"

详细内容：[前端 · 表情包](/guide/frontend/emoticons.md)

更新兼容 [OwO 格式](https://github.com/DIYgod/OwO)，支持 URL 动态加载。

设置为 `false` 关闭表情包功能。

:::warning 请替换 CDN 资源
JS DELIVR 在中国大陆的 [ICP 牌照已被吊销](https://github.com/jsdelivr/jsdelivr/issues/18348#issuecomment-997777996)。
:::

## 界面

### placeholder

**评论框占位字符**

- 类型：`String`
- 默认值：`"键入内容..."`

### noComment

**评论为空时显示字符**

- 类型：`String`
- 默认值：`"「此时无声胜有声」"`

### sendBtn

**发送按钮文字**

- 类型：`String`
- 默认值：`"发送评论"`

### editorTravel

**评论框穿梭**

- 类型：`Boolean`
- 默认值：`true`

设置为 `true` 当回复评论时，评论框移动到待回复评论位置之后，而不是固定不动。

### darkMode

**夜间模式**

- 类型：`Boolean|"auto"`
- 默认值：`false`

当 Artalk 被 new 时会读取该值，并根据该值选择是否开启夜间模式（可与博客主题配合使用）。

代码动态修改 darkMode：

```js
artalk.setDarkMode(true)
```

> 参考代码：“[index.html](https://github.com/ArtalkJS/Artalk/blob/master/ui/artalk/index.html#L97-L150)”

可设置为 `"auto"`，Artalk 将监听 `(prefers-color-scheme: dark)` 根据用户操作系统判断自动切换夜间模式。

### flatMode

**平铺模式**

- 类型：`Boolean|"auto"`
- 默认值：`"auto"`

默认 `"auto"` 仅小尺寸屏幕设备自动开启「平铺」模式 (屏幕宽度 < 768px 时)

设置 `true` 评论以「平铺模式」形式显示

设置 `false` 评论以「层级嵌套」形式显示

### nestMax

**最大嵌套层数**

- 类型：`Number`
- 默认值：`2`

评论「层级嵌套」模式的最大嵌套层数。

### nestSort

**嵌套评论的排序规则**

- 类型：`"DATE_ASC"|"DATE_DESC"|"VOTE_UP_DESC"`
- 默认值：`"DATE_ASC"`

嵌套评论的子评论默认以「日期升序 (新评的论在末尾)」排列。

## 功能

### pvEl

**页面浏览量 (PV) 绑定元素**

- 类型：`String`
- 默认值：`"#ArtalkPV"`

你可以在页面任意位置，放置 HTML 标签：`<span id="ArtalkPV"></span>`

当 Artalk 完成加载时展示页面的浏览量。

该项填入绑定元素的 Selector，默认为 `#ArtalkPV`。

### countEl

**评论数绑定元素**

- 类型：`String`
- 默认值：`"#ArtalkCount"`

你可以在页面任意位置，放置 HTML 标签：`<span id="ArtalkCount"></span>` 显示当前页面的评论数。

::: tip

pvEl 和 countEl 元素标签都可以设置 `data-page-key` 属性值，来指定显示某个页面的统计数目，例如：`<span id="ArtalkCount" data-page-key="/t/1.html"></span>`

详情参考：[浏览量统计](./pv.md#显示多个页面的浏览量)

:::

### vote

**投票按钮**

- 类型：`Boolean`
- 默认值：`true`

启用评论投票功能 (赞同 / 反对)。

### voteDown

**反对按钮**

- 类型：`Boolean`
- 默认值：`false`

反对的投票按钮（默认隐藏）。

### uaBadge

**显示用户的 UserAgent 信息徽标**

- 类型：`Boolean`
- 默认值：`false`

### listSort

**评论排序功能**

- 类型：`Boolean`
- 默认值：`true`

鼠标移到评论列表左上角「n 条评论」的位置，显示悬浮下拉层，可以让评论列表按照「最新 / 最热 / 最早 / 作者」等规则排序筛选显示。

### imgUpload

**图片上传功能**

- 类型：`Boolean`
- 默认值：`true`

该配置项自动跟随后端，当后端图片上传功能关闭时，仅管理员会显示图片上传按钮。

### imgUploader

**图片上传器**

- 类型：`(file: File) => Promise<string>`
- 默认值：`undefined`

自定义图片上传器，例如：

```js
Artalk.init({
  imgUploader: async (file) => {
    const form = new FormData()
    form.set('file', file)

    const imgUrl = await fetch('https://api.example.org/upload', {
      method: 'POST',
      body: form,
    })

    return imgUrl
  },
})
```

### preview

**编辑器实时预览功能**

- 类型：`Boolean`
- 默认值：`true`

显示编辑器的「预览」按钮。

## 头像

```js
gravatar: {
  mirror: '<Gravatar 镜像地址>',
  default: 'mp',
}
```

### gravatar.mirror

**Gravatar 镜像地址**

- 类型：`String`
- 默认值：`"https://sdn.geekzu.org/avatar/"`

如果你觉得 Gravatar 头像加载速度不理想，可以尝试替换。

例如：

> Cravatar：https://cravatar.cn/avatar/
>
> V2EX：https://cdn.v2ex.com/gravatar/
>
> 极客族：https://sdn.geekzu.org/avatar/
>
> loli：https://gravatar.loli.net/avatar/

### gravatar.params

**Gravatar API 参数**

- 类型：`String`
- 默认值：`"d=mp&s=240"`

例如，你可以通过该配置项设置默认头像 (`d=mp`) 和头像尺寸 (`s=240`)。

参考：[Gravatar API 文档](http://cn.gravatar.org/site/implement/images/)

该配置项格式为 HTTP Query。

::: warning 更新注意

v2.5.5 已废弃 `gravatar.default` 配置项，请使用 `gravatar.params` 替代。

:::

### avatarURLBuilder

**头像链接生成器**

- 类型：`(comment: CommentData) => string`
- 默认值：`undefined`

自定义用户头像图片链接生成，例如：

```js
Artalk.init({
  avatarURLBuilder: (comment) => {
    return `/api/avatar?email=${comment.email_encrypted}`
  },
})
```

## 评论分页

```js
pagination: {
  pageSize: 20,   // 每页评论数
  readMore: true, // 加载更多 or 分页条
  autoLoad: true, // 自动加载 (加载更多)
}
```

### pagination.readMore

**加载更多模式**

- 类型：`Boolean`
- 默认值：`true`

设置 `true` 为 “**加载更多**” 模式

设置 `false` 为 “**分页条**” 模式

### pagination.autoLoad

**滚动到底部自动加载**

- 类型：`Boolean`
- 默认值：`true`

（需同时开启“加载更多”模式，将 `readMore` 设置为 `true`）

### pagination.pageSize

**每次请求获取数量**

- 类型：`Number`
- 默认值：`20`

## 内容限高

超过设定高度的内容将被隐藏，并显示“阅读更多”按钮。

```js
heightLimit: {
  content: 300, // 评论内容限高
  children: 400, // 子评论区域限高
  scrollable: false, // 限高滚动
}
```

### heightLimit.content

**评论内容限高**

- 类型：`Number`
- 默认值：`300`

> 当值为 0 时，关闭限高

### heightLimit.children

**子评论区域限高**

- 类型：`Number`
- 默认值：`400`

### heightLimit.scrollable

**限高区域滚动**

- 类型：`Boolean`
- 默认：`false`

允许限高区域滚动。

## 版本检测

### versionCheck

**版本检测**

- 类型：`Boolean`
- 默认：`true`

当前端和后端版本不兼容时，显示警告提示框。

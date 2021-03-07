<p align="center">
<img src="./docs/banner-2021.png" alt="Artalk" width="100%">
</p>

# [Artalk](https://artalk.js.org)
> 一款简洁有趣的自托管评论系统

[![](https://img.shields.io/npm/v/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![](https://img.shields.io/npm/dt/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![](https://img.shields.io/github/last-commit/ArtalkJS/Artalk/master.svg?style=flat-square)](https://github.com/ArtalkJS/Artalk/commits/master)
[![](https://img.shields.io/github/issues-raw/ArtalkJS/Artalk.svg?style=flat-square)](https://github.com/ArtalkJS/Artalk/issues)
[![](https://img.shields.io/github/issues-pr-raw/ArtalkJS/Artalk.svg?style=flat-square)](https://github.com/ArtalkJS/Artalk/pulls)
[![](https://img.shields.io/travis/com/ArtalkJS/Artalk?style=flat-square)](https://travis-ci.com/ArtalkJS/Artalk)
[![](https://img.shields.io/github/license/ArtalkJS/Artalk.svg?style=flat-square)](./LICENSE)

## 特性
- 轻量简洁 (~23kB gzipped)
- 有趣有爱
- 自托管
- Markdown
- 表情自定
- 滑稽表情包
- 管理员密码，防冒名
- 验证码，提交频率限制
- 通知中心，邮件提醒
- 仅管理员可评论
- 无限层级回复
- 滚动加载更多
- 内容自动保存
- 暗黑模式
- 评论折叠
- 一页多个评论
- TypeScript

[查看 DEMO](https://artalk.js.org)

## 基本食用方法

> 前端资源下载：[Artalk.js](./dist/Artalk.js) | [Artalk.css](./dist/Artalk.css)

1. 部署 Artalk 的后端，传送门：[PHP API](https://github.com/ArtalkJS/Artalk-API-PHP) / Go API
2. 前端页面配置 Artalk 如下：

```html
<!DOCTYPE html>
<html>
<head>
  <!-- ... -->
  <link href="dist/Artalk.css" rel="stylesheet">
</head>
<body>
  <div id="ArtalkComments"></div>
  <!-- ... -->
  <script src="dist/Artalk.js"></script>
  <script>
  new Artalk({
    el: '#ArtalkComments', // 元素选择
    placeholder: '来啊，快活啊 ( ゜- ゜)', // 占位符
    noComment: '快来成为第一个评论的人吧~', // 无评论时显示
    pageKey: '[页面唯一标识]',
    serverUrl: '[后端程序URL]',
    readMore: { // 阅读更多配置
      pageSize: 15, // 每次请求获取评论数
      autoLoad: true // 滚动到底部自动加载
    }
  });
  </script>
</body>
</html>
```

前端更多 QuickStart 栗子，请参考 [/example/](./example/) 目录

## 一些进阶的操作

<details>

<summary>点我给你看</summary>

### 自定义头像 Gravatar 镜像源

Artalk 依赖于 [Gravatar](https://gravatar.com) 服务，但 Gravatar 在部分地区可能会出现连接问题。

可通过以下配置解决：

```js
new Artalk({
  gravatar: {
    cdn: 'https://cdn.v2ex.com/gravatar/'
  }

  // ... 你的其他配置
})
```

目前默认使用 v2ex 的镜像，如果您有更好的，欢迎反馈~~

### 默认头像

Gravatar 默认头像，参考：[传送门](https://cn.gravatar.com/site/implement/images/#default-image)

```js
new Artalk({
  defaultAvatar: 'mp',
  
  // ... 你的其他配置
})
```

### 开启暗黑模式

以下给出简单的栗子，可以结合博客主题的暗黑模式食用：

```html
<button onclick="switchDarkMode()">暗黑模式按钮</button>

<script>
var artalk = new Artalk({ // ① 暴露 artalk 变量以供调用
  // ... 各种配置

  darkMode: false,
  // ↑ ② 当 Artalk 初始化时，是否立刻开启暗黑模式
  // ↑ 这里可以直接读取当前主题的模式
})

// ③ 动态设置 Artalk 的暗黑模式
let isDarkMode = false // 读取当前你主题的模式
artalk.ui.setDarkMode(darkMode)

// ④ 你主题 暗黑模式切换按钮 点击时的触发操作
function switchDarkMode() {
  let isDarkMode = true // ...
  artalk.ui.setDarkMode(darkMode)
}
</script>
```

独立开发新的暗黑模式可以参考：[我的栗子](https://github.com/ArtalkJS/Artalk/blob/master/index.html#L88) | [如何和操作系统的暗黑模式同步？](https://stackoverflow.com/questions/50840168/how-to-detect-if-the-os-is-in-dark-mode-in-browsers)

### 自定义表情包

表情包配置格式参考：[emoticons.json](/src/assets/emoticons.json)

```js
// ↓ 首先将表情包数据加载并存到变量中
let eData = {
  // ...
}

new Artalk({
  emoticons: eData

  // ... 你的其他配置
})
```

前端更多配置项详见 [/types/artalk-config.d.ts](./types/artalk-config.d.ts)

</details>

## 开发

<details>

<summary>点我给你看</summary>

```bash
git clone https://github.com/ArtalkJS/Artalk.git
cd Artalk
yarn install

# Dev
yarn run dev

# Build
yarn run build
```

Made with ♥

</details>

## License
[GPL-2.0](./LICENSE)

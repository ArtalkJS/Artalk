# 表情包

编辑后端 `artalk.yml` 配置文件并修改配置项 `frontend.emoticons` 填入表情包列表 URL：

```yaml
frontend:
  emoticons: https://raw.githubusercontent.com/ArtalkJS/Emoticons/master/grps/default.json
```

## 表情包预设

Artalk 社区提供许多表情包预设，你能够挑选几套喜欢的表情包，仅需简单的配置，轻松添加到你的评论系统中，前往仓库：[@ArtalkJS/Emoticons](https://github.com/ArtalkJS/Emoticons)。

## 格式支持

### OwO 格式

[OwO](https://github.com/DIYgod/OwO) 是作者 [@DIYgod](https://github.com/DIYgod) 开发的开源 JS 插件，能让输入框快速拥有插入表情符号的功能。

Artalk 的表情包功能借鉴了其优秀的设计，并且 Artalk 原生适配和兼容 OwO 格式的表情包数据文件，示例如下：

```yaml
frontend:
  emoticons: https://raw.githubusercontent.com/DIYgod/OwO/master/demo/OwO.json
  # 直接食用 OwO 格式的表情包 ↑↑
```

社区也有许多可以直接食用的 OwO 格式表情包资源，例如：[@2X-ercha/Twikoo-Magic](https://github.com/2X-ercha/Twikoo-Magic)。

### Artalk 表情包列表文件标准格式

Artalk 除了支持 OwO 格式外，还内置支持一种标准的表情包列表文件格式：

```js
;[
  {
    name: '颜表情',
    type: 'emoticon', // 字符类型
    items: [
      { key: 'Hi', val: '|´・ω・)ノ' },
      { key: '开心', val: 'ヾ(≧∇≦*)ゝ' },
      //...
    ],
  },
  {
    name: '滑稽',
    type: 'image', // 图片类型
    items: [
      {
        key: '原味稽',
        val: '<图片 URL>',
      },
      //...
    ],
  },
]
```

## 前端配置

::: warning

前端的配置可能会被后端所覆盖，详情参考：[在后端控制前端](/guide/backend/fe-control.html)

:::

在前端配置表情包列表，例如：

```js
Artalk.init({
  // 默认表情包列表，动态引入 ↓↓
  emoticons:
    'https://raw.githubusercontent.com/ArtalkJS/Emoticons/master/grps/default.json',
})
```

## 关闭表情包

你可以将 `emoticons` 设置为 `false` 来禁用表情包功能：

```js
Artalk.init({
  emoticons: false,
})
```

## 加载方式

### 动态加载

将 `emoticons` 属性设置为表情包数据文件的 URL，当打开表情包列表时，Artalk 会动态引入。

```js
Artalk.init({
  emoticons: '<表情包数据文件 URL>',
})
```

远程的表情包文件支持 Artalk、OwO 格式，且支持嵌套、混合加载。

### 静态加载

相较于动态引入，可以将表情包列表对象，作为 Artalk 配置，静态保存在页面的 JS 代码中，避免动态加载：

```js
Artalk.init({
  emoticons: [
    {
      name: '颜表情',
      type: 'emoticon', // 字符类型
      items: [
        { key: 'Hi', val: '|´・ω・)ノ' },
        { key: '开心', val: 'ヾ(≧∇≦*)ゝ' },
        //...
      ],
    },
    {
      name: '滑稽',
      type: 'image', // 图片类型
      items: [
        {
          key: '原味稽',
          val: '<图片 URL>',
        },
        //...
      ],
    },
  ],
})
```

### 混合加载

Artalk 支持 **动态**、**静态** 混合加载，例如：

```js
Artalk.init({
  emoticons: [
    // 动态加载
    'https://raw.githubusercontent.com/DIYgod/OwO/master/demo/OwO.json', // OwO 格式表情包
    'https://raw.githubusercontent.com/qwqcode/huaji/master/huaji.json',
    // 静态加载
    {
      name: '表情包名字',
      type: 'emoticon', // 字符类型
      items: [
        { key: '去吧大师球', val: '(╯°A°)╯︵○○○' },
        //...
      ],
    },
  ],
})
```

### 嵌套引入

Artalk 支持远程表情包资源中**嵌套引入**另外的表情包资源，例如：

```js
Artalk.init({
  emoticons: ['https://example.org/表情包.json'],
})
```

文件 `表情包.json` 中的数据：

```json
[
  "https://example.org/其他表情包.json",
  //...
  {
    // Artalk 格式、OwO 格式
    //...
  }
]
```

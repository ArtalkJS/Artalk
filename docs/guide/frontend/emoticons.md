# 表情包

你可以在前端配置表情包列表，例如：

```js
new Artalk({
    // 默认表情包列表，动态引入 ↓↓
    emoticons: "https://raw.githubusercontent.com/ArtalkJS/Emoticons/master/grps/default.json",
})
```

## 表情包预设

Artalk 社区提供许多表情包预设，你能够挑选几套喜欢的表情包，仅需简单的配置，轻松添加到你的评论系统中，前往仓库：[“@ArtalkJS/Emoticons”](https://github.com/ArtalkJS/Emoticons)。

## 格式支持

### Artalk 标准格式

```js
[{
    "name": "颜表情",
    "type": "emoticon", // 字符类型
    "items": [
        { "key": "Hi", "val": "|´・ω・)ノ" },
        { "key": "开心", "val": "ヾ(≧∇≦*)ゝ" },
        //...
    ]
}, {
    "name": "滑稽",
    "type": "image", // 图片类型
    "items": [
        {
            "key": "原味稽",
            "val": "<图片 URL>"
        },
        //...
    ]
}]
```


### OwO 格式

[OwO](https://github.com/DIYgod/OwO) 是作者 [@DIYgod](https://github.com/DIYgod) 开发的 JS 插件，能让输入框快速拥有插入表情符号的功能。

Artalk 的表情包功能灵感也源于此，并且 Artalk 适配和兼容 OwO 格式的表情包数据文件，示例如下：

```js
new Artalk({
    emoticons: "https://raw.githubusercontent.com/DIYgod/OwO/master/demo/OwO.json",
    // 直接食用 OwO 格式的表情包 ↑↑
})
```

社区也有许多可以直接食用的 OwO 格式表情包资源，例如：[@2X-ercha/Twikoo-Magic](https://github.com/2X-ercha/Twikoo-Magic)。

## 装载方式

### 动态引入

将 `emoticons` 属性设置为表情包数据文件的 URL，当打开表情包列表时，Artalk 会动态引入。

```js
new Artalk({
    emoticons: "<表情包数据文件 URL>",
})
```

远程的表情包文件支持 Artalk、OwO 格式，且支持嵌套、混合装载。

### 静态装载

相较于动态引入，可以将表情包列表对象，作为 Artalk 配置，静态保存在页面的 JS 代码中，避免动态加载：

```js
new Artalk({
    emoticons: [{
        "name": "颜表情",
        "type": "emoticon", // 字符类型
        "items": [
            { "key": "Hi", "val": "|´・ω・)ノ" },
            { "key": "开心", "val": "ヾ(≧∇≦*)ゝ" },
            //...
        ]
    }, {
        "name": "滑稽",
        "type": "image", // 图片类型
        "items": [
            {
                "key": "原味稽",
                "val": "<图片 URL>"
            },
            //...
        ]
    }],
})
```

### 混合装载

Artalk 支持 **动态**、**静态** 混合装载，例如：

```js
new Artalk({
    emoticons: [
        // 动态装载
        "https://raw.githubusercontent.com/DIYgod/OwO/master/demo/OwO.json", // OwO 格式表情包
        "https://raw.githubusercontent.com/qwqcode/huaji/master/huaji.json",
        // 静态装载
        {
            "name": "表情包名字",
            "type": "emoticon", // 字符类型
            "items": [
                { "key": "去吧大师球", "val": "(╯°A°)╯︵○○○" },
                //...
            ]
        }
    ]
})
```

### 嵌套装载

Artalk 支持远程表情包资源中**嵌套引入**另外的表情包资源，例如：

```js
new Artalk({
    emoticons: [
        "https://example.org/表情包.json"
    ]
})
```

文件 `表情包.json` 中的数据：

```json
[
    "https://example.org/其他表情包.json",
    //...
    { // Artalk 格式、OwO 格式
        //...
    }
]
```

### 关闭表情包功能

你可以将 `emoticons` 设置为 `false` 来禁用表情包功能：

```js
new Artalk({
    emoticons: false
})
```

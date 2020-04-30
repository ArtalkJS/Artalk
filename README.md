<p align="center">
<img src="./docs/banner.png" alt="Artalk" width="200">
</p>

# [Artalk](https://artalk.js.org)
> 一款简洁有趣的自托管评论系统

[![](https://img.shields.io/npm/v/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![](https://img.shields.io/npm/dt/artalk.svg?style=flat-square)](https://www.npmjs.com/package/artalk)
[![](https://img.shields.io/github/last-commit/qwqcode/Artalk/master.svg?style=flat-square)](https://github.com/qwqcode/Artalk/commits/master)
[![](https://img.shields.io/github/issues-raw/qwqcode/Artalk.svg?style=flat-square)](https://github.com/qwqcode/Artalk/issues)
[![](https://img.shields.io/github/issues-pr-raw/qwqcode/Artalk.svg?style=flat-square)](https://github.com/qwqcode/Artalk/pulls)
[![](https://img.shields.io/github/license/qwqcode/Artalk.svg?style=flat-square)](./LICENSE)
[![](https://img.shields.io/badge/%24-donate-%23ff69b4.svg?style=flat-square)](#捐助)

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
- 无限层级回复
- 滚动自动加载
- 评论打分
- TypeScript

[查看 DEMO](https://artalk.js.org)

## 食用

> 生产环境版下载：[Artalk.js](https://raw.githubusercontent.com/qwqcode/Artalk/master/dist/Artalk.js) | [Artalk.css](https://raw.githubusercontent.com/qwqcode/Artalk/master/dist/Artalk.css)

1. 部署 Artalk 的后端，传送门：[PHP API](https://github.com/qwqcode/Artalk-API-PHP) / Go API / Node API / Python API
2. 前端页面配置 Artalk 如下：

```html
<!DOCTYPE html>
<html>
<head>
  <!-- ... -->
  <link href="dist/Artalk.css" rel="stylesheet">
  <script type="text/javascript" src="dist/Artalk.js"></script>
</head>
<body>
  <div id="ArtalkComments"></div>
  <script>
  new Artalk({
    el: '#ArtalkComments',
    placeholder: '来啊，快活啊 (/ω＼)',
    noComment: '快来成为第一个评论的人吧~',
    defaultAvatar: 'mp',
    pageKey: '[页面唯一标识]',
    serverUrl: '[后端程序URL]',
    readMore: {
      pageSize: 15, // 每次请求获取评论数
      autoLoad: true // 滚动到底部自动加载
    }
  });
  </script>
</body>
</html>
```

更多配置项详见 [types/artalk-config.d.ts](./types/artalk-config.d.ts)

## 开发

```bash
git clone https://github.com/qwqcode/Artalk.git
cd Artalk
yarn install

# Dev
yarn run dev

# Build
yarn run build
```

Made with ♥ by QWQAQ

## 捐助

如果您愿意给予我一点点支持与帮助，可以通过下面的方式向我捐助，爱你~~~ (/ω＼)

| Alipay | Wechat | 
| :------: | :------: | 
| <img width="150" src="https://raw.githubusercontent.com/qwqcode/donate-qwqaq/master/docs/donate/alipay.png"> | <img width="150" src="https://raw.githubusercontent.com/qwqcode/donate-qwqaq/master/docs/donate/wechat.png"> | 

捐助者的名字将保存于 [捐助者名单](https://github.com/qwqcode/donate-qwqaq)

最后，我再次对您致以我最最最最最诚挚的感谢！

## License
[GPL-2.0](./LICENSE)

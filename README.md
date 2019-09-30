<p align="center">
<img src="./docs/banner.png" alt="Artalk" width="200">
</p>

# Artalk
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
- 评论打分
- TypeScript

[查看 DEMO](https://artalk.js.org)

## 食用

1. 部署 Artalk 的后端，[PHP 版传送门](https://github.com/qwqcode/ArtalkServerPhp)
2. 前端页面配置 Artalk

```html
<!DOCTYPE html>
<html>
<head>
  <!-- ... -->
  <link href="dist/Artalk.css" rel="stylesheet">
  <script type="text/javascript" src="dist/Artalk.js"></script>
</head>
<body>
  <!-- ... -->
  <div id="ArtalkComments"></div>
  <!-- ... -->
  <script>
  $(document).ready(function () {
    var artalk = new Artalk({
      el: '#ArtalkComments',
      placeholder: '来啊，快活啊 (/ω＼)',
      defaultAvatar: 'mp',
      pageKey: '[页面唯一标识]',
      serverUrl: '[后端程序URL]'
    });
  });
  </script>
</body>
</html>
```

生产环境版下载：[Artalk.js](https://raw.githubusercontent.com/qwqcode/Artalk/master/dist/Artalk.js) | [Artalk.css](https://raw.githubusercontent.com/qwqcode/Artalk/master/dist/Artalk.css)

## 开发

Made with ♥ by QWQAQ

```bash
git clone https://github.com/qwqcode/Artalk.git
cd Artalk
yarn install

# Dev
yarn run dev

# Build
yarn run build
```

## 捐助

如果您觉得我的项目对您有帮助，并且您愿意给予我一点小小的支持，您可以通过以下方式向我捐助，这样可以维持项目持续地发展，非常感谢！(/ω＼)

| Alipay | Wechat | 
| :------: | :------: | 
| <img width="150" src="https://raw.githubusercontent.com/qwqcode/donate-qwqaq/master/docs/donate/alipay.png"> | <img width="150" src="https://raw.githubusercontent.com/qwqcode/donate-qwqaq/master/docs/donate/wechat.png"> | 

捐助者的名字将保存于 [捐助者名单](https://github.com/qwqcode/donate-qwqaq)

最后，我再次对您致以诚挚的感谢！

## License
[GPL-2.0](./LICENSE)

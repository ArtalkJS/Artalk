# @artalkjs/plugin-katex

食用方法

```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Document</title>

  <!-- 引入 katex -->
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.15.3/dist/katex.min.css">
  <script defer src="https://cdn.jsdelivr.net/npm/katex@0.15.3/dist/katex.min.js"></script>

  <!-- 引入 Artalk -->
  <link href="https://cdn.jsdelivr.net/npm/artalk/dist/Artalk.css" rel="stylesheet">
  <script src="https://cdn.jsdelivr.net/npm/artalk/dist/Artalk.js"></script>

  <!-- 引入 @artalkjs/plugin-katex -->
  <script defer src="https://cdn.jsdelivr.net/npm/@artalkjs/plugin-katex/dist/artalk-plugin-katex.js"></script>
</head>
<body>
  
<!-- Artalk -->
<div id="Comments"></div>
<script>
  new Artalk({
    el:        '#Comments',
    pageKey:   '<页面链接>',
    pageTitle: '<页面标题>',
    server:    '<后端地址>',
    site:      '<站点名称>',
  })
</script>

</body>
</html>
```

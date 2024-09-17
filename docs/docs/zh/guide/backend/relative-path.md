# 相对 / 绝对路径

Artalk 支持解析相对路径，因此你可以在前端页面进行如下配置：

```js
Artalk.init({
  site: '举个栗子站点', // 你的站点名
  pageKey: '/relative-path/xx.html', // 使用相对路径
})
```

建议页面使用相对路径，因为这为日后的「站点迁移」需求创建条件。

然后，你需要在侧边栏「[控制中心](../frontend/sidebar.md#控制中心)」-「站点」找到 “举个栗子站点”，修改站点 URL。

![](/images/relative-path/1.png)

之后，所有相对路径都会「基于这个站点 URL」，例如：

```bash
"/relative-path/xx.html"
         ↓ 解析为
"https://设定的举个栗子站点URL.xxx/relative-path/xx.html"
```

## 解析后的 URL 用途

站点 URL + 页面 相对路径 将用于：

- **邮件通知**中的回复评论链接
- **侧边栏**快速跳转到某条评论
- **控制中心**页面管理打开页面
- **获取页面标题**等信息时使用

## 配置多个站点 URL 的情况

你可能需要配置站点的多个 URL 来允许 Referer 和跨域。

「[控制中心](../frontend/sidebar.md#控制中心)」-「站点」修改站点 URL 支持为站点添加多个 URL，用英文逗号 `,` 分隔开每个 URL 即可。

**当站点存在多个 URL 时**，「相对路径」会基于多个 URL 中的「**第一个**」URL。

## 使用绝对路径的情况

区别于使用相对路径，你可以使用绝对路径，例如前端这样配置：

```js
Artalk.init({
  pageKey: 'https://your_domain.com/relative-path/xx.html', // 使用绝对路径
})
```

这时后端不会去解析该地址，邮件、侧边栏等地方都是直接使用 `pageKey` 这个绝对路径来定位页面。

# 开发说明

## 参考资源

📖 [开发者贡献指南 (CONTRIBUTING.md)](https://github.com/ArtalkJS/Artalk/blob/master/CONTRIBUTING.md)

🔬 <a href="/http-api.html" target="_blank">HTTP API 文档 (OpenAPI)</a>

🔖 [贡献者契约行为准则 (CODE_OF_CONDUCT.md)](https://github.com/ArtalkJS/Artalk/blob/master/CODE_OF_CONDUCT.md)

## 配置文档

- [后端配置](../guide/backend/config.md)
- [前端配置](../guide/frontend/config.md)

## 补充说明

由于 Artalk 正处于开发阶段，使用此文档中 `API`、`Event` 前请务必检查时效性。

- `API` 部分参考源码
  - [@ArtalkJS/Artalk - src/api/](https://github.com/ArtalkJS/Artalk/tree/master/ui/artalk/src/api)
  - [@ArtalkJS/Artalk - server/server.go](https://github.com/ArtalkJS/Artalk/blob/master/server/server.go)
- `UI` 及 `Event` 部分参考源码
  - [@ArtalkJS/Artalk - src/artalk.ts](https://github.com/ArtalkJS/Artalk/blob/master/ui/artalk/src/artalk.ts)
  - [@ArtalkJS/Artalk - src/types/event.ts](https://github.com/ArtalkJS/Artalk/blob/master/ui/artalk/src/types/event.ts)

通过 Artalk 提供的 `API` 和 `Event`，你可以实现很多高级功能，比如编写评论管理机器人、评论提醒推送插件等。Artalk 并不为此提供技术指导，但鼓励你参考此处的文档自行定制。

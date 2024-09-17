# 兼容性

## 浏览器兼容性

Artalk 客户端生产环境代码的构建通过 esbuild 执行，构建目标为 `es2015`。但请注意，构建仅会处理语法转义遵循 `es2015` 标准，不会处理 API 的兼容性（不包含任何 polyfill）。

如果你的项目需要支持较旧的浏览器和设备，请考虑在捆绑应用程序中包含全局 polyfill。我们推荐使用 [Polyfill.io](https://polyfill.io/) 服务，它会根据用户的 User-Agent 动态返回所需的 polyfill。另外，你也可以使用 [core-js](https://github.com/zloirock/core-js) 或 [Babel](https://babeljs.io/) 来手动添加 polyfill。

以下是 Artalk 使用的现代功能的列表：

- [All ECMAScript 2015 (ES6) features](https://compat-table.github.io/compat-table/es6/)
- [Fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API)
- [AbortController](https://developer.mozilla.org/en-US/docs/Web/API/AbortController)
- [Intersection Observer](https://developer.mozilla.org/en-US/docs/Web/API/Intersection_Observer_API)

## 服务器兼容性

Artalk 服务端程序基于 Golang 开发，支持 Linux、Windows 和 macOS 等操作系统。当前 Artalk 的 Golang 版本可通过 [go.mod](https://github.com/ArtalkJS/Artalk/blob/master/go.mod#L3) 查看。关于 Golang 的最低操作系统要求，可参考 [Go Wiki](https://go.dev/wiki/MinimumRequirements)。

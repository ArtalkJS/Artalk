# Artalk Plugin Kit

Artalk Plugin Kit 是 Artalk 的插件开发工具包，提供了一系列的工具函数和组件，帮助开发者更方便的开发插件。

## 安装

```bash
pnpm add artalk
pnpm add @artalk/plugin-kit -D
```

工具包中提供了 Vite 的集成插件，使用 Vite 开发能够开箱即用地构建 Artalk 插件，简化 Vite 的配置。

```js
// vite.config.js
import { ViteArtalkPluginKit } from '@artalk/plugin-kit'

export default {
  plugins: [ViteArtalkPluginKit()],
}
```

## 开发

在开发插件前，你需要在 `package.json` 设置插件名字，修改 `name` 字段，以 `artalk-plugin-` 开头。

插件入口文件默认为 `src/main.ts`，你需要在入口文件导出名为 `ArtalkDemoPlugin` 的 `ArtalkPlugin` 类型的对象。

```ts
// src/main.ts
import type { ArtalkPlugin } from 'artalk'

export const ArtalkDemoPlugin: ArtalkPlugin = (ctx) => {
  ctx.on('mounted', () => {
    // Your plugin code here
  })
}
```

执行 `pnpm dev` 开发插件，Vite 将会启动开发服务器。浏览器访问 ViteArtalkPluginKit 插件提供的内置 Artalk 调试页面，该页面已自动注入并启用了你当前正在开发的插件。

## 构建

执行 `pnpm build` 构建插件，构建产物为 `dist/artalk-plugin-demo.js`，以及 `.mjs`、`.cjs`、`.d.ts` 等文件。

ViteArtalkPluginKit 将帮助你自动生成合并的 `.d.ts` 文件，并自动在 `package.json` 中填充 `exports`、`types` 等字段。

ViteArtalkPluginKit 会进行一系列代码检查，以确保 Artalk 的插件符合开发规范。

## 发布

执行 `pnpm publish` 发布插件，你可以在 `package.json` 中配置 `"files": ["dist"]` 只推送需要发布的构建产物。

## 配置

可以通过配置 `artalkInitOptions` 控制 Artalk 的初始化参数。

```ts
// vite.config.js
import { ViteArtalkPluginKit } from '@artalk/plugin-kit'

export default {
  plugins: [
    ViteArtalkPluginKit({
      artalkInitOptions: {
        // Your Artalk init options here
      },
    }),
  ],
}
```

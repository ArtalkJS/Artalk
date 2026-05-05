# 2026-05 依赖升级说明

本文档记录了本次大规模依赖升级所涉及的全部变更，涵盖前端工具链、ESLint 生态、TypeScript 编译配置、构建工具迁移及相关文档与 CI 的同步更新。

---

## 一、前端依赖升级

### 核心工具链

| 包 | 旧版本 | 新版本 |
|----|--------|--------|
| `typescript` | 5.6.3 | 6.0.3 |
| `vite` | 5.4.9 | 8.0.10 |
| `vitest` | ^2.1.3 | ^4.1.5 |
| `rollup` | 4.24.0 | 4.60.3 |
| `sass` | 1.80.3 | 1.99.0 |
| `terser` | 5.36.0 | 5.46.2 |
| `postcss` | 8.4.47 | 8.5.14 |
| `autoprefixer` | 10.4.20 | 10.5.0 |
| `tsx` | ^4.19.1 | ^4.21.0 |
| `cross-env` | ^7.0.3 | ^10.1.0 |
| `prettier` | 3.3.3 | 3.8.3 |

### ESLint 生态

| 包 | 旧版本 | 新版本 |
|----|--------|--------|
| `eslint` | ^9.13.0 | ^10.3.0 |
| `typescript-eslint` | ^8.10.0 | ^8.59.1 |
| `@typescript-eslint/eslint-plugin` | 8.10.0 | 8.59.1 |
| `eslint-plugin-vue` | ^9.29.0 | ^10.9.0 |
| `eslint-plugin-import-x` | ^4.3.1 | ^4.16.2 |
| `eslint-plugin-compat` | ^6.0.1 | ^7.0.2 |
| `eslint-plugin-react-hooks` | 5.1.0-beta | 7.1.1 |
| `eslint-config-prettier` | 9.1.0 | 10.1.8 |
| `eslint-import-resolver-typescript` | 3.6.3 | 4.4.4 |
| `vue-eslint-parser` | ^9.4.3 | ^10.4.0 |
| `globals` | ^15.11.0 | ^17.6.0 |

### Vite 插件

| 包 | 旧版本 | 新版本 |
|----|--------|--------|
| `vite-plugin-dts` | 4.2.4 | 5.0.0 |
| `vite-plugin-checker` | 0.8.0 | 0.13.0 |

### UI 子包

| 包 | 子包 | 旧版本 | 新版本 |
|----|------|--------|--------|
| `marked` | `ui/artalk` (dep + peer) | ^14.1.3 | ^18.0.3 |
| `vue` | artalk-sidebar, docs | ^3.5.12 | ^3.5.33 |
| `vue-router` | artalk-sidebar | ^4.4.5 | ^5.0.6 |
| `vue-i18n` | artalk-sidebar | ^10.0.4 | ^11.4.0 |
| `pinia` | artalk-sidebar | ^2.2.4 | ^3.0.4 |
| `vue-tsc` | artalk-sidebar | ^2.1.6 | ^3.2.8	 |
| `@vitejs/plugin-vue` | artalk-sidebar | ^5.1.4 | ^6.0.6 |
| `unplugin-auto-import` | artalk-sidebar | ^0.18.3 | ^21.0.0 |
| `unplugin-vue-components` | artalk-sidebar | ^0.27.4 | ^32.0.0 |
| `react` / `react-dom` | docs/landing | ^18.3.1 | ^19.2.5 |
| `@vitejs/plugin-react-swc` | docs/landing | ^3.7.1 | ^4.3.0 |
| `react-i18next` | docs/landing | ^15.0.3 | ^17.0.6 |
| `i18next` | docs/landing | ^23.16.0 | ^26.0.8 |
| `katex` | plugin-katex | ^0.16.11 | ^0.16.45 |
| `lightgallery` | plugin-lightbox | ^2.7.2 | ^2.9.0 |
| `vitepress` | docs/docs | 1.4.1 | 1.6.4 |
| `@redocly/cli` | docs/swagger | 1.25.7 | 2.30.3 |
| `solid-js` | plugin-auth | ^1.9.2 | ^1.9.12 |

---

## 二、包管理器与运行时要求更新

**`package.json` 根配置：**

```jsonc
// 旧
"packageManager": "pnpm@9.10.0"

// 新
"packageManager": "pnpm@10.33.2",
"engines": {
  "node": ">=22.19.0",
  "pnpm": ">=10.33.2"
}
```

新增 `pnpm.peerDependencyRules.allowedVersions` 以抑制 peer dep 版本警告：

```json
"pnpm": {
  "peerDependencyRules": {
    "allowedVersions": {
      "eslint-plugin-import>eslint": "10",
      "eslint-plugin-react>eslint": "10",
      "artalk>marked": "18"
    }
  }
}
```

新增 `prepare` 脚本，`pnpm install` 后自动构建本地 workspace 包：

```json
"prepare": "pnpm build:plugin-kit && pnpm build:eslint-plugin"
```

本地 workspace 引用方式也由 npm 发布版改为 workspace 直引：

```jsonc
// 旧
"@artalk/plugin-kit": "^1.0.7",
"eslint-plugin-artalk": "^1.0.2",

// 新
"@artalk/plugin-kit": "workspace:*",
"eslint-plugin-artalk": "workspace:*",
```

---

## 三、构建工具迁移：tsup → tsdown

`eslint-plugin-artalk` 和 `@artalk/plugin-kit` 从 `tsup`（基于 esbuild）迁移到 `tsdown`（基于 rolldown）。

### 迁移对比

| 项目 | tsup | tsdown |
|------|------|--------|
| 构建器内核 | esbuild | rolldown |
| `?raw` 文件插件 | `esbuild-plugin-raw`，API：`setup/onResolve/onLoad` | 内联 rolldown 插件，API：`resolveId/load` |
| externals | `external: [...]` | `deps.neverBundle: [...]` |
| 产物扩展名 | `.js` / `.d.ts` | `.mjs`/`.cjs` / `.d.mts`/`.d.cts` |
| 插件类型 | `esbuild.Plugin` | `TsdownPlugin`（从 `tsdown` 直接导入） |

### 产物入口更新

由于产物扩展名变化，`package.json` 的 `main`/`types`/`exports` 字段同步更新：

```jsonc
// 旧
"main": "dist/main.js",
"types": "dist/main.d.ts",
"exports": {
  "require": { "types": "./dist/main.d.cjs" },
  "default": { "types": "./dist/main.d.ts", "default": "./dist/main.js" }
}

// 新
"main": "dist/main.cjs",
"types": "dist/main.d.cts",
"exports": {
  "require": { "types": "./dist/main.d.cts", "default": "./dist/main.cjs" },
  "default": { "types": "./dist/main.d.mts", "default": "./dist/main.mjs" }
}
```

---

## 四、eslint-plugin-artalk 重构

### 问题背景

`eslint@10` 内置 `@eslint/core@1.2.1`，而 `typescript-eslint@8` 依赖 `@eslint/core@0.7.0`，两者的 `RuleContext` API 不兼容，导致 `TSESLint.RuleModule` 无法直接赋值给 `eslint` 的 `Rule.RuleModule`。

### 解决方案

完全弃用 `ESLintUtils.RuleCreator`，改为直接使用 `eslint` 原生的 `Rule.RuleModule` 类型构造规则，无需类型桥接或强制转换。

**`src/helper.ts`（重写）：**

```ts
import type { Rule } from 'eslint'

export function createRule(def: {
  name: string
  meta: Rule.RuleMetaData
  defaultOptions?: unknown[]
  create: Rule.RuleModule['create']
}): Rule.RuleModule {
  return {
    meta: { ...def.meta, docs: { ...def.meta.docs, url: `${docsBaseUrl}${def.name}` } },
    create: def.create,
  }
}
```

**`src/main.ts`（重写）：**

改为 default-only 导出，避免 rolldown MIXED_EXPORTS 警告：

```ts
import type { ESLint } from 'eslint'

const plugin: ESLint.Plugin & { configs: typeof configs } = {
  rules: { 'artalk-plugin': artalkPlugin },
  configs,
}

export default plugin
```

---

## 五、TypeScript 编译配置重构

### `tsconfig.base.json`

- 新增 `"composite": true`：支持项目引用（Project References），使 ESLint `projectService` 能正确解析各子包
- 移除 `"noEmit": true`（与 `composite` 不兼容）
- 移除 `"experimentalDecorators": true`（TypeScript 5+ 已默认支持）
- 移除 `"allowSyntheticDefaultImports": true`（`esModuleInterop: true` 已隐含）

### 新增根级 `tsconfig.json`

新增仅供工具链（ESLint、IDE）使用的根 tsconfig，通过 `references` 聚合所有子包：

```json
{
  "extends": "./tsconfig.base.json",
  "compilerOptions": { "noEmit": true, "allowJs": true },
  "files": ["eslint.config.mjs", "vitest.workspace.ts", ...],
  "references": [
    { "path": "./ui/artalk" },
    { "path": "./ui/artalk-sidebar" },
    ...
  ]
}
```

### 各子包 `tsconfig.json` 变更

需要 composite 支持但由 vite 构建（不 emit）的包（`artalk-sidebar`、`docs/landing`）：

```jsonc
// 不能用 noEmit，改为
"emitDeclarationOnly": true,
"declarationDir": ".tsbuildinfo-dts"
```

并将 `.tsbuildinfo-dts/` 加入 `.gitignore`。

`ui/artalk/tsconfig.json` 新增：

```json
"include": ["./src", "./types", "./tests", "./scripts", "package.json"],
"files": ["./package.json"]
```

原因：`composite: true` 下所有 import 的文件必须被 tsconfig 显式涵盖；`import { version } from '~/package.json'` 的 JSON 文件不会被目录 glob 匹配，需手动列出。同时 `vite-plugin-dts` 检查 `include` 而 `tsc` 检查 `files`，因此两处均需声明。

`ui/artalk/vite.config.ts` 为 `vite-plugin-dts` 添加：

```ts
compilerOptions: { composite: false }
```

避免 dts 插件在生成声明文件时触发 composite 文件列表检查（TS6307）。

---

## 六、ESLint 配置重构（`eslint.config.mjs`）

| 变更 | 旧 | 新 |
|------|----|----|
| 配置入口 | `eslintTs.config()` | `defineConfig()` from `'eslint/config'` |
| typescript-eslint 配置展开 | `...eslintTs.configs.recommended` | `eslintTs.configs.recommended`（无需展开） |
| vue 配置展开 | `...pluginVue.configs['flat/recommended']` | `pluginVue.configs['flat/recommended']` |
| `@typescript-eslint` 插件 | 手动注册 `pluginTS` | 由 `typescript-eslint` 内部管理，移除手动注册 |
| TypeScript 项目解析 | `project: tsProjects`（数组 glob） | `projectService: true`（语言服务模式） |
| import-x resolver | `project: tsProjects` | `project: './tsconfig.json'`（单一根 tsconfig） |
| React 版本 | `version: '18.3'` | `version: '19.2'` |
| typed linting | `recommended` | `recommendedTypeChecked` |

---

## 七、依赖分类修复

### `ui/plugin-katex`

- `katex`：从 `dependencies` 移至 `devDependencies`（vite 配置中已 external，不打包进产物，应由用户提供）
- `@types/katex`：从 `dependencies` + `peerDependencies` 完全移除（katex >= 0.16 已自带 `types/katex.d.ts`，`@types/katex` 为旧版遗留包）

### `ui/eslint-plugin-artalk`

- 新增 `eslint: "^10.3.0"` 到 `devDependencies`（peerDep 已声明但开发依赖缺失）

### `ui/plugin-kit`

- 新增 `vite: "^8.0.10"` 到 `devDependencies`（源码中 `import type { ViteDevServer, Plugin } from 'vite'`，peerDep 已声明但开发依赖缺失）

### `ui/plugin-auth`、`ui/plugin-katex`、`ui/plugin-lightbox`

- 新增 `"@artalk/plugin-kit": "workspace:^"` 到各自 `devDependencies`（这些插件基于 plugin-kit 开发，之前遗漏了显式声明）

### 根 `package.json`

- `@microsoft/api-extractor` 保留在根级（由 `vite-plugin-dts` 的 `bundleTypes: true` 调用，属于根级工具链依赖）

---

## 八、vitest 配置迁移

`vitest.workspace.ts` 从 `defineWorkspace` API 迁移到 `defineConfig` + `test.projects`（vitest v4 新 API）：

```ts
// 旧
import { defineWorkspace } from 'vitest/config'
export default defineWorkspace(['./ui/artalk/vitest.config.ts'])

// 新
import { defineConfig } from 'vitest/config'
export default defineConfig({
  test: { projects: ['./ui/artalk/vitest.config.ts'] },
})
```

---

## 九、GitHub Actions CI 更新

所有使用 Node.js 和 pnpm 的工作流同步更新版本：

| 文件 | Node.js | pnpm |
|------|---------|------|
| `test-frontend.yml` | `[18, 20]` → `[22]` | `9.10.0` → `10.33.2` |
| `build-ui.yml` | `20.x` → `22.x` | `9.10.0` → `10.33.2` |
| `test-docs.yml` | `20.x` → `22.x` | `9.10.0` → `10.33.2` |
| `docs-cn.yml` | `20` → `22` | `9.10.0` → `10.33.2` |
| `build-nightly.yml` | `20.x` → `22.x` | `9.10.0` → `10.33.2` |
| `release.yml` | `20.x` → `22.x` | — |

---

## 十、文档更新

- `CONTRIBUTING.md`（英文）：Node.js 要求 `>=20.17.0` → `>=22.19.0`，pnpm `>=9.10.0` → `>=10.33.2`，修复 `make build-fronted` 拼写错误为 `make build-frontend`
- `docs/docs/zh/develop/contributing.md`（中文）：同步上述版本要求更新
- `docs/swagger/redocly.yaml`：更新 Redocly 配置适配 CLI v2 的新格式
- `docs/swagger/package.json`：`swagger:serve` 命令 `npx @redocly/cli preview-docs` → `redocly preview`

---

## 十一、Go 后端依赖更新

`go.mod` / `go.sum` 同步更新了若干后端 Go 模块依赖（具体版本见 `go.mod` diff）。

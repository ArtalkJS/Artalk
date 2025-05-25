# 开发者贡献指南

本指南帮助开发者了解如何在 Artalk 项目中贡献代码、文档、测试和反馈。我们欢迎所有形式的贡献，无论是大型功能开发、小型修复、文档改进还是测试。我们相信每一个贡献都是推动项目前进的重要一环。

本指南包含多语言版本：[English](https://github.com/ArtalkJS/Artalk/blob/master/CONTRIBUTING.md) • [简体中文](https://artalk.js.org/develop/contributing.html)

## 克隆仓库

克隆仓库：

```sh
git clone https://github.com/ArtalkJS/Artalk.git
```

建议你先 fork 仓库，然后克隆 fork 后的仓库。

导航到目录：

```sh
cd Artalk
```

## 开发环境

为了开发 Artalk 的前端和后端，请安装以下工具：

- [Node.js](https://nodejs.org/en/) (>= 20.17.0)
- [PNPM](https://pnpm.io/) (>= 9.10.0)
- [Go](https://golang.org/) (>= 1.22)
- [Docker](https://www.docker.com/) (>= 20.10.0)（可选）
- [Docker Compose](https://docs.docker.com/compose/) (>= 1.29.0)（可选）

## 开发工作流程概述

开发工作流程包括以下步骤：

1. **搭建完整的开发实例**：

   - 导航到 Artalk 仓库目录。
   - 执行 `make dev` 在端口 `23366` 上运行后端。
   - 执行 `pnpm dev` 在端口 `5173` 上运行前端。
   - 可选地，执行 `pnpm dev:sidebar` 在端口 `23367` 上运行侧边栏前端。

2. **访问前端进行开发**：

   - 打开浏览器，访问 `http://localhost:5173` 进行开发和测试。

3. **测试**：

   - 使用 `make test` 运行后端单元测试。
   - 使用 `pnpm test` 运行 UI 单元测试。
   - 使用 `make test-frontend-e2e` 运行 UI 端到端测试。

4. **构建前端代码**：

   - 当你对前端代码进行更改后，使用 `pnpm build:all` 构建完整的前端程序。
   - JavaScript 和 CSS 代码将在 `ui/artalk/dist` 目录中找到。
   - 前端代码由 GitHub Actions 自动部署到 NPM。

5. **构建后端代码**：

   - 当你对后端代码进行更改后，注意后端程序会嵌入前端代码。
   - 在构建后端之前，运行 `make build-frontend`（它运行 `scripts/build-frontend.sh`）。这将把嵌入的前端主程序和侧边栏前端程序放入 `/public` 目录。
   - 然后，使用 `make build` 构建后端程序。

更多探索：

- **开发 Artalk 插件或主题**：
  请参考 [Artalk 插件开发指南](https://artalk.js.org/develop/plugin.html)。

- **Makefile**：
  你可以探索 `Makefile` 中的过程和命令。

- **GitHub 上的自动化 CI**：
  GitHub 上的自动化 CI 流水线位于 `.github/workflows` 目录中，负责构建前端、后端和 Docker 镜像，发布 Docker 镜像到 Docker Hub，处理 npm 包发布，并创建 GitHub Releases，同时测试前端和后端。

以下部分提供了更详细的开发工作流程信息。

## 项目结构

Artalk 是一个 monorepo 项目，这意味着所有源代码都在一个仓库中。尽管如此，前端和后端保持分离。前端部分位于 `./ui` 目录中。

### 目录概览

- `bin/`：编译的二进制文件。（此目录被 git 忽略）
- `cmd/`：命令行工具的源代码。
- `conf/`：示例配置文件。
- `data/`：本地数据。（此目录被 git 忽略）
- `docs/`：文档站点的源代码。
- `i18n/`：翻译文件。
- `local/`：本地文件。（此目录被 git 忽略）
- `internal/`：内部包的 Go 源代码。
- `public/`：静态文件。构建的前端文件将被复制到这里。
- `scripts/`：开发和构建的脚本。
- `server/`：后端的 Go 源代码。
- `ui/`：前端的 UI 源代码。
- `.github/`：GitHub Actions 工作流。
- `.vscode/`：VSCode 设置。

## 后端开发

### 运行后端

1. **创建配置文件**：

   - 将 `./conf/artalk.example.yml` 复制到根目录。
   - 重命名为 `artalk.yml`。
   - 根据需要修改文件。

2. **启动后端程序**：

   - 运行 `make dev` 启动后端，端口为 `23366`。
   - 通过 `http://localhost:23366` 访问。
   - 建议保留默认端口进行测试。
   - （替代）使用 `ARGS="version" make dev` 传递启动参数（默认是 `ARGS="server"`）。

`make dev` 命令会带有调试符号构建后端，这方便使用 GDB 进行调试。

如果你使用 VSCode，可以使用 `F5` 键（或运行按钮）直接启动调试后端程序。

### 构建后端

1. **获取最新的 Git 标签**：
   执行 `git fetch --tags` 获取最新的 Git 标签信息。此标签将用作后端应用程序的版本号。或者，通过环境变量指定：`VERSION="v1.0.0"` 和 `COMMIT_HASH="66128e"`。

2. **构建前端**：
   运行 `make build-frontend` 构建前端，因为后端应用程序会嵌入前端代码。

3. **构建后端**：
   执行 `make build` 构建后端程序。

构建的二进制文件将位于 `./bin` 目录中。

## 前端开发（UI）

### 运行前端

1. **安装依赖**：
   运行 `pnpm install` 安装所有必要的依赖。

2. **启动开发服务器**：
   执行 `pnpm dev --port 5173` 启动前端开发服务器。

3. **访问前端**：
   前端应用程序默认将在端口 `5173` 上运行。在浏览器中访问 `http://localhost:5173`。

4. **侧边栏开发**：
   要运行侧边栏部分的前端，请使用 `pnpm dev:sidebar`。

端口 `5173` 和 `23367` 用于前端开发和测试。这些端口被包括在后端程序的 `ATK_TRUSTED_DOMAINS` 配置中。

### 构建前端

要构建前端，请运行以下命令：

```sh
make build-frontend
```

编译的 JavaScript 和 CSS 文件将位于 `/public` 目录中，嵌入到后端程序中。

更多详细信息，请参考 `scripts/build-frontend.sh` 脚本。

### 发布前端

核心的 Artalk 客户端代码通过 GitHub Actions 自动部署到 NPM。部署过程在 `.github/workflows/build-ui.yml` 文件中定义。

如果你在 monorepo 中编写 Artalk 插件，它不会被自动部署。你需要手动发布插件。要发布插件，请首先确保在 `package.json` 文件中更新版本号。然后，运行以下命令：

```sh
pnpm publish --access public
```

#### 检查已发布版本

`pnpm check:publish` 脚本旨在验证项目中的所有包是否在 npm 上发布了最新版本。它会自动跳过私有包，并允许使用 `-F` 选项进行筛选以检查特定包。

检查所有包：

```sh
pnpm check:publish
```

这将启动检查存储库中所有公共包的过程，以确保它们在 npm 上是最新的。

### 开发和调试 `@artalk/plugin-kit` 和 `eslint-plugin-artalk`

这个 monorepo 默认会从 NPM 安装最新版本的 `@artalk/plugin-kit` 和 `eslint-plugin-artalk` 发布包。要开发和调试这些包，请在仓库的根目录下运行以下命令：

```sh
pnpm link --global --dir ui/eslint-plugin-artalk
pnpm link --global --dir ui/plugin-kit

pnpm link eslint-plugin-artalk
pnpm link @artalk/plugin-kit
```

执行这些命令后，`eslint-plugin-artalk` 和 `@artalk/plugin-kit` 包将链接到本地开发环境中的 monorepo 工作空间。你可以在 `ui/eslint-plugin-artalk` 和 `ui/plugin-kit` 目录中修改源代码，这些更改会自动反映在 monorepo 工作空间中。如果你希望取消链接这些包，请运行以下命令：

```sh
pnpm unlink eslint-plugin-artalk
pnpm unlink @artalk/plugin-kit

pnpm uninstall --global eslint-plugin-artalk
pnpm uninstall --global @artalk/plugin-kit
```

更多详情请参考 [pnpm link 文档](https://pnpm.io/cli/link)。

## Docker 开发

### 构建 Docker 镜像

要构建 Docker 镜像，请运行以下命令：

```sh
docker build -t artalk:TAG .
```

将 `TAG` 替换为所需的标签名（例如 `latest`）。

### 跳过前端构建

如果你已经在 Docker 容器外部构建了前端，可以跳过容器内的前端构建过程以加快构建速度。使用以下命令：

```sh
docker build --build-arg SKIP_UI_BUILD=true -t artalk:latest .
```

更多详细信息，请参考 `Dockerfile`。

## 测试

### 测试后端

要测试后端 Go 程序，可以使用以下命令：

- **单元测试**：
  运行 `make test` 执行单元测试。结果将输出在终端。

- **测试覆盖率**：
  运行 `make test-coverage` 检查代码测试覆盖率。

- **HTML 覆盖率报告**：
  运行 `make test-coverage-html` 生成并浏览 HTML 格式的测试覆盖率报告。

- **运行特定测试**：
  要运行特定测试，请使用带有环境变量的命令：
  ```sh
  TEST_PATHS="./server/..." make test
  ```

### 测试前端

前端测试包括单元测试和端到端（E2E）测试。

- **单元测试**：
  前端的单元测试使用 Vitest 进行。运行 `pnpm test` 开始单元测试。

- **端到端测试**：
  E2E 测试使用 Playwright 进行。要启动 E2E 测试，运行：
  ```sh
  make test-frontend-e2e
  ```

### 持续集成测试

前端和后端的测试都是自动化的，并将在 Git pull 请求和构建过程中使用 CircleCI 或 GitHub Actions 执行。

## 配置

模板配置文件位于 `/conf` 目录中，模板以多种语言注释，命名格式为 `artalk.example.[lang].yml`（例如，`artalk.example.zh-CN.yml`）。Artalk 将解析这些模板配置文件以生成设置界面、环境变量名称和文档。为了提高性能，程序运行期间将缓存部分数据以消除解析时间。

在 `internal/config/config.go` 文件中，有一个名为 `Config` 的结构体定义，该结构体用于将 yml 文件解析为 Go 结构体。该结构体包含更精确的类型定义。如果需要添加、修改或删除配置项，必须同时修改 `/conf` 目录中的配置文件模板和 `Config` 结构体。

当你修改配置文件后，请运行以下命令以更新配置数据缓存：

```sh
make update-conf
```

要更新环境变量文档，可以运行以下命令：

```sh
make update-conf-docs
```

## 文档

文档由三个部分组成，每个部分与特定包相关：**指南文档**、**着陆页**和 **Swagger API 文档**。下面是每个包的相关命令及其输出的总结。

### 指南文档 (`docs`)

- **目录**：`docs/docs/guide`
- **开发命令**：`pnpm -F docs dev:docs`
- **构建命令**：`pnpm -F docs build:docs`
- **输出目录**：`docs/docs/.vitepress/dist`

### 着陆页 (`docs-landing`)

- **目录**：`docs/landing`
- **开发命令**：`pnpm -F docs-landing dev:landing`
- **构建命令**：`pnpm -F docs-landing build:landing`
- **输出目录**：`docs/landing/dist`

### Swagger API 文档 (`docs-swagger`)

- **目录**：`docs/swagger`
- **构建命令**：`make update-swagger`
- **输出目录**：`docs/swagger`

Swagger 定义位于后端代码的 `/server` 目录中。修改 Swagger 定义后，运行构建命令将更新 Swagger 文档并在 `/ui/artalk/src/api/v2.ts` 生成 HTTP 客户端代码。

### 综合构建

- **构建所有包的命令**：`pnpm build:docs`
- **综合输出目录**：`docs/docs/.vitepress/dist`

此命令构建并组合所有三个包到单个目录中，准备发布。

### 环境变量文档

- **文件**：`docs/docs/guide/env.md`
- **更新命令**：`make update-conf-docs`

此命令读取位于 `/conf` 的配置模板文件以更新环境变量文档。在修改配置模板文件后运行此命令。

## 翻译 (i18n)

如果你编写了新功能或进行了修复/重构，请使用以下命令通过解析源代码中的 `i18n.T` 函数调用来增量生成翻译文件：

```sh
make update-i18n
```

如果你不是程序员，但希望帮助改进翻译，可以直接编辑 `/i18n` 目录中的翻译文件，然后提交 pull 请求。

## 结尾

感谢你耐心阅读本指南，并感谢你对 Artalk 的兴趣和支持。我们深知，一个开源项目的成功依赖于每一位开发者的贡献。无论是通过代码、文档、测试还是反馈，你的参与是推动项目持续进步的动力。我们诚挚邀请你加入我们，共同改进和增强 Artalk。

通过合作，我们可以创建一个更加健壮和多功能的平台，造福整个社区。你的独特见解和专业知识是无价的，只有共同努力，我们才能实现显著的进步。无论你是经验丰富的开发者，还是刚刚起步的新手，我们的项目中总有你可以贡献的一席之地。

如果你有任何问题或需要帮助，请随时在我们的 GitHub 页面上的 issue 部分提问：<https://github.com/ArtalkJS/Artalk/issues>。再次感谢你的支持，我们期待与你合作。

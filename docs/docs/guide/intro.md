# <span class="wave">👋</span> Hello Friend

Artalk 是一款简单易用但功能丰富的评论系统，你可以开箱即用地部署并置入任何博客、网站、Web 应用。

![](https://user-images.githubusercontent.com/22412567/141147152-df30a0ff-bf41-42ee-9958-4722206a7a90.png)

## 功能亮点

**轻量设计**

前端采用 TypeScript (Vanilla JS)，轻量级，无冗余依赖，仅 ~40KB (gzipped)。

后端采用 Golang 重制 (Artalk v2)，跨平台，体积小巧，五脏俱全，快速部署。

**“麻雀虽小，五脏俱全”**

Artalk 的功能包括但不限于：

<!-- features -->
* [侧边栏](https://artalk.js.org/guide/frontend/sidebar.html): 快速管理、直观浏览
* [社交登录](https://artalk.js.org/guide/frontend/auth.html): 通过社交账号快速登录
* [邮件通知](https://artalk.js.org/guide/backend/email.html): 多种发送方式、邮件模板
* [多元推送](https://artalk.js.org/guide/backend/admin_notify.html): 多种推送方式、通知模版
* [站内通知](https://artalk.js.org/guide/frontend/sidebar.html): 红点标记、提及列表
* [验证码](https://artalk.js.org/guide/backend/captcha.html): 多种验证类型、频率限制
* [评论审核](https://artalk.js.org/guide/backend/moderator.html): 内容检测、垃圾拦截
* [图片上传](https://artalk.js.org/guide/backend/img-upload.html): 自定义上传、支持图床
* [Markdown](https://artalk.js.org/guide/intro.html): 支持 Markdown 语法
* [表情包](https://artalk.js.org/guide/frontend/emoticons.html): 兼容 OwO，快速集成
* [多站点](https://artalk.js.org/guide/backend/multi-site.html): 站点隔离、集中管理
* [管理员](https://artalk.js.org/guide/backend/multi-site.html): 密码验证、徽章标识
* [页面管理](https://artalk.js.org/guide/frontend/sidebar.html): 快速查看、标题一键跳转
* [浏览量统计](https://artalk.js.org/guide/frontend/pv.html): 轻松统计网页浏览量
* [层级结构](https://artalk.js.org/guide/frontend/config.html#nestmax): 嵌套分页列表、滚动加载
* [评论投票](https://artalk.js.org/guide/frontend/config.html#vote): 赞同或反对评论
* [评论排序](https://artalk.js.org/guide/frontend/config.html#listsort): 多种排序方式，自由选择
* [评论搜索](https://artalk.js.org/guide/frontend/sidebar.html): 快速搜索评论内容
* [评论置顶](https://artalk.js.org/guide/frontend/sidebar.html): 重要消息置顶显示
* [仅看作者](https://artalk.js.org/guide/frontend/config.html): 仅显示作者的评论
* [评论跳转](https://artalk.js.org/guide/intro.html): 快速跳转到引用的评论
* [自动保存](https://artalk.js.org/guide/frontend/config.html): 输入内容防丢功能
* [IP 属地](https://artalk.js.org/guide/frontend/ip-region.html): 用户 IP 属地展示
* [数据迁移](https://artalk.js.org/guide/transfer.html): 自由迁移、快速备份
* [图片灯箱](https://artalk.js.org/guide/frontend/lightbox.html): 图片灯箱快速集成
* [图片懒加载](https://artalk.js.org/guide/frontend/img-lazy-load.html): 延迟加载图片，优化体验
* [Latex](https://artalk.js.org/guide/frontend/latex.html): Latex 公式解析集成
* [夜间模式](https://artalk.js.org/guide/frontend/config.html#darkmode): 夜间模式切换
* [扩展插件](https://artalk.js.org/develop/plugin.html): 创造更多可能性
* [多语言](https://artalk.js.org/guide/frontend/i18n.html): 多国语言切换
* [命令行](https://artalk.js.org/guide/backend/config.html): 命令行操作管理能力
* [API 文档](https://artalk.js.org/http-api.html): 提供 OpenAPI 格式文档
* [程序升级](https://artalk.js.org/guide/backend/update.html): 版本检测，一键升级
<!-- /features -->

穷举不是我们的特长，更多有趣的功能期待你来探索！

**“Unlimited Blade Works”**

Artalk 正在持续成长，创意由你发挥，价值由你赋予！

不论是 Vue、React、Svelte 的前端项目，还是 WordPress、Typecho、Hexo 等博客系统，都可以快速引入 Artalk，结合诸位的聪明才智，我们相信 Artalk 能够自如应对各种业务场景。

## 用户体验

我们相信优雅的设计能带来良好的用户体验，良好的用户体验能帮助项目走得更远。

「平凡而不平庸的设计」倍受专业 UI 设计师青睐的设计工具 Figma 这次在 Artalk 的重新设计中也帮了大忙。我们预先使用 Figma 设计人性化、现代化的界面，再编写前端样式使其自然融合至现代化的网站中，简约清新的界面由此诞生。此外，我们还设计了用户身份认证徽章、评论平铺 / 无限嵌套模式、评论分页等样式，同时兼顾不同尺寸的设备，在有限的空间体验无限的内容。

「崩溃就在一瞬间」对于不加优化的评论系统，用户每次评论可能需要反复输入个人信息，发生意外状况时辛苦键入的见解还可能完全丢失。要知道，成年人的崩溃只在一瞬间。为解决这些痛点，Artalk 借助浏览器缓存自动填充用户信息、自动保存评论数据，让用户以最少的成本发表见解。

「丰富站点表情，重燃评论热情」千篇一律的表情包可能容易使访客丧失评论的热情，于是 Artalk 自带一套精心挑选的滑稽表情包。除此之外，Artalk 也支持自定义图片表情。

「你所热爱的，就是你的生活？」用户体验不仅仅就访客而言，对于站点管理者，Artalk 也不乏人性化的设计。通过侧边栏集成管理[控制中心](./frontend/sidebar.md#控制中心)，管理员用户可以方便快捷地管理名下多个站点，所有数据通过规范化 API 交流并且异步处理，减少数据处理阻塞，降低服务资源占用。针对可能出现的垃圾评论，Artalk 支持自动拦截，降低管理者工作强度，也还站点以清净。

我们希望 Artalk 不仅能实现评论系统应有的基础功能，更能成为搭建 **知识传播者和知识学习者交流思想** 桥梁的媒介，让知识不再局限于文本，帮助知识传播者创造其应有的价值。

## 社区理念

“**化繁为简，简而不凡**”

Artalk 的目标是在尽量 **简洁** 的前提下，实现 **丰富** 的功能。

2018 年 10 月 2 日，Artalk 的 [第一行代码](https://github.com/ArtalkJS/Artalk/commit/66128e2c8d9a8ac00a8d1498ff0ec035a7727daf) 被推送至 GitHub，直至 2021 年 10 月 20 日，才发布了 v2 版本。由于团队成员较少且开发者时间并不充裕，项目整体发展较慢。我们非常需要社区的力量，无论是为项目反馈 Bug，还是提供新功能的创意，我们都十分期待。

Artalk 社区是包容开放的社区，我们欢迎不同水平的人员帮助 / 参与项目开发。如果你是入门新手，除了积极学习项目相关知识外，你也可以尝试体验已有 Artalk 部署，在使用中寻找、确认 Artalk 的不足之处，复现、总结后在相关项目的 [Issues](https://github.com/ArtalkJS/Artalk/issues)、[Discussions](https://github.com/ArtalkJS/Artalk/discussions) 中发表相关讨论，帮助开发者更好地定位问题、更快地做出优化。如果你是颇有技术的开发人员，你可以在 [@ArtalkJS](https://github.com/ArtalkJS) 找到项目的所有源码，结合此文档，我们相信这也许不难理解。无论是优化前后端结构、实现全新功能还是编写社区项目，我们都期待 Artalk 汇入新鲜血液。

“More action, less talk”，Artalk 社区不欢迎无意义的争论，我们希望社区成员和谐共处、为社区发展出谋划策。在提出问题前，你应当读过《[提问的智慧](https://lug.ustc.edu.cn/wiki/doc/smart-questions/)》，这可能决定了你最终是否能得到有用的回答。在表达观点前，你应当具备基本的礼仪，比如保持平和的态度、使用得体的语言，切忌恶语相向、冷嘲热讽、不尊重他人信仰和立场等。

我们作为开源精神的推崇者以及实践者，希望我们所创造的自由软件，都应该被自由的使用，自由的研究，自由的更改和自由的分享。本项目主程序使用 [MIT](https://github.com/ArtalkJS/Artalk/blob/master/LICENSE) 协议开源，文档使用 [CC](https://creativecommons.org/licenses/by-nc-sa/4.0/deed.zh) 许可协议。
::: tip 立即为社区贡献力量？

- 浏览开发者资料（[开发文档](../develop/index.md) / [CONTRIBUTING.md](https://github.com/ArtalkJS/Artalk/blob/master/CONTRIBUTING.md)）
- 维护 Artalk 后端（代码仓库 [@ArtalkJS/Artalk:/](https://github.com/ArtalkJS/Artalk)）
- 维护 Artalk 前端（代码仓库 [@ArtalkJS/Artalk:/ui](https://github.com/ArtalkJS/Artalk/tree/master/ui)）
- 完善 Artalk 文档（代码仓库 [@ArtalkJS/Artalk:/docs](https://github.com/ArtalkJS/Artalk/tree/master/docs)）
- 翻译 多语言 i18n（前往查看 [多语言说明](./frontend/i18n.html)）
- 改进数据迁移工具（代码仓库 [@ArtalkJS/Artransfer](https://github.com/ArtalkJS/Artransfer)）
- 分享你的想法创意（下方留言 / [Discussions](https://github.com/ArtalkJS/Artalk/discussions)）
- 编写相关社区项目（扩展插件 / 部署教程等）

:::

## 写在结尾

至此，相信你已经了解 Artalk 的基本情况。无论你是否选择 Artalk，我们都十分感谢你对 Artalk 的关注。如果 Artalk 尚未满足你的需求，希望你能提出一针见血的建议帮助 Artalk 成长。

欢迎使用 Artalk，

起飞！🛫️

# 安全防范

## 用户身份凭证存储

RESTful API 具有无状态 (Stateless) 这一特性，尽管 Artalk 的 API 并不属于 RESTful，但保证无状态这一特性能让 API 更灵活易用，相反使用 Session 会话技术在处理多站点、跨域等操作的时候尤为繁琐，对于普通用户来说，服务器配置也成为麻烦，不当的配置会让用户的安全受到威胁。

在各种权衡考虑之下，Artalk 目前还是依然使用 localStorage 保存用户凭证数据，因为使用 localStorage 能保证 API 的无状态性，也无需考虑因使用 Cookie 而带来的各种麻烦，而且使用 Cookie 并不能换来多大的安全性提升。无论是 Cookie 还是 localStorage 都存在有一些缺陷。

Cookie 有 CSRF 安全隐患，而无状态 API 不容易有。如果使用 Cookie，就需要在任何有身份验证的地方，额外考虑对 CSRF 安全问题的防范，虽然可以使用全局 CSRF_TOKEN，但这无疑增加了 API 和客户端的复杂程度，开发也很难保证永远不会遗漏一些细节。所以，我认为使用 Cookie 既引入了一些额外的安全隐患，并且程序的开发难度也变大了，为什么要这样折磨自己呢？

还有一点是 Cookie 这种被广泛用作判断用户身份状态的技术特别容易被跟踪滥用，尽管可信的现代浏览器有很多内部的安全策略，但如果用户操作系统安装了一些恶意流氓软件，拥有较高权限，你无法保证某些软件窃取 Cookie 不是轻而易举的事 (特指某些大厂)。

那，Cookie 有啥优点？使用 Cookie 的主要优点在于能防止客户端 (浏览器) 通过 JS 对用户凭证数据的「直接」访问，阻止恶意 XSS 攻击脚本「直接」窃取数据的可能性。注意，这里用了「直接」这一个修饰词。虽然使用 Cookie (httpOnly) 不能通过 JS 读取凭证数据，但 XSS 攻击脚本依然可以在客户端执行 fetch 代码向 API 发起请求，这些请求依然是携带了 Cookie 的，是具有用户身份的 API 访问。也就是说，使用 Cookie (httpOnly) 并不能从根源上解决安全问题，不管你是 Cookie 还是 localStorage，只要网站存在 XSS 漏洞，攻击者成功注入了 XSS 脚本，都会对用户的安全造成威胁。所以使用 Cookie (httpOnly) 防 XSS 我认为是个「伪命题」。

XSS 是属于非常严重的问题了。所以无论使用何种方式存放敏感数据，XSS 防范才是重中之重。

对于我们用户来说，不应随意引入不明的 JS 脚本；不应在不了解代码意图的情况下到浏览器控制台执行 JS 代码；不应随意点击来历不明的链接；尽量使用新开的隐身模式窗口打开不明链接。

Artalk 在努力保护用户的安全。参考 [OWASP 安全备忘录](https://cheatsheetseries.owasp.org/cheatsheets/DOM_based_XSS_Prevention_Cheat_Sheet.html) 中提到的内容，Artalk 的开发对于来自用户输入的任何数据始终持以怀疑的态度，我们严格遵循以下准则防范 XSS 攻击：

- 不直接或间接地调用 `innerHTML` 输出未经处理的用户输入数据，而是使用 `innerText`。
- 需要额外注意将用户输入数据带入程序上下文时的处理，防止非法数据在页面被执行。例如不要在 `setAttribute` 时直接引用来自用户的输入值。
- 注意类型合法性判断、做好类型的转换。例如 URL 以 `javascript:` 前缀开头应被视为非法。
- 不要去调用一些高风险的内置函数，因为会带来潜在的安全隐患。例如：`eval()`。
- 注意来自内置 API 的数据，应任被视为不可信的。例如：`location.hash.split("#")[1]`。

## 跨域请求安全性

参考：[Cross-Origin Resource Sharing (CORS) ｜ MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)

【待补充】

## API 设计安全性

【待补充】

## 写在结尾

Artalk 的功能在不断完善，同时也在不断向提高安全性的目标努力，但在努力的同时还是难免发生一些疏忽和遗漏，若遇问题或有好的建议，欢迎反馈与指正。

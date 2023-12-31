# 安全防范

## 程序运行安全性

有以下一些安全措施可以提高服务器的安全性：

1. 确保 Artalk 运行在非 `root` 用户或 Docker 容器中，让 Artalk 运行在较低的权限以及隔离的环境中，以降低程序被攻击后的危害。
2. 保持 Artalk 的版本为最新，新版可能修复了一些安全漏洞，以及提升了安全性。
3. 确保 Artalk 运行在 HTTPS 协议下，以防止数据被中间人窃取。
4. 设置强度足够的密码，并配置好验证码等安全措施，以防止密码被破解。
5. 设置 `trusted_domains`，以防止恶意的跨域请求。
6. 定期备份数据库，以防止数据丢失。

## 程序设计安全性

### 用户身份凭证存储

Artalk 在努力保护用户的安全。参考 [OWASP 安全备忘录](https://cheatsheetseries.owasp.org/cheatsheets/DOM_based_XSS_Prevention_Cheat_Sheet.html) 中提到的内容，Artalk 的开发对于来自用户输入的任何数据始终持以怀疑的态度，我们严格遵循以下准则防范 XSS 攻击：

- 不直接或间接地调用 `innerHTML` 输出未经处理的用户输入数据，而是使用 `innerText`。
- 需要额外注意将用户输入数据带入程序上下文时的处理，防止非法数据在页面被执行。例如不要在 `setAttribute` 时直接引用来自用户的输入值。
- 注意类型合法性判断、做好类型的转换。例如 URL 以 `javascript:` 前缀开头应被视为非法。
- 不要去调用一些高风险的内置函数，因为会带来潜在的安全隐患。例如：`eval()`。
- 注意来自内置 API 的数据，应任被视为不可信的。例如：`location.hash.split("#")[1]`。

### 跨域请求安全性

参考：[Cross-Origin Resource Sharing (CORS) ｜ MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)

【待补充】

### API 设计安全性

【待补充】

## 写在结尾

Artalk 的功能在不断完善，同时也在不断向提高安全性的目标努力，但难免存在可能的问题，如果您发现了 Artalk 的安全漏洞，请通过 [GitHub Issues](https://github.com/ArtalkJS/Artalk/issues) 或者邮箱 artalkjs@gmail.com 快速向我们取得联系，我们会尽最大的努力在第一时间解决！感谢您的支持！

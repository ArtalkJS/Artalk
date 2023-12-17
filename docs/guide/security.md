# 安全防范

## 用户身份凭证存储

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

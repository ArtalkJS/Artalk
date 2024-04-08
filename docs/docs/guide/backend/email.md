# 邮件通知

Artalk 支持 SMTP 协议、阿里云邮件推送、调用系统 sendmail 命令等方式向用户发送邮件通知。

你可以在[控制中心](/guide/frontend/sidebar.md#控制中心)找到「设置」界面修改此配置。

## 配置文件

完整的 `email` 配置如下：

```yaml
# 邮件通知
email:
  enabled: false # 总开关
  send_type: smtp # 发送方式 [smtp, ali_dm, sendmail]
  send_name: '{{reply_nick}}' # 发信人昵称
  send_addr: example@qq.com # 发信人地址
  mail_subject: '[{{site_name}}] 您收到了来自 @{{reply_nick}} 的回复'
  mail_tpl: default # 邮件模板文件
  smtp:
    host: smtp.qq.com
    port: 587
    username: example@qq.com
    password: ''
  ali_dm: # https://help.aliyun.com/document_detail/29444.html
    access_key_id: '' # 阿里云颁发给用户的访问服务所用的密钥 ID
    access_key_secret: '' # 用于加密的密钥
    account_name: example@example.com # 管理控制台中配置的发信地址
```

### 选择发件方式

配置项 `enabled` 启用邮件，`send_type` 用于选择发送方式，可选：`smtp`, `ali_dm`, `sendmail`。

```yaml
email:
  enabled: true
  send_type: smtp # 发送方式
  # 省略其他配置...
  smtp:
    # SMTP 配置...
  ali_dm:
    # 阿里云推送配置...
```

### SMTP 配置

```yaml
# 邮件通知
email:
  enabled: true
  send_type: smtp # 选择 smtp
  smtp:
    host: smtp.qq.com
    port: 587
    username: example@qq.com
    password: ''
```

### 阿里云推送配置

```yaml
email:
  enabled: true
  send_type: ali_dm # 选择 ali_dm
  ali_dm:
    access_key_id: '' # 阿里云颁发给用户的访问服务所用的密钥 ID
    access_key_secret: '' # 用于加密的密钥
    account_name: example@example.com # 管理控制台中配置的发信地址
```

可参考：[阿里云官方文档](https://help.aliyun.com/document_detail/29444.html)

## 评论回复

邮件中会有一个评论回复按钮，该链接指向前端给定的页面 PageKey，若你提供的 `pageKey` 配置项为页面的「相对路径」，你需要在「[控制中心](../frontend/sidebar.md#控制中心)」-「站点」为你的站点设置一个 URL：

<img src="/images/sidebar/site_url.png" width="400px">

## 邮件模板

### 模板变量

你可以在 `mail_subject` 和 `mail_subject_to_admin` 以及邮件模板文件中使用模板变量：

```yaml
email:
  # 省略其他配置...
  mail_subject: '[{{site_name}}] 您收到了来自 @{{reply_nick}} 的回复'
  mail_subject_to_admin: '[{{site_name}}] 您的文章 "{{page_title}}" 有新回复'
  mail_tpl: default # 邮件模板文件
```

变量是 “Mustache” 的语法，`双大括号` + `变量名` 的形式即可输出一个变量：

**基本内容变量**

```yaml
{{content}}        # 评论内容
{{link_to_reply}}  # 回复链接
{{nick}}           # 评论者昵称
{{page_title}}     # 页面标题
{{page_url}}       # 页面 PageKey (URL)
{{reply_content}}  # 回复对象的内容
{{reply_nick}}     # 回复对象的昵称f
{{site_name}}      # 站点名
{{site_url}}       # 站点 URL
```

::: details 查看其他变量

```yaml
# 评论创建者
{{comment.badge_color}}
{{comment.badge_name}}
{{comment.content}}
{{comment.content_raw}}
{{comment.date}}
{{comment.datetime}}
{{comment.email}}
{{comment.email_encrypted}}
{{comment.id}}
{{comment.is_allow_reply}}
{{comment.is_collapsed}}
{{comment.is_pending}}
{{comment.link}}
{{comment.nick}}
{{comment.page.admin_only}}
{{comment.page.id}}
{{comment.page.key}}
{{comment.page.site_name}}
{{comment.page.title}}
{{comment.page.url}}
{{comment.page.vote_down}}
{{comment.page.vote_up}}
{{comment.page_key}}
{{comment.page_title}}
{{comment.rid}}
{{comment.site.first_url}}
{{comment.site.id}}
{{comment.site.name}}
{{comment.site.urls.0}}
{{comment.site.urls_raw}}
{{comment.site_name}}
{{comment.time}}
{{comment.ua}}
{{comment.visible}}
{{comment.vote_down}}
{{comment.vote_up}}

# 父评论（评论创建者回复的评论）
{{parent_comment.badge_color}}
{{parent_comment.badge_name}}
{{parent_comment.content}}
{{parent_comment.content_raw}}
{{parent_comment.date}}
{{parent_comment.datetime}}
{{parent_comment.email}}
{{parent_comment.email_encrypted}}
{{parent_comment.id}}
{{parent_comment.is_allow_reply}}
{{parent_comment.is_collapsed}}
{{parent_comment.is_pending}}
{{parent_comment.link}}
{{parent_comment.nick}}
{{parent_comment.page.admin_only}}
{{parent_comment.page.id}}
{{parent_comment.page.key}}
{{parent_comment.page.site_name}}
{{parent_comment.page.title}}
{{parent_comment.page.url}}
{{parent_comment.page.vote_down}}
{{parent_comment.page.vote_up}}
{{parent_comment.page_key}}
{{parent_comment.page_title}}
{{parent_comment.rid}}
{{parent_comment.site.first_url}}
{{parent_comment.site.id}}
{{parent_comment.site.name}}
{{parent_comment.site.urls}}
{{parent_comment.site.urls_raw}}
{{parent_comment.site_name}}
{{parent_comment.time}}
{{parent_comment.ua}}
{{parent_comment.visible}}
{{parent_comment.vote_down}}
{{parent_comment.vote_up}}
```

:::

::: details 查看数据样例

```json
{
  "content": "<p>测试内容</p>\n",
  "link_to_reply": "https://artalk.js.org/?atk_comment=8100&atk_notify_key=44a9b2f08312565fba47c716df9d177f",
  "nick": "用户名",
  "page_title": "Artalk",
  "page_url": "https://artalk.js.org/",

  "reply_content": "<p>回复者内容</p>\n",
  "reply_nick": "回复者",
  "site_name": "ArtalkDemo",
  "site_url": "http://localhost:3000/",

  "comment.badge_color": "",
  "comment.badge_name": "",
  "comment.content": "<p>回复者内容</p>\n",
  "comment.content_raw": "回复者内容",
  "comment.date": "2021-11-22",
  "comment.datetime": "2021-11-22 22:22:42",
  "comment.email": "replyer@example.com",
  "comment.email_encrypted": "249898bd50e0febc5799485cf10b345a",
  "comment.id": 8100,
  "comment.is_allow_reply": true,
  "comment.is_collapsed": false,
  "comment.is_pending": false,
  "comment.link": "",
  "comment.nick": "回复者",
  "comment.page.admin_only": false,
  "comment.page.id": 75,
  "comment.page.key": "https://artalk.js.org/",
  "comment.page.site_name": "ArtalkDemo",
  "comment.page.title": "Artalk",
  "comment.page.url": "https://artalk.js.org/",
  "comment.page.vote_down": 0,
  "comment.page.vote_up": 0,
  "comment.page_key": "https://artalk.js.org/",
  "comment.page_title": "Artalk",
  "comment.rid": 8099,
  "comment.site.first_url": "http://localhost:3000/",
  "comment.site.id": 2,
  "comment.site.name": "ArtalkDemo",
  "comment.site.urls.0": "http://localhost:3000/",
  "comment.site.urls_raw": "http://localhost:3000/",
  "comment.site_name": "ArtalkDemo",
  "comment.time": "22:22:42",
  "comment.ua": "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
  "comment.visible": false,
  "comment.vote_down": 0,
  "comment.vote_up": 0,

  "parent_comment.badge_color": "",
  "parent_comment.badge_name": "",
  "parent_comment.content": "<p>测试内容</p>\n",
  "parent_comment.content_raw": "测试内容",
  "parent_comment.date": "2021-11-22",
  "parent_comment.datetime": "2021-11-22 22:21:17",
  "parent_comment.email": "test@example.com",
  "parent_comment.email_encrypted": "55502f40dc8b7c769880b10874abc9d0",
  "parent_comment.id": 8099,
  "parent_comment.is_allow_reply": true,
  "parent_comment.is_collapsed": false,
  "parent_comment.is_pending": false,
  "parent_comment.link": "https://qwqaq.com",
  "parent_comment.nick": "用户名",
  "parent_comment.page.admin_only": false,
  "parent_comment.page.id": 75,
  "parent_comment.page.key": "https://artalk.js.org/",
  "parent_comment.page.site_name": "ArtalkDemo",
  "parent_comment.page.title": "Artalk",
  "parent_comment.page.url": "https://artalk.js.org/",
  "parent_comment.page.vote_down": 0,
  "parent_comment.page.vote_up": 0,
  "parent_comment.page_key": "https://artalk.js.org/",
  "parent_comment.page_title": "Artalk",
  "parent_comment.rid": 0,
  "parent_comment.site.first_url": "http://localhost:3000/",
  "parent_comment.site.id": 2,
  "parent_comment.site.name": "ArtalkDemo",
  "parent_comment.site.urls.0": "http://localhost:3000/",
  "parent_comment.site.urls_raw": "http://localhost:3000/",
  "parent_comment.site_name": "ArtalkDemo",
  "parent_comment.time": "22:21:17",
  "parent_comment.ua": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36",
  "parent_comment.visible": false,
  "parent_comment.vote_down": 0,
  "parent_comment.vote_up": 0
}
```

:::

### 自定义模板

你可以将 `mail_tpl` 设置为一个「具体的文件路径」，使用外部的自定义邮件模板。

例如，将 `mail_tpl` 配置为 `"/root/Artalk/data/mail_tpl/your_email_template.html"`

```yaml
email:
  mail_tpl: /root/Artalk/data/mail_tpl/your_email_template.html
  # 其他配置省略...
```

那么，在这个路径应该有一个文件：

```html
<div>
  <p>Hi, {{nick}}：</p>
  <p>您在 “{{page_title}}” 收到了回复：</p>
  <div>
    <div>@{{reply_nick}}:</div>
    <div>{{reply_content}}</div>
  </div>
  <p><a href="{{link_to_reply}}" target="_blank">回复消息 »</a></p>
</div>
```

Artalk 内置许多预设的邮件模板，例如 `mail_tpl: "default"` 使用的就是：[@ArtalkJS/Artalk:/internal/template/email_tpl/default.html](https://github.com/ArtalkJS/Artalk/blob/master/internal/template/email_tpl/default.html)

## 发向管理员的邮件

邮件通知目标为管理员和普通用户，你可通过如下配置，为发向管理员的邮件设定不同的标题：

```yaml
admin_notify:
  enabled: true
  mail_subject: '[{{site_name}}] 您的文章「{{page_title}}」有新回复'
```

注：旧版 `email.mail_subject_to_admin` 配置项已弃用，请使用以上替代。

不局限于邮件，Artalk 支持多种方式向管理员发送通知，参考：[多元推送](./admin_notify.md#邮件通知)。

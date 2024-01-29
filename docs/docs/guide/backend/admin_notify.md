# 多元推送

Artalk 支持通过多元推送功能以多种方式发送管理员通知。

支持 **Telegram**、**飞书**、**钉钉**、**Bark**、**Slack**、**LINE**，并且多种方式可以同时启用。

你可以在[控制中心](/guide/frontend/sidebar.md#控制中心)找到「设置」界面修改此配置。

## 配置文件

完整的 `admin_notify` 配置如下：

::: details 点击显示

```yaml
# 多元推送
admin_notify:
  # 通知模版
  notify_tpl: "default"
  noise_mode: false
  # 邮件通知管理员
  email:
    enabled: true # 当使用其他推送方式时，可以关闭管理员邮件通知
    mail_subject: "[{{site_name}}] 您的文章「{{page_title}}」有新回复"
    mail_tpl: ""
  # Telegram
  telegram:
    enabled: false
    api_token: ""
    receivers:
      - 7777777
  # 飞书
  lark:
    enabled: false
    webhook_url: ""
  # 钉钉
  ding_talk:
    enabled: false
    token: ""
    secret: ""
  # Bark
  bark:
    enabled: false
    server: "http://day.app/xxxxxxx/"
  # Slack
  slack:
    enabled: false
    oauth_token: ""
    receivers:
      - "CHANNEL_ID"
  # LINE
  line:
    enabled: false
    channel_secret: ""
    channel_access_token: ""
    receivers:
      - "USER_ID_1"
      - "GROUP_ID_1"
  # WebHook
  webhook:
    enabled: false
    url: ""
```

:::

## 邮件通知

通过邮件的方式向管理员发送消息通知。

```yaml
admin_notify:
  enabled: true # 当使用其他推送方式时，可以关闭管理员邮件通知
  mail_subject: "[{{site_name}}] 您的文章「{{page_title}}」有新回复"
  mail_tpl: ""
```

当使用其他推送方式时，可以关闭管理员邮件通知。

在这之前，你需要配置全局邮件发送功能：[参考此处](./email.md)。

### 管理员邮件模板

- 配置项 `mail_subject` 为发送给管理员的邮件标题。
- 配置项 `mail_tpl` 为发送给管理员的邮件选用特定的[邮件模板](./email.md#邮件模板)（填写邮件模板文件路径）。
  
  (当该项留空时，将继承 `email.mail_tpl` 配置项)

## Telegram

```yaml
admin_notify:
  # Telegram
  telegram:
    enabled: true
    api_token: ""
    receivers:
      - 7777777
```

- `api_token`：TG Bot 的 API Token。
- `receivers`：消息接受者的数字 ID，可设置多个。

### 创建 TG Bot

搜索 `@BotFather` 回复 `/newbot` 并按提示创建新的 TG 机器人。

![](/images/notify/tg-1.png)

标红的文字就是你之后需要在 Artalk 配置中填入的 `api_token`。

配置中的 `receivers` 填入需要接受消息的账号数字 ID，可以搜索机器人 `@RawDataBot` 获取如图：

![](/images/notify/tg-2.png)

详情可参考：[Bots: An introduction for developers - Telegram](https://core.telegram.org/bots)

::: tip

鉴于复杂的网络环境，如需使用代理，请在 Artalk 启动之前配置环境变量，例如：

```sh
export https_proxy=http://127.0.0.1:7890
```

:::

## 飞书

```yaml
admin_notify:
  # 飞书
  lark:
    enabled: true
    webhook_url: ""
```

- `webhook_url`：填入创建群组机器人时得到的 WebHook 地址。

### 创建群组机器人

点击顶部的加号，创建一个新的群组：

![](/images/notify/lark-1.png)

找到右侧的「群设置」-「群机器人」- 点击「添加机器人」- 选择「自定义机器人」并按照提示创建。

<img src="/images/notify/lark-2.png" width="700px">

复制如上图的 WebHook 地址，并修改 Artalk 的 `webhook_url` 配置即可。

<img src="/images/notify/lark-3.png" width="400px">

可参考：[飞书帮助中心文档](https://www.feishu.cn/hc/zh-CN/articles/360024984973)

## 钉钉

```yaml
admin_notify:
  # 钉钉
  ding_talk:
    enabled: true
    token: ""
    secret: ""
```

可参考：[钉钉开放文档](https://open.dingtalk.com/document/robots/custom-robot-access)

## Bark

```yaml
admin_notify:
  # Bark
  bark:
    enabled: true
    server: "http://day.app/xxxxxxx/"
```

[Bark](https://github.com/Finb/Bark) 是一款开源的 iOS App，并且[支持自托管](https://github.com/Finb/bark-server)，你能使用 Bark 轻松地推送消息给你的 iOS 设备。

你可以在 App Store 搜索下载，并获得需要填入 Artalk 的 `server` 配置项：


<img src="/images/notify/bark.png" width="700px">


## Slack

```yaml
admin_notify:
  # Slack
  slack:
    enabled: true
    oauth_token: ""
    receivers:
      - "CHANNEL_ID"
```

## LINE

```yaml
admin_notify:
  # LINE
  line:
    enabled: true
    channel_secret: ""
    channel_access_token: ""
    receivers:
      - "USER_ID_1"
      - "GROUP_ID_1"
```

## 通知模版

```yaml
admin_notify:
  notify_tpl: "default"
```

配置项 `admin_notify.notify_tpl` 可设置为自定义通知模版「文件路径」，默认的通知模版为：

```
@{{reply_nick}}:

{{reply_content}}

{{link_to_reply}}
```

可用变量和邮件模板相同，可参考：[邮件模版](./email.md#邮件模板)

## 待审评论仍然发送通知 `notify_pending`

```yaml
admin_notify:
  notify_pending: false
```

`notify_pending` 默认为关闭状态，当该项设置为 `false` 时，待审评论不会发送通知。你可以在控制中心查看所有待审核的评论。

## 嘈杂模式 `noise_mode`

```yaml
admin_notify:
  noise_mode: false
```

`noise_mode` 默认为关闭状态，当该项设置为 `false` 时，站内仅向管理员回复的消息会发送通知，例如「普通用户 A」回复「普通用户 B」，这两个用户之间的通讯不会通知管理员。

注：当 `moderator.pending_default` 为 `true` 时，noise_mode 为始终开启状态。

## WebHook 回调

开启 WebHook 后，创建新评论将以 **POST** 方式携带 `application/json` 类型的 Body 数据请求设定的 WebHook 地址。

你可以编写自己的 Server 端代码，处理来自 Artalk 的请求。

**Artalk 配置文件**

```yaml
admin_notify:
    webhook:
      enabled: true
      url: "http://localhost:8080/"
```

**Body 数据内容**

`application/json` 类型

|Key|描述|类型|备注|
| - | - | - | - |
|`notify_subject`|通知标题     |String| 对应 admin_notify.notify_subject 配置项 |
|`notify_body`   |通知内容     |String| 根据 admin_notify.notify_tpl 模版渲染 |
|`comment`       |评论内容     |Object| 新创建的评论数据对象 |
|`parent_comment`|评论回复的目标|Object| 如果是根节点评论值为 null |

**Body 数据样本**

```js
{
  "notify_subject": "",
  "notify_body": "@测试用户:\n\n测试内容\n\nhttps://127.0.0.1/index.html?atk_comment=1057",
  "comment": {
    "id": 1057,
    "content": "测试内容",
    "user_id": 226,
    "nick": "测试用户",
    "email_encrypted": "654236c1e78i4c09a17c4869c9d43910",
    "link": "https://qwqaq.com",
    "ua": "Mozilla/5.0 (Macintosh; Intel Mac OS X 12_4_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Safari/537.36",
    "date": "2022-05-23 17:00:23",
    "is_collapsed": false,
    "is_pending": false,
    "is_pinned": false,
    "is_allow_reply": false,
    "rid": 0,
    "badge_name": "",
    "badge_color": "",
    "visible": true,
    "vote_up": 0,
    "vote_down": 0,
    "page_key": "/index.html",
    "page_url": "https://127.0.0.1/index.html",
    "site_name": "ArtalkDocs"
  },
  "parent_comment": null
}
```

**Node.js express 处理示例**

```js
const express = require('express');

const app = express();

app.use(express.json()); // Use JSON middleware

app.post('/', function(request, response){
  console.log(request.body);

  const notifySubject = request.body.notify_subject
  const notifyBody    = request.body.notify_body
  console.log(notifySubject, notifyBody);

  response.send(request.body);
});

app.listen(8080);
```

**Node.js http 处理示例**

```js
const http = require("http");

const requestListener = function (req, res) {
  // receive json request
  let body = "";
  req.on("data", function (data) {
    body += data;
  });
  req.on("end", function () {
    let json = "";
    try {
      json = JSON.parse(body);
    } catch {}

    // do something with json
    console.log(json);
    res.end();
  });

  res.writeHead(200);
  res.end("Hello, World!");
};

const server = http.createServer(requestListener);
server.listen(8080);
```

**PHP Laravel 处理示例**

```php
Route::get('/', function (Request $request) {
    $data = $request->json()->all();
    $notify_subject = $data["notify_subject"];
    $notify_body    = $data["notify_body"];
});
```

**Golang net/http 处理示例**

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
)

type ArtalkNotify struct {
	NotifySubject string      `json:"notify_subject"`
	NotifyBody    string      `json:"notify_body"`
	Comment       interface{} `json:"comment"`
	ParentComment interface{} `json:"parent_comment"`
}

func webhookHandler(rw http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)
    var notify ArtalkNotify
    err := decoder.Decode(&notify)
    if err != nil {
        panic(err)
    }
    log.Println(notify.NotifyBody)
}

func main() {
    http.HandleFunc("/webhook", webhookHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```


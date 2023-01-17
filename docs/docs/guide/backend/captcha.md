# 验证码

ArtalkGo 内置图片验证码功能，你可以配置操作频率限制，当超过限度时激活验证码。

此外，你也可以接入[极验](https://www.geetest.com/)，拥有一个滑动验证码。

完整的 `captcha` 配置如下：

```yaml
# 验证码
captcha:
  enabled: true    # 总开关
  always: false    # 总是需要验证码
  action_limit: 3  # 激活验证码所需操作次数
  action_reset: 60 # 重置操作计数器超时 (单位：s, 设为 -1 不重置)
  # Geetest 极验
  geetest: # https://www.geetest.com
    enabled: false
    captcha_id: ""
    captcha_key: ""
```

- **always**：当该项为 `true` 时，总是需要输入验证码。
- **action_limit**：激活评论所需的操作次数。
- **action_reset**：当时间超过该值时会重置操作计数器，单位为秒，设为 `-1` 将永不重置。

注：当 `always` 开启时，`action_limit` 和 `action_reset` 配置将失效。

## 配置举例

#### 例 1

在 60s 时间范围内，当操作次数超过 3 次，将一直被要求输入验证码：

```yaml
captcha:
  action_limit: 3
  action_reset: 60
```

在 60s 后将自动重置计数器，即重新获得 3 次不用输入验证码的机会。

#### 例 2

无论多少时间范围内，这个 IP 地址操作次数只要超过 5 次时，将一直被要求输入验证码：

```yaml
captcha:
  action_limit: 5
  action_reset: -1
```

#### 例 3

总是要求输入验证码，无论这个 IP 操作多少次：

```yaml
captcha:
  always: true
```

## 操作的定义

一个 IP 地址的一次「评论、投票、图片上传、密码验证」都算作一次「操作」。

## Geetest 极验

ArtalkGo 支持接入 [Geetest 极验](https://www.geetest.com/adaptive-captcha) 第四代「行为验」，启用极验后，验证码将切换为滑动验证码。

你需要在官网注册账号，并申请获得 `captcha_id` 和 `captcha_key`，并填入配置文件：

```yaml
captcha:
  # 省略其他配置...
  geetest:
    enabled: true
    captcha_id: ""
    captcha_key: ""
```


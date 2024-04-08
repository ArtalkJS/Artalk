# 评论审核

Artalk 支持多种评论过滤方式以拦截垃圾评论，支持通过 Akismet、腾讯云内容安全、阿里云内容安全 等在线 API 实现垃圾评论识别和拦截，也可通过本地离线关键词词库过滤评论，结合[验证码](./captcha.md)功能进一步提高评论审核强度，让垃圾评论没有容身之地。

你可以在[控制中心](/guide/frontend/sidebar.md#控制中心)找到「设置」界面修改此配置。

## 配置文件

完整的 `moderator` 配置如下：

```yaml
# 评论审核
moderator:
  pending_default: false # 发表新评论默认为 “待审状态”
  api_fail_block: false # 垃圾检测 API 请求错误仍然拦截
  # akismet.com 反垃圾
  akismet_key: ''
  # 腾讯云文本内容安全 (tms)
  tencent: # https://cloud.tencent.com/document/product/1124/64508
    enabled: false
    secret_id: ''
    secret_key: ''
    region: ap-guangzhou
  # 阿里云内容安全
  aliyun: # https://help.aliyun.com/document_detail/28417.html
    enabled: false
    access_key_id: ''
    access_key_secret: ''
    region: cn-shanghai
  # 关键词词库过滤
  keywords:
    enabled: false
    pending: false # 匹配成功设为待审状态
    files: # 支持多个词库文件
      - './data/词库_1.txt'
    file_sep: "\n" # 词库文件内容分割符
    replac_to: 'x' # 替换字符
```

## 默认待审模式

开启发表新评论默认为 “待审状态”：

```yaml
moderator:
  pending_default: true
```

## Akismet

[Akismet](https://akismet.com/) 是 WordPress 提供的面向全球范围的老牌垃圾拦截 API，通常对一些英文的垃圾评论十分凑效。Akismet 提供了 Personal 免费版本，适用于个人博客站点。

![](/images/akismet/1.png)

你能在 [Akismet 官网](https://akismet.com/) 轻松地申请 `akismet_key`，并填入配置文件中，即可启用 Akismet 垃圾拦截。

![](/images/akismet/2.png)

```yaml
moderator:
  akismet_key: your_key
```

## 腾讯云文本内容安全

可参考：[腾讯云文档](https://cloud.tencent.com/document/product/1124/64508)

开通「文本内容安全」后，在「访问管理」-「API 密钥管理」新增具有权限的 Secret，然后填入配置：

```yaml
moderator:
  tencent:
    enabled: true
    secret_id: ''
    secret_key: ''
    region: ap-guangzhou
```

## 阿里云内容安全

可参考：[阿里云文档](https://help.aliyun.com/document_detail/28417.html)

开通「阿里云内容安全」后，阿里云后台创建 Access Key 并填入配置：

```yaml
moderator:
  aliyun:
    enabled: true
    access_key_id: ''
    access_key_secret: ''
    region: cn-shanghai
```

## 关键词词库过滤

如果你不想依赖于远程 API，可以在本地配置导入词库文件，让 Artalk 根据词语来检测垃圾评论：

```yaml
moderator:
  keywords:
    enabled: true
    pending: false # 匹配成功设为待审状态
    files: # 支持多个词库文件
      - ./data/词库_1.txt
    file_sep: "\n" # 词库文件内容分割符
    replac_to: 'x' # 替换字符
```

- **pending**：当成功匹配时，是否将评论设为待审核状态。
- **files**：词库文件。允许多个文件，Artalk 启动时会合并词库。
- **file_sep**：词库文件内容分割符。例如：文件中每行一个词语，该项配置 `\n`。
- **replac_to**：替换字符。例如：该项设置为 `x`，你可以将 `pending` 设置为 `false`，评论自动过审，但匹配到的词语会被替换为 `x`，例如 `fxxk`、`xxxx`。

注：`replac_to` 不建议使用 `*` 星号，应为它和 Markdown 的加粗语法冲突。

## 使用验证码

你可以开启 Artalk 的验证码功能，支持图片和滑动验证码，[参考此处](./captcha.md)。

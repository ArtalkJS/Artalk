# Comment Moderation

Artalk supports multiple comment filtering methods to intercept spam comments. It offers integration with online APIs like Akismet, Tencent Cloud Content Safety, Alibaba Cloud Content Safety, and local offline keyword libraries. You can further enhance comment moderation by combining it with the [captcha](./captcha.md) feature, making it hard for spam comments to get through.

You can modify this configuration in the settings interface of the [Dashboard](../frontend/sidebar.md#settings), or configure it through [configuration files](./config.md#comment-moderation-moderator) or [environment variables](../env.md#comment-moderation).

## Configuration File

The complete `moderator` configuration is as follows:

```yaml
# Comment Moderation
moderator:
  pending_default: false # New comments default to "pending review" status
  api_fail_block: false # Intercept comments even if the spam detection API request fails
  # akismet.com anti-spam
  akismet_key: ''
  # Tencent Cloud Text Content Safety (tms)
  tencent: # https://cloud.tencent.com/document/product/1124/64508
    enabled: false
    secret_id: ''
    secret_key: ''
    region: ap-guangzhou
  # Alibaba Cloud Content Safety
  aliyun: # https://help.aliyun.com/document_detail/28417.html
    enabled: false
    access_key_id: ''
    access_key_secret: ''
    region: cn-shanghai
  # Keyword library filtering
  keywords:
    enabled: false
    pending: false # Set to pending review if matched
    files: # Support multiple keyword files
      - './data/keyword_1.txt'
    file_sep: "\n" # Keyword file content separator
    replace_to: 'x' # Replacement character
```

## Default Pending Mode

Enable new comments to default to "pending review" status:

```yaml
moderator:
  pending_default: true
```

## Akismet

[Akismet](https://akismet.com/) is a globally recognized anti-spam API provided by WordPress, particularly effective against English spam comments. Akismet offers a free Personal version suitable for personal blog sites.

![](/images/akismet/1.png)

You can easily apply for an `akismet_key` on the [Akismet official website](https://akismet.com/) and fill it in the configuration file to enable Akismet spam filtering.

![](/images/akismet/2.png)

```yaml
moderator:
  akismet_key: your_key
```

## Tencent Cloud Text Content Safety

Refer to: [Tencent Cloud Documentation](https://cloud.tencent.com/document/product/1124/64508)

After enabling "Text Content Safety", add a Secret with permissions in "Access Management" - "API Key Management", and fill in the configuration:

```yaml
moderator:
  tencent:
    enabled: true
    secret_id: ''
    secret_key: ''
    region: ap-guangzhou
```

## Alibaba Cloud Content Safety

Refer to: [Alibaba Cloud Documentation](https://help.aliyun.com/document_detail/28417.html)

After enabling "Alibaba Cloud Content Safety", create an Access Key in the Alibaba Cloud backend and fill in the configuration:

```yaml
moderator:
  aliyun:
    enabled: true
    access_key_id: ''
    access_key_secret: ''
    region: cn-shanghai
```

## Keyword Library Filtering

If you prefer not to rely on remote APIs, you can configure and import keyword files locally to let Artalk detect spam comments based on keywords:

```yaml
moderator:
  keywords:
    enabled: true
    pending: false # Set to pending review if matched
    files: # Support multiple keyword files
      - ./data/keyword_1.txt
    file_sep: "\n" # Keyword file content separator
    replace_to: 'x' # Replacement character
```

- **pending**: Whether to set the comment to pending review status upon successful match.
- **files**: Keyword files. Multiple files are allowed, and Artalk will merge the keyword libraries upon startup.
- **file_sep**: Keyword file content separator. For example, if each line in the file contains a keyword, set this item to `\n`.
- **replace_to**: Replacement character. For example, if this item is set to `x`, you can set `pending` to `false`, allowing the comment to pass review automatically, but the matched keywords will be replaced with `x`, resulting in `fxxk` or `xxxx`.

Note: It is recommended not to use `*` asterisk as `replace_to` because it conflicts with the Markdown bold syntax.

## Using Captcha

You can enable Artalk's captcha feature, supporting image and slider captchas, [refer here](./captcha.md).

# 多元通知

Artalk v2.1.8+ 配置项：

- `notify` 弃用并变更为 `admin_notify`
- `email.mail_subject_to_admin` 弃用并变更为 `admin_notify.email.mail_subject`

```yaml
# 管理员多元推送
admin_notify:
  # 通知模版
  notify_tpl: default
  noise_mode: false
  # 邮件通知管理员
  email:
    enabled: true # 当使用其他推送方式时，可以关闭管理员邮件通知
    mail_subject: '[{{site_name}}] 您的文章「{{page_title}}」有新回复'
```

请参考新文档：[多元推送](./admin_notify.md)

# Email Notifications

Artalk supports sending email notifications to users through various methods such as SMTP protocol, Alibaba Cloud Email Service, and system Sendmail command invocation.

You can modify this configuration in the settings interface of the [Dashboard](../frontend/sidebar.md#settings), or configure it through [configuration files](./config.md#email-notifications-email) or [environment variables](../env.md#email-notifications).

## Configuration File

The complete `email` configuration is as follows:

```yaml
# Email Notifications
email:
  enabled: false # Master Switch
  send_type: smtp # Sending Method [smtp, ali_dm, sendmail]
  send_name: '{{reply_nick}}' # Sender's Nickname
  send_addr: example@qq.com # Sender's Address
  mail_subject: '[{{site_name}}] You have received a reply from @{{reply_nick}}'
  mail_tpl: default # Email Template File
  smtp:
    host: smtp.qq.com
    port: 587
    username: example@qq.com
    password: ''
  ali_dm: # https://help.aliyun.com/document_detail/29444.html
    access_key_id: '' # Access Key ID issued by Alibaba Cloud
    access_key_secret: '' # Secret key for encryption
    account_name: example@example.com # Sending address configured in the management console
```

### Choosing a Sending Method

The configuration item `enabled` enables email notifications, and `send_type` is used to select the sending method. Options are: `smtp`, `ali_dm`, `sendmail`.

```yaml
email:
  enabled: true
  send_type: smtp # Sending Method
  # Other configurations omitted...
  smtp:
    # SMTP Configuration...
  ali_dm:
    # Alibaba Cloud Push Configuration...
```

### SMTP Configuration

```yaml
# Email Notifications
email:
  enabled: true
  send_type: smtp # Selecting smtp
  smtp:
    host: smtp.qq.com
    port: 587
    username: example@qq.com
    password: ''
```

### Alibaba Cloud Push Configuration

```yaml
email:
  enabled: true
  send_type: ali_dm # Selecting ali_dm
  ali_dm:
    access_key_id: '' # Access Key ID issued by Alibaba Cloud
    access_key_secret: '' # Secret key for encryption
    account_name: example@example.com # Sending address configured in the management console
```

Refer to: [Alibaba Cloud Official Documentation](https://help.aliyun.com/document_detail/29444.html)

## Comment Replies

The email will include a comment reply button, linking to the given PageKey on the frontend. If your `pageKey` configuration item is a "relative path" of the page, you need to set a URL for your site in the "[Dashboard](../frontend/sidebar.md#dashboard)" - "Site":

<img src="/images/sidebar/site_url.png" width="400px">

## Email Templates

### Template Variables

You can use template variables in `mail_subject` and `mail_subject_to_admin` as well as in the email template file:

```yaml
email:
  # Other configurations omitted...
  mail_subject: '[{{site_name}}] You have received a reply from @{{reply_nick}}'
  mail_subject_to_admin: '[{{site_name}}] Your article "{{page_title}}" has a new reply'
  mail_tpl: default # Email Template File
```

Variables follow the "Mustache" syntax, outputting a variable with `double braces` + `variable name`:

**Basic Content Variables**

```yaml
{{content}}        # Comment content
{{link_to_reply}}  # Reply link
{{nick}}           # Commenter's nickname
{{page_title}}     # Page title
{{page_url}}       # Page PageKey (URL)
{{reply_content}}  # Content of the reply
{{reply_nick}}     # Nickname of the reply target
{{site_name}}      # Site name
{{site_url}}       # Site URL
```

::: details View Other Variables

```yaml
# Comment Creator
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

# Parent Comment (Comment that the creator is replying to)
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

::: details View Data Sample

```json
{
  "content": "<p>Test content</p>\n",
  "link_to_reply": "https://artalk.js.org/?atk_comment=8100&atk_notify_key=44a9b2f08312565fba47c716df9d177f",
  "nick": "Username",
  "page_title": "Artalk",
  "page_url": "https://artalk.js.org/",

  "reply_content": "<p>Replier's content</p>\n",
  "reply_nick": "Replier",
  "site_name": "ArtalkDemo",
  "site_url": "http://localhost:3000/",

  "comment.badge_color": "",
  "comment.badge_name": "",
  "comment.content": "<p>Replier's content</p>\n",
  "comment.content_raw": "Replier's content",
  "comment.date": "2021-11-22",
  "comment.datetime": "2021-11-22 22:22:42",
  "comment.email": "replyer@example.com",
  "comment.email_encrypted": "249898bd50e0febc5799485cf10b345a",
  "comment.id": 8100,
  "comment.is_allow_reply": true,
  "comment.is_collapsed": false,
  "comment.is_pending": false,
  "comment.link": "",
  "comment.nick": "Replier",
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
  "parent_comment.content": "<p>Test content</p>\n",
  "parent_comment.content_raw": "Test content",
  "parent_comment.date": "2021-11-22",
  "parent_comment.datetime": "2021-11-22 22:21:17",
  "parent_comment.email": "test@example.com",
  "parent_comment.email_encrypted": "55502f40dc8b7c769880b10874abc9d0",
  "parent_comment.id": 8099,
  "parent_comment.is_allow_reply": true,
  "parent_comment.is_collapsed": false,
  "parent_comment.is_pending": false,
  "parent_comment.link": "https://qwqaq.com",
  "parent_comment.nick": "Username",
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

### Custom Templates

You can set `mail_tpl` to a "specific file path" to use an external custom email template.

For example, set `mail_tpl` to `"/root/Artalk/data/mail_tpl/your_email_template.html"`

```yaml
email:
  mail_tpl: /root/Artalk/data/mail_tpl/your_email_template.html
  # Other configurations omitted...
```

Then, there should be a file at this path:

```html
<div>
  <p>Hi, {{nick}}：</p>
  <p>You have received a reply in “{{page_title}}”:</p>
  <div>
    <div>@{{reply_nick}}:</div>
    <div>{{reply_content}}</div>
  </div>
  <p><a href="{{link_to_reply}}" target="_blank">Reply to the message »</a></p>
</div>
```

Artalk includes many preset email templates, such as `mail_tpl: "default"`, which uses: [@ArtalkJS/Artalk:/internal/template/email_tpl/default.html](https://github.com/ArtalkJS/Artalk/blob/master/internal/template/email_tpl/default.html)

## Emails to Administrators

Email notifications target both administrators and regular users. You can set different subjects for emails sent to administrators via the following configuration:

```yaml
admin_notify:
  enabled: true
  mail_subject: '[{{site_name}}] Your article "{{page_title}}" has a new reply'
```

Note: The old `email.mail_subject_to_admin` configuration item has been deprecated. Please use the above instead.

Artalk supports various methods to send notifications to administrators, not limited to emails. Refer to: [Diverse Push Notifications](./admin_notify.md#email-notifications).

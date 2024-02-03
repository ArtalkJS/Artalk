# 后端配置

## 轻松配置

推荐在侧边栏 “[控制中心](/guide/frontend/sidebar.md)” 通过图形界面修改配置，无需手动编辑配置文件。

## 指定配置文件路径

Artalk 默认以工作目录下的 `artalk.yml` 作为配置文件，可使用参数 `-c` 来指定具体文件：

```bash
artalk -c ./conf.yml
```

## 通过环境变量配置

Artalk 读取以 `ATK_` 为前缀的环境变量，并且全部大写，子节点用单个下划线表示，配置名含下划线请用双下划线表示：

  - `_` (单下划线) 转为 `.` 表示子节点
  - `__` (双下划线) 转为 `_` 表示配置名的下划线

e.g.

```bash
ATK_TIMEZONE             -> timezone
ATK_LOGIN__TIMEOUT       -> login_timeout
ATK_SITE__DEFAULT        -> site_default
ATK_DB_TYPE              -> db.type
ATK_DB_TABLE__PREFIX     -> db.table_prefix
ATK_CACHE_REDIS_USERNAME -> cache.redis.username
ATK_ADMIN_USERS_0_NAME   -> admin_users[0].name
```

## 获取模版配置文件

可参考一份「完整的配置文件」：[artalk.example.zh-CN.yml](https://github.com/ArtalkJS/Artalk/blob/master/conf/artalk.example.zh-CN.yml)

#### 使用 gen 命令生成配置文件

Artalk 提供 `gen` 命令，你可以快速生成一份新的配置文件：

```bash
artalk gen conf -lang zh-CN ./artalk.yml
```

#### 命令行下载配置文件

```bash
wget -O artalk.yml https://raw.githubusercontent.com/ArtalkJS/Artalk/master/conf/artalk.example.zh-CN.yml
```

## 加密密钥 `app_key`

在 Artalk 启动之前，你需要配置一个 `app_key` 用于对网站内容进行安全加密：

```yaml
app_key: "<任意的字符>"
```

## 语言 `locale`

设置 Artalk 的语言。遵循 Unicode BCP 47 规范，该项默认为 "zh-CN" (简体中文)。

```yml
locale: "zh-CN"
```

详情参考：[多语言](../frontend/i18n.md)

## 数据库 `db`

Artalk 支持连接多种数据库，支持 SQLite、MySQL、PostgreSQL、SQL Server 配置如下：

#### SQLite

SQLite 是轻型数据库，使用单个文件存储数据，无需额外运行程序，尤其适合小型站点，例如个人博客。

```yaml
db:
  type: "sqlite"
  file: "./data/artalk.db"
```

#### MySQL / PostgreSQL / SQL Server

修改 `type` 为你的数据库类型：

```yaml
db:
  type: "mysql"      # sqlite, mysql, pgsql, mssql
  name: "artalk"     # 数据库名
  host: "localhost"  # 地址
  port: "3306"       # 端口
  user: "root"       # 账号
  password: ""       # 密码
  charset: "utf8mb4" # 编码格式
  table_prefix: ""   # 表前缀 (例如："atk_")
  ssl: false         # 启用 SSL
  prepare_stmt: true # 预编译语句
```

数据表将在 Artalk 启动时自动完成创建，无需额外操作。

#### 数据库连接字符串 (DSN)

如有需要，你还可以手动配置 `db.dsn` 来指定数据库连接字符串，例如：

```yaml
db:
  type: "mysql"
  dsn: "mysql://myuser:mypassword@localhost:3306/mydatabase?tls=skip-verify"
```

更多内容参考：[@go-sql-driver/mysql:README.md](https://github.com/go-sql-driver/mysql)

## 管理员 `admin_users`

你需要配置管理员账户，这样才能通过「[控制中心](../frontend/sidebar.md)」对站点内容进行管理。

Artalk 支持多站点，你可以创建多个管理员账户，为其分配站点，让你的朋友们共用同一个后端程序。

详情参考：[管理员 × 多站点](/guide/backend/multi-site.md)

## 可信域名 `trusted_domains`

```yaml
trusted_domains:
  - "https://前端使用域名A.com"
  - "https://前端使用域名B.com"
```

配置该项能限制来自列表外的 Referer 和跨域请求。

:::tip

你需要将「使用该后端的前端」URL 地址加入可信域名列表中，

若非默认 80/443 端口需额外附带端口号，例如：`https://example.com:8080`

:::

在侧边栏[控制中心](../frontend/sidebar.md#控制中心)「站点」选项卡 - 选择站点「修改 URL」，填入站点 URL 也具有相同的效果；添加多个 URL 可使用 `","` 英文逗号分隔，修改后请重启 Artalk。

可以将其关闭：

```yaml
trusted_domains:
  - "*"
```

::: danger

但并不建议这样做，关闭后将存在潜在安全风险，例如可能遭受 CSRF 跨域攻击。

:::

细节：`trusted_domains` 配置项实际上是对响应标头：

- `Access-Control-Allow-Origin` 的控制 (参考：[W3C Cross-Origin Resource Sharing](https://fetch.spec.whatwg.org/#http-cors-protocol))
- `Referer` 的限制 (参考：[Referer - HTTP | MDN](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Referer))

CSRF 跨域攻击防范措施参考：[OWASP 安全备忘单](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)

## 默认站点 `site_default`

如果你觉得大概是不会用到 Artalk 的多站点功能，可以直接将该项配置为你的站点名，例如：

```yaml
site_default: "Artalk 官网"
```

然后在前端直接使用这个站点名：

```js
Artalk.init({ site: "Artalk 官网" })
```

这样，你就无需在侧边栏的[控制中心](../frontend/sidebar.md#控制中心)手动创建站点。

## 前端配置 `frontend`

增加 `frontend` 字段内容可以在后端控制前端的配置，详情可参考：[在后端控制前端](/guide/backend/fe-control)。

## 邮件通知 `email`

配置邮件通知，让回复通过邮件的形式通知目标用户，你可以自定义邮件发送者名称、标题、模版等。

详情参考：[后端 · 邮件通知](/guide/backend/email.md)

## 多元推送 `admin_notify`

你可以配置多种消息发送方式，例如飞书、Telegram 等，当收到新的评论时通知管理员。

详情参考：[后端 · 多元推送](/guide/backend/admin_notify.md)

## 评论审核 `moderator`

配置评论审核来自动拦截垃圾评论。

详情参考：[后端 · 评论审核](/guide/backend/moderator.md)

## 验证码 `captcha`

支持图片、滑动验证码，通过验证码对请求频率进行限制。

详情参考：[后端 · 验证码](/guide/backend/captcha.md)

## 高速缓存 `cache`

为节省内存资源占用，缓存默认关闭。如果你对网站性能有较高要求，请手动开启。你还可以连接外部缓存服务器，支持 Redis 和 Memcache。

```yaml
cache:
  enabled: true   # 启用缓存 (默认关闭)
  type: "builtin" # 支持 redis, memcache, builtin (自带缓存)
  expires: 30     # 缓存过期时间 (单位：分钟)
  warm_up: false  # 程序启动时预热缓存
  server: ""      # 连接缓存服务器 (例如："localhost:6379")
```

- **warm_up**：缓存预热功能。设置为 `true`，在 Artalk 启动时会立刻对数据库内容进行全面缓存，如果你的评论数据较多，多达上万条，启动时间可能会延长。
- **type**：缓存类型，默认为 `builtin`。可选：`redis`, `memcache`, `builtin`。

注：如果在 Artalk 程序外部修改数据库内容，需要刷新 Artalk 缓存才能更新。

---

Redis 身份认证、数据库配置：

```yaml
# 缓存
cache:
  # 省略其他配置项...
  redis:
    network: "tcp" # 连接方式 (tcp 或 unix)
    username: ""   # 用户名
    password: ""   # 密码
    db: 0          # 使用零号数据库
```

<!-- 技术细节：[Artalk 缓存机制 时序图.png](/images/artalk/artalk-cache.png) -->
<!-- ![](/images/artalk/artalk-cache.png) -->

## 监听地址 `host`

Artalk 的默认 HTTP 端口为 23366，你可以在配置文件中指定：

```yaml
host: "0.0.0.0"
port: 23366
```

配置 `host` 监听地址为 `0.0.0.0` 将 Artalk 服务暴露到全网可访问范围，

如果你只想让 Artalk 仅本地能够访问，可将 `host` 配置为 `127.0.0.1`。

命令行下启动 Artalk 时，可以携带 `--host` 和 `--port` 参数分别对地址和端口进行指定，例如：

```bash
artalk server --host 127.0.0.1 --port 8080
```

## 加密传输 `ssl`

```yaml
ssl:
  enabled: true
  cert_path: ""
  key_path: ""
```

你可以配置该项，让 HTTP 升级为 HTTPS，通过 SSL 协议加密传输数据。

- `cert_path`：SSL 证书公钥文件路径。
- `key_path`：SSL 证书私钥文件路径。

你也可以直接反向代理 Artalk 本地服务器，然后在例如 Nginx 启用 HTTPS。

## 时区配置 `timezone`

```yaml
timezone: "Asia/Shanghai"
```

该值填写你所在地时区，对应 IANA 数据库时区名，参考：[Wikipedia](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) / [RFC-6557](https://www.rfc-editor.org/rfc/rfc6557.html)。

```
UTC+08:00   Asia/Shanghai
UTC+09:00   Asia/Tokyo
UTC-07:00   America/Los_Angeles
UTC-04:00   America/New_York
```

## 登录超时 `login_timeout`

该值设定管理员账户登录 JWT 令牌的有效期，单位：秒。

例如，3 天有效：

```yaml
login_timeout: 259200
```

## 日志配置 `log`

打开日志后，系统错误等信息将被记录到设定的文件中。

```yaml
log:
  enabled: true # 总开关
  filename: "./data/artalk.log" # 日志文件路径
```

## 调试模式 `debug`

将 `debug` 配置为 `true` 启用调试模式。

```yaml
debug: true
```

## 工作目录 `-w` 参数

Artalk 在不指定工作目录的情况下，会使用「程序启动时的目录」作为工作目录。

```bash
pwd  # 显示当前目录路径
```

使用参数 `-w` 来指定工作目录，它通常是一个「绝对路径」，例如：

```bash
artalk -w /root/artalk -c ./conf.yml
```

注：`-c` 的相对路径会基于 `-w` 的路径，Artalk 此时会读取 `/root/artalk/conf.yml` 作为配置文件。

其次，在「配置文件中」使用的「相对路径」，也会基于「工作目录」。

例如 `conf.yml` 中有这样的配置：

```yaml
test_file: "./data/artalk.log"
```

将读取 `/root/artalk/data/artalk.log`。

::: tip

配置文件相关代码：[/internal/config/config.go](https://github.com/ArtalkJS/Artalk/blob/master/internal/config/config.go)

前往：[前端配置](/guide/frontend/config.md)
:::

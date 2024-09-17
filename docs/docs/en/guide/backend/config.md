# Configuration File

## Easy Configuration

It is recommended to modify the configuration through the graphical interface in the sidebar "[Dashboard](../frontend/sidebar.md)" without manually editing the configuration file.

## Environment Variables

You can configure via environment variables, refer to the documentation: [Environment Variables](../env.md).

## Specifying the Configuration File Path

By default, Artalk uses `artalk.yml` in the working directory as the configuration file. You can specify a specific file using the `-c` parameter:

```bash
artalk -c ./conf.yml
```

## Obtaining a Template Configuration File

You can refer to a "complete configuration file": [artalk.example.zh-CN.yml](https://github.com/ArtalkJS/Artalk/blob/master/conf/artalk.example.zh-CN.yml)

### Generate Configuration File Using the `gen` Command

Artalk provides the `gen` command, allowing you to quickly generate a new configuration file:

```bash
artalk gen conf -lang zh-CN ./artalk.yml
```

### Download Configuration File via Command Line

```bash
wget -O artalk.yml https://raw.githubusercontent.com/ArtalkJS/Artalk/master/conf/artalk.example.zh-CN.yml
```

## Encryption Key `app_key`

Before starting Artalk, you need to configure an `app_key` for securely encrypting site content:

```yaml
app_key: <any random characters>
```

## Language `locale`

Set the language for Artalk. Follows the Unicode BCP 47 standard, with the default being "zh-CN" (Simplified Chinese).

```yaml
locale: zh-CN
```

For details, refer to: [Multilingual](../frontend/i18n.md)

## Database `db`

Artalk supports connecting to various databases, including SQLite, MySQL, PostgreSQL, and SQL Server. The configurations are as follows:

#### SQLite

SQLite is a lightweight database that stores data in a single file, requiring no additional running programs. It is especially suitable for small sites, such as personal blogs.

```yaml
db:
  type: sqlite
  file: ./data/artalk.db
```

#### MySQL / PostgreSQL / SQL Server

Modify `type` to your database type:

```yaml
db:
  type: mysql # sqlite, mysql, pgsql, mssql
  name: artalk # Database name
  host: localhost # Address
  port: 3306 # Port
  user: root # Username
  password: '' # Password
  charset: utf8mb4 # Charset
  table_prefix: '' # Table prefix (e.g., "atk_")
  ssl: false # Enable SSL
  prepare_stmt: true # Prepared statements
```

The data tables will be automatically created when Artalk starts, no additional operation is required.

#### Database Connection String (DSN)

If needed, you can manually configure `db.dsn` to specify the database connection string, for example:

```yaml
db:
  type: mysql
  dsn: mysql://myuser:mypassword@localhost:3306/mydatabase?tls=skip-verify
```

For more details, refer to: [@go-sql-driver/mysql:README.md](https://github.com/go-sql-driver/mysql)

## Server `http`

```yaml
http:
  # Request body size limit (unit: MB)
  body_limit: 100
  # Proxy header name (when using CDN, fill in `X-Forwarded-For` to get the user's real IP)
  proxy_header: ""
```

## Admin Users `admin_users`

You need to configure administrator accounts to manage site content through the "[Dashboard](../frontend/sidebar.md)".

Artalk supports multiple sites, allowing you to create multiple administrator accounts and assign sites to them, enabling your friends to share the same backend program.

For details, refer to: [Admin Users × Multi-site](./multi-site.md)

## Trusted Domains `trusted_domains`

```yaml
trusted_domains:
  - https://frontend.domainA.com
  - https://frontend.domainB.com
```

Only the domains in the list are allowed to access the backend API, restricting cross-domain requests from external domains.

::: tip

You need to add the URL address of the "frontend page" to the trusted domains list.

If it's not the default 80/443 port, include the port number, e.g., `https://example.com:8080`.

:::

In the sidebar [Dashboard](../frontend/sidebar.md#控制中心), under the "Site" tab - select the site and "Modify URL", fill in the site URL to achieve the same setting effect; to add multiple URLs, separate them with commas `,`. After modification, manually restart Artalk.

The default site address `ATK_SITE_URL` configuration item will also be automatically added to the trusted domains list.

Additionally, you can configure trusted domains via the environment variable `ATK_TRUSTED_DOMAINS` when starting, for example:

```bash
ATK_TRUSTED_DOMAINS="https://a.com https://b.org" artalk server
```

Note that environment variables will override the settings in the configuration file.

Technical details: `trusted_domains` configuration controls cross-origin requests by setting the HTTP response header `Access-Control-Allow-Origin` (refer to: [W3C Cross-Origin Resource Sharing](https://fetch.spec.whatwg.org/#http-cors-protocol) / [OWASP Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)).

## Default Site `site_default`

If you think you won't use the multi-site feature of Artalk, you can directly set this item to your site name, for example:

```yaml
site_default: Artalk Official Website
```

Then use this site name in the frontend:

```js
Artalk.init({ site: 'Artalk Official Website' })
```

This way, you don't need to manually create the site in the sidebar [Dashboard](../frontend/sidebar.md#控制中心).

## Frontend Configuration `frontend`

In Artalk's configuration file `artalk.yml`, you can configure the `frontend` field to control the frontend interface, for example:

```yaml
frontend:
  placeholder: Type content...
  noComment: "Silence is golden"
  sendBtn: Send Comment
  emoticons: 'https://raw.githubusercontent.com/ArtalkJS/Emoticons/master/grps/default.json'
  # ----- Omitted -------
  # Keep the same as the frontend configuration item name
```

A complete backend `frontend` field configuration file can be referred to: [artalk.example.zh-CN.yml](https://github.com/ArtalkJS/Artalk/blob/master/conf/artalk.example.zh-CN.yml)

## Email Notifications `email`

Configure email notifications to notify the target user of replies via email. You can customize the sender name, title, template, etc.

For details, refer to: [Backend · Email Notifications](./email.md)

## Multi-channel Notifications `admin_notify`

You can configure various message sending methods, such as Feishu, Telegram, etc., to notify administrators when new comments are received.

For details, refer to: [Backend · Multi-channel Notifications](./admin_notify.md)

## Comment Moderation `moderator`

Configure comment moderation to automatically intercept spam comments.

For details, refer to: [Backend · Comment Moderation](./moderator.md)

## Captcha `captcha`

Supports image and slide captchas to limit request frequency.

For details, refer to: [Backend · Captcha](./captcha.md)

## Cache `cache`

To save memory resources, caching is disabled by default. If you have high performance requirements for your site, enable it manually. You can also connect to external cache servers, supporting Redis and Memcache.

```yaml
cache:
  enabled: true # Enable cache (default is off)
  type: builtin # Supports redis, memcache, builtin (built-in cache)
  expires: 30 # Cache expiration time (unit: minutes)
  warm_up: false # Warm up cache on program startup
  server: '' # Connect to cache server (e.g., "localhost:6379")
```

- **warm_up**: Cache warm-up feature. Set to `true`, it will immediately cache all database content when Artalk starts. If you have a large number of comments, this may extend the startup time.
- **type**: Cache type, defaults to `builtin`. Options: `redis`, `memcache`, `builtin`.

Note: If you modify the database content outside the Artalk program, you need to refresh the Artalk cache to update it.

---

Redis authentication and database configuration:

```yaml
# Cache
cache:
  # Omit other configurations...
  redis:
    network: tcp # Connection method (tcp or unix)
    username: '' # Username
    password: '' # Password
    db: 0 # Use database 0
```

<!-- Technical details: [Artalk Cache Mechanism Sequence Diagram.png](/images/artalk/artalk-cache.png) -->
<!-- ![](/images/artalk/artalk-cache.png) -->

## Listening Address `host`

The default HTTP port for Artalk is 23366. You can specify it in the configuration file:

```yaml
host: '0.0.0.0'
port: 23366
```

Configuring the `host` listening address to `0.0.0.0` will expose the Artalk service to the entire network.

If you want Artalk to be accessible only locally, configure `host` to `127.0.0.1`.

When starting Artalk via the command line, you can specify the address and port using the `--host` and `--port` parameters, for example:

```bash
artalk server --host 127.0.0.1 --port 8080
```

## Encrypted Transmission `ssl`

```yaml
ssl:
  enabled: true
  cert_path: ''
  key_path: ''
```

You can configure this item to upgrade HTTP to HTTPS, encrypting data transmission via SSL.

- `cert_path`: Path to the SSL certificate public key file.
- `key_path`: Path to the SSL certificate private key file.

You can also directly reverse proxy the Artalk local server and then enable HTTPS in, for example, Nginx.

## Time Zone Configuration `timezone`

```yaml
timezone: Asia/Shanghai
```

Fill in the time zone where you are located, corresponding to the IANA database time zone name, refer to: [Wikipedia](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) / [RFC-6557](https://www.rfc-editor.org/rfc/rfc6557.html).

```
UTC+08:00   Asia/Shanghai
UTC+09:00   Asia/Tokyo
UTC-07:00   America/Los_Angeles
UTC-04:00   America/New_York
```

## Login Timeout `login_timeout`

This value sets the validity period of the administrator account login JWT token, in seconds.

For example, valid for 3 days:

```yaml
login_timeout: 259200
```

## Log Configuration `log`

When logging is enabled, system errors and other information will be recorded in the specified file.

```yaml
log:
  enabled: true # Master switch
  filename: ./data/artalk.log # Log file path
```

Logs are in JSON format, which can be viewed using the `tail` command combined with the `jq` tool in Linux:

```bash
tail -f ./data/artalk.log | jq
```

If you are a Grafana user, you can use Loki and Promtail to collect and query logs.

## Debug Mode `debug`

Set `debug` to `true` to enable debug mode.

```yaml
debug: true
```

## Working Directory `-w` Parameter

If the working directory is not specified, Artalk will use the "directory where the program is started" as the working directory.

```bash
pwd  # Display the current directory path
```

Use the `-w` parameter to specify the working directory, which is usually an "absolute path", for example:

```bash
artalk -w /root/artalk -c ./conf.yml
```

Note: The relative path of `-c` will be based on the path of `-w`, and Artalk will read `/root/artalk/conf.yml` as the configuration file.

Additionally, the "relative path" used in the "configuration file" will also be based on the "working directory".

For example, if the configuration in `conf.yml` is:

```yaml
test_file: ./data/artalk.log
```

It will read `/root/artalk/data/artalk.log`.

::: tip

Configuration file related code: [/internal/config/config.go](https://github.com/ArtalkJS/Artalk/blob/master/internal/config/config.go)

Go to: [Frontend Configuration](../frontend/config.md)
:::

# 反向代理

## Caddy

**/etc/caddy/Caddyfile**

```nginx
artalk.your_domain.com {
  tls /etc/caddy/ssl/cert.pem /etc/caddy/ssl/cert.key

  reverse_proxy http://localhost:23366 {
    header_up X-Forwarded-For {header.X-Forwarded-For}
  }
}
```

执行重载命令：

```sh
sudo systemctl reload caddy
```

## Nginx

假定：

- 你想绑定的域名是：`artalk.your_domain.com`
- Artalk 本地地址：`http://localhost:23366`

以 Ubuntu 20.04 为例：

创建站点配置文件：

```bash
sudo vim /etc/nginx/sites-available/artalk.your_domain.com
```

编辑反向代理配置文件：

```nginx
server
{
  listen 80;
  listen [::]:80;

  server_name artalk.your_domain.com;

  location / {
    proxy_redirect off;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_pass http://localhost:23366/;
  }
}
```

创建软链接启用站点：

```bash
sudo ln -s /etc/nginx/sites-available/artalk.your_domain.com /etc/nginx/sites-enabled/
```

验证配置文件是否有效：

```bash
sudo nginx -t
```

配置没有问题，重启 Nginx：

```bash
sudo systemctl restart nginx
```

配置前端：

```js
Artalk.init({ server: 'http://artalk.your_domain.com' })
```

::: tip
你还可以再套一层 CDN，然后加上 SSL

注意配置文件权限，以及反代目标 URL 可访问性

尤其是运行在 Docker 容器内的 artalk，注意检查 IP 和端口是否能够被 Nginx 正常访问
:::

## Apache

需要启用反代模块 `mod_proxy.c`

```apache
<VirtualHost *:80>
    ServerName your_domain.xxx
    ServerAlias

    RewriteEngine On
    RewriteCond %{QUERY_STRING} transport=polling         [NC]
    RewriteRule /(.*)           http://localhost:23366/$1 [P]

    <IfModule mod_proxy.c>
        ProxyRequests Off
        SSLProxyEngine on
        ProxyPass / http://localhost:23366/
        ProxyPassReverse / http://localhost:23366/
    </IfModule>
</VirtualHost>
```

## 宝塔面板

首先创建一个站点 (例如 `artalk.your_domain.com`)，然后点击站点的「设置」：

![](/images/baota-proxy/1.png)

打开「反向代理」选项卡，点击「添加反向代理」，「目标 URL」填写 `http://localhost:端口号`（端口号与 Artalk 端口对应），「发送域名」填写 `$host`，如图：

![](/images/baota-proxy/2.png)

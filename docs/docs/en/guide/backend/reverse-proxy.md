# Reverse Proxy

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

Reload the Caddy configuration:

```sh
sudo systemctl reload caddy
```

## Nginx

Assumptions:

- The domain you want to bind is: `artalk.your_domain.com`
- Artalk local address: `http://localhost:23366`

For Ubuntu 20.04:

Create the site configuration file:

```bash
sudo vim /etc/nginx/sites-available/artalk.your_domain.com
```

Edit the reverse proxy configuration file:

```nginx
server {
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

Create a symbolic link to enable the site:

```bash
sudo ln -s /etc/nginx/sites-available/artalk.your_domain.com /etc/nginx/sites-enabled/
```

Check if the configuration file is valid:

```bash
sudo nginx -t
```

If the configuration is correct, restart Nginx:

```bash
sudo systemctl restart nginx
```

Configure the frontend:

```js
Artalk.init({ server: 'http://artalk.your_domain.com' })
```

::: tip
You can also add a layer of CDN and SSL.

Pay attention to the configuration file permissions and the accessibility of the reverse proxy target URL.

Especially if Artalk is running inside a Docker container, make sure the IP and port can be accessed by Nginx.
:::

## Apache

You need to enable the reverse proxy module `mod_proxy.c`

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

## BaoTa Panel (宝塔面板)

First, create a site (e.g., `artalk.your_domain.com`), then click the site "Settings":

![](/images/baota-proxy/1.png)

Open the "Reverse Proxy" tab, click "Add Reverse Proxy", fill in the "Target URL" with `http://localhost:port` (where the port corresponds to the Artalk port), and fill in the "Send Domain" with `$host`, as shown in the figure:

![](/images/baota-proxy/2.png)

## Getting the Accurate IP Address

When using a reverse proxy server, you need to configure the proxy headers to get the user's accurate IP address. Refer to the [IP Region](../frontend/ip-region.md#获取准确的-ip-地址) documentation for more details.

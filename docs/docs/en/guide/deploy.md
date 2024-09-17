# 📦 Program Deployment

## Docker Deployment

It is recommended to deploy using Docker. Pre-install the [Docker Engine](https://docs.docker.com/engine/install/) and execute the command to create the container:

```bash
docker run -d \
    --name artalk \
    -p 8080:23366 \
    -v $(pwd)/data:/data \
    -e "TZ=America/New_York" \
    -e "ATK_LOCALE=en" \
    -e "ATK_SITE_DEFAULT=Artalk Blog" \
    -e "ATK_SITE_URL=https://example.com" \
    artalk/artalk-go
```

Execute the command to create an admin account:

```bash
docker exec -it artalk artalk admin
```

Enter `http://artalk.example.com:8080` in your browser to access the Artalk admin login interface.

Incorporate the embedded front-end JS and CSS resources of the Artalk program into the webpage and initialize Artalk:

<!-- prettier-ignore-start -->

```html
<!-- CSS -->
<link href="http://artalk.example.com:8080/dist/Artalk.css" rel="stylesheet" />

<!-- JS -->
<script src="http://artalk.example.com:8080/dist/Artalk.js"></script>

<!-- Artalk -->
<div id="Comments"></div>
<script>
Artalk.init({
  el:        '#Comments',                // Selector of the bound element
  pageKey:   '/post/1',                  // Permalink
  pageTitle: 'About Introducing Artalk', // Page title (leave blank to auto-fetch)
  server:    'http://artalk.example.com:8080',  // Backend address
  site:      'Artalk Blog',            // Your site name
})
</script>
```
<!-- prettier-ignore-end -->

Enter the admin username and email in the comment box, and the "Console" button will appear at the bottom right of the comment box.

In the console, you can [configure the comment system](./backend/config.md) according to your preferences or [migrate comments to Artalk](./transfer.md).

🥳 You have successfully completed the Artalk deployment!

## Standard Deployment

1. Download the program package from [GitHub Release](https://github.com/ArtalkJS/Artalk/releases)
2. Extract it with `tar -zxvf artalk_version_system_architecture.tar.gz`
3. Run `./artalk server`
4. Configure

   ```js
   Artalk.init({ server: 'http://artalk.example.com:23366' })
   ```

Advanced operations:

- [Daemon Process (Systemd, Supervisor)](./backend/daemon.md)
- [Reverse Proxy (Caddy, Nginx, Apache)](./backend/reverse-proxy.md)
- [Self-compile (Build from the latest code)](../develop/contributing.md)

## Compose Deployment

**compose.yaml**

```yaml
version: '3.8'
services:
  artalk:
    container_name: artalk
    image: artalk/artalk-go
    restart: unless-stopped
    ports:
      - 8080:23366
    volumes:
      - ./data:/data
    environment:
      - TZ=America/New_York
      - ATK_LOCALE=en
      - ATK_SITE_DEFAULT=Artalk Blog
      - ATK_SITE_URL=https://your_domain
```

Create the container:

```bash
docker-compose up -d
```

::: details Common Compose Commands

```bash
docker-compose restart  # Restart the container
docker-compose stop     # Stop the container
docker-compose down     # Remove the container
docker-compose pull     # Update the image
docker-compose exec artalk bash # Enter the container
```

:::

Refer to the documentation: [Docker](./backend/docker.md) / [Environment Variables](./env.md)

## Linux Distributions

**Arch Linux (AUR)**:

```bash
paru -S artalk
```

**NixOS**:

```bash
nix-env -iA nixpkgs.artalk
```

**Termux**:

```bash
pkg install artalk
```

[![Packaging status](https://repology.org/badge/vertical-allrepos/artalk.svg)](https://repology.org/project/artalk/versions)

## CDN Resources

::: tip Latest Version of Artalk

The current latest version number of the Artalk front-end is: **:ArtalkVersion:**

To upgrade the front-end, simply replace the numeric part of the version number in the URL.
:::

The back-end program of Artalk integrates front-end JS and CSS files. Please pay attention to version compatibility when using public CDN resources.

**CDNJS**

> <https://cdnjs.cloudflare.com/ajax/libs/artalk/:ArtalkVersion:/Artalk.js>
>
> <https://cdnjs.cloudflare.com/ajax/libs/artalk/:ArtalkVersion:/Artalk.css>

::: details View More
**SUSTech Mirrors (Domestic)**

> <https://mirrors.sustech.edu.cn/cdnjs/ajax/libs/artalk/:ArtalkVersion:/Artalk.js>
>
> <https://mirrors.sustech.edu.cn/cdnjs/ajax/libs/artalk/:ArtalkVersion:/Artalk.css>

**UNPKG**

> <https://unpkg.com/artalk@:ArtalkVersion:/dist/Artalk.js>
>
> <https://unpkg.com/artalk@:ArtalkVersion:/dist/Artalk.css>

**JSDELIVR**

> <https://cdn.jsdelivr.net/npm/artalk@:ArtalkVersion:/dist/Artalk.js>
>
> <https://cdn.jsdelivr.net/npm/artalk@:ArtalkVersion:/dist/Artalk.css>

:::

## Node Projects

Install Artalk:

```bash
npm install artalk
```

Import Artalk:

```js
import 'artalk/dist/Artalk.css'
import Artalk from 'artalk'

Artalk.init({
  // ...
})
```

For more references:

- [Embedding in Blog Documentation](../develop/import-blog.md)
- [Embedding in Framework Documentation](../develop/import-framework.md)
- [Front-end Configuration](./frontend/config.md)
- [Front-end API](../develop/fe-api.md)

## Data Import

Import data from other comment systems: [Data Migration](./transfer.md)

## ArtalkLite

Consider the lightweight version [ArtalkLite](./frontend/artalk-lite.md): smaller and more streamlined.

## Development Environment

Refer to: [Developer Guide](https://github.com/ArtalkJS/Artalk/blob/master/CONTRIBUTING.md)

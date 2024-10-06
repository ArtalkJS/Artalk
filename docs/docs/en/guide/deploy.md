# 📦 Getting Started

This guide will help you deploy Artalk on your server. Once deployed, you can integrate the Artalk client into your website or blog to enable comment functionality.

## Docker

Here is a simple guide for deploying the Artalk **Server** and **Client**.

### Launch Server

It is recommended to use Docker for deployment. Pre-install the [Docker Engine](https://docs.docker.com/engine/install/) and create a working directory, then run the command to launch a container in the background:

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

(Tip: We also provide a [Docker Compose](#docker-compose) configuration file).

Run the command to create an admin account:

```bash
docker exec -it artalk artalk admin
```

Enter `http://artalk.example.com:8080` in your browser to access the Artalk Dashboard.

### Integrate Client

Import the embedded client JS and CSS resources from the Artalk server to your webpage and initialize the Artalk client:

```html
<!-- CSS -->
<link href="http://artalk.example.com:8080/dist/Artalk.css" rel="stylesheet" />

<!-- JS -->
<script src="http://artalk.example.com:8080/dist/Artalk.js"></script>

<!-- Artalk -->
<div id="Comments"></div>
<script>
Artalk.init({
  el:        '#Comments',                       // Selector of the bound element
  pageKey:   '/post/1',                         // Permalink
  pageTitle: 'About Introducing Artalk',        // Page title (leave blank to auto-fetch)
  server:    'http://artalk.example.com:8080',  // Artalk server address
  site:      'Artalk Blog',                     // Your site name
})
</script>
```

Enter the admin username and email in the comment box, and the "Dashboard" button will appear at the bottom right of the comment box.

In the Dashboard, you can [configure the comment system](./backend/config.md) to your liking or [migrate comments to Artalk](./transfer.md).

🥳 You have successfully deployed Artalk!

## Binary

1. Download the program package from [GitHub Release](https://github.com/ArtalkJS/Artalk/releases)
2. Extract it with `tar -zxvf artalk_version_system_architecture.tar.gz`
3. Run `./artalk server`
4. Configure and initialize the client on your webpage:

   ```js
   Artalk.init({ server: 'http://artalk.example.com:23366' })
   ```

Advanced operations:

- [Daemon Process (Systemd, Supervisor)](./backend/daemon.md)
- [Reverse Proxy (Caddy, Nginx, Apache)](./backend/reverse-proxy.md)
- [Self-compile (Build from the latest code)](https://github.com/ArtalkJS/Artalk/blob/master/CONTRIBUTING.md)

## Go Install

If you have the Golang toolchain installed, you can run the following commands to compile and install the latest version of Artalk:

```bash
go install github.com/artalkjs/artalk/v2@latest
```

Then run the server:

```bash
artalk server
```

The client integration steps are [here](#integrate-client).

## Linux Distributions

**Fedora (Copr)**：

```bash
dnf install 'dnf-command(copr)'
dnf copr enable @artalk/artalk
dnf install artalk
```

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

[![Copr status](https://img.shields.io/badge/dynamic/json?color=blue&label=Fedora%20Copr&query=builds.latest.source_package.version&url=https%3A%2F%2Fcopr.fedorainfracloud.org%2Fapi_3%2Fpackage%3Fownername%3D%40artalk%26projectname%3Dartalk%26packagename%3Dartalk%26with_latest_build%3DTrue)](https://copr.fedorainfracloud.org/coprs/g/artalk/artalk/)

## Docker Compose

Create a working directory and edit the `docker-compose.yml` file:

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

Start the container to launch Artalk server:

```bash
docker-compose up -d
```

The client integration steps are [here](#integrate-client).

::: details Common Compose Commands

```bash
docker-compose restart  # Restart the container
docker-compose stop     # Stop the container
docker-compose down     # Remove the container
docker-compose pull     # Update the image
docker-compose exec artalk bash # Enter the container
```

:::

More information: [Docker](./backend/docker.md) / [Environment Variables](./env.md)

## Client

If you have a frontend or Node.js web project, this following guide will help you integrate the Artalk **Client** into your web project.

### Client Installation

Install via npm:

::: code-group

```sh [npm]
npm install artalk
```

```sh [yarn]
yarn add artalk
```

```sh [pnpm]
pnpm add artalk
```

```sh [bun]
bun add artalk
```

:::

(or see [CDN Resources](#client-cdn-resources)).

### Client Usage

Import Artalk in your web project:

```js
import 'artalk/Artalk.css'
import Artalk from 'artalk'

const artalk = Artalk.init({
  el:        document.querySelector('#Comments'), // DOM element where you want to mount
  pageKey:   '/post/1',                           // Permalink
  pageTitle: 'About Introducing Artalk',          // Page title
  server:    'http://artalk.example.com:8080',    // Artalk server address
  site:      'Artalk Blog',                       // Site name
})
```

For more references:

- [Import to Blog Documentation](../develop/import-blog.md)
- [Import to Frameworks Documentation](../develop/import-framework.md)
- [Client API Reference](../develop/fe-api.md)
- [UI (Client) Configuration](./frontend/config.md)

### Client CDN Resources

To use Artalk in the browser via CDN, for modern browsers (which support [ES modules](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules)), you can use [esm.sh](https://esm.sh) or [esm.run](https://esm.run):

::: code-group

```html [ES Module]
<body>
  <div id="Comments"></div>

  <script type="module">
    import Artalk from 'https://esm.sh/artalk@:ArtalkVersion:'

    Artalk.init({
      el: '#Comments',
    })
  </script>
</body>
```

```html [Legacy]
<body>
  <div id="Comments"></div>

  <script src="https://cdnjs.cloudflare.com/ajax/libs/artalk/:ArtalkVersion:/Artalk.js"></script>
  <script>
    Artalk.init({
      el: '#Comments',
    })
  </script>
</body>
```

:::

Don't forget to include the CSS file:

```html
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/artalk/:ArtalkVersion:/Artalk.css">
```

::: tip Latest Version of Artalk

The latest Artalk client version is: `:ArtalkVersion:`

To upgrade the client, replace the numeric part of the version number in the URL.
:::

Note: The Artalk server had embedded the client resources. Ensure version compatibility when using public CDN resources.

::: details Alternative CDNs

**CDNJS**

```
https://cdnjs.cloudflare.com/ajax/libs/artalk/:ArtalkVersion:/Artalk.js
```

```
https://cdnjs.cloudflare.com/ajax/libs/artalk/:ArtalkVersion:/Artalk.css
```

**JSDELIVR**

```
https://cdn.jsdelivr.net/npm/artalk@:ArtalkVersion:/dist/Artalk.js
```

```
https://cdn.jsdelivr.net/npm/artalk@:ArtalkVersion:/dist/Artalk.css
```

**UNPKG**

```
https://unpkg.com/artalk@:ArtalkVersion:/dist/Artalk.js
```

```
https://unpkg.com/artalk@:ArtalkVersion:/dist/Artalk.css
```

:::

## Data Migration

Import data from other comment systems: [Data Migration](./transfer.md).

## ArtalkLite

ArtalkLite is a lightweight version of the Artalk client. See: [ArtalkLite](./frontend/artalk-lite.md).

## Development Environment

Please refer to: [Developer Guide](https://github.com/ArtalkJS/Artalk/blob/master/CONTRIBUTING.md).

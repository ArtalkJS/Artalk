# 守护进程

## Systemd

前提条件：

- Artalk 二进制文件，可在 [GitHub Release](https://github.com/ArtalkJS/Artalk/releases) 下载
- `systemctl --version` 232 或更新的版本
- `sudo` 管理员权限

移动 Artalk 到 `$PATH`，例如：

```bash
sudo mv artalk /usr/bin/
```

测试是否有效：

```bash
artalk version
```

创建名为 artalk 的用户组：

```bash
sudo groupadd --system artalk
```

创建一个名为 artalk 的用户，并且拥有一个可写的 home 目录：

```bash
sudo useradd --system \
    --gid artalk \
    --create-home \
    --home-dir /var/lib/artalk \
    --shell /usr/sbin/nologin \
    --comment "Artalk server" \
    artalk
```

如果你有 Artalk 的配置文件，请确保刚刚创建的 artalk 用户有权限可以读取。

创建服务文件：

```bash
sudo vim /etc/systemd/system/artalk.service
```

```ini
[Unit]
Description=Artalk
Documentation=https://artalk.js.org
After=network.target network-online.target
Requires=network-online.target

[Service]
Type=simple
User=artalk
Group=artalk
ExecStart=/usr/bin/artalk server -w /var/lib/artalk -c /etc/artalk/artalk.yml
ExecReload=/bin/kill -s HUP $MAINPID
ExecStop=/bin/kill -s QUIT $MAINPID
TimeoutStopSec=5s
LimitNOFILE=1048576
LimitNPROC=512
PrivateTmp=true
ProtectSystem=full
AmbientCapabilities=CAP_NET_ADMIN CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target
```

请仔细检查 `ExecStart` 和 `ExecReload`。 确保二进制文件的位置和程序启动参数是正确的

例如：更改 `-c` 参数的路径来指定配置文件，`-w` 参数可以更改工作目录。

请注意配置文件中的所有相对路径都是基于工作目录，例如配置文件中的 `./data/` 文件夹，如果启动参数 `-w /var/lib/artalk`，将读取 `/var/lib/artalk/data/` 目录。请确保该目录中的文件对于创建的 `artalk` 账户有权限可读写。

保存服务文件后，你可以设置自动启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now artalk
```

验证服务是否正常运行：

```bash
systemctl status artalk
```

一些常用的命令：

- 启动服务：`systemctl start artalk`
- 停止服务：`systemctl stop artalk`
- 查看状态：`systemctl status artalk`
- 查看日志：`journalctl -u artalk --no-pager | less +G`

## Tmux

tmux 将创建一个持续的命令行会话，在 SSH 或 tty 断开后保持在后台。

Note: 服务器关闭或重启后，tmux 会话将被清除，需要手动重新运行程序。

1. 创建会话 `tmux new -s artalk`
2. 运行程序 `./artalk server`

恢复接入会话：`tmux attach -t artalk`

查看所有会话：`tmux ls`

## Supervisor

以宝塔面板举例：打开「软件商店」，搜索并安装「Supervisor管理器」：

![](/images/baota-supervisor/0.png)

安装后，打开插件，点击「添加守护程序」：

![](/images/baota-supervisor/1.png)

> - 启动用户：`root` 或其他
> - 运行目录：点击右侧图标，选择 Artalk 所在目录
> - 启动命令：`./artalk server`

## Docker

更新 Docker 容器的 [Restart 策略](https://docs.docker.com/config/containers/start-containers-automatically/#use-a-restart-policy) 以达到进程守护效果。

```bash
docker update --restart=unless-stopped artalk
```

## Docker Compose

在 `docker-compose.yml` 文件给服务添加 `restart: unless-stopped` 策略：

```diff
version: '3'
services:
  artalk:
+   restart: unless-stopped
```

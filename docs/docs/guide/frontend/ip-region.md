# IP 属地

Artalk 内置用户 IP 属地展示功能，并且你可以设置显示的精度：精确到市、省。

该功能默认关闭，你可在 Artalk 控制中心设置中启用 IP 属地展示功能。

## IP 属地数据库

在开启 IP 属地展示功能之前，你需要下载一个数据库文件：

- [GitHub 下载](https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.xdb) / [镜像站下载](https://gh-proxy.com/?q=https%3A%2F%2Fgithub.com%2Flionsoul2014%2Fip2region%2Fblob%2Fmaster%2Fdata%2Fip2region.xdb) (境内推荐)

下载后请手动放置到 `./data/` 目录下，文件命名为：`ip2region.xdb`

## 精度设置

你可在设置中找到该配置项。

| 显示精度   | 描述             | 示例       |
| ---------- | ---------------- | ---------- |
| `province` | 精确到省（默认） | `四川`     |
| `city`     | 精确到城市       | `四川成都` |
| `country`  | 精确到国家、地区 | `中国`     |

配置文件：

```yaml
# IP 属地
ip_region:
  # 启用 IP 属地展示
  enabled: false
  # 数据文件路径 (.xdb 格式)
  db_path: ./data/ip2region.xdb
  # 显示精度 ["province", "city", "country"]
  precision: province
```

## 获取准确的 IP 地址

如果你正在使用 CDN 或者 Nginx 等可信的反向代理服务器，那么你需要在「设置」-「服务器」选项 -「代理标头名 (`http.proxy_header`)」填写包含用户真实 IP 的请求头字段名，如：`X-Real-IP`（为了安全，该字段默认为空）。修改后，请手动重启 Artalk 服务以生效。

否则 Artalk 将无法获取到用户真实 IP 地址（如果使用了 Docker，可能获取到的 IP 始终是 172.17.0.X，这是 Docker 虚拟网卡的 IP）。

## 隐私权

Artalk 评论将记录用户的 `IP` 和 `User-Agent` 数据，此类数据有关用户隐私权，请在你的网站隐私政策中声明，并提示用户评论将会收集隐私数据。

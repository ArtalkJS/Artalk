# IP 属地

Artalk 内置用户 IP 属地展示功能，并且你可以设置显示的精度：精确到市、省。

该功能默认关闭，你可在 Artalk 控制中心设置中启用 IP 属地展示功能。

## IP 属地数据库

Artalk 支持以下两种数据库格式，选择其中一种即可。已有配置和数据库文件可以继续使用。

### ip2region（`.xdb`）

- [GitHub 下载（IPv4）](https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region_v4.xdb) / [镜像站下载](https://gh-proxy.com/?q=https%3A%2F%2Fgithub.com%2Flionsoul2014%2Fip2region%2Fblob%2Fmaster%2Fdata%2Fip2region_v4.xdb) (境内推荐)

下载后请手动放置到 `./data/` 目录下，文件命名为：`ip2region.xdb`

Artalk 兼容 ip2region 上游的旧版（`国家|区域|省份|城市|运营商`）和新版（`国家|省份|城市|运营商|国家代码`）数据。已有数据库无需改名或转换。

ip2region 的 IPv4 和 [IPv6 数据库](https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region_v6.xdb) 分开提供。当前 `db_path` 只加载一个文件；需要同时覆盖 IPv4 和 IPv6 时，可选择下面的标准版 IPDB。

### qqwry 标准版（`.ipdb`）

[metowolf/qqwry.ipdb](https://github.com/metowolf/qqwry.ipdb) 的标准版提供中文国家、省、市字段，一个文件同时支持 IPv4 和 IPv6。

- [下载标准版](https://cdn.jsdelivr.net/npm/qqwry.ipdb/qqwry.ipdb) / [备用下载](https://unpkg.com/qqwry.ipdb/qqwry.ipdb)

将文件保存为 `./data/qqwry.ipdb`，并将 `ip_region.db_path` 设置为该路径。请保留 `.ipdb` 扩展名，Artalk 会据此选择 IPDB 读取器，无需额外设置数据库类型。不要使用 `qqwry.raw.ipdb` 原版，它没有独立、完整的省市字段，不适用于这里的精度设置。

IPDB 读取器会将整个文件载入内存。建议保持 `cache.enabled: true`，以便复用已加载的数据；关闭缓存后，每次查询都会重新加载文件。替换数据库文件后，请重启 Artalk 以刷新内存中的数据。

## 精度设置

你可在设置中找到该配置项。

| 显示精度   | 描述             | 示例       |
| ---------- | ---------------- | ---------- |
| `province` | 精确到省（默认） | `四川`     |
| `city`     | 精确到城市       | `四川成都` |
| `country`  | 精确到国家、地区 | `中国`     |

省份数据缺失时，`province` 精度回退到国家或地区，不显示更精确的城市。

配置文件：

```yaml
# IP 属地
ip_region:
  # 启用 IP 属地展示
  enabled: false
  # 数据文件路径（.xdb 或标准版 .ipdb）
  db_path: ./data/ip2region.xdb
  # 显示精度 ["province", "city", "country"]
  precision: province
```

## 获取准确的 IP 地址

如果你正在使用 CDN 或者 Nginx 等可信的反向代理服务器，那么你需要在「设置」-「服务器」选项 -「代理标头名 (`http.proxy_header`)」填写包含用户真实 IP 的请求头字段名，如：`X-Real-IP`（为了安全，该字段默认为空）。修改后，请手动重启 Artalk 服务以生效。

否则 Artalk 将无法获取到用户真实 IP 地址（如果使用了 Docker，可能获取到的 IP 始终是 172.17.0.X，这是 Docker 虚拟网卡的 IP）。

## 隐私权

Artalk 评论将记录用户的 `IP` 和 `User-Agent` 数据，此类数据有关用户隐私权，请在你的网站隐私政策中声明，并提示用户评论将会收集隐私数据。

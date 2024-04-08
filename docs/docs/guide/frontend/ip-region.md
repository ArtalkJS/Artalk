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
| `city`     | 精确到城市       | `浙江杭州` |
| `country`  | 精确到国家、地区 | `美国`     |

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

## 隐私权

Artalk 评论将记录用户的 `IP` 和 `User-Agent` 数据，此类数据有关用户隐私权，请在你的网站隐私政策中声明，并提示用户评论将会收集隐私数据。

# IP Region

Artalk has a built-in feature to display the geolocation of user IPs, with adjustable precision levels: city or province.

This feature is disabled by default. You can enable the IP region display feature in the Artalk Dashboard settings.

## IP Region Database

Before enabling the IP region display feature, you need to download a database file:

- [GitHub Download (IPv4)](https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region_v4.xdb) / [Mirror Download](https://gh-proxy.com/?q=https%3A%2F%2Fgithub.com%2Flionsoul2014%2Fip2region%2Fblob%2Fmaster%2Fdata%2Fip2region_v4.xdb) (recommended within China)

After downloading, manually place it in the `./data/` directory, and name the file: `ip2region.xdb`.

Artalk supports both the legacy (`country|area|province|city|ISP`) and current (`country|province|city|ISP|country code`) upstream ip2region data layouts. Existing databases do not need to be renamed or converted.

Currently, `db_path` can load only one database. The upstream [IPv6 database](https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region_v6.xdb) can also be used, but it does not contain IPv4 data and cannot replace the IPv4 database to cover both address families.

## Precision Settings

You can find this configuration item in the settings.

| Precision Level | Description       | Example       |
| --------------- | ----------------- | ------------- |
| `province`      | Province (default) | `Sichuan`     |
| `city`          | City              | `Sichuan Chengdu` |
| `country`       | Country/Region    | `China`       |

When province data is missing, `province` precision falls back to the country or region without exposing the more precise city.

Configuration file:

```yaml
# IP Region
ip_region:
  # Enable IP region display
  enabled: false
  # Data file path (.xdb format)
  db_path: ./data/ip2region.xdb
  # Display precision ["province", "city", "country"]
  precision: province
```

## Obtaining the Correct IP Address

If you are using a CDN or a trusted reverse proxy server like Nginx, you need to specify the request header field containing the user's real IP in the "Settings" - "Server" option - "Proxy Header Name (`http.proxy_header`)", such as `X-Real-IP` (for security, this field is empty by default). After modification, please manually restart the Artalk service to take effect.

Otherwise, Artalk will not be able to obtain the user's real IP address (if using Docker, the IP obtained may always be 172.17.0.X, which is the IP of the Docker virtual network card).

## Privacy Policy

Artalk comments will record users' `IP` and `User-Agent` data. Since such data pertains to user privacy, please declare this in your website's privacy policy and inform users that privacy data will be collected when they comment.

# IP Region

Artalk has a built-in feature to display the geolocation of user IPs, with adjustable precision levels: city or province.

This feature is disabled by default. You can enable the IP region display feature in the Artalk Dashboard settings.

## IP Region Database

Artalk supports the following two database formats. Choose one; existing configurations and database files remain supported.

### qqwry standard edition (`.ipdb`, recommended)

The standard edition of [metowolf/qqwry.ipdb](https://github.com/metowolf/qqwry.ipdb) provides Chinese country, province, and city fields, with IPv4 and IPv6 support in one file.

- [Download standard edition](https://cdn.jsdelivr.net/npm/qqwry.ipdb/qqwry.ipdb) / [Alternative download](https://unpkg.com/qqwry.ipdb/qqwry.ipdb)

Save the file as `./data/qqwry.ipdb` and set `ip_region.db_path` to that path. Keep the `.ipdb` extension: Artalk uses it to select the IPDB reader without an additional database type setting. Do not use the raw edition, `qqwry.raw.ipdb`, which lacks the separate, complete province and city fields required for precision settings.

The IPDB reader loads the entire file into memory on the first query and reuses it for subsequent queries, regardless of the `cache.enabled` setting. Restart Artalk after replacing an IPDB file to load the new data.

### ip2region (`.xdb`)

- [GitHub Download (IPv4)](https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region_v4.xdb) / [Mirror Download](https://gh-proxy.com/?q=https%3A%2F%2Fgithub.com%2Flionsoul2014%2Fip2region%2Fblob%2Fmaster%2Fdata%2Fip2region_v4.xdb) (recommended within China)

After downloading, manually place it in the `./data/` directory, and name the file: `ip2region.xdb`.

Artalk supports both the legacy (`country|area|province|city|ISP`) and current (`country|province|city|ISP|country code`) upstream ip2region data layouts. Existing databases do not need to be renamed or converted.

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
  # Data file path (.xdb or standard .ipdb)
  db_path: ./data/qqwry.ipdb
  # Display precision ["province", "city", "country"]
  precision: province
```

## Obtaining the Correct IP Address

If you are using a CDN or a trusted reverse proxy server like Nginx, you need to specify the request header field containing the user's real IP in the "Settings" - "Server" option - "Proxy Header Name (`http.proxy_header`)", such as `X-Real-IP` (for security, this field is empty by default). After modification, please manually restart the Artalk service to take effect.

Otherwise, Artalk will not be able to obtain the user's real IP address (if using Docker, the IP obtained may always be 172.17.0.X, which is the IP of the Docker virtual network card).

## Privacy Policy

Artalk comments will record users' `IP` and `User-Agent` data. Since such data pertains to user privacy, please declare this in your website's privacy policy and inform users that privacy data will be collected when they comment.

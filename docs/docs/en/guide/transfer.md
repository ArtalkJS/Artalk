# 🛬 Data Migration

## Data Bundle

The data bundle (Artrans = Art + Ran) is the Artalk persistent data storage standard format.

::: details Artran Format Definition

We define it as follows: each comment data (Object) is called an Artran, and multiple comment data together form Artrans (Array type).

```json
{
  "id": "123",
  "rid": "233",
  "content": "Hello Artalk",
  "ua": "Artalk/6.6",
  "ip": "233.233.233.233",
  "created_at": "2021-10-28 20:50:15 +0800",
  "updated_at": "2021-10-28 20:50:15 +0800",
  "is_collapsed": "false",
  "is_pending": "false",
  "vote_up": "666",
  "vote_down": "0",
  "nick": "qwqcode",
  "email": "qwqcode@github.com",
  "link": "https://qwqaq.com",
  "password": "",
  "badge_name": "Administrator",
  "badge_color": "#FF716D",
  "page_key": "https://artalk.js.org/guide/transfer.html",
  "page_title": "Data Migration",
  "page_admin_only": "false",
  "site_name": "Artalk",
  "site_urls": "http://localhost:3000/demo/,https://artalk.js.org"
}
```

We call a JSON array Artrans, and each Object item within the array is an Artran (without the s).

:::

## Conversion Tool

Use the following tool to convert comment data from other formats to Artrans, and then import them into Artalk. [Open in a new window](https://artransfer.netlify.app)

<Artransfer />

::: tip

The following text provides various methods to obtain source data for reference; if you encounter any issues, please submit feedback via [issue](https://github.com/ArtalkJS/Artransfer/issues).

:::

## Data Import

Data files converted to the `.artrans` format can be imported into Artalk:

- **Dashboard Import**: You can find the "Migration" tab in the "[Dashboard](./frontend/sidebar.md#control-center)" and follow the prompts to import Artrans.
- **Command Line Import**: Refer to [Command Line Import](#command-line-import).

## Obtaining Source Data

### Typecho

**Install Plugin to Obtain Artrans**

We provide an Artrans export plugin:

1. Click "[here](https://github.com/ArtalkJS/Artrans-Typecho/releases/download/v1.0.0/ArtransExporter.zip)" to download the plugin and "unzip" it into the Typecho directory `/usr/plugins/`.
2. Go to the Typecho backend "Console - Plugins" to enable the "ArtransExporter" plugin.
3. Go to "Console - Export Comments (Artrans)" to export all Typecho comments in Artrans format.

**Direct Database Connection to Obtain Artrans**

If your blog is no longer active but the database still exists, you can use our command line tool that supports direct connection to the Typecho database.

[Download Artransfer-CLI](https://github.com/ArtalkJS/Artransfer-CLI/releases), unzip the package, and execute:

```sh
./artransfer typecho \
    --db="mysql" \
    --host="localhost" \
    --port="3306" \
    --user="root" \
    --password="123456" \
    --name="typecho_database_name"
```

After execution, you will get a file in Artrans format:

```sh
> ls
typecho-20220424-202246.artrans
```

Note: It supports connecting to various databases. For more details, refer to [here](https://github.com/ArtalkJS/Artransfer-CLI).

### WordPress

Go to the WordPress backend "Tools - Export", check "All Content", and export the file. You can then use the [conversion tool](#conversion-tool) for conversion.

![](/images/transfer/wordpress.png)

### Valine

Go to the [LeanCloud backend](https://console.leancloud.cn/) to export the comment data file in JSON format, then use the [conversion tool](#conversion-tool) for conversion.

![](/images/transfer/leancloud.png)

### Waline

Using the LeanCloud database, Waline can refer to the above method for Valine as their formats are similar.

For independently deployed Waline, download [Artransfer-CLI](https://github.com/ArtalkJS/Artransfer-CLI/releases) to connect to the local database for export. Execute the command line:

```bash
./artransfer waline \
    --db="mysql" \
    --host="localhost" \
    --port="3306" \
    --user="root" \
    --password="123456" \
    --name="waline_database_name" \
    --table-prefix="wl_"
```

You will obtain a data file in Artrans format, then [import into Artalk](#how-to-import-artrans).

Note: It supports connecting to various databases. For more details, refer to [here](https://github.com/ArtalkJS/Artransfer-CLI).

### Disqus

Go to the [Disqus backend](https://disqus.com/admin), find "Moderation - Export" and click to export. Disqus will send a `.gz` compressed package to your email. After extracting, you will get a `.xml` data file, which you can then use the [conversion tool](#conversion-tool) to convert to Artrans.

![](/images/transfer/disqus.png)

### Commento

You can export the data file in JSON format from the Commento backend, then use the [conversion tool](#conversion-tool) for conversion.

[Image, to be supplemented...]

### Twikoo

[Twikoo](https://twikoo.js.org/) is a comment system developed based on Tencent Cloud. Go to the [Tencent Cloud backend](https://console.cloud.tencent.com/tcb) to export the comment data file in JSON format, then use the [conversion tool](#conversion-tool) for conversion.

<img src="/images/transfer/tencent-tcb.png" style="max-width: 480px;">

### Artalk v1 (Old PHP Backend)

[Artalk v1](https://github.com/ArtalkJS/ArtalkPHP) is the old backend of Artalk, written in PHP. The new backend has fully transitioned to Golang with a redesigned data table structure. Upgrading to the new version requires using the [conversion tool](#conversion-tool) for conversion.

Old version data path: `/data/comments.data.json`

## Command Line Import

Execute `artalk import -h` to view the help documentation.

```bash
./artalk import [parameters...] [filename]
```

Set import parameters via the `-p` flag on the command line:

```bash
./artalk import -p '{ "target_site_name": "Site", "target_site_url": "https://xx.com", "json_file": "", "url_resolver": true }' ./artrans.json
```

If importing via the web backend, you can fill the JSON in the text box:

```json
{
  "target_site_name": "Site",
  "target_site_url": "https://xx.com",
  "json_file": "path on the server",
  "url_resolver": true
}
```

Artalk import parameters:

|        Parameter        | Type    | Description                                                                                          |
| :---------------------: | ------- | ---------------------------------------------------------------------------------------------------- |
| `target_site_name`      | String  | Name of the import site                                                                              |
| `target_site_url`       | String  | URL of the import site                                                                               |
| `url_resolver`          | Boolean | URL resolver, default is off. Re-generates the `page_key` based on the `target_site_url` as the new `page_key` for comments |
| `url_keep_domain`       | Boolean | Default is off. Whether to keep the original domain part of the URL. If off, removes the domain part of `pageKey`. When `url_resolver` is on, `url_keep_domain` is also enabled |
| `json_file`             | String  | Path to the JSON data file                                                                           |
| `json_data`             | String  | Content of the JSON data string                                                                      |
| `assumeyes`             | Boolean | Execute directly without confirmation `y/n`                                                          |

## Data Backup

You can find the "Migration" tab in the "[Dashboard](./frontend/sidebar.md#dashboard)" on the front end, and export comment data in Artrans format.

### Command Line Backup

Export: `artalk export ./artrans`

Import: `artalk import ./artrans`

### Advanced Usage

Execute `artalk export` to directly "standard output", and perform "pipe" or "output redirection" operations, for example:

```bash
artalk export | gzip -9 | ssh username@remote_ip "cat > ~/backup/artrans.gz"
```

## Conclusion

We currently support converting data from Typecho, WordPress, Valine, Waline, Disqus, Commento, Twikoo, etc., to Artrans. However, considering the diversity of comment systems, although we have adapted the above types of data, many are still not compatible. If you happen to be using an unsupported comment system, besides waiting for official Artalk support, you can also try to understand the Artrans data format and write your own tools for importing and exporting comment data. If you think your tool is well-written, we would be happy to include it, allowing us to create a tool that can freely switch between different comment systems together.

Visit: [Artransfer Migration Tool Code Repository](https://github.com/ArtalkJS/Artransfer)

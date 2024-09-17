# Image Upload

Artalk provides an image upload feature with options to limit image size, upload frequency, and more. You can also integrate with UpGit to upload images to image hosting services.

You can modify this configuration in the settings interface of the [Dashboard](../frontend/sidebar.md#settings), or configure it through [configuration files](./config.md#image-upload-img_upload) or [environment variables](../env.md#image-upload).

## Configuration File

The complete `img_upload` configuration is as follows:

```yaml
# Image Upload
img_upload:
  enabled: true # Master Switch
  path: ./data/artalk-img/ # Image storage path
  max_size: 5 # Image size limit (Unit: MB)
  public_path: null # Specify the base URL for image links (default is "/static/images/")
  # Use UpGit to upload images to GitHub or image hosting services
  upgit:
    enabled: false # Enable UpGit
    exec: upgit -c <upgit configuration file path> -t /artalk-img
    del_local: true # Delete local images after uploading
```

## Using UpGit to Upload to Image Hosting Services

[UpGit](https://github.com/pluveto/upgit) supports uploading images to various image hosting services or code repositories such as GitHub, Gitee, Tencent Cloud COS, Qiniu Cloud, UpYun, SM.MS, and more.

First, download UpGit and complete the configuration for your target image hosting service according to the [README.md](https://github.com/pluveto/upgit).

Then, add UpGit to the system's environment variables by adding the following to `~/.bashrc`:

```bash
export PATH=$PATH:/path/to/upgit
```

(Or move it directly to `/usr/bin`)

Finally, fill in the UpGit startup parameters in Artalk's `img_upload.upgit` field:

```yaml
upgit:
  enabled: true # Enable UpGit
  exec: upgit -c <upgit configuration file path> -t /artalk-img
  del_local: true # Delete local images after uploading
```

::: warning Update Notice
Starting from version `v2.8.4`, to enhance security, Artalk no longer allows specifying the UpGit executable file path. Please add it to the system's environment variables. :)
:::

### Mounting UpGit with Docker

If you are deploying Artalk with Docker, you can mount the UpGit executable to the container:

```bash
docker run -d --name artalk -v /path/to/upgit:/usr/bin/upgit -v /path/to/artalk:/app/data -p 8080:23366 artalk
```

## Upload Frequency Limit

The frequency limit follows the `captcha` configuration. When the limit is exceeded, a captcha will be prompted.

Refer to: [Backend · Captcha](./captcha.md)

## Path

`img_upload.path` is the "local storage directory" path for uploaded image files. This directory will be mapped by Artalk to be accessible at:

```
http://<backend address>/static/images/
```

## Public Path

The default value for `img.public_path` is: `/static/images/`

When this item is a "relative path", for example: `/static/images/`, the HTML tag for the uploaded image on the frontend will be:

```html
<img src="http://<backend address>/static/images/1.png" />
```

Note: The `<backend address>` is configured in the frontend `conf.server`.

When this item is a "complete URL path", for example: `https://cdn.github.com/img/`, the image tag will be:

```html
<img src="https://cdn.github.com/img/1.png" />
```

Tip: This configuration can be used in scenarios such as load balancing.

## Custom Upload API on the Frontend

The frontend provides the `imgUploader` configuration option, allowing you to customize the API for image upload requests, for example:

```js
Artalk.init({
  imgUploader: async (file) => {
    const form = new FormData()
    form.set('file', file)

    const imgUrl = await fetch('https://api.example.org/upload', {
      method: 'POST',
      body: form,
    })

    return imgUrl
  },
})
```

Refer to: [Frontend Configuration Documentation](../frontend/config.md#imguploader)

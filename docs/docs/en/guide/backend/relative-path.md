# Resolve Relative Path

Artalk supports parsing relative paths, allowing you to configure it on the frontend as follows:

```js
Artalk.init({
  site: 'Example Site', // Your site name
  pageKey: '/relative-path/xx.html', // Using relative path
})
```

Using relative paths is recommended because it facilitates future "site migration" needs.

Next, in the sidebar, navigate to "Dashboard" - "Site", find "Example Site", and modify the site URL.

![](/images/relative-path/1.png)

Afterward, all relative paths will be "based on this site URL". For example:

```bash
"/relative-path/xx.html"
         ↓ resolves to
"https://set-example-site-url.xxx/relative-path/xx.html"
```

## Uses of Resolved URLs

The site URL + page relative path will be used for:

- **Email notifications** containing links to reply to comments
- **Sidebar** quick jumps to specific comments
- **Dashboard** page management to open pages
- **Fetching page titles** and other information

## Configuring Multiple Site URLs

You might need to configure multiple URLs for a site to allow referer and cross-origin requests.

In "Dashboard" - "Site", you can modify the site URL to support adding multiple URLs by separating each URL with a comma `,`.

**When a site has multiple URLs**, the "relative path" will be based on the "first" URL in the list.

## Using Absolute Paths

As opposed to using relative paths, you can use absolute paths. For example, configure the frontend like this:

```js
Artalk.init({
  pageKey: 'https://your_domain.com/relative-path/xx.html', // Using absolute path
})
```

In this case, the backend will not resolve this address. The absolute path specified in `pageKey` will be directly used for locating the page in emails, sidebar, etc.

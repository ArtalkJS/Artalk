# Emoticons

You can modify this configuration in the settings interface of the [Dashboard](./sidebar.md#设置), or through [configuration files](../backend/config.md#界面配置-frontend) and [environment variables](../env.md#界面配置).

In the configuration file, you can fill in the URL of the emoticon list under `frontend.emoticons`:

```yaml
frontend:
  emoticons: https://raw.githubusercontent.com/ArtalkJS/Emoticons/master/grps/default.json
```

## Emoticon Presets

The Artalk community offers many emoticon presets. You can choose a few favorite sets and easily add them to your comment system with simple configurations. Visit the repository: [@ArtalkJS/Emoticons](https://github.com/ArtalkJS/Emoticons).

## Supported Formats

### OwO Format

[OwO](https://github.com/DIYgod/OwO) is an open-source JS plugin developed by [@DIYgod](https://github.com/DIYgod) that allows quick insertion of emoticons into the input box.

Artalk's emoticon feature draws inspiration from its excellent design and natively supports and is compatible with OwO format emoticon data files, as shown below:

```yaml
frontend:
  emoticons: https://raw.githubusercontent.com/DIYgod/OwO/master/demo/OwO.json
  # Directly use OwO format emoticons ↑↑
```

The community also has many OwO format emoticon resources that can be used directly, for example: [@2X-ercha/Twikoo-Magic](https://github.com/2X-ercha/Twikoo-Magic).

### Artalk Emoticon List File Standard Format

In addition to supporting the OwO format, Artalk also natively supports a standard format for emoticon list files:

```js
[
  {
    name: 'Emoji',
    type: 'emoticon', // Text type
    items: [
      { key: 'Hi', val: '|´・ω・)ノ' },
      { key: 'Happy', val: 'ヾ(≧∇≦*)ゝ' },
      //...
    ],
  },
  {
    name: 'Funny',
    type: 'image', // Image type
    items: [
      {
        key: 'Original Funny',
        val: '<Image URL>',
      },
      //...
    ],
  },
]
```

## Front-end Configuration

To configure emoticon lists on the front end, for example:

```js
Artalk.init({
  // Default emoticon list, dynamically imported ↓↓
  emoticons:
    'https://raw.githubusercontent.com/ArtalkJS/Emoticons/master/grps/default.json',
})
```

::: tip

Configuration file and dashboard emoticon settings have a higher priority than front-end configuration.

:::

## Disabling Emoticons

You can disable the emoticon feature by setting `emoticons` to `false`:

```js
Artalk.init({
  emoticons: false,
})
```

## Loading Methods

### Dynamic Loading

Set the `emoticons` attribute to the URL of the emoticon data file. When the emoticon list is opened, Artalk will dynamically import it.

```js
Artalk.init({
  emoticons: '<Emoticon Data File URL>',
})
```

Remote emoticon files support both Artalk and OwO formats, and support nested and mixed loading.

### Static Loading

Compared to dynamic importing, you can statically save the emoticon list object in the JS code of the page to avoid dynamic loading:

```js
Artalk.init({
  emoticons: [
    {
      name: 'Emoji',
      type: 'emoticon', // Text type
      items: [
        { key: 'Hi', val: '|´・ω・)ノ' },
        { key: 'Happy', val: 'ヾ(≧∇≦*)ゝ' },
        //...
      ],
    },
    {
      name: 'Funny',
      type: 'image', // Image type
      items: [
        {
          key: 'Original Funny',
          val: '<Image URL>',
        },
        //...
      ],
    },
  ],
})
```

### Mixed Loading

Artalk supports a mix of **dynamic** and **static** loading, for example:

```js
Artalk.init({
  emoticons: [
    // Dynamic loading
    'https://raw.githubusercontent.com/DIYgod/OwO/master/demo/OwO.json', // OwO format emoticons
    'https://raw.githubusercontent.com/qwqcode/huaji/master/huaji.json',
    // Static loading
    {
      name: 'Emoticon Name',
      type: 'emoticon', // Text type
      items: [
        { key: 'Go Master Ball', val: '(╯°A°)╯︵○○○' },
        //...
      ],
    },
  ],
})
```

### Nested Import

Artalk supports **nested import** of other emoticon resources within remote emoticon resources, for example:

```js
Artalk.init({
  emoticons: ['https://example.org/emoticons.json'],
})
```

Data in `emoticons.json`:

```json
[
  "https://example.org/other-emoticons.json",
  //...
  {
    // Artalk format, OwO format
    //...
  }
]
```

This allows you to efficiently manage and load multiple emoticon sets with minimal configuration.

# Plugin Development

## Developing with Vite

We provide an integrated Vite plugin, enabling you to effortlessly build Artalk plugins with out-of-the-box Vite configuration.

Reference Documentation: [@artalk/plugin-kit](https://github.com/ArtalkJS/Artalk/blob/master/ui/plugin-kit/README.md).

Additionally, we offer a template repository that you can fork to develop your plugins: [artalk-plugin-sample](https://github.com/ArtalkJS/artalk-plugin-sample).

Utilizing Vite in combination with the frontend ecosystem, you can develop Artalk plugins using any framework you prefer, such as Vue, React, Svelte, or SolidJS.

## `Artalk.use`

Extend Artalk's functionality using `Artalk.use`.

```js
import Artalk from 'artalk'

Artalk.use((ctx) => {
  ctx.on('list-loaded', () => {
    console.log('Comment list loaded')
    ctx.getCommentNodes().forEach((node) => {
      node.getEl().style.background = 'red'
    })
  })
})

Artalk.use((ctx) => {
  let el = null

  ctx.on('mounted', () => {
    el = document.createElement('div')
    document.body.appendChild(el)
    console.log('Artalk instance mounted')
  })

  ctx.on('unmounted', () => {
    el.remove()
    console.log('Artalk instance unmounted')
  })
})

const artalk = Artalk.init({ ... })
```

Please note:

- `Artalk.use` must be called before `Artalk.init` to take effect.
- Do not rely on the order of plugin loading; instead, listen for events to execute plugin logic.
- The `created` event is triggered once all plugins are loaded.
- The `mounted` event is triggered when the instance and configuration are fully loaded.
- The `unmounted` event is triggered when the Artalk instance is destroyed.
- The `updated` event is triggered when the configuration changes.

Within the use function, you can access Artalk's Context object via `ctx`.

### Configurable Plugin Options

In TypeScript, you can define a plugin with options using `ArtalkPlugin<T>`.

```ts
export interface DemoPluginOptions {
  foo?: string
}

export const ArtalkDemoPlugin: ArtalkPlugin<DemoPluginOptions> = (ctx, options = {}) => {
  console.log(options.foo)
}
```

Pass options when calling `Artalk.use`:

```ts
import { ArtalkDemoPlugin } from 'artalk-plugin-demo'

Artalk.use(ArtalkDemoPlugin, { foo: 'bar' })
```

## Context API

The Context object contains Artalk's contextual information.

| Member               | Description             |
| -------------------- | ----------------------- |
| `ctx.getEl`          | Get container element   |
| `ctx.getConf`        | Get configuration       |
| `ctx.updateConf`     | Update configuration    |
| `ctx.watchConf`      | Watch configuration     |
| `ctx.setDarkMode`    | Set dark mode           |
| `ctx.getApi`         | Get API client object   |
| --                   | --                      |
| `ctx.fetch`          | Fetch comment data      |
| `ctx.reload`         | Reload comment list     |
| `ctx.getComments`    | Get all comment data    |
| `ctx.getCommentNodes`| Get all comment nodes   |
| --                   | --                      |
| `ctx.on`             | Add event listener      |
| `ctx.off`            | Remove event listener   |
| `ctx.trigger`        | Trigger event           |
| --                   | --                      |
| `ctx.get`            | Get dependency          |
| `ctx.inject`         | Inject dependency       |

::: warning
The Context API is currently unstable and subject to change. Please refer to the CHANGELOG when upgrading.
:::

Reference: [@artalk/src/types/context.ts](https://github.com/ArtalkJS/Artalk/blob/master/ui/artalk/src/types/context.ts)

Artalk includes comprehensive TypeScript type definitions, allowing you to explore the API through editor autocompletion.

## Sample Plugins

Below is a list of externally maintained Artalk plugins:

| Plugin                                                                                       | Description             |
| -------------------------------------------------------------------------------------------- | ----------------------- |
| [artalk-plugin-sample](https://github.com/ArtalkJS/artalk-plugin-sample)                     | Sample Plugin           |
| [@artalk/plugin-katex](https://github.com/ArtalkJS/Artalk/tree/master/ui/plugin-katex)       | LaTeX Formula Plugin    |
| [@artalk/plugin-auth](https://github.com/ArtalkJS/Artalk/tree/master/ui/plugin-auth)         | Auth Plugin (SolidJS)   |
| [@artalk/plugin-lightbox](https://github.com/ArtalkJS/Artalk/tree/master/ui/plugin-lightbox) | Basic Image Lightbox    |

Additionally, many plugins are implemented within Artalk itself, which you can reference for developing your own plugins:

[@ArtalkJS/Artalk - src/plugins](https://github.com/ArtalkJS/Artalk/tree/master/ui/artalk/src/plugins)

## Backend Plugin Development

【TODO】

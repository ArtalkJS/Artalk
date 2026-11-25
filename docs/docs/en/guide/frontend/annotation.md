# Highlight Annotations

The `@artalk/plugin-annotation` plugin adds text-anchored comments to Artalk. Readers select text in your article, create an annotation, and the discussion stays attached to that exact passage.

## Install

```bash
pnpm add @artalk/plugin-annotation
```

## Enable

```ts
import Artalk from 'artalk'
import { ArtalkAnnotationPlugin } from '@artalk/plugin-annotation'

Artalk.use(ArtalkAnnotationPlugin)

const artalk = Artalk.init({
  el: '#comments',
  server: 'https://artalk.example.com',
  site: 'MySite',
  pageKey: location.pathname,
  annotation: {
    enabled: true,
    root: 'article',
    hideInList: true,
    allowGuest: true,
  },
})
```

Or load the plugin with script tags:

```html
<script src="https://artalk.example.com/dist/Artalk.js"></script>
<script src="https://artalk.example.com/dist/plugins/artalk-plugin-annotation.js"></script>
<script>
  Artalk.init({
    el: '#comments',
    server: 'https://artalk.example.com',
    site: 'MySite',
    annotation: { enabled: true, root: 'article' },
  })
</script>
```

## Options

| Option | Type | Default | Description |
| --- | --- | --- | --- |
| `enabled` | `boolean` | `false` | Enable highlight annotations |
| `root` | `string \| HTMLElement` | `'article'` | Element or selector that allows text selection |
| `hideInList` | `boolean` | `true` | Hide annotation roots from the normal comment list |
| `allowGuest` | `boolean` | `true` | Allow guests to create annotations |

The plugin also accepts `excludeSelector` and `panelContainer` through `Artalk.use` options.

## How it works

- An annotation is a root comment with a W3C `TextQuoteSelector` anchor.
- Replies reuse the normal comment tree.
- The normal comment list excludes annotation roots by default (`hideInList`).
- Moderation, anti-spam, notifications and permissions are shared with normal comments.
- If the article text changes and the anchor cannot be relocated, the annotation is marked as **orphan** and can be manually relocated.

## SPA

After a route change, update the page key:

```ts
artalk.update({ pageKey: location.pathname })
```

The plugin clears old highlights and reloads the current page.

## Limitations

- Only text selections are supported (no image or region annotations).
- The annotation panel is scoped to the current page.

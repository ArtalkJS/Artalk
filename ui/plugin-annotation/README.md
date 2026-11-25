# @artalk/plugin-annotation

Text-anchored comments for [Artalk](https://artalk.js.org). Readers select text in an article, create an annotation, and the discussion stays attached to that exact passage.

## Install

```bash
pnpm add artalk @artalk/plugin-annotation
```

## Usage

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

## Options

| Option            | Type                    | Default               | Description                                        |
| ----------------- | ----------------------- | --------------------- | -------------------------------------------------- |
| `enabled`         | `boolean`               | `false`               | Enable highlight annotations                       |
| `root`            | `string \| HTMLElement` | `'article'`           | Element or selector that allows text selection     |
| `hideInList`      | `boolean`               | `true`                | Hide annotation roots from the normal comment list |
| `allowGuest`      | `boolean`               | `true`                | Allow guests to create annotations                 |
| `excludeSelector` | `string`                | `pre, code, svg, ...` | Elements excluded from selection                   |
| `panelContainer`  | `string \| HTMLElement` | `document.body`       | Container for the Shadow DOM panel                 |

## Behavior

- An annotation is a root comment with a W3C `TextQuoteSelector` anchor.
- Replies reuse the normal Artalk comment tree.
- The normal comment list excludes annotation roots when `hideInList` is enabled.
- Moderation, anti-spam, notifications and permissions are shared with normal comments.
- Unresolvable anchors are marked as orphan and can be manually relocated.

## SPA

Call `artalk.update({ pageKey: location.pathname })` after route changes.

## License

MIT

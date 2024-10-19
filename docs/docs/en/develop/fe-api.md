# Frontend API

Here are some of the main Artalk client APIs for operating the Artalk instance in frontend code.

For more detailed API information, you can browse the [TypeDoc documentation](https://artalk.js.org/typedoc/).

## Install Dependencies

Use your preferred package manager to install Artalk:

```bash
pnpm add artalk
# or using npm
npm install artalk
```

## Create Instance `init`

Calling this function creates and returns an instantiated Artalk object for subsequent operations.

<!-- prettier-ignore-start -->

```js
import Artalk from 'artalk'

const artalk = Artalk.init({
  el:        '#Comments',
  pageKey:   '/post/1',
  pageTitle: 'About Introducing Artalk',
  server:    'http://your_domain:8080',
  site:      'Artalk Blog',
})
```

<!-- prettier-ignore-end -->

Calling this function will asynchronously send requests to the backend to:

1. Fetch frontend configuration
2. Fetch comment list

The `mounted` event is triggered after configuration and plugins are loaded.

The `list-loaded` event is triggered after the comment list is loaded.

> Note: The frontend UI configuration may be overridden. See: [UI Configuration](../guide/frontend/config.md).

## Release Resources `destroy`

Destroy the Artalk instance to release resources.

```js
const artalk = Artalk.init({ ... })

artalk.destroy()
```

Calling this function triggers the `unmounted` event. Releasing resources involves deleting DOM nodes mounted by Artalk and removing Artalk event listeners. After this function completes, the Artalk instance will no longer be usable.

In frameworks like Vue/React, make sure to call this function when the component is destroyed to prevent memory leaks:

```ts
import { onUnmounted } from 'vue'

onUnmounted(() => {
  artalk.destroy()
})
```

## Update Configuration `update`

Modify the current Artalk configuration.

```js
const artalk = Artalk.init({ ... })

artalk.update({
  // new configuration...
})
```

Calling this function will trigger the `updated` event.

Updating the configuration will not automatically refresh the comment list. You can continue to call the `artalk.reload()` function as needed.

Note that this function is a member method of an instantiated object, not a global function.

Frontend configuration reference: [Frontend Configuration](../guide/frontend/config.md)

## Reload `reload`

Reload the Artalk comment list data from the backend.

```js
const artalk = Artalk.init({ ... })

artalk.reload()
```

The `list-load` event is triggered before loading the list, and the `list-loaded` event is triggered after the loading is complete.

## Listen to Events `on`

Add Artalk event listeners.

```js
const artalk = Artalk.init({ ... })

artalk.on('list-loaded', (comments) => {
  // 'comments' contains all comments on the current page (not just the comments fetched in this request)
  console.log('Comments have been loaded: ', comments)
})
```

Refer to: [Frontend Event](./event.md).

## Remove Event Listeners `off`

Remove Artalk event listeners.

```js
const artalk = Artalk.init({ ... })

const handler = () => {
  alert('Comments have been loaded')
}

artalk.on('list-loaded', handler)
artalk.off('list-loaded', handler)
```

## Trigger Events `trigger`

Trigger Artalk events.

```js
const artalk = Artalk.init({ ... })

artalk.trigger('list-loaded', [...])
```

## Load Plugins `use`

Load plugins using the `Artalk.use` function. Note that this function is a global function, not a method of an instantiated object. Plugins loaded through this function will be effective for all Artalk instances.

```js
import Artalk from 'artalk'

Artalk.use((ctx) => {
  ctx.editor.setContent("Hello World")
})

const artalk = Artalk.init({ ... })
```

See: [Plugin Development](./plugin.md)

## Dark Mode `setDarkMode`

Toggle dark mode, which can be called in conjunction with the blog theme, such as when a user clicks a dark mode toggle button.

```js
const artalk = Artalk.init({ ... })

artalk.setDarkMode(true)
```

You can also set the initial dark mode by passing the `darkMode` parameter when calling `Artalk.init`.

```js
const artalk = Artalk.init({
  // ...other configurations
  darkMode: 'auto',
})
```

## View Count Widget `loadCountWidget`

See: [Page View Statistics](../guide/frontend/pv.md)

## Get Configuration `getConf`

Get the current Artalk configuration.

```js
const artalk = Artalk.init({ ... })

const conf = artalk.getConf()
```

## Get Mounted Element `getEl`

Get the DOM element that Artalk is currently mounted to.

```js
const artalk = Artalk.init({ ... })

const el = artalk.getEl()
```

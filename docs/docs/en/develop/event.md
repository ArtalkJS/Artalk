# Frontend Event

## Basic Events

| Event Name          | Description                     |
| ------------------- | ------------------------------- |
| `created`           | After initialization            |
| `mounted`           | After data loading              |
| `updated`           | After data update               |
| `unmounted`         | After destruction               |
| `list-fetch`        | During comment list request     |
| `list-fetched`      | After comment list request      |
| `list-load`         | Before comment loading          |
| `list-loaded`       | After comment loading           |
| `list-failed`       | When comment loading fails      |
| `list-goto-first`   | When comment list resets        |
| `list-reach-bottom` | When comment list reaches bottom|
| `comment-inserted`  | After comment insertion         |
| `comment-updated`   | After comment update            |
| `comment-deleted`   | After comment deletion          |
| `comment-rendered`  | After comment node rendering    |
| `notifies-updated`  | When unread messages change     |
| `list-goto`         | During comment jump             |
| `page-loaded`       | After page data update          |
| `editor-submit`     | When editor submits             |
| `editor-submitted`  | After editor submission         |
| `user-changed`      | When local user data changes    |
| `sidebar-show`      | When sidebar is shown           |
| `sidebar-hide`      | When sidebar is hidden          |

Event declaration code: [@ArtalkJS/Artalk - src/types/event.ts](https://github.com/ArtalkJS/Artalk/blob/master/ui/artalk/src/types/event.ts)

## Adding Event Listeners

```js
artalk.on('list-loaded', () => {
  alert('Comments have been loaded')
})
```

## Removing Event Listeners

```js
let foo = function () {
  /* do something */
}

artalk.on('list-loaded', foo)
artalk.off('list-loaded', foo)
```

# Compatibility

## Browser Compatibility

The Artalk client production environment code is built using esbuild with a target of `es2015`. However, please note that the build process only handles syntax transpilation to comply with the `es2015` standard and does not address API compatibility (no polyfills are included).

If your project needs to support older browsers and devices, consider including global polyfills in your bundled application. We recommend using the [Polyfill.io](https://polyfill.io/) service, which dynamically returns the necessary polyfills based on the user's User-Agent. Alternatively, you can manually add polyfills using [core-js](https://github.com/zloirock/core-js) or [Babel](https://babeljs.io/).

Here is a list of modern features used by Artalk:

- [All ECMAScript 2015 (ES6) features](https://compat-table.github.io/compat-table/es6/)
- [Fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API)
- [AbortController](https://developer.mozilla.org/en-US/docs/Web/API/AbortController)
- [Intersection Observer](https://developer.mozilla.org/en-US/docs/Web/API/Intersection_Observer_API)

## Server Compatibility

The Artalk server program is developed in Golang and supports operating systems such as Linux, Windows, and macOS. The current Golang version for Artalk can be found in the [go.mod](https://github.com/ArtalkJS/Artalk/blob/master/go.mod#L3). For the minimum operating system requirements for Golang, refer to the [Go Wiki](https://go.dev/wiki/MinimumRequirements).

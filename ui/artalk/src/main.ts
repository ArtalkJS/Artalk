import Artalk from './artalk'
import type * as ArtalkType from './types'
import { Defaults } from './defaults'

export type * from './types'
export { ArtalkType, Defaults }

export default Artalk

// Manually expose the static methods from the Artalk class.
// Directly exporting static methods from a class is not supported,
// and it's important to maintain consistency in the library's API.
// To ensure compatibility across different environments, such as CommonJS and browser IIFE,
// this approach allows us to use `Artalk.init()` instead of `Artalk.default.init()`.
export const init = Artalk.init
export const use = Artalk.use
export const loadCountWidget = Artalk.loadCountWidget

import Artalk from './artalk'
import type * as ArtalkType from './types'

export type * from './types'
export { ArtalkType }

export default Artalk

// Expose the static methods from the Artalk class
// because direct export of static methods is not supported
// for adapting to different environments like CommonJS and browser IIFE
// for example, we can use `Artalk.init()` rather than `Artalk.default.init()`
// therefore, we need to manually expose the static methods in the Artalk class.
export const init = Artalk.init
export const use = Artalk.use
export const loadCountWidget = Artalk.loadCountWidget

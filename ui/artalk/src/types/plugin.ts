import { Context } from './context'

export type ArtalkPlugin<T = any> = (ctx: Context, options?: T) => void

import { ContextApi } from './context'

export type ArtalkPlugin<T = any> = (ctx: ContextApi, options?: T) => void

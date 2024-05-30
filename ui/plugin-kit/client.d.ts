import Artalk, { ArtalkPlugin } from 'artalk'

export {}

declare global {
  interface Window {
    Artalk?: typeof Artalk
    ArtalkPlugins?: { [name: string]: ArtalkPlugin }
  }
}

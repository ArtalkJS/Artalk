import { ArtalkPlugin } from '.'

export {}

declare global {
  interface Window {
    ArtalkPlugins?: { [name: string]: ArtalkPlugin }
  }
}

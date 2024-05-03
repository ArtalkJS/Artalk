import type { ArtalkPlugin } from 'artalk'

export {}

declare global {
  interface Window {
    ArtalkPlugins?: { [name: string]: ArtalkPlugin }
  }
}

import Artalk from 'artalk'

export {}

declare global {
  interface Window {
    Artalk?: typeof Artalk
  }
}

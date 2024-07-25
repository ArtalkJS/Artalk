import type lightGallery from 'lightgallery'
import type PhotoSwipeLightbox from 'photoswipe'
import type { ArtalkLightboxPluginOptions } from './main'

const winKeys = {
  lightGallery: 'lightGallery',
  lightBox: 'lightbox',
  fancyBox: 'Fancybox',
  photoSwipe: 'PhotoSwipeLightbox',
}

let deps: {
  lightGallery?: typeof lightGallery
  lightBox?: any
  fancyBox?: any
  photoSwipe?: typeof PhotoSwipeLightbox
}

export async function getDeps(opts?: ArtalkLightboxPluginOptions) {
  if (deps) return deps
  deps = {}
  for (const key in winKeys) {
    try {
      if (typeof window !== 'undefined' && window[winKeys[key]]) {
        deps[key] = window[winKeys[key]]
      } else {
        deps[key] = await opts?.[key]?.lib?.()
      }
      if (deps[key].default) deps[key] = deps[key].default
      if (deps[key]) return deps
    } catch {
      void 0
    }
  }
  return deps
}

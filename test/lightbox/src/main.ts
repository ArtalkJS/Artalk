import Artalk from 'artalk'
import 'artalk/Artalk.css'
import { ArtalkLightboxPlugin } from '@artalk/plugin-lightbox'

// @ts-ignore
import jQuery from 'jquery'

// @ts-ignore
window.$ = jQuery

const currentLib = new URLSearchParams(window.location.search).get('lib')
const libs = ['lightGallery', 'lightBox', 'fancyBox', 'photoSwipe']

document.querySelector<HTMLDivElement>('#app')!.innerHTML = `
  <div class="lib-switcher">
    ${libs
      .map(
        (lib) => `
      <button onclick="window.location
        .search = '?lib=${lib}'"
        class="${lib === currentLib ? 'active' : ''}">
        ${lib}
      </button>
    `,
      )
      .join('')}
  </div>
  <div class="artalk-container"></div>
`

if (currentLib === 'lightGallery') {
  import('lightgallery/css/lightgallery.css')

  Artalk.use(ArtalkLightboxPlugin, {
    lightGallery: {
      lib: () => import('lightgallery'),
    },
  })
} else if (currentLib === 'lightBox') {
  import('lightbox2/dist/css/lightbox.min.css')

  Artalk.use(ArtalkLightboxPlugin, {
    lightBox: {
      // @ts-ignore
      lib: () => import('lightbox2'),
    },
  })
} else if (currentLib === 'fancyBox') {
  import('fancybox/dist/css/jquery.fancybox.css')

  Artalk.use(ArtalkLightboxPlugin, {
    fancyBox: {
      // @ts-ignore
      lib: () => import('fancybox'),
    },
  })
} else if (currentLib === 'photoSwipe') {
  import('photoswipe/style.css')

  Artalk.use(ArtalkLightboxPlugin, {
    photoSwipe: {
      lib: () => import('photoswipe/lightbox'),
      pswpModule: () => import('photoswipe'),
    },
  })
}

Artalk.init({
  el: '.artalk-container',
  server: 'http://localhost:23366',
  site: 'ArtalkDocs',
  pageKey: 'https://artalk.js.org/guide/intro.html',
})

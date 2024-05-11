import { marked } from 'marked'
import type { ArtalkConfig } from '@/types'
import { renderCode } from './highlight'

export interface RendererOptions {
  imgLazyLoad: ArtalkConfig['imgLazyLoad']
}

export function getRenderer(options: RendererOptions) {
  const renderer = new marked.Renderer()
  renderer.link = markedLinkRenderer(renderer, renderer.link)
  renderer.code = markedCodeRenderer()
  renderer.image = markedImageRenderer(renderer, renderer.image, options)
  return renderer
}

const markedLinkRenderer =
  (renderer: any, orgLinkRenderer: Function) =>
  (href: string, title: string, text: string): string => {
    const getLinkOrigin = (link: string) => {
      try {
        return new URL(link).origin
      } catch {
        return ''
      }
    }
    const isSameOriginLink = getLinkOrigin(href) === window.location.origin
    const html = orgLinkRenderer.call(renderer, href, title, text)
    return html.replace(
      /^<a /,
      `<a target="_blank" ${!isSameOriginLink ? `rel="noreferrer noopener nofollow"` : ''} `,
    )
  }

const markedCodeRenderer =
  () =>
  (block: string, lang: string | undefined): string => {
    // Colorize the block only if the language is known to highlight.js
    const realLang = !lang ? 'plaintext' : lang
    let colorized = block
    if ((window as any).hljs) {
      if (realLang && (window as any).hljs.getLanguage(realLang)) {
        colorized = (window as any).hljs.highlight(realLang, block).value
      }
    } else {
      colorized = renderCode(block)
    }

    return (
      `<pre rel="${realLang}">\n` +
      `<code class="hljs language-${realLang}">${colorized.replace(/&amp;/g, '&')}</code>\n` +
      `</pre>`
    )
  }

const markedImageRenderer =
  (renderer: any, orgImageRenderer: Function, { imgLazyLoad }: RendererOptions) =>
  (href: string, title: string | null, text: string): string => {
    const html = orgImageRenderer.call(renderer, href, title, text)
    if (!imgLazyLoad) return html
    if (imgLazyLoad === 'native' || (imgLazyLoad as any) === true)
      return html.replace(/^<img /, '<img class="lazyload" loading="lazy" ')
    if (imgLazyLoad === 'data-src')
      return html.replace(/^<img /, '<img class="lazyload" ').replace('src=', 'data-src=')
    return html
  }

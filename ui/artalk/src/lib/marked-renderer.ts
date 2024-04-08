import { marked as libMarked } from 'marked'
import { renderCode } from './highlight'

export function getRenderer() {
  const renderer = new libMarked.Renderer()
  renderer.link = markedLinkRenderer(renderer, renderer.link)
  renderer.code = markedCodeRenderer()
  return renderer
}

export const markedLinkRenderer =
  (renderer: any, orgLinkRenderer: Function) =>
  (href: string, title: string, text: string): string => {
    const localLink = href?.startsWith(`${window.location.protocol}//${window.location.hostname}`)
    const html = orgLinkRenderer.call(renderer, href, title, text)
    return html.replace(
      /^<a /,
      `<a target="_blank" ${!localLink ? `rel="noreferrer noopener nofollow"` : ''} `,
    )
  }

export const markedCodeRenderer =
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

import { marked as libMarked } from 'marked'
import insane from 'insane'
import hanabi from 'hanabi'
import Context from '~/types/context'

let instance: (typeof libMarked)|undefined

export type TMarked = typeof libMarked

/** Get Marked instance */
export function getInstance() {
  return instance
}

/** 初始化 marked */
export function initMarked() {
  try { if (!libMarked.name) return } catch { return }

  const renderer = new libMarked.Renderer()
  const orgLinkRenderer = renderer.link
  renderer.link = (href, title, text) => {
    const localLink = href?.startsWith(`${window.location.protocol}//${window.location.hostname}`);
    const html = orgLinkRenderer.call(renderer as any, href, title, text);
    return html.replace(/^<a /, `<a target="_blank" ${!localLink ? `rel="noreferrer noopener nofollow"` : ''} `);
  }

  renderer.code = (block, lang) => {
    // Colorize the block only if the language is known to highlight.js
    const realLang = (!lang ? 'plaintext' : lang)
    let colorized = block
    if ((window as any).hljs) {
      if (realLang && (window as any).hljs.getLanguage(realLang)) {
        colorized = (window as any).hljs.highlight(realLang, block).value
      }
    } else {
      colorized = hanabi(block)
    }

    return `<pre rel="${realLang}">\n`
      + `<code class="hljs language-${realLang}">${colorized.replace(/&amp;/g, '&')}</code>\n`
      + `</pre>`
  }

  // @see https://github.com/markedjs/marked/blob/4afb228d956a415624c4e5554bb8f25d047676fe/src/Tokenizer.js#L329
  const nMarked = libMarked
  libMarked.setOptions({
    renderer,
    pedantic: false,
    gfm: true,
    breaks: true,
    smartLists: true,
    smartypants: true,
    xhtml: false,
    sanitize: false,
    silent: true,
  })

  instance = nMarked
}

/** 解析 markdown */
export default function marked(ctx: Context, src: string): string {
  let markedContent = getInstance()?.parse(src)
  if (!markedContent) {
    // 无 Markdown 模式简单处理
    markedContent = src
      // .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
      .replace(/```\s*([^]+?.*?[^]+?[^]+?)```/g, (_, code) => `<pre><code>${hanabi(code)}</code></pre>`)
      // .replace(/`([^`]+?)`/g, '<code>$1</code>')
      .replace(/!\[(.*?)\]\((.*?)\)/g, (_, alt, imgSrc) => `<img src="${imgSrc}" alt="${alt}" />`)
      .replace(/\[(.*?)\]\((.*?)\)/g, (_, text, link) => `<a href="${link}" target="_blank">${text}</a>`)
      .replace(/\n/g, '<br>')
  }

  // @link https://github.com/markedjs/marked/discussions/1232
  // @link https://gist.github.com/lionel-rowe/bb384465ba4e4c81a9c8dada84167225
  let dest = insane(markedContent, {
    allowedClasses: {},
    // @refer CVE-2018-8495
    // @link https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2018-8495
    // @link https://leucosite.com/Microsoft-Edge-RCE/
    // @link https://medium.com/@knownsec404team/analysis-of-the-security-issues-of-url-scheme-in-pc-from-cve-2018-8495-934478a36756
    allowedSchemes: [
      'http', 'https', 'mailto',
      'data' // for support base64 encoded image (安全性有待考虑)
    ],
    allowedTags: [
      'a', 'abbr', 'article', 'b', 'blockquote', 'br', 'caption', 'code', 'del', 'details', 'div', 'em',
      'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'hr', 'i', 'img', 'ins', 'kbd', 'li', 'main', 'mark',
      'ol', 'p', 'pre', 'section', 'span', 'strike', 'strong', 'sub', 'summary', 'sup', 'table',
      'tbody', 'td', 'th', 'thead', 'tr', 'u', 'ul'
    ],
    allowedAttributes: {
      '*': ['title', 'accesskey'],
      a: ['href', 'name', 'target', 'aria-label', 'rel'],
      img: ['src', 'alt', 'title', 'atk-emoticon', 'aria-label'],
      // for code highlight
      code: ['class'],
      span: ['class', 'style'],
    },
    filter: node => {
      // allow hljs style
      const allowed = [
        [ 'code', /^hljs\W+language-(.*)$/ ],
        [ 'span', /^(hljs-.*)$/ ]
      ]
      allowed.forEach(([ tag, reg ]) => {
        if (
          node.tag === tag
          && !!node.attrs.class
          && !(reg as RegExp).test(node.attrs.class)
        ) {
          delete node.attrs.class
        }
      })

      // allow <span> set color sty
      if (node.tag === 'span' && !!node.attrs.style
          && !/^color:(\W+)?#[0-9a-f]{3,6};?$/i.test(node.attrs.style)) {
        delete node.attrs.style
      }

      return true
    }
  })

  // 内容替换器
  ctx.markedReplacers.forEach((replacer) => {
    if (typeof replacer === 'function') dest = replacer(dest)
  })

  return dest
}

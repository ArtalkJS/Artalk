import { marked as libMarked, MarkedOptions } from 'marked'

import { sanitize } from './sanitizer'
import { renderCode } from './highlight'
import { getRenderer } from './marked-renderer'

type Replacer = (raw: string) => string
export type TMarked = typeof libMarked

let instance: TMarked | undefined
let replacers: Replacer[] = []

const markedOptions: MarkedOptions = {
  gfm: true,
  breaks: true,
  async: false,
}

/** Get Marked instance */
export function getInstance() {
  return instance
}

export function setReplacers(arr: Replacer[]) {
  replacers = arr
}

/** 初始化 marked */
export function initMarked() {
  try {
    if (!libMarked.name) return
  } catch {
    return
  }

  // @see https://github.com/markedjs/marked/blob/4afb228d956a415624c4e5554bb8f25d047676fe/src/Tokenizer.js#L329
  libMarked.setOptions({
    renderer: getRenderer(),
    ...markedOptions,
  })

  instance = libMarked
}

/** 解析 markdown */
export default function marked(src: string): string {
  let markedContent = getInstance()?.parse(src) as string
  if (!markedContent) {
    // 无 Markdown 模式简单处理
    markedContent = simpleMarked(src)
  }

  let dest = sanitize(markedContent)

  // 内容替换器
  replacers.forEach((replacer) => {
    if (typeof replacer === 'function') dest = replacer(dest)
  })

  return dest
}

function simpleMarked(src: string) {
  return (
    src
      // .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
      .replace(
        /```\s*([^]+?.*?[^]+?[^]+?)```/g,
        (_, code) => `<pre><code>${renderCode(code)}</code></pre>`,
      )
      // .replace(/`([^`]+?)`/g, '<code>$1</code>')
      .replace(/!\[(.*?)\]\((.*?)\)/g, (_, alt, imgSrc) => `<img src="${imgSrc}" alt="${alt}" />`)
      .replace(
        /\[(.*?)\]\((.*?)\)/g,
        (_, text, link) => `<a href="${link}" target="_blank">${text}</a>`,
      )
      .replace(/\n/g, '<br>')
  )
}

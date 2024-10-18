import { Marked } from 'marked'
import type { MarkedOptions } from 'marked'

import { sanitize } from './sanitizer'
import { renderCode } from './highlight'
import { getRenderer } from './marked-renderer'
import type { Config } from '@/types'

type Replacer = (raw: string) => string

let instance: Marked | undefined
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

export interface MarkedInitOptions {
  markedOptions: Config['markedOptions']
  imgLazyLoad: Config['imgLazyLoad']
}

/** 初始化 marked */
export function initMarked(options: MarkedInitOptions) {
  try {
    if (!Marked.name) return
  } catch {
    return
  }

  // @see https://github.com/markedjs/marked/blob/4afb228d956a415624c4e5554bb8f25d047676fe/src/Tokenizer.js#L329
  const marked = new Marked()
  marked.setOptions({
    renderer: getRenderer({
      imgLazyLoad: options.imgLazyLoad,
    }),
    ...markedOptions,
    ...options.markedOptions,
  })

  instance = marked
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

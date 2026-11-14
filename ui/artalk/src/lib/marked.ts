import { Marked } from 'marked'
import type { MarkedOptions } from 'marked'

import { sanitize } from './sanitizer'
import { renderCode } from './highlight'
import { getRenderer } from './marked-renderer'
import type { Config } from '@/types'

type Replacer = (raw: string) => string

export interface MarkedInitOptions {
  markedOptions?: Config['markedOptions']
  imgLazyLoad: Config['imgLazyLoad']
}

/** Each widget owns its parser. Standalone callers keep the legacy exports below. */
export function createMarkdown(getReplacers: () => Replacer[] | undefined = () => []) {
  let instance: Marked | undefined
  let defaults: MarkedOptions
  let configured: MarkedOptions = {}
  const rendererOptions = { imgLazyLoad: false as Config['imgLazyLoad'] }

  return {
    getInstance: () => instance,
    configure(options: MarkedInitOptions) {
      // Artalk Lite can run without the optional external Marked library.
      try {
        if (!Marked.name) return
      } catch {
        return
      }
      rendererOptions.imgLazyLoad = options.imgLazyLoad
      if (!instance) {
        instance = new Marked()
        instance.setOptions({
          renderer: getRenderer(rendererOptions),
          gfm: true,
          breaks: true,
          async: false,
        })
        defaults = { ...instance.defaults }
      }
      const next = options.markedOptions || {}
      // Reset removed config options, without replacing plugin extensions or
      // renderers when unrelated config changes.
      for (const key of Object.keys(configured)) {
        if (!(key in next)) instance.setOptions({ [key]: defaults[key] })
      }
      for (const key of Object.keys(next)) {
        if (next[key] === configured[key]) continue
        let value = next[key]
        // Marked writes parser state to these objects. A config object may be
        // reused by multiple widgets, so give each parser its own copy.
        if (value && ['renderer', 'tokenizer', 'hooks'].includes(key)) {
          value = Object.assign(Object.create(Object.getPrototypeOf(value)), value)
        }
        instance.setOptions({ [key]: value })
      }
      configured = next
    },
    render(src: string): string {
      const html = instance?.parse(src) as string | undefined
      let dest = sanitize(html || simpleMarked(src))
      getReplacers()?.forEach((replace) => {
        if (typeof replace === 'function') dest = replace(dest)
      })
      return dest
    },
  }
}

let replacers: Replacer[] = []
let standalone = createMarkdown(() => replacers)

export const getInstance = () => standalone.getInstance()
export function setReplacers(arr: Replacer[]) {
  replacers = arr
}
export function initMarked(options: MarkedInitOptions) {
  standalone = createMarkdown(() => replacers)
  standalone.configure(options)
}
export default function marked(src: string): string {
  return standalone.render(src)
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

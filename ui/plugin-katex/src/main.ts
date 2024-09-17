import { ArtalkPlugin } from 'artalk'
import katex from 'katex'
import type { TokenizerExtension } from 'marked'

if (import.meta.env.DEV) {
  // only dev mode, because some client may not support or require css import
  import('katex/dist/katex.min.css')
}

const inlineMathStart = /\$.*?\$/
const inlineMathReg = /^\$(.*?)\$/
const blockMathReg = /^(?:\s{0,3})\$\$((?:[^\n]|\n[^\n])+?)\n{0,1}\$\$/

export const ArtalkKatexPlugin: ArtalkPlugin = (ctx) => {
  // Using placeholder to replace the tex expressions in order to bypass the HTML sanitization
  let curtIdx = 0
  const next = () => `__atk_katex_id_${curtIdx++}__`
  const texExprs: {
    [key: string]: { isBlock: boolean; tex: string }
  } = {}

  const getPlaceholder = (tex: string, isBlock: boolean) => {
    const key = next()
    texExprs[key] = {
      isBlock,
      tex,
    }
    return key
  }

  const blockMathExtension: TokenizerExtension = {
    name: 'blockMath',
    level: 'block',
    tokenizer: (src: string) => {
      const cap = blockMathReg.exec(src)

      if (cap) {
        return {
          type: 'html',
          raw: cap[0],
          text: getPlaceholder(cap[1], true),
        }
      }

      return undefined
    },
  }

  const inlineMathExtension: TokenizerExtension = {
    name: 'inlineMath',
    level: 'inline',
    start: (src: string) => {
      const idx = src.search(inlineMathStart)
      return idx !== -1 ? idx : src.length
    },
    tokenizer: (src: string) => {
      const cap = inlineMathReg.exec(src)

      if (cap) {
        return {
          type: 'html',
          raw: cap[0],
          text: getPlaceholder(cap[1], false),
        }
      }

      return undefined
    },
  }

  ctx.on('mounted', () => {
    const markedInstance = ctx.getMarked() // must be called after `mounted` event
    if (!markedInstance) {
      console.error('[artalk-plugin-katex] no marked instance found in artalk context')
      return
    }

    if (typeof katex === 'undefined') {
      console.error(
        '[artalk-plugin-katex] katex not found, please make sure you have imported katex in your project',
      )
      return
    }

    markedInstance.use({
      extensions: [blockMathExtension, inlineMathExtension],
    })

    ctx.updateConf({
      markedReplacers: [
        (text) => {
          text = text.replace(/(__atk_katex_id_\d+__)/g, (_match, key) => {
            const { tex, isBlock } = texExprs[key]

            try {
              return katex.renderToString(tex, {
                displayMode: isBlock,
              })
            } catch (e) {
              console.error('[artalk-plugin-katex] failed to render katex:', e)
              return `<code>${e}</code>`
            }
          })

          return text
        },
      ],
    })
  })
}

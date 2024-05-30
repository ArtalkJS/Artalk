import { ArtalkPlugin } from 'artalk'
import katex from 'katex'

if (import.meta.env.DEV) {
  import('katex/dist/katex.min.css')
}

export const ArtalkKatexPlugin: ArtalkPlugin = (ctx) => {
  ctx.on('mounted', () => {
    // @see https://github.com/markedjs/marked/issues/1538#issuecomment-575838181
    const markedInstance = ctx.getMarked() // must be called after `mounted` event
    if (!markedInstance) {
      console.error('[artalk-plugin-katex] no marked instance found')
      return
    }

    let i = 0
    const nextID = () => `__atk_katex_id_${i++}__`
    const mathExpressions: {
      [key: string]: { type: 'block' | 'inline'; expression: string }
    } = {}

    function replaceMathWithIds(text: string) {
      // Allowing newlines inside of `$$...$$`
      text = text.replace(/\$\$([\s\S]+?)\$\$/g, (_match, expression) => {
        const id = nextID()
        mathExpressions[id] = { type: 'block', expression }
        return id
      })

      // Not allowing newlines or space inside of `$...$`
      text = text.replace(/\$([^\n]+?)\$/g, (_match, expression) => {
        const id = nextID()
        mathExpressions[id] = { type: 'inline', expression }
        return id
      })

      return text
    }

    // Marked render
    const renderer = new markedInstance.Renderer() as any

    const orgListitem = renderer.listitem
    const orgParagraph = renderer.paragraph
    const orgTablecell = renderer.tablecell
    const orgCodespan = renderer.codespan
    const orgText = renderer.text

    renderer.listitem = (text: string, task: boolean, checked: boolean) =>
      orgListitem(replaceMathWithIds(text), task, checked)
    renderer.paragraph = (text: string) => orgParagraph(replaceMathWithIds(text))
    renderer.tablecell = (content: string, flags: any) =>
      orgTablecell(replaceMathWithIds(content), flags)
    renderer.codespan = (code: string) => orgCodespan(replaceMathWithIds(code))
    renderer.text = (text: string) => orgText(replaceMathWithIds(text)) // Inline level, maybe unneeded

    ctx.updateConf({
      markedReplacers: [
        (text) => {
          text = text.replace(/(__atk_katex_id_\d+__)/g, (_match, capture) => {
            const v = mathExpressions[capture]
            const type = v.type
            let expression = v.expression

            // replace <br/> tag to \n
            expression = expression.replace(/<br\s*\/?>/gm, '\n')

            return katex.renderToString(expression, {
              displayMode: type === 'block',
            })
          })

          return text
        },
      ],
    })

    markedInstance.use({
      renderer,
    })
  })
}

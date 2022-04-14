/* eslint-disable guard-for-in */
/* eslint-disable no-restricted-syntax */
import Artalk from 'artalk'
import katex from 'katex'

Artalk.Use((ctx) => {
  const renderer: any = {}

  function renderMathsExpression(expr) {
    if (expr[0] === '$' && expr[expr.length - 1] === '$') {
      let displayStyle = false
      expr = expr.substr(1, expr.length - 2)
      if (expr[0] === '$' && expr[expr.length - 1] === '$') {
        displayStyle = true
        expr = expr.substr(1, expr.length - 2)
      }
      let html: any = null
      try {
        html = katex.renderToString(expr)
      } catch (e) {
        console.error(e)
      }
      if (displayStyle && html) {
        html = html.replace(/class="katex"/g, 'class="katex katex-block" style="display: block;"')
      }
      return html
    }

    return null
  }

  renderer.paragraph = (text) => {
    const blockRegex = /\$\$[^$]*\$\$/g
    const inlineRegex = /\$[^$]*\$/g
    const blockExprArray = text.match(blockRegex)
    const inlineExprArray = text.match(inlineRegex)

    for (const i in blockExprArray) {
      const expr = blockExprArray[i]
      const result = renderMathsExpression(expr)
      text = text.replace(expr, result)
    }

    for (const i in inlineExprArray) {
      const expr = inlineExprArray[i]
      const result = renderMathsExpression(expr)
      text = text.replace(expr, result)
    }

    return text
  }

  ctx.markedInstance.use({ renderer });
})

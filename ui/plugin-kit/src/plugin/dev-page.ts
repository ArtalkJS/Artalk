import type { Connect, HtmlTagDescriptor } from 'vite'
import DevHTML from '../../../artalk/index.html?raw'
import { RUNTIME_PATH } from './runtime-helper'

function patchArtalkInitOptionsInHTML(html: string, options: Record<string, any>) {
  const jsonStr = JSON.stringify(options, null, 8).slice(2, -2)
  const optionsStr = jsonStr ? `${jsonStr},` : ''
  html = html.replace(/(Artalk\.init\({[\s\S]*?)(\n\s*\}\))/, (match, p1, p2) => {
    const cleanedP1 = p1.replace(/,\s*$/, '') // remove last comma
    return `${cleanedP1},\n${optionsStr}${p2}`
  })
  return html
}

export function hijackIndexPage(
  middlewares: Connect.Server,
  transformIndexHtml: (url: string, html: string) => Promise<string>,
  artalkInitOptions: Record<string, any>,
) {
  middlewares.use(async (req, res, next) => {
    // only handle the index page
    let [url] = (req.originalUrl || '').split('?')
    if (url.endsWith('/')) url += 'index.html'
    if (url !== '/index.html') return next()

    try {
      let html = DevHTML.replace(/import.*?from.*/, '') // remove imports
      html = patchArtalkInitOptionsInHTML(html, artalkInitOptions)
      html = await transformIndexHtml('/', html)
      res.end(html)
    } catch (e) {
      console.log(e)
      res.statusCode = 500
      res.end(e)
    }

    return undefined
  })
}

export function getInjectHTMLTags(entrySrc: string): HtmlTagDescriptor[] {
  return [
    // artalk code
    {
      tag: 'script',
      attrs: { type: 'module' },
      children: `
        import Artalk from '/node_modules/artalk'
        import '/node_modules/artalk/dist/Artalk.css'
        window.Artalk = Artalk
      `,
      injectTo: 'head',
    },

    // plugin-kit-runtime code
    {
      tag: 'script',
      attrs: { type: 'module' },
      children: `
        import { inject } from "${RUNTIME_PATH}";
        inject({
          config: ${JSON.stringify({ test: 'Hello World' })},
        });
      `,
    },

    // entry code
    {
      tag: 'script',
      attrs: { type: 'module', src: entrySrc },
      children: '',
      injectTo: 'head',
    },
  ]
}

import insane from 'insane'

const insaneOptions = {
  allowedClasses: {},
  // @refer CVE-2018-8495
  // @link https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2018-8495
  // @link https://leucosite.com/Microsoft-Edge-RCE/
  // @link https://medium.com/@knownsec404team/analysis-of-the-security-issues-of-url-scheme-in-pc-from-cve-2018-8495-934478a36756
  allowedSchemes: [
    'http',
    'https',
    'mailto',
    'data', // for support base64 encoded image (安全性有待考虑)
  ],
  allowedTags: [
    'a',
    'abbr',
    'article',
    'b',
    'blockquote',
    'br',
    'caption',
    'code',
    'del',
    'details',
    'div',
    'em',
    'h1',
    'h2',
    'h3',
    'h4',
    'h5',
    'h6',
    'hr',
    'i',
    'img',
    'ins',
    'kbd',
    'li',
    'main',
    'mark',
    'ol',
    'p',
    'pre',
    'section',
    'span',
    'strike',
    'strong',
    'sub',
    'summary',
    'sup',
    'table',
    'tbody',
    'td',
    'th',
    'thead',
    'tr',
    'u',
    'ul',
  ],
  allowedAttributes: {
    '*': ['title', 'accesskey'],
    a: ['href', 'name', 'target', 'aria-label', 'rel'],
    img: ['src', 'alt', 'title', 'atk-emoticon', 'aria-label', 'data-src', 'class', 'loading'],
    // for code highlight
    code: ['class'],
    span: ['class', 'style'],
  },
  filter: (node) => {
    // class whitelist
    const allowed = [
      ['code', /^hljs\W+language-(.*)$/],
      ['span', /^(hljs-.*)$/],
      ['img', /^lazyload$/],
    ]
    allowed.forEach(([tag, reg]) => {
      if (node.tag === tag && !!node.attrs.class && !(reg as RegExp).test(node.attrs.class)) {
        delete node.attrs.class
      }
    })

    // allow <span> set color sty
    if (
      node.tag === 'span' &&
      !!node.attrs.style &&
      !/^color:(\W+)?#[0-9a-f]{3,6};?$/i.test(node.attrs.style)
    ) {
      delete node.attrs.style
    }

    return true
  },
}

export function sanitize(content: string): string {
  // @link https://github.com/markedjs/marked/discussions/1232
  // @link https://gist.github.com/lionel-rowe/bb384465ba4e4c81a9c8dada84167225
  return insane(content, insaneOptions)
}

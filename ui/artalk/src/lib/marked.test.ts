import { beforeEach, describe, expect, it } from 'vitest'
import renderMarkdown, { getInstance, initMarked, setReplacers } from './marked'

describe('Marked integration', () => {
  beforeEach(() => {
    initMarked({ markedOptions: {}, imgLazyLoad: 'native' })
    setReplacers([])
  })

  it('renders and sanitizes links, code, images, and lists', () => {
    const html = renderMarkdown(`
[external](https://example.com)

~~~ts
const answer = 42
~~~

![alt](https://example.com/image.png)

- first
- second
`)

    expect(html).toContain('target="_blank"')
    expect(html).toContain('rel="noreferrer noopener nofollow ugc"')
    expect(html).toContain('<code class="hljs language-ts">')
    expect(html).toContain('class="lazyload" loading="lazy"')
    expect(html).toContain('<ul>')
    expect(html).toContain('<li>first</li>')
    expect(html).toContain('<li>second</li>')
  })

  it('supports the tokenizer extension shape used by the KaTeX plugin', () => {
    const delimiter = '$'
    getInstance()?.use({
      extensions: [
        {
          name: 'inlineMathFixture',
          level: 'inline',
          start: (src) => src.indexOf(delimiter),
          tokenizer: (src) => {
            const raw = `${delimiter}x${delimiter}`
            if (!src.startsWith(raw)) return undefined
            return { type: 'html', raw, text: '<span>math</span>' }
          },
        },
      ],
    })

    expect(renderMarkdown(`${delimiter}x${delimiter}`)).toContain('<span>math</span>')
  })
})

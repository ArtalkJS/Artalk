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

  it.each(['native', 'data-src'] as const)(
    'filters data URLs through the comment/preview renderer with %s images',
    (imgLazyLoad) => {
      initMarked({ markedOptions: {}, imgLazyLoad })
      const html = renderMarkdown(`
[link](data:text/html;base64,PHNjcmlwdD5hbGVydCgxKTwvc2NyaXB0Pg==)

<a href="data:text/html,hello">raw HTML link</a>

![unsafe](data:image/svg+xml,%3Csvg%3E)

![safe](data:image/png;base64,AA==)
`)
      expect(html).not.toContain('data:text/html')
      expect(html).not.toContain('data:image/svg+xml')
      expect(html).toContain(
        `${imgLazyLoad === 'data-src' ? 'data-src' : 'src'}="data:image/png;base64,AA=="`,
      )
      expect(html).toContain('raw HTML link</a>')
    },
  )
})

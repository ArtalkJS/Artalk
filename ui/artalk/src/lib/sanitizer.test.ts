import { describe, expect, it } from 'vitest'
import { sanitize } from './sanitizer'

describe('HTML sanitizer', () => {
  it('keeps allowed tags, attributes, classes, styles, and schemes', () => {
    const html = sanitize(`
      <a href="mailto:test@example.com" target="_blank" class="drop">mail</a>
      <a href="https://example.com">web</a>
      <code class="hljs language-ts">code</code>
      <span class="hljs-keyword" style="color: #abc;">keyword</span>
      <img src="data:image/png;base64,AA==" class="lazyload" loading="lazy" alt="image">
    `)

    expect(html).toContain('<a href="mailto:test@example.com" target="_blank">mail</a>')
    expect(html).toContain('<a href="https://example.com">web</a>')
    expect(html).toContain('<code class="hljs language-ts">code</code>')
    expect(html).toContain('<span class="hljs-keyword" style="color: #abc;">keyword</span>')
    expect(html).toContain(
      '<img src="data:image/png;base64,AA==" class="lazyload" loading="lazy" alt="image">',
    )
  })

  it.each([
    'data:text/html;base64,PHNjcmlwdD5hbGVydCgxKTwvc2NyaXB0Pg==',
    'data:text/plain,hello',
    'data:image/png;base64,AA==',
    'data:,hello',
    'DATA:text/html,hello',
    'data&#58;text/html,hello',
    'd&#x09;ata:text/html,hello',
    ' &#x0a;data:text/html,hello',
  ])('removes data URLs from links: %s', (url) => {
    expect(sanitize(`<a href="${url}" target="_blank">link</a>`)).toBe(
      '<a target="_blank">link</a>',
    )
  })

  describe.each(['src', 'data-src'])('image %s', (attribute) => {
    it.each(['png', 'jpeg', 'gif', 'webp', 'avif', 'bmp'])(
      'preserves inline %s images',
      (format) => {
        const url = `data:image/${format};base64,AA==`
        expect(sanitize(`<img ${attribute}="${url}" alt="image">`)).toBe(
          `<img ${attribute}="${url}" alt="image">`,
        )
      },
    )

    it.each(['data:image/png,%89PNG', 'data:IMAGE/PNG;BASE64,AA=='])(
      'preserves percent-encoded or mixed-case image URLs: %s',
      (url) => {
        expect(sanitize(`<img ${attribute}="${url}">`)).toBe(`<img ${attribute}="${url}">`)
      },
    )

    it.each([
      'data:text/html,hello',
      'data:text/javascript,alert(1)',
      'data:application/javascript,alert(1)',
      'data:image/svg+xml,%3Csvg%20onload=alert(1)%3E',
      'data:,hello',
      'data:image/png+xml,hello',
      'data:image/png;unknown,hello',
      'data:image/png;base64',
      'DATA:text/html,hello',
      'd&#x09;ata:text/html,hello',
      'javascript:alert(1)',
      'java&#x0a;script:alert(1)',
      'vbscript:msgbox(1)',
      'blob:https://example.com/id',
    ])('removes non-raster or malformed data URLs: %s', (url) => {
      expect(sanitize(`<img ${attribute}="${url}" alt="image">`)).toBe('<img alt="image">')
    })
  })

  it.each([
    'http://example.com',
    'https://example.com',
    'mailto:test@example.com',
    '/path',
    '//example.com/path',
    '#reply',
    '?page=2',
  ])('preserves ordinary links: %s', (url) => {
    expect(sanitize(`<a href="${url}">link</a>`)).toBe(`<a href="${url}">link</a>`)
  })

  it('removes executable markup and attributes outside the per-tag allowlist', () => {
    const html = sanitize(`
      <script>alert(1)</script>
      <a href="java&#x0a;script:alert(1)" style="color: #fff">bad link</a>
      <img src="javascript:alert(1)" href="https://example.com" onerror="alert(1)">
      <span class="not-highlight" style="background: url(javascript:alert(1))">text</span>
    `)

    expect(html).not.toContain('script')
    expect(html).not.toContain('href=')
    expect(html).not.toContain('src=')
    expect(html).not.toContain('style=')
    expect(html).not.toContain('class=')
    expect(html).not.toContain('onerror')
  })

  it('keeps only explicitly allowlisted ARIA and data attributes', () => {
    const html = sanitize(`
      <a aria-label="permalink" aria-hidden="true" data-test="drop">link</a>
      <img aria-label="preview" data-src="https://example.com/image.png" data-test="drop">
    `)

    expect(html).toContain('<a aria-label="permalink">link</a>')
    expect(html).toContain('<img aria-label="preview" data-src="https://example.com/image.png">')
    expect(html).not.toContain('aria-hidden')
    expect(html).not.toContain('data-test')
  })

  it('handles the former insane parser ReDoS payload', () => {
    const payload = `<b>foo</b><foo bar= ${' -=""'.repeat(50)}`
    expect(sanitize(payload)).toContain('<b>foo</b>')
  }, 1000)
})

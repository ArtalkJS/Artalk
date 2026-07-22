import { describe, expect, it } from 'vitest'
import { sanitize } from './sanitizer'

describe('HTML sanitizer', () => {
  it('keeps the existing tags, attributes, classes, styles, and schemes', () => {
    const html = sanitize(`
      <a href="mailto:test@example.com" target="_blank" class="drop">mail</a>
      <a href="data:text/plain,hello">data</a>
      <code class="hljs language-ts">code</code>
      <span class="hljs-keyword" style="color: #abc;">keyword</span>
      <img src="data:image/png;base64,AA==" class="lazyload" loading="lazy" alt="image">
    `)

    expect(html).toContain('<a href="mailto:test@example.com" target="_blank">mail</a>')
    expect(html).toContain('<a href="data:text/plain,hello">data</a>')
    expect(html).toContain('<code class="hljs language-ts">code</code>')
    expect(html).toContain('<span class="hljs-keyword" style="color: #abc;">keyword</span>')
    expect(html).toContain(
      '<img src="data:image/png;base64,AA==" class="lazyload" loading="lazy" alt="image">',
    )
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

import { marked as libMarked } from 'marked'
import insane from 'insane'
import hanabi from 'hanabi'
import Context from '~/types/context'

export function createElement<E extends HTMLElement = HTMLElement>(htmlStr: string = ''): E {
  const div = document.createElement('div')
  div.innerHTML = htmlStr.trim()
  return (div.firstElementChild || div) as E
}

export function getHeight(el: HTMLElement) {
  return parseFloat(getComputedStyle(el, null).height.replace('px', ''))
}

export function htmlEncode(str: string) {
  const temp = document.createElement('div')
  temp.innerText = str
  const output = temp.innerHTML
  return output
}

export function htmlDecode(str: string) {
  const temp = document.createElement('div')
  temp.innerHTML = str
  const output = temp.innerText
  return output
}

export function getQueryParam(name: string) {
  const match = RegExp(`[?&]${name}=([^&]*)`).exec(window.location.search);
  return match && decodeURIComponent(match[1].replace(/\+/g, ' '));
}

export function getOffset(el: HTMLElement) {
  const rect = el.getBoundingClientRect()
  return {
    top: rect.top + window.scrollY,
    left: rect.left + window.scrollX
  }
}

export function padWithZeros(vNumber: number, width: number) {
  let numAsString = vNumber.toString()
  while (numAsString.length < width) {
    numAsString = `0${numAsString}`
  }
  return numAsString
}

export function dateFormat(date: Date) {
  const vDay = padWithZeros(date.getDate(), 2)
  const vMonth = padWithZeros(date.getMonth() + 1, 2)
  const vYear = padWithZeros(date.getFullYear(), 2)
  // var vHour = padWithZeros(date.getHours(), 2);
  // var vMinute = padWithZeros(date.getMinutes(), 2);
  // var vSecond = padWithZeros(date.getSeconds(), 2);
  return `${vYear}-${vMonth}-${vDay}`
}

export function timeAgo(date: Date, ctx: Context) {
  try {
    const oldTime = date.getTime()
    const currTime = new Date().getTime()
    const diffValue = currTime - oldTime

    const days = Math.floor(diffValue / (24 * 3600 * 1000))
    if (days === 0) {
      // 计算相差小时数
      const leave1 = diffValue % (24 * 3600 * 1000) // 计算天数后剩余的毫秒数
      const hours = Math.floor(leave1 / (3600 * 1000))
      if (hours === 0) {
        // 计算相差分钟数
        const leave2 = leave1 % (3600 * 1000) // 计算小时数后剩余的毫秒数
        const minutes = Math.floor(leave2 / (60 * 1000))
        if (minutes === 0) {
          // 计算相差秒数
          const leave3 = leave2 % (60 * 1000) // 计算分钟数后剩余的毫秒数
          const seconds = Math.round(leave3 / 1000)
          return `${seconds} ${ctx.$t('seconds')}`
        }
        return `${minutes} ${ctx.$t('minutes')}`
      }
      return `${hours} ${ctx.$t('hours')}`
    }
    if (days < 0) return ctx.$t('now')

    if (days < 8) {
      return `${days} ${ctx.$t('days')}`
    }

    return dateFormat(date)
  } catch (error) {
    console.error(error)
    return ' - '
  }
}

/** 所有图片加载完毕后执行 */
export function onImagesLoaded($container: HTMLElement, event: Function) {
  if (!$container) return
  const images = $container.getElementsByTagName('img')
  if (!images.length) return
  let loaded = images.length
  for (let i = 0; i < images.length; i++) {
    if (images[i].complete) {
      loaded--
    } else {
      // eslint-disable-next-line @typescript-eslint/no-loop-func
      images[i].addEventListener('load', () => {
        loaded--
        if (loaded === 0) event()
      })
    }

    if (loaded === 0) event()
  }
}

export function getGravatarURL(ctx: Context, emailMD5: string) {
  return `${(ctx.conf.gravatar.mirror).replace(/\/$/, '')}/${emailMD5}?d=${encodeURIComponent(ctx.conf.gravatar.default)}&s=80`
}

export function sleep(ms: number) {
  return new Promise(resolve => { setTimeout(resolve, ms) });
}

/** 版本号比较（a < b :-1 | 0 | b < a :1） */
export function versionCompare(a: string, b: string) {
  const pa = a.split('.')
  const pb = b.split('.')
  for (let i = 0; i < 3; i++) {
      const na = Number(pa[i])
      const nb = Number(pb[i])
      if (na > nb) return 1
      if (nb > na) return -1
      if (!Number.isNaN(na) && Number.isNaN(nb)) return 1
      if (Number.isNaN(na) && !Number.isNaN(nb)) return -1
  }
  return 0
}

/** 初始化 marked */
export function initMarked(ctx: Context) {
  try { if (!libMarked.name) return } catch { return }

  const renderer = new libMarked.Renderer()
  const orgLinkRenderer = renderer.link
  renderer.link = (href, title, text) => {
    const localLink = href?.startsWith(`${window.location.protocol}//${window.location.hostname}`);
    const html = orgLinkRenderer.call(renderer as any, href, title, text);
    return html.replace(/^<a /, `<a target="_blank" ${!localLink ? `rel="noreferrer noopener nofollow"` : ''} `);
  }

  renderer.code = (block, lang) => {
    // Colorize the block only if the language is known to highlight.js
    const realLang = (!lang ? 'plaintext' : lang)
    let colorized = block
    if ((window as any).hljs) {
      if (realLang && (window as any).hljs.getLanguage(realLang)) {
        colorized = (window as any).hljs.highlight(realLang, block).value
      }
    } else {
      colorized = hanabi(block)
    }

    return `<pre rel="${realLang}">\n`
      + `<code class="hljs language-${realLang}">${colorized.replace(/&amp;/g, '&')}</code>\n`
      + `</pre>`
  }

  // @see https://github.com/markedjs/marked/blob/4afb228d956a415624c4e5554bb8f25d047676fe/src/Tokenizer.js#L329
  const nMarked = libMarked
  libMarked.setOptions({
    renderer,
    pedantic: false,
    gfm: true,
    breaks: true,
    smartLists: true,
    smartypants: true,
    xhtml: false,
    sanitize: false,
    silent: true,
  })

  ctx.markedInstance = nMarked
}

/** 解析 markdown */
export function marked(ctx: Context, src: string): string {
  let markedContent = ctx.markedInstance?.parse(src)
  if (!markedContent) {
    // 无 Markdown 模式简单处理
    markedContent = src
      // .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
      .replace(/```\s*([^]+?.*?[^]+?[^]+?)```/g, (_, code) => `<pre><code>${hanabi(code)}</code></pre>`)
      // .replace(/`([^`]+?)`/g, '<code>$1</code>')
      .replace(/!\[(.*?)\]\((.*?)\)/g, (_, alt, imgSrc) => `<img src="${imgSrc}" alt="${alt}" />`)
      .replace(/\[(.*?)\]\((.*?)\)/g, (_, text, link) => `<a href="${link}" target="_blank">${text}</a>`)
      .replace(/\n/g, '<br>')
  }

  // @link https://github.com/markedjs/marked/discussions/1232
  // @link https://gist.github.com/lionel-rowe/bb384465ba4e4c81a9c8dada84167225
  let dest = insane(markedContent, {
    allowedClasses: {},
    // @refer CVE-2018-8495
    // @link https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2018-8495
    // @link https://leucosite.com/Microsoft-Edge-RCE/
    // @link https://medium.com/@knownsec404team/analysis-of-the-security-issues-of-url-scheme-in-pc-from-cve-2018-8495-934478a36756
    allowedSchemes: [
      'http', 'https', 'mailto',
      'data' // for support base64 encoded image (安全性有待考虑)
    ],
    allowedTags: [
      'a', 'abbr', 'article', 'b', 'blockquote', 'br', 'caption', 'code', 'del', 'details', 'div', 'em',
      'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'hr', 'i', 'img', 'ins', 'kbd', 'li', 'main', 'mark',
      'ol', 'p', 'pre', 'section', 'span', 'strike', 'strong', 'sub', 'summary', 'sup', 'table',
      'tbody', 'td', 'th', 'thead', 'tr', 'u', 'ul'
    ],
    allowedAttributes: {
      '*': ['title', 'accesskey'],
      a: ['href', 'name', 'target', 'aria-label', 'rel'],
      img: ['src', 'alt', 'title', 'atk-emoticon', 'aria-label'],
      // for code highlight
      code: ['class'],
      span: ['class', 'style'],
    },
    filter: node => {
      // allow hljs style
      const allowed = [
        [ 'code', /^hljs\W+language-(.*)$/ ],
        [ 'span', /^(hljs-.*)$/ ]
      ]
      allowed.forEach(([ tag, reg ]) => {
        if (
          node.tag === tag
          && !!node.attrs.class
          && !(reg as RegExp).test(node.attrs.class)
        ) {
          delete node.attrs.class
        }
      })

      // allow <span> set color sty
      if (node.tag === 'span' && !!node.attrs.style
          && !/^color:(\W+)?#[0-9a-f]{3,6};?$/i.test(node.attrs.style)) {
        delete node.attrs.style
      }

      return true
    }
  })

  // 内容替换器
  ctx.markedReplacers.forEach((replacer) => {
    if (typeof replacer === 'function') dest = replacer(dest)
  })

  return dest
}

/** 获取修正后的 UserAgent */
export async function getCorrectUserAgent() {
  const uaRaw = navigator.userAgent
  if (!(navigator as any).userAgentData || !(navigator as any).userAgentData.getHighEntropyValues) {
    return uaRaw
  }

  // @link https://web.dev/migrate-to-ua-ch/
  // @link https://web.dev/user-agent-client-hints/
  const uaData = (navigator as any).userAgentData
  let uaGot: any = null
  try {
    uaGot = await uaData.getHighEntropyValues(["platformVersion"])
  } catch (err) { console.error(err); return uaRaw }
  const majorPlatformVersion = Number(uaGot.platformVersion.split('.')[0])

  if (uaData.platform === "Windows") {
    if (majorPlatformVersion >= 13) { // @link https://docs.microsoft.com/en-us/microsoft-edge/web-platform/how-to-detect-win11
      // @date 2022-4-29
      // Win 11 样本："Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36"
      return uaRaw.replace(/Windows NT 10.0/, 'Windows NT 11.0')
    }
  }
  if (uaData.platform === "macOS") {
    if (majorPlatformVersion >= 11) { // 11 => BigSur
      // @date 2022-4-29
      // (Intel Chip) macOS 12.3 样本："Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36"
      // (Arm Chip) macOS 样本："Mozilla/5.0 (Macintosh; ARM Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.2 Safari/605.1.15"
      return uaRaw.replace(/(Mac OS X \d+_\d+_\d+|Mac OS X)/, `Mac OS X ${uaGot.platformVersion.replace(/\./g, '_')}`)
    }
  }

  return uaRaw
}

/** 是否为完整的 URL */
export function isValidURL(urlRaw: string) {
  // @link https://www.rfc-editor.org/rfc/rfc3986
  let url: URL
  try {
    url = new URL(urlRaw)
  } catch (_) { return false }
  return url.protocol === "http:" || url.protocol === "https:"
}

/** 获取基于 conf.server 的 URL */
export function getURLBasedOnApi(ctx: Context, path: string) {
  return getURLBasedOn(ctx.conf.server, path)
}

/** 获取基于某个 baseURL 的 URL */
export function getURLBasedOn(baseURL: string, path: string) {
  return `${baseURL.replace(/\/$/, '')}/${path.replace(/^\//, '')}`
}

/**
 * Performs a deep merge of `source` into `target`.
 * Mutates `target` only but not its objects and arrays.
 *
 * @author inspired by [jhildenbiddle](https://stackoverflow.com/a/48218209).
 * @link https://gist.github.com/ahtcx/0cd94e62691f539160b32ecda18af3d6
 */
export function mergeDeep(target: any, source: any) {
  const isObject = (obj: any) => obj && typeof obj === 'object'

  if (!isObject(target) || !isObject(source)) {
    return source
  }

  Object.keys(source).forEach(key => {
    const targetValue = target[key]
    const sourceValue = source[key]

    if (Array.isArray(targetValue) && Array.isArray(sourceValue)) {
      target[key] = targetValue.concat(sourceValue)
    } else if (isObject(targetValue) && isObject(sourceValue)) {
      target[key] = mergeDeep({ ...targetValue}, sourceValue)
    } else {
      target[key] = sourceValue
    }
  })

  return target
}

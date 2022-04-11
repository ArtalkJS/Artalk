import { marked as libMarked } from 'marked'
import insane from 'insane'
import hanabi from 'hanabi'
import Context from '../context'

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

export function timeAgo(date: Date) {
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
          return `${seconds} 秒前`
        }
        return `${minutes} 分钟前`
      }
      return `${hours} 小时前`
    }
    if (days < 0) return '刚刚'

    if (days < 8) {
      return `${days} 天前`
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
  return `${(ctx.conf.gravatar?.mirror || '').replace(/\/$/, '')}/${emailMD5}?d=${encodeURIComponent(ctx.conf.gravatar?.default || '')}&s=80`
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

let markedInstance: typeof libMarked
export function marked(ctx: Context, src: string): string {
  if (!markedInstance) {
    const renderer = new libMarked.Renderer()
    const linkRenderer = renderer.link
    renderer.link = (href, title, text) => {
      const localLink = href?.startsWith(`${window.location.protocol}//${window.location.hostname}`);
      const html = linkRenderer.call(renderer as any, href, title, text);
      return html.replace(/^<a /, `<a target="_blank" ${!localLink ? `rel="noreferrer noopener nofollow"` : ''} `);
    }

    // @see https://github.com/markedjs/marked/blob/4afb228d956a415624c4e5554bb8f25d047676fe/src/Tokenizer.js#L329
    const nMarked = libMarked
    libMarked.setOptions({
      renderer,
      highlight: (code) => hanabi(code),
      pedantic: false,
      gfm: true,
      breaks: true,
      smartLists: true,
      smartypants: true,
      xhtml: false,
      sanitize: true,
      sanitizer: (html) => insane(html, {
        ...insane.defaults,
        allowedAttributes: {
          ...insane.defaults.allowedAttributes,
          img: ['src', 'atk-emoticon']
        },
      }),
      silent: true,
    })

    markedInstance = nMarked
  }

  return markedInstance.parse(src)
}

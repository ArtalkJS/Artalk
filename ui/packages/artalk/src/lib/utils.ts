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
          if (seconds < 10) return ctx.$t('now')
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
  const { mirror, params } = ctx.conf.gravatar
  return `${mirror.replace(/\/$/, '')}/${emailMD5}?${params.replace(/^\?/, '')}`
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

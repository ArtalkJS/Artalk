const win = (window as any) || {}
const nav = navigator || {}

function Detect (userAgent?: string) {
  const u = String(userAgent || nav.userAgent)
  const dest = {
    os: '',
    osVersion: '',
    engine: '',
    browser: '',
    device: '',
    language: '',
    version: ''
  }

  // 内核
  const engineMatch = {
    Trident: u.includes('Trident') || u.includes('NET CLR'),
    Presto: u.includes('Presto'),
    WebKit: u.includes('AppleWebKit'),
    Gecko: u.includes('Gecko/')
  }

  // 浏览器
  const browserMatch = {
    Safari: u.includes('Safari'),
    Chrome: u.includes('Chrome') || u.includes('CriOS'),
    IE: u.includes('MSIE') || u.includes('Trident'),
    Edge: u.includes('Edge') || u.includes('Edg'),
    Firefox: u.includes('Firefox') || u.includes('FxiOS'),
    'Firefox Focus': u.includes('Focus'),
    Chromium: u.includes('Chromium'),
    Opera: u.includes('Opera') || u.includes('OPR'),
    Vivaldi: u.includes('Vivaldi'),
    Yandex: u.includes('YaBrowser'),
    Kindle: u.includes('Kindle') || u.includes('Silk/'),
    360: u.includes('360EE') || u.includes('360SE'),
    UC: u.includes('UC') || u.includes(' UBrowser'),
    QQBrowser: u.includes('QQBrowser'),
    QQ: u.includes('QQ/'),
    Baidu: u.includes('Baidu') || u.includes('BIDUBrowser'),
    Maxthon: u.includes('Maxthon'),
    Sogou: u.includes('MetaSr') || u.includes('Sogou'),
    LBBROWSER: u.includes('LBBROWSER'),
    '2345Explorer': u.includes('2345Explorer'),
    TheWorld: u.includes('TheWorld'),
    MIUI: u.includes('MiuiBrowser'),
    Quark: u.includes('Quark'),
    Qiyu: u.includes('Qiyu'),
    Wechat: u.includes('MicroMessenger'),
    Taobao: u.includes('AliApp(TB'),
    Alipay: u.includes('AliApp(AP'),
    Weibo: u.includes('Weibo'),
    Douban: u.includes('com.douban.frodo'),
    Suning: u.includes('SNEBUY-APP'),
    iQiYi: u.includes('IqiyiApp'),
  }

  // 系统或平台
  const osMatch = {
    Windows: u.includes('Windows'),
    Linux: u.includes('Linux') || u.includes('X11'),
    'macOS': u.includes('Macintosh'),
    Android: u.includes('Android') || u.includes('Adr'),
    Ubuntu: u.includes('Ubuntu'),
    FreeBSD: u.includes('FreeBSD'),
    Debian: u.includes('Debian'),
    'Windows Phone': u.includes('IEMobile') || u.includes('Windows Phone'),
    BlackBerry: u.includes('BlackBerry') || u.includes('RIM'),
    MeeGo: u.includes('MeeGo'),
    Symbian: u.includes('Symbian'),
    iOS: u.includes('like Mac OS X'),
    'Chrome OS': u.includes('CrOS'),
    WebOS: u.includes('hpwOS'),
  }

  // 设备
  const deviceMatch = {
    Mobile: u.includes('Mobi') || u.includes('iPh') || u.includes('480'),
    Tablet: u.includes('Tablet') || u.includes('Pad') || u.includes('Nexus 7')
  }

  // 修正
  if (deviceMatch.Mobile) {
    deviceMatch.Mobile = !(u.includes('iPad'))
  } else if (browserMatch.Chrome && u.includes('Edg')) {
    // Chrome 内核的 Edge
    browserMatch.Chrome = false
    browserMatch.Edge = true
  } else if (win.showModalDialog && win.chrome) {
    browserMatch.Chrome = false
    browserMatch['360'] = true
  }

  // 默认设备
  dest.device = 'PC'

  // 语言
  dest.language = (() => {
    const g = ((nav as any).browserLanguage || nav.language)
    const arr = g.split('-')
    if (arr[1]) arr[1] = arr[1].toUpperCase()
    return arr.join('_')
  })()

  // 应用判断数据
  const hash = {
    engine: engineMatch,
    browser: browserMatch,
    os: osMatch,
    device: deviceMatch,
  }
  Object.entries(hash).forEach(([type, match]) => {
    Object.entries(match).forEach(([name, result]) => {
      if (result === true) dest[type] = name
    })
  })

  // 系统版本信息
  const osVersion = {
    Windows: () => {
      const v = u.replace(/^.*Windows NT ([\d.]+);.*$/, '$1')
      const wvHash = {
        '6.4': '10',
        '6.3': '8.1',
        '6.2': '8',
        '6.1': '7',
        '6.0': 'Vista',
        '5.2': 'XP',
        '5.1': 'XP',
        '5.0': '2000',
        '10.0': '10',
        '11.0': '11' // 自定的，不是微软官方的判断方法
      }
      return wvHash[v] || v
    },
    Android: () => u.replace(/^.*Android ([\d.]+);.*$/, '$1'),
    iOS: () => u.replace(/^.*OS ([\d_]+) like.*$/, '$1').replace(/_/g, '.'),
    Debian: () => u.replace(/^.*Debian\/([\d.]+).*$/, '$1'),
    'Windows Phone': () => u.replace(/^.*Windows Phone( OS)? ([\d.]+);.*$/, '$2'),
    'macOS': () => u.replace(/^.*Mac OS X ([\d_]+).*$/, '$1').replace(/_/g, '.'),
    WebOS: () => u.replace(/^.*hpwOS\/([\d.]+);.*$/, '$1')
  }

  dest.osVersion = ''
  if (osVersion[dest.os]) {
    dest.osVersion = osVersion[dest.os]()
    if (dest.osVersion === u) {
      dest.osVersion = ''
    }
  }

  // 浏览器版本信息
  const version = {
    Safari: () => u.replace(/^.*Version\/([\d.]+).*$/, '$1'),
    Chrome: () => u.replace(/^.*Chrome\/([\d.]+).*$/, '$1').replace(/^.*CriOS\/([\d.]+).*$/, '$1'),
    IE: () => u.replace(/^.*MSIE ([\d.]+).*$/, '$1').replace(/^.*rv:([\d.]+).*$/, '$1'),
    Edge: () => u.replace(/^.*(Edge|Edg|Edg[A-Z]{1})\/([\d.]+).*$/, '$2'),
    Firefox: () => u.replace(/^.*Firefox\/([\d.]+).*$/, '$1').replace(/^.*FxiOS\/([\d.]+).*$/, '$1'),
    'Firefox Focus': () => u.replace(/^.*Focus\/([\d.]+).*$/, '$1'),
    Chromium: () => u.replace(/^.*Chromium\/([\d.]+).*$/, '$1'),
    Opera: () => u.replace(/^.*Opera\/([\d.]+).*$/, '$1').replace(/^.*OPR\/([\d.]+).*$/, '$1'),
    Vivaldi: () => u.replace(/^.*Vivaldi\/([\d.]+).*$/, '$1'),
    Yandex: () => u.replace(/^.*YaBrowser\/([\d.]+).*$/, '$1'),
    Kindle: () => u.replace(/^.*Version\/([\d.]+).*$/, '$1'),
    Maxthon: () => u.replace(/^.*Maxthon\/([\d.]+).*$/, '$1'),
    QQBrowser: () => u.replace(/^.*QQBrowser\/([\d.]+).*$/, '$1'),
    QQ: () => u.replace(/^.*QQ\/([\d.]+).*$/, '$1'),
    Baidu: () => u.replace(/^.*BIDUBrowser[\s/]([\d.]+).*$/, '$1'),
    UC: () => u.replace(/^.*UC?Browser\/([\d.]+).*$/, '$1'),
    Sogou: () => u.replace(/^.*SE ([\d.X]+).*$/, '$1').replace(/^.*SogouMobileBrowser\/([\d.]+).*$/, '$1'),
    '2345Explorer': () => u.replace(/^.*2345Explorer\/([\d.]+).*$/, '$1'),
    TheWorld: () => u.replace(/^.*TheWorld ([\d.]+).*$/, '$1'),
    MIUI: () => u.replace(/^.*MiuiBrowser\/([\d.]+).*$/, '$1'),
    Quark: () => u.replace(/^.*Quark\/([\d.]+).*$/, '$1'),
    Qiyu: () => u.replace(/^.*Qiyu\/([\d.]+).*$/, '$1'),
    Wechat: () => u.replace(/^.*MicroMessenger\/([\d.]+).*$/, '$1'),
    Taobao: () => u.replace(/^.*AliApp\(TB\/([\d.]+).*$/, '$1'),
    Alipay: () => u.replace(/^.*AliApp\(AP\/([\d.]+).*$/, '$1'),
    Weibo: () => u.replace(/^.*weibo__([\d.]+).*$/, '$1'),
    Douban: () => u.replace(/^.*com.douban.frodo\/([\d.]+).*$/, '$1'),
    Suning: () => u.replace(/^.*SNEBUY-APP([\d.]+).*$/, '$1'),
    iQiYi: () => u.replace(/^.*IqiyiVersion\/([\d.]+).*$/, '$1'),
  }

  dest.version = ''
  if (version[dest.browser]) {
    dest.version = version[dest.browser]()
    if (dest.version === u) {
      dest.version = ''
    }
  }

  // 简化版本号
  /* if (_this.osVersion.indexOf('.')) {
    _this.osVersion = _this.osVersion.substring(0, _this.osVersion.indexOf('.'))
  } */
  if (dest.version.indexOf('.')) {
    dest.version = dest.version.substring(0, dest.version.indexOf('.'))
  }

  // 修正
  if (dest.os === 'iOS' && u.includes('iPad')) {
    dest.os = 'iPadOS'
  } else if (dest.browser === 'Edge' && !u.includes('Edg')) {
    dest.engine = 'EdgeHTML'
  } else if (dest.browser === 'MIUI') {
    dest.os = 'Android'
  } else if (dest.browser === 'Chrome' && Number(dest.version) > 27) {
    dest.engine = 'Blink'
  } else if (dest.browser === 'Opera' && Number(dest.version) > 12) {
    dest.engine = 'Blink'
  } else if (dest.browser === 'Yandex') {
    dest.engine = 'Blink'
  } else if (dest.browser === undefined) {
    dest.browser = 'Unknow App'
  }

  return dest
}

export default Detect

export function escapeHtml(value: string): string {
  return value.replace(/[&<>"']/g, (ch) => {
    switch (ch) {
      case '&':
        return '&amp;'
      case '<':
        return '&lt;'
      case '>':
        return '&gt;'
      case '"':
        return '&quot;'
      default:
        return '&#39;'
    }
  })
}

export function relativeTime(date: string, locale: string): string {
  const time = Date.parse(date)
  if (Number.isNaN(time)) return date
  const diff = Date.now() - time
  const minute = 60_000
  const hour = 60 * minute
  const day = 24 * hour
  const zh = locale.toLowerCase().startsWith('zh')
  if (diff < minute) return zh ? '刚刚' : 'just now'
  if (diff < hour) {
    const n = Math.floor(diff / minute)
    return zh ? `${n} 分钟前` : `${n}m ago`
  }
  if (diff < day) {
    const n = Math.floor(diff / hour)
    return zh ? `${n} 小时前` : `${n}h ago`
  }
  if (diff < 30 * day) {
    const n = Math.floor(diff / day)
    return zh ? `${n} 天前` : `${n}d ago`
  }
  return new Date(time).toLocaleDateString(zh ? 'zh-CN' : 'en-US')
}

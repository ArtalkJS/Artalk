import ArtalkSidebar from './artalk-sidebar'

const searchParams = new URLSearchParams(document.location.search)
const pageKey = searchParams.get('pageKey') || ''
const site = searchParams.get('site') || ''
const user = JSON.parse(searchParams.get('user') || '{}')
const darkMode = searchParams.get('darkMode') === '1'
;(window as any).referer = searchParams.get('referer')

const artalkSidebar = new ArtalkSidebar({
  el: '#ArtalkSidebar',
  server: (import.meta.env.MODE === 'development') ? 'http://localhost:23366/api' : '../api',
  pageKey,
  site,
  darkMode,
}, user)

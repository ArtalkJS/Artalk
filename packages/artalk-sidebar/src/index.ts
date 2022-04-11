import ArtalkSidebar from './artalk-sidebar'

const searchParams = new URLSearchParams(document.location.search)
const pageKey = searchParams.get('pageKey') || ''
const site = searchParams.get('site') || ''
const user = JSON.parse(searchParams.get('user') || '{}')

const artalkSidebar = new ArtalkSidebar({
  el: '#ArtalkSidebar',
  server: (['localhost', '127.0.0.1'].includes(window.location.hostname)) ? 'http://localhost:23366/api' : '/api',
  pageKey,
  site
}, user)

export interface SidebarLayer {
  onUserChanged: () => void
  show: (conf?: SidebarShowPayload) => void
  hide: () => void
}

export interface SidebarShowPayload {
  view?: 'comments' | 'sites' | 'pages' | 'transfer'
}

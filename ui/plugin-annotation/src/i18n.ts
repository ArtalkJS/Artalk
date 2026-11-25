export interface Messages {
  panelTitle: string
  create: string
  reply: string
  resolve: string
  reopen: string
  orphan: string
  relocate: string
  empty: string
  deleteConfirm: string
  quote: string
  resolved: string
  filterAll: string
  filterMine: string
  filterResolved: string
  filterPending: string
  searchPlaceholder: string
  loadMore: string
  sortPosition: string
  sortNewest: string
  submit: string
  cancel: string
  required: string
  loadFail: string
  loginRequired: string
}

const en: Messages = {
  panelTitle: 'Annotations',
  create: 'Annotate',
  reply: 'Reply',
  resolve: 'Resolve',
  reopen: 'Reopen',
  orphan: 'Original text changed',
  relocate: 'Relocate',
  empty: 'No annotations yet',
  deleteConfirm: 'Delete this annotation and all replies?',
  quote: 'Quoted text',
  resolved: 'Resolved',
  filterAll: 'All',
  filterMine: 'Mine',
  filterResolved: 'Resolved',
  filterPending: 'Pending',
  searchPlaceholder: 'Search annotations',
  loadMore: 'Load more',
  sortPosition: 'Document order',
  sortNewest: 'Newest',
  submit: 'Publish',
  cancel: 'Cancel',
  required: 'Please write something',
  loadFail: 'Failed to load annotations',
  loginRequired: 'Please fill in your name and email first',
}

const zhCN: Messages = {
  panelTitle: '划线批注',
  create: '批注',
  reply: '回复',
  resolve: '已解决',
  reopen: '重新打开',
  orphan: '原文已变更',
  relocate: '重新定位',
  empty: '还没有划线批注',
  deleteConfirm: '删除这条划线批注及其所有回复？',
  quote: '引用原文',
  resolved: '已解决',
  filterAll: '全部',
  filterMine: '我的',
  filterResolved: '已解决',
  filterPending: '待审核',
  searchPlaceholder: '搜索划线批注',
  loadMore: '加载更多',
  sortPosition: '按原文顺序',
  sortNewest: '最新',
  submit: '发布',
  cancel: '取消',
  required: '请填写内容',
  loadFail: '划线批注加载失败',
  loginRequired: '请先填写昵称和邮箱',
}

export function messages(locale: string): Messages {
  return locale.toLowerCase().startsWith('zh') ? zhCN : en
}

import { defineLocaleExternal } from './external'

export default defineLocaleExternal('zh-TW', {
  /* Editor */
  placeholder: '鍵入內容...',
  noComment: '「此時無聲勝有聲」',
  send: '發送評論',
  save: '保存評論',
  nick: '昵稱',
  email: '郵箱',
  link: '網址',
  emoticon: '表情',
  preview: '預覽',
  image: '圖片',
  uploadFail: '上傳失敗',
  commentFail: '評論失敗',
  restoredMsg: '內容已自動恢復',
  onlyAdminCanReply: '僅管理員可評論',
  uploadLoginMsg: '填入你的名字郵箱才能上傳哦',

  /* List */
  counter: '{count} 條評論',
  sortLatest: '最新',
  sortOldest: '最早',
  sortBest: '最熱',
  sortAuthor: '作者',
  openComment: '打開評論',
  closeComment: '關閉評論',
  listLoadFailMsg: '無法獲取評論列表數據',
  listRetry: '點擊重新獲取',
  loadMore: '加載更多',

  /* Comment */
  admin: '管理員',
  reply: '回覆',
  voteUp: '贊同',
  voteDown: '反對',
  voteFail: '投票失敗',
  readMore: '閱讀更多',
  actionConfirm: '確認操作',
  collapse: '摺疊',
  collapsed: '已摺疊',
  collapsedMsg: '該評論已被系統或管理員摺疊',
  expand: '展開',
  approved: '已審',
  pending: '待審',
  pendingMsg: '審核中，僅本人可見。',
  edit: '編輯',
  editCancel: '取消編輯',
  delete: '刪除',
  deleteConfirm: '確認刪除',
  pin: '置頂',
  unpin: '取消置頂',

  /* Time */
  seconds: '秒前',
  minutes: '分鐘前',
  hours: '小時前',
  days: '天前',
  now: '剛剛',

  /* Checker */
  adminCheck: '鍵入密碼來驗證管理員身份：',
  captchaCheck: '鍵入驗證碼繼續：',
  confirm: '確認',
  cancel: '取消',

  /* Sidebar */
  msgCenter: '通知中心',
  ctrlCenter: '控制中心',

  /* General */
  frontend: '前端',
  backend: '後端',
  loading: '加載中',
  loadFail: '加載失敗',
  editing: '修改中',
  editFail: '修改失敗',
  deleting: '刪除中',
  deleteFail: '刪除失敗',
  reqGot: '請求響應',
  reqAborted: '請求超時或意外終止',
})

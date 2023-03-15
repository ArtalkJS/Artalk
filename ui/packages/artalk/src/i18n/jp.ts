import { defineLocaleExternal } from './external'

export default defineLocaleExternal('jp', {
  /* エディタ */
  placeholder: 'コメントを残す',
  noComment: 'コメントなし',
  send: '送信',
  save: '保存',
  nick: 'ニックネーム',
  email: 'メール',
  link: 'ウェブサイト',
  emoticon: '絵文字',
  preview: 'プレビュー',
  image: '画像',
  uploadFail: 'アップロードに失敗しました',
  commentFail: 'コメントに失敗しました',
  restoredMsg: 'コンテンツが復元されました',
  onlyAdminCanReply: '管理者のみが返信できます',
  uploadLoginMsg: 'アップロードするには、名前とメールアドレスを入力してください',

  /* リスト */
  counter: '{count} コメント',
  sortLatest: '最新のもの',
  sortOldest: '最も古い',
  sortBest: 'ベスト',
  sortAuthor: '著者',
  openComment: 'コメントを開く',
  closeComment: 'コメントを閉じる',
  listLoadFailMsg: 'コメントの読み込みに失敗しました',
  listRetry: 'クリックで再試行',
  loadMore: 'もっと読み込む',

  /* コメント */
  admin: '管理者',
  reply: '返信',
  voteUp: 'アップ',
  voteDown: 'ダウン',
  voteFail: '投票に失敗しました',
  readMore: '続きを読む',
  actionConfirm: '確認',
  collapse: '折りたたむ',
  collapsed: '折りたたまれた',
  collapsedMsg: 'このコメントは折りたたまれました',
  expand: '展開',
  approved: '承認されました',
  pending: '保留中',
  pendingMsg: '保留中、コメント投稿者にのみ表示されます',
  edit: '編集',
  editCancel: '編集をキャンセル',
  delete: '削除',
  deleteConfirm: '確認',
  pin: 'ピン',
  unpin: 'ピン外し',

  /* 時間 */
  seconds: '秒前',
  minutes: '分前',
  hours: '時間前',
  days: '数日前',
  now: 'たった今',

  /* チェッカー */
  adminCheck: '管理者パスワードを入力してください:',
  captchaCheck: '続けるにはCAPTCHAを入力してください:',
  confirm: '確認',
  cancel: '取消',

  /* サイドバー */
  msgCenter: 'メッセージ',
  ctrlCenter: 'コンソール',

  /* 一般 */
  frontend: 'フロントエンド',
  backend: 'バックエンド',
  loading: '読み込み中',
  loadFail: '読み込みに失敗しました',
  editing: '編集中',
  editFail: '編集に失敗しました',
  deleting: '削除中',
  deleteFail: '削除に失敗しました',
  reqGot: 'リクエストを取得しました',
  reqAborted: 'リクエストがタイムアウトしたか、予期せず終了した',
}, ['jp-JP'])

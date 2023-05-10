import type { CommentData } from '~/types/artalk-data'
import Editor from './editor'
import User from '../lib/user'

export interface SubmitManager {
  editor: Editor
  do: () => Promise<void>
  customs: CustomSubmit[]
  registerCustom(c: CustomSubmit): void
}

export function createSubmitManager(editor: Editor) {
  const m: SubmitManager = {
    editor,
    do: () => Do(m),
    customs: [],
    registerCustom: (c) => { m.customs.push(c) }
  }

  return m
}

async function Do(m: SubmitManager) {
  if (m.editor.getFinalContent().trim() === '') {
    m.editor.focus()
    return
  }

  const custom = m.customs.find(o => o.activeCond())

  m.editor.ctx.trigger('editor-submit')
  m.editor.showLoading()

  try {
    // pre submit
    if (custom?.pre) custom.pre()

    let nComment: CommentData

    // submit request
    if (custom?.req) nComment = await custom.req()
    else nComment = await reqAdd(m)

    // post submit
    if (custom?.post) custom.post(nComment)
    else postSubmitAdd(m, nComment)
  } catch (err: any) {
    // submit error
    console.error(err)
    m.editor.showNotify(`${m.editor.$t('commentFail')}，${err.msg || String(err)}`, 'e')
    return
  } finally {
    m.editor.hideLoading()
  }

  m.editor.reset() // 复原编辑器
  m.editor.ctx.trigger('editor-submitted')
}

interface CustomSubmit {
  activeCond: () => void
  pre?: () => void
  req?: () => Promise<CommentData>
  post?: (nComment: CommentData) => void
}

// ================== Submit Add =======================
function getSubmitAddParams(m: SubmitManager) {
  const { nick, email, link } = User.data
  const conf = m.editor.ctx.conf
  const reply = m.editor.getReplyManager()?.comment

  return {
    content: m.editor.getFinalContent(),
    nick, email, link,
    rid: (!reply) ? 0 : reply.id,
    page_key: (!reply) ? conf.pageKey : reply.page_key,
    page_title: (!reply) ? conf.pageTitle : undefined,
    site_name: (!reply) ? conf.site : reply.site_name
  }
}

async function reqAdd(m: SubmitManager) {
  const nComment = await m.editor.ctx.getApi().comment.add({
    ...getSubmitAddParams(m)
  })
  return nComment
}

function postSubmitAdd(m: SubmitManager, commentNew: CommentData) {
  // 回复不同页面的评论，跳转到新页面
  const reply = m.editor.getReplyManager()
  const conf = m.editor.ctx.conf
  if (!!reply?.comment && reply.comment.page_key !== conf.pageKey) {
    window.open(`${reply.comment.page_url}#atk-comment-${commentNew.id}`)
  }

  m.editor.ctx.insertComment(commentNew)
}

import Editor from './editor'
import User from '../lib/user'
import { createMover } from './mover'
import { createPlugManager } from './plug-manager'
import { createReplyManager } from './reply'
import { createEditModeManager } from './edit-mode'
import { createSubmitManager } from './submit'

const EditorFuncs = {
  localStorage, header, textarea, submitBtn, submitManager,
  plugs, mover, replyManager, editModeManager,
}

export default function initEditorFuncs(editor: Editor) {
  Object.entries(EditorFuncs).forEach(([k, init]) => {
    init(editor)
  })
}

// ================== Editor Functions ======================

function localStorage(editor: Editor) {
  const localContent = window.localStorage.getItem('ArtalkContent') || ''
  if (localContent.trim() !== '') {
    editor.showNotify(editor.$t('restoredMsg'), 'i')
    editor.setContent(localContent)
  }

  editor.getUI().$textarea.addEventListener('input', () => (editor.saveToLocalStorage()))
}

function header(editor: Editor) {
  const $inputs = editor.getHeaderInputEls()
  Object.entries($inputs).forEach(([key, $input]) => {
    $input.value = User.data[key] || ''
    $input.addEventListener('input', () => onHeaderInput(editor, key, $input))

    // 设置 Placeholder
    $input.placeholder = `${editor.$t(key as any)}`
  })
}

function onHeaderInput(editor: Editor, field: string, $input: HTMLInputElement) {
  if (editor.isEditMode) return // 评论编辑模式，不修改个人信息

  User.update({
    [field]: $input.value.trim()
  })

  editor.getPlugs()?.triggerHeaderInputEvt(field, $input)
}

function textarea(editor: Editor) {
  // 占位符
  editor.getUI().$textarea.placeholder = editor.ctx.conf.placeholder || editor.$t('placeholder')

  // 修复按下 Tab 输入的内容
  editor.getUI().$textarea.addEventListener('keydown', (e) => {
    const keyCode = e.keyCode || e.which

    if (keyCode === 9) {
      e.preventDefault()
      editor.insertContent('\t')
    }
  })

  // 输入框高度随内容而变化
  editor.getUI().$textarea.addEventListener('input', () => {
    editor.adjustTextareaHeight()
  })
}

function submitBtn(editor: Editor) {
  editor.refreshSendBtnText()
  editor.getUI().$submitBtn.addEventListener('click', () => (editor.submit()))
}

function plugs(editor: Editor) {
  editor.setPlugs(createPlugManager(editor))
}

function mover(editor: Editor) {
  if (!editor.conf.editorTravel) return
  editor.setMover(createMover(editor))
}

function replyManager(editor: Editor) {
  editor.setReplyManager(createReplyManager(editor))
}

function editModeManager(editor: Editor) {
  editor.setEditModeManager(createEditModeManager(editor))
}

function submitManager(editor: Editor) {
  editor.setSubmitManager(createSubmitManager(editor))
}

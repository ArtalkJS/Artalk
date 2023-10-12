import Editor from './editor'
import User from '../lib/user'
import { createMover } from './mover'
import { createPlugManager } from './plug-manager'
import { createReplyManager } from './reply'
import { createEditModeManager } from './edit-mode'
import { createSubmitManager } from './submit'

/** Editor Enabled Functions */
const EditorFuncs = {
  localStorage, header, textarea,
  submitBtn, submitManager,
  mover, reply,
  plugs,
}

export default function initEditorFuncs(editor: Editor) {
  // storage the unmount funcs
  const unmountFuncs: (() => void)[] = []

  Object.entries(EditorFuncs).forEach(([k, init]) => {
    const unmountFunc = init(editor)
    if (unmountFunc) unmountFuncs.push(unmountFunc)
  })

  return unmountFuncs
}

// ================== Editor Functions ======================

function localStorage(editor: Editor) {
  // load editor content from localStorage when init
  const localContent = window.localStorage.getItem('ArtalkContent') || ''
  if (localContent.trim() !== '') {
    editor.showNotify(editor.$t('restoredMsg'), 'i')
    editor.setContent(localContent)
  }

  // save editor content to localStorage when input
  const onEditorInput = () => {
    editor.saveToLocalStorage()
  }

  // bind the event
  editor.getUI().$textarea.addEventListener('input', onEditorInput)

  return () => {
    // unmount the event
    editor.getUI().$textarea.removeEventListener('input', onEditorInput)
  }
}

function header(editor: Editor) {
  const $inputs = editor.getHeaderInputEls()
  const inputFuncs: {[name: string]: () => void} = {}

  // the input event
  const onInput = ($input: HTMLInputElement, key: string) => () => {
    if (editor.isEditMode) return // 评论编辑模式，不修改个人信息

    User.update({ [key]: $input.value.trim() })
    editor.getPlugs()?.triggerHeaderInputEvt(key, $input)
  }

  // set placeholder and sync header input value
  Object.entries($inputs).forEach(([key, $input]) => {
    $input.placeholder = `${editor.$t(key as any)}`
    $input.value = User.data[key] || ''
  })

  // bind the event
  Object.entries($inputs).forEach(([key, $input]) => {
    $input.addEventListener('input', inputFuncs[key] = onInput($input, key))
  })

  return () => {
    // unmount the event
    Object.entries($inputs).forEach(([key, $input]) => {
      $input.removeEventListener('input', inputFuncs[key])
    })
  }
}

function textarea(editor: Editor) {
  // 占位符
  editor.getUI().$textarea.placeholder = editor.ctx.conf.placeholder || editor.$t('placeholder')

  // 按下 Tab 输入内容，而不是失去焦距
  const onKeydown = (e: KeyboardEvent) => {
    const keyCode = e.keyCode || e.which

    if (keyCode === 9) {
      e.preventDefault()
      editor.insertContent('\t')
    }
  }

  // 输入框高度随内容而变化
  const onInput = () => {
    editor.adjustTextareaHeight()
  }

  // bind the event
  editor.getUI().$textarea.addEventListener('keydown', onKeydown)
  editor.getUI().$textarea.addEventListener('input', onInput)

  return () => {
    // unmount the event
    editor.getUI().$textarea.removeEventListener('keydown', onKeydown)
    editor.getUI().$textarea.removeEventListener('input', onInput)
  }
}

function submitBtn(editor: Editor) {
  editor.refreshSendBtnText()

  const onClick = () => {
    editor.submit()
  }

  editor.getUI().$submitBtn.addEventListener('click', onClick)

  return () => {
    editor.getUI().$submitBtn.removeEventListener('click', onClick)
  }
}

function submitManager(editor: Editor) {
  editor.setSubmitManager(createSubmitManager(editor))

  // Note:
  // the EditModeManger depends on the submitManager,
  // so must init the submitManager first.
  // createEditModeManager will make a side-effect on submitManager (call `submitManger.registerCustom`)
  // side-effect is not a good practice, to be improved.
  editor.setEditModeManager(createEditModeManager(editor))
}

function mover(editor: Editor) {
  if (!editor.conf.editorTravel) return
  editor.setMover(createMover(editor))
}

function reply(editor: Editor) {
  editor.setReplyManager(createReplyManager(editor))
}

function plugs(editor: Editor) {
  editor.setPlugs(createPlugManager(editor))
}

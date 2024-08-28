import * as Utils from '../lib/utils'
import EditorHTML from './editor.html?raw'

const Sel = {
  $header: '.atk-header',
  $name: '.atk-header [name="name"]',
  $email: '.atk-header [name="email"]',
  $link: '.atk-header [name="link"]',
  $textareaWrap: '.atk-textarea-wrap',
  $textarea: '.atk-textarea',
  $bottom: '.atk-bottom',
  $submitBtn: '.atk-send-btn',
  $notifyWrap: '.atk-notify-wrap',
  $bottomLeft: '.atk-bottom-left',
  $stateWrap: '.atk-state-wrap',
  $plugBtnWrap: '.atk-plug-btn-wrap',
  $plugPanelWrap: '.atk-plug-panel-wrap',
}

export interface EditorUI extends Record<keyof typeof Sel, HTMLElement> {
  $el: HTMLElement
  $name: HTMLInputElement
  $email: HTMLInputElement
  $link: HTMLInputElement
  $textarea: HTMLTextAreaElement
  $submitBtn: HTMLButtonElement
  $sendReplyBtn?: HTMLElement
  $editCancelBtn?: HTMLElement
}

export function render() {
  const $el = Utils.createElement(EditorHTML)
  const ui = { $el }
  Object.entries(Sel).forEach(([k, sel]) => {
    ui[k] = $el.querySelector(sel)
  })
  return ui as EditorUI
}

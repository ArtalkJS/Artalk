import Context from '~/types/context'
import * as Utils from '../lib/utils'

interface ActionBtnConf {
  /** 按钮文字 (动态/静态) */
  text: (() => string) | string

  /** 仅管理员可用 */
  adminOnly?: boolean

  /** 确认操作 */
  confirm?: boolean

  /** 确认时提示文字 */
  confirmText?: string
}

/**
 * 通用操作按钮
 */
export default class ActionBtn {
  private ctx: Context
  private conf: ActionBtnConf
  public $el: HTMLElement

  public isLoading = false // 正在加载

  public msgRecTimer?: number // 消息显示复原定时器
  public msgRecTimerFunc?: Function // 消息显示复原操作
  public get isMessaging() { return !!this.msgRecTimer } // 消息正在显示

  public isConfirming = false // 正在确认
  public confirmRecTimer?: number // 确认消息复原定时器

  /** 构造函数 */
  constructor(ctx: Context, conf: ActionBtnConf|string|(() => string)) {
    this.ctx = ctx
    this.$el = Utils.createElement(`<span class="atk-common-action-btn"></span>`)

    this.conf = (typeof conf !== 'object') ? ({ text: conf }) : conf
    this.$el.innerText = this.getText()

    // 仅管理员可操作
    if (this.conf.adminOnly) this.$el.setAttribute('atk-only-admin-show', '')
  }

  /** 将按钮装载到指定元素 */
  public appendTo(dom: HTMLElement) {
    dom.append(this.$el)
    return this
  }

  /** 获取按钮文字（动态/静态） */
  private getText() {
    return (typeof this.conf.text === 'string') ? this.conf.text : this.conf.text()
  }

  /** 设置点击事件 */
  public setClick(func: Function) {
    this.$el.onclick = (e) => {
      e.stopPropagation() // 防止穿透

      // 按钮处于加载状态，禁止点击
      if (this.isLoading) {
        return
      }

      // 操作确认
      if (this.conf.confirm && !this.isMessaging) {
        const confirmRestore = () => {
          this.isConfirming = false
          this.$el.classList.remove('atk-btn-confirm')
          this.$el.innerText = this.getText()
        }

        if (!this.isConfirming) {
          this.isConfirming = true
          this.$el.classList.add('atk-btn-confirm')
          this.$el.innerText = this.conf.confirmText || this.ctx.$t('actionConfirm')
          this.confirmRecTimer = window.setTimeout(() => confirmRestore(), 5000)
          return
        }

        if (this.confirmRecTimer) window.clearTimeout(this.confirmRecTimer)
        confirmRestore()
      }

      // 立刻释放掉原有的定时器（当显示消息后，用户立刻点击时）
      if (this.msgRecTimer) {
        this.fireMsgRecTimer()
        this.clearMsgRecTimer()
        return
      }

      // 触发设定的点击事件
      func()
    }
  }

  /** 文字刷新（动态/静态） */
  public updateText(text?: (() => string) | string) {
    if (text) this.conf.text = text
    this.setLoading(false)
    this.$el.innerText = this.getText()
  }

  /** 设置加载状态 */
  public setLoading(value: boolean, loadingText?: string) {
    if (this.isLoading === value) return
    this.isLoading = value
    if (value) {
      this.$el.classList.add('atk-btn-loading')
      this.$el.innerText = loadingText || `${this.ctx.$t('loading')}...`
    } else {
      this.$el.classList.remove('atk-btn-loading')
      this.$el.innerText = this.getText()
    }
  }

  /** 错误消息 */
  public setError(text: string) {
    this.setMsg(text, 'atk-btn-error')
  }

  /** 警告消息 */
  public setWarn(text: string) {
    this.setMsg(text, 'atk-btn-warn')
  }

  /** 成功消息 */
  public setSuccess(text: string) {
    this.setMsg(text, 'atk-btn-success')
  }

  /** 设置消息 */
  public setMsg(text: string, className?: string, duringTime?: number, after?: Function) {
    this.setLoading(false)
    if (className) this.$el.classList.add(className)
    this.$el.innerText = text

    // 设定复原 timer
    this.setMsgRecTimer(() => {
      this.$el.innerText = this.getText()
      if (className) this.$el.classList.remove(className)
      if (after) after()
    }, duringTime || 2500) // 消息默认显示持续 2500s
  }

  /** 设置消息复原操作定时器 */
  private setMsgRecTimer(func: Function, duringTime: number) {
    this.fireMsgRecTimer()
    this.clearMsgRecTimer()

    this.msgRecTimerFunc = func
    this.msgRecTimer = window.setTimeout(() => {
      func()
      this.clearMsgRecTimer()
    }, duringTime)
  }

  /** 立刻触发器复原定时器 */
  private fireMsgRecTimer() {
    if (this.msgRecTimerFunc) this.msgRecTimerFunc()
  }

  /** 仅清除 timer */
  private clearMsgRecTimer() {
    if (this.msgRecTimer) window.clearTimeout(this.msgRecTimer)
    this.msgRecTimer = undefined
    this.msgRecTimerFunc = undefined
  }
}

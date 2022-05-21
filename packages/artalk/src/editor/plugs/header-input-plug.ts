import Editor from '../editor'
import EditorPlug from './editor-plug'

export default class HeaderInputPlug extends EditorPlug {
  public static Name = 'headerInput'

  public constructor(editor: Editor) {
    super(editor)

    this.registerHeaderInputEvt((key, $input) => {
      if (key === 'nick' || key === 'email') {
        this.fetchUserInfo()
      }
    })

    // Link URL 自动补全协议
    this.editor.$link.addEventListener('change', () => {
      const link = this.editor.$link.value.trim()
      if (!!link && !/^(http|https):\/\//.test(link)) {
        this.editor.$link.value = `https://${link}`
        this.ctx.user.update({ link: this.editor.$link.value })
      }
    })
  }

  queryUserInfo = {
    timeout: <number|null>null,
    abortFunc: <(() => void)|null>null
  }

  /** 远程获取用户数据 */
  fetchUserInfo() {
    this.ctx.user.logout()

    // 获取用户信息
    if (this.queryUserInfo.timeout) window.clearTimeout(this.queryUserInfo.timeout) // 清除待发出的请求
    if (this.queryUserInfo.abortFunc) this.queryUserInfo.abortFunc() // 之前发出未完成的请求立刻中止

    this.queryUserInfo.timeout = window.setTimeout(() => {
      this.queryUserInfo.timeout = null // 清理

      const {req, abort} = this.ctx.getApi().userGet(
        this.ctx.user.data.nick, this.ctx.user.data.email
      )
      this.queryUserInfo.abortFunc = abort
      req.then(data => {
        if (!data.is_login) {
          this.ctx.user.logout()
        }

        // 未读消息更新
        this.ctx.updateNotifies(data.unread)

        // 若用户为管理员，执行登陆操作
        if (this.ctx.user.checkHasBasicUserInfo() && !data.is_login && data.user?.is_admin) {
          this.showLoginDialog()
        }

        // 自动填入 link
        if (data.user && data.user.link) {
          this.editor.$link.value = data.user.link
          this.ctx.user.update({ link: data.user.link })
        }
      })
      .catch(() => {})
      .finally(() => {
        this.queryUserInfo.abortFunc = null // 清理
      })
    }, 400) // 延迟执行，减少请求次数
  }

  showLoginDialog() {
    this.ctx.checkAdmin({
      onSuccess: () => {
      }
    })
  }
}

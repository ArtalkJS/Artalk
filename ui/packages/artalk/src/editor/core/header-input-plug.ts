import Editor from '../editor'
import User from '../../lib/user'
import EditorPlug from '../editor-plug'

export default class HeaderInputPlug extends EditorPlug {
  constructor(editor: Editor) {
    super(editor)

    this.kit.useHeaderInput((key, $input) => {
      if (key === 'nick' || key === 'email') {
        this.fetchUserInfo()
      }
    })

    const onLinkInputChange = () => this.onLinkInputChange()

    this.kit.useMounted(() => {
      this.editor.getUI().$link.addEventListener('change', onLinkInputChange)
    })
    this.kit.useUnmounted(() => {
      this.editor.getUI().$link.addEventListener('change', onLinkInputChange)
    })
  }

  private queryUserInfo = {
    timeout: <number|null>null,
    abortFunc: <(() => void)|null>null
  }

  /** 远程获取用户数据 */
  private fetchUserInfo() {
    User.logout()

    // 获取用户信息
    if (this.queryUserInfo.timeout) window.clearTimeout(this.queryUserInfo.timeout) // 清除待发出的请求
    if (this.queryUserInfo.abortFunc) this.queryUserInfo.abortFunc() // 之前发出未完成的请求立刻中止

    this.queryUserInfo.timeout = window.setTimeout(() => {
      this.queryUserInfo.timeout = null // 清理

      const {req, abort} = this.editor.ctx.getApi().user.userGet(
        User.data.nick, User.data.email
      )
      this.queryUserInfo.abortFunc = abort
      req.then(data => {
        if (!data.is_login) {
          User.logout()
        }

        // 未读消息更新
        this.editor.ctx.updateNotifies(data.unread)

        // 若用户为管理员，执行登陆操作
        if (User.checkHasBasicUserInfo() && !data.is_login && data.user?.is_admin) {
          // 显示登录窗口
          this.editor.ctx.checkAdmin({
            onSuccess: () => {}
          })
        }

        // 自动填入 link
        if (data.user && data.user.link) {
          this.editor.getUI().$link.value = data.user.link
          User.update({ link: data.user.link })
        }
      })
      .catch(() => {})
      .finally(() => {
        this.queryUserInfo.abortFunc = null // 清理
      })
    }, 400) // 延迟执行，减少请求次数
  }

  private onLinkInputChange() {
    // Link URL 自动补全协议
    const link = this.editor.getUI().$link.value.trim()
    if (!!link && !/^(http|https):\/\//.test(link)) {
      this.editor.getUI().$link.value = `https://${link}`
      User.update({ link: this.editor.getUI().$link.value })
    }
  }
}

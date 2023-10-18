import User from '@/lib/user'
import $t from '@/i18n'
import EditorPlug from './_plug'
import PlugKit from './_kit'

export default class HeaderUser extends EditorPlug {
  constructor(kit: PlugKit) {
    super(kit)

    const onInput = ({ $input, field }: { $input: HTMLInputElement, field: string }) => {
      if (this.kit.useEditor().getState() === 'edit')
        return // TODO prevent execute when editing, since update comment.user not support

      // update user data
      User.update({ [field]: $input.value.trim() })

      // remote fetch user info
      if (field === 'nick' || field === 'email')
        this.fetchUserInfo() // must after update user data, since fetchUserInfo() will use User.data
    }

    this.kit.useMounted(() => {
      Object.entries(this.kit.useEditor().getHeaderInputEls())
        .forEach(([key, $input]) => {
          // set placeholder
          $input.placeholder = `${$t(key as any)}`

          // sync header values from User.data
          $input.value = User.data[key] || ''
        })

      // bind events
      this.kit.useEvents().on('header-input', onInput)
    })

    this.kit.useUnmounted(() => {
      this.kit.useEvents().off('header-input', onInput)
    })
  }

  private query = {
    timer: <number|null>null,
    abortFn: <(() => void)|null>null
  }

  /**
   * Fetch user info from server
   */
  private fetchUserInfo() {
    User.logout() // clear login status

    if (this.query.timer) window.clearTimeout(this.query.timer) // clear the not executed timeout task
    if (this.query.abortFn) this.query.abortFn() // abort the last request (if request is pending not finished)

    this.query.timer = window.setTimeout(() => {
      this.query.timer = null // clear the timer (clarify the timer is executing)

      const {req, abort} = this.kit.useApi().user.userGet(
        User.data.nick, User.data.email
      )
      this.query.abortFn = abort
      req.then(data => this.onUserInfoFetched(data))
        .catch((err) => {})
        .finally(() => {
          this.query.abortFn = null // clear the abort function (clarify the request is finished)
        })
    }, 400) // delay to reduce request
  }

  /**
   * Function called when user info fetched
   *
   * @param data The response data from server
   */
  private onUserInfoFetched(
    data: any // TODO fix type
  ) {
    // If api response is not login, logout
    if (!data.is_login) User.logout()

    // Update unread notifies
    this.kit.useGlobalCtx().updateNotifies(data.unread)

    // If user is admin and not login,
    if (User.checkHasBasicUserInfo() && !data.is_login && data.user?.is_admin) {
      // then show login window
      this.kit.useGlobalCtx().checkAdmin({
        onSuccess: () => {}
      })
    }

    // Auto fill user link from server
    if (data.user && data.user.link) {
      this.kit.useUI().$link.value = data.user.link
      User.update({ link: data.user.link })
    }
  }
}

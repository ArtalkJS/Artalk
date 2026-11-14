import EditorPlugin from './_plug'
import type PlugKit from './_kit'
import $t from '@/i18n'
import type { UserInfoApiResponseData } from '@/types'

export default class HeaderUser extends EditorPlugin {
  constructor(kit: PlugKit) {
    super(kit)

    const onInput = ({ $input, field }: { $input: HTMLInputElement; field: string }) => {
      if (this.kit.useEditor().getState() === 'edit') return // TODO: prevent execute when editing, since update comment.user not support

      // update user data
      this.kit.useUser().update({ [field]: $input.value.trim() })

      // remote fetch user info
      if (field === 'name' || field === 'email') this.fetchUserInfo() // must after update user data, since fetchUserInfo() will use User.data
    }

    const placeholders = {
      name: $t('nick'),
      email: $t('email'),
      link: $t('link'),
    }

    this.kit.useMounted(() => {
      Object.entries(this.kit.useEditor().getHeaderInputEls()).forEach(([key, $input]) => {
        // set placeholder
        $input.placeholder = placeholders[key]

        // sync header values from User.data
        $input.value = this.kit.useUser().getData()[key] || ''
      })

      // bind events
      this.kit.useEvents().on('header-input', onInput)
    })

    this.kit.useUnmounted(() => {
      this.kit.useEvents().off('header-input', onInput)
      if (this.query.timer) window.clearTimeout(this.query.timer)
      this.query.controller?.abort()
    })
  }

  private query = {
    timer: <number | null>null,
    controller: <AbortController | null>null,
  }

  /**
   * Fetch user info from server
   */
  private fetchUserInfo() {
    this.kit.useUser().logout() // clear login status

    if (this.query.timer) window.clearTimeout(this.query.timer) // clear the not executed timeout task
    this.query.controller?.abort()

    this.query.timer = window.setTimeout(() => {
      this.query.timer = null // clear the timer (clarify the timer is executing)

      const api = this.kit.useApi()
      const controller = (this.query.controller = new AbortController())
      const identity = JSON.stringify(this.kit.useUser().getData())
      api.user
        .getUser({ ...api.getUserFields() }, { signal: controller.signal })
        .then(({ data }) => {
          if (
            !controller.signal.aborted &&
            identity === JSON.stringify(this.kit.useUser().getData())
          )
            this.onUserInfoFetched(data)
        })
        .catch(() => {})
        .finally(() => {
          if (this.query.controller === controller) this.query.controller = null
        })
    }, 400) // delay to reduce request
  }

  /**
   * Function called when user info fetched
   *
   * @param data The response data from server
   */
  private onUserInfoFetched(data: UserInfoApiResponseData) {
    // If api response is not login, logout
    if (!data.is_login) this.kit.useUser().logout()

    // Update unread notifies
    this.kit.useData().updateNotifies(data.notifies)

    // If user is admin and not login,
    if (this.kit.useUser().checkHasBasicUserInfo() && !data.is_login && data.user?.is_admin) {
      // then show login window
      this.kit
        .useCheckers()
        .checkAdmin({})
        .catch(() => {})
    }

    // Auto fill user link from server
    if (data.user && data.user.link) {
      this.kit.useUI().$link.value = data.user.link
      this.kit.useUser().update({ link: data.user.link })
    }
  }
}

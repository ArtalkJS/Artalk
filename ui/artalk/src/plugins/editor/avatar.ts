import EditorPlugin from './_plug'
import type PlugKit from './_kit'
import * as Utils from '@/lib/utils'

/**
 * Editor avatar
 *
 * Renders the guest avatar on the left of the editor and refreshes it from the
 * Gravatar service whenever the email field changes.
 */
export default class Avatar extends EditorPlugin {
  constructor(kit: PlugKit) {
    super(kit)

    const onHeaderChange = ({ field }: { field: string }) => {
      if (field === 'email') this.refresh()
    }

    this.kit.useMounted(() => {
      this.refresh()
      this.kit.useEvents().on('header-change', onHeaderChange)
    })

    this.kit.useUnmounted(() => {
      this.kit.useEvents().off('header-change', onHeaderChange)
    })
  }

  private async refresh() {
    const email = (this.kit.useUI().$email.value || '').trim().toLowerCase()
    const emailHash = email ? await sha256(email) : ''

    const { mirror, params } = this.kit.useConf().gravatar
    this.kit.useUI().$avatarImg.src = Utils.getGravatarURL({ mirror, params, emailHash })
  }
}

async function sha256(str: string): Promise<string> {
  const buf = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(str))
  return Array.from(new Uint8Array(buf))
    .map((b) => b.toString(16).padStart(2, '0'))
    .join('')
}

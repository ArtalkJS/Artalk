import EditorPlugin from './_plug'
import type PlugKit from './_kit'

export default class HeaderLink extends EditorPlugin {
  constructor(kit: PlugKit) {
    super(kit)

    const onLinkChange = ({ field }: { field: string }) => {
      if (field === 'link') this.onLinkInputChange()
    }

    // bind events
    this.kit.useMounted(() => {
      this.kit.useEvents().on('header-change', onLinkChange)
    })

    this.kit.useUnmounted(() => {
      this.kit.useEvents().off('header-change', onLinkChange)
    })
  }

  private onLinkInputChange() {
    // auto and force add protocol prefix for user input link
    const link = this.kit.useUI().$link.value.trim()
    if (!!link && !/^(http|https):\/\//.test(link)) {
      this.kit.useUI().$link.value = `https://${link}`
      this.kit.useUser().update({ link: this.kit.useUI().$link.value })
    }
  }
}

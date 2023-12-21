import type { ArtalkConfig, ContextApi } from '~/types'

export default abstract class Component {
  public $el!: HTMLElement
  public readonly conf: ArtalkConfig

  public constructor(
    public ctx: ContextApi
  ) {
    this.conf = ctx.conf
  }
}

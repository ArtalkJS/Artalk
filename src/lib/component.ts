import ArtalkConfig from '~/types/artalk-config'
import Context from '../context'

export default class Component {
  public $el!: HTMLElement

  public ctx: Context
  public readonly conf: ArtalkConfig

  public constructor(ctx: Context) {
    this.ctx = ctx
    this.conf = ctx.conf
  }
}

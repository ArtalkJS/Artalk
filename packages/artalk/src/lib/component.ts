import ArtalkConfig from '~/types/artalk-config'
import Context from '~/types/context'
import { I18n } from '../i18n'

export default abstract class Component {
  public $el!: HTMLElement

  public ctx: Context
  public readonly conf: ArtalkConfig

  public constructor(ctx: Context) {
    this.ctx = ctx
    this.conf = ctx.conf
  }

  public $t(key: keyof I18n, args: {[key: string]: string} = {}): string {
    return this.ctx.$t(key, args)
  }
}

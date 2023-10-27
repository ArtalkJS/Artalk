import type { ArtalkConfig, ContextApi } from '~/types'
import { I18n } from '../i18n'

export default abstract class Component {
  public $el!: HTMLElement

  public ctx: ContextApi
  public readonly conf: ArtalkConfig

  public constructor(ctx: ContextApi) {
    this.ctx = ctx
    this.conf = ctx.conf
  }

  public $t(key: keyof I18n, args: {[key: string]: string} = {}): string {
    return this.ctx.$t(key, args)
  }
}

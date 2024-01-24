import type { ContextApi } from '@/types'

export default abstract class Component {
  public $el!: HTMLElement
  public get conf() {
    return this.ctx.conf
  }

  public constructor(
    public ctx: ContextApi
  ) {}
}

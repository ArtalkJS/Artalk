import { IPgHolderConf, TPgMode } from '.'
import Adaptors from './adaptors'

/**
 * 分页方式持有者（调度器）
 */
export default class PgHolder {
  private conf: IPgHolderConf

  constructor(conf: IPgHolderConf) {
    this.conf = conf
    this.init()
  }

  public getAdaptor() {
    return Adaptors[this.conf.mode]
  }

  public init() {
    const adaptor = this.getAdaptor()
    const [instance, el] = adaptor.createInstance(this.conf)
    adaptor.instance = instance
    adaptor.el = el
    this.conf.list.$el.append(adaptor.el)
  }

  public setLoading(val: boolean) {
    this.getAdaptor().setLoading(val)
  }

  public update(offset: number, total: number) {
    this.getAdaptor().update(offset, total)
  }

  public getEl() {
    return this.getAdaptor().el
  }

  public showErr(msg: string) {
    const that = this.getAdaptor()
    const func = that.showErr
    if (func) func.bind(that)(msg)
  }

  public setMode(val: TPgMode) {
    if (val !== this.conf.mode) {
      this.getEl().remove()
      this.conf.mode = val
      this.init()
    }
  }

  public next() {
    this.getAdaptor().next()
  }
}

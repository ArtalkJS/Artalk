import { IPgHolderConf, TPgMode } from '../index'
import PaginationAdaptor from './pagination'
import ReadMoreAdaptor from './read-more'

// 分页方式适配器
const Adaptors: Record<TPgMode, IPgAdaptor<any>> = {
  'pagination': PaginationAdaptor,
  'read-more': ReadMoreAdaptor
}

export default Adaptors

export interface IPgAdaptor<T> {
  instance: T
  el: HTMLElement
  createInstance(conf: IPgHolderConf): [T, HTMLElement]
  setLoading(val: boolean): void
  update(offset: number, total: number): void
  next(): void
  showErr?(msg: string): void
}

import { LayerWrap } from './wrap'
import type { LayerManager as ILayerManager } from '@/types'

export class LayerManager implements ILayerManager {
  private wrap = new LayerWrap()

  getEl() {
    return this.wrap.getWrap()
  }

  create(name: string, el?: HTMLElement) {
    return this.wrap.createItem(name, el)
  }

  destroy() {
    this.wrap.destroy()
  }
}

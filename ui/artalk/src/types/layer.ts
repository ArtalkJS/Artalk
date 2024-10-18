export interface Layer {
  show: () => void
  hide: () => void
  destroy: () => void
  setOnAfterHide(func: () => void): void
  setAllowMaskClose(allow: boolean): void
  getAllowMaskClose(): boolean
  getEl: () => HTMLElement
}

export interface LayerManager {
  create(name: string, el?: HTMLElement): Layer
  destroy(): void
  getEl(): HTMLElement
}

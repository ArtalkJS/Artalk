import type { ArtalkConfig, ContextApi } from '@/types'

export function watchConf<T extends (keyof ArtalkConfig)[]>(
  ctx: ContextApi,
  keys: T,
  effect: (conf: Pick<ArtalkConfig, T[number]>) => void,
): void {
  const deepEqual = (a: any, b: any) => JSON.stringify(a) === JSON.stringify(b)
  const val = () => {
    const conf = ctx.getConf()
    const res: any = {}
    keys.forEach((key) => {
      res[key] = conf[key]
    })
    return res
  }
  let lastVal: any = null
  const handler = () => {
    const newVal = val()
    const isDiff = lastVal == null || !deepEqual(lastVal, newVal)
    // only trigger when specific keys changed
    if (isDiff) {
      lastVal = newVal
      effect(newVal)
    }
  }
  ctx.on('mounted', handler)
  ctx.on('updated', handler)
}

import type { Config, EventManager } from '@/types'

export function watchConf<T extends (keyof Config)[]>({
  keys,
  effect,
  getConf,
  getEvents,
}: {
  keys: T
  effect: (conf: Pick<Config, T[number]>) => void
  getConf: () => Config
  getEvents: () => EventManager
}): void {
  const deepEqual = (a: any, b: any) => JSON.stringify(a) === JSON.stringify(b)
  const val = () => {
    const conf = getConf()
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
  getEvents().on('mounted', handler)
  getEvents().on('updated', handler)
}

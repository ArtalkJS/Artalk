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
    const isDiff =
      lastVal == null || keys.some((key) => !isConfigValueEqual(lastVal[key], newVal[key]))
    // only trigger when specific keys changed
    if (isDiff) {
      lastVal = newVal
      effect(newVal)
    }
  }
  getEvents().on('mounted', handler)
  getEvents().on('updated', handler)
  getEvents().on('unmounted', () => {
    getEvents().off('mounted', handler)
    getEvents().off('updated', handler)
  })
}

// Compare plain config objects and arrays by content. Functions and other
// objects (such as Markdown renderers and DOM nodes) are compared by reference.
function isConfigValueEqual(a: unknown, b: unknown): boolean {
  if (Object.is(a, b)) return true
  if (!a || !b || typeof a !== 'object' || typeof b !== 'object') return false

  if (Array.isArray(a) || Array.isArray(b)) {
    if (!Array.isArray(a) || !Array.isArray(b) || a.length !== b.length) return false
  } else if (
    Object.getPrototypeOf(a) !== Object.prototype ||
    Object.getPrototypeOf(b) !== Object.prototype
  ) {
    return false
  }

  const first = a as Record<string, unknown>
  const second = b as Record<string, unknown>
  const keys = Object.keys(first)
  return (
    keys.length === Object.keys(second).length &&
    keys.every(
      (key) =>
        Object.prototype.hasOwnProperty.call(second, key) &&
        isConfigValueEqual(first[key], second[key]),
    )
  )
}

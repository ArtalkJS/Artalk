/**
 * Performs a deep merge of objects and returns new object.
 * Does not modify objects (immutable) and merges arrays via concatenation.
 *
 * @param objects - Objects to merge
 * @returns New object with merged key/values
 */
export function mergeDeep<T>(...objects: any[]): T {
  const isObject = (obj: any) => obj && typeof obj === 'object' && obj.constructor === Object

  return objects.reduce((prev, obj) => {
    Object.keys(obj ?? {}).forEach((key) => {
      // Avoid prototype pollution
      if (key === '__proto__' || key === 'constructor' || key === 'prototype') {
        return
      }

      const pVal = prev[key]
      const oVal = obj[key]

      if (Array.isArray(pVal) && Array.isArray(oVal)) {
        prev[key] = pVal.concat(...oVal)
      } else if (isObject(pVal) && isObject(oVal)) {
        prev[key] = mergeDeep(pVal, oVal)
      } else {
        prev[key] = oVal
      }
    })

    return prev
  }, {})
}

Promise.withResolvers ??= function <T>() {
  let resolve: PromiseWithResolvers<T>['resolve']
  let reject: PromiseWithResolvers<T>['reject']
  const promise = new Promise<T>((res, rej) => {
    resolve = res
    reject = rej
  })
  return { promise, resolve: resolve!, reject: reject! }
}

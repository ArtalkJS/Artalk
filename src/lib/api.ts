import { ListData } from '~/types/artalk-data'
import Context from '../Context'


export default class Api {
  private ctx: Context
  private serverURL: string

  constructor (ctx: Context) {
    this.ctx = ctx
    this.serverURL = ctx.conf.serverUrl
  }

  public async Get(offset: number): Promise<ListData> {
    const params = new FormData()
    params.append('page_key', this.ctx.conf.pageKey)
    params.append('limit', String(this.ctx.conf.readMore?.pageSize || 15))
    params.append('offset', String(offset))

    return timeoutPromise(4000, fetch(`${this.serverURL}/api/get`, {
      method: 'POST',
      body: params,
    })).then(async (response) => {
        const data: ListData = await response.json()
        return data
    })
  }
}

/** TODO: 我靠，一个 timeout，都要丑陋的实现 */
function timeoutPromise<T>(ms: number, promise: Promise<T>): Promise<T> {
  return new Promise((resolve, reject) => {
    const timeoutId = setTimeout(() => {
      reject(new Error("promise timeout"))
    }, ms);
    promise.then(
      (res) => {
        clearTimeout(timeoutId);
        resolve(res);
      },
      (err) => {
        clearTimeout(timeoutId);
        reject(err);
      }
    );
  })
}

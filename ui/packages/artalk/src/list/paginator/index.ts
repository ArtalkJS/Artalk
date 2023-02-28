import ListLite from '../list-lite'
import PgHolder from './holder'

export type TPgMode = 'pagination'|'read-more'

export interface IPgHolderConf {
  list: ListLite
  mode: TPgMode
  total: number
  pageSize: number
}

export default PgHolder

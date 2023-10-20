import List from '../list'
import PgHolder from './holder'

export type TPgMode = 'pagination'|'read-more'

export interface IPgHolderConf {
  list: List
  mode: TPgMode
  total: number
  pageSize: number
}

export default PgHolder

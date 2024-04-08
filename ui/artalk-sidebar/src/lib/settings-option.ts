import YAML from 'yaml'
type Pair = YAML.Pair<YAML.Scalar<any>, YAML.Scalar<any> & YAML.YAMLMap<any, any>>

export interface OptionNode {
  name: string
  path: string
  level: number
  default?: string | number | boolean
  selector?: string[]
  type: 'string' | 'number' | 'boolean' | 'object' | 'array'
  title: string
  subTitle?: string
  items?: OptionNode[]
}

function extractItemComment(item: Pair, index: number, parentPair?: Pair): string {
  let comment = ''
  if (index === 0 && parentPair) comment = parentPair?.value?.commentBefore || ''
  else comment = item?.key?.commentBefore || ''
  return comment
}

export function getTree(yamlObj: YAML.Document.Parsed): OptionNode {
  const tree: OptionNode = {
    name: '',
    path: '',
    title: '',
    level: 0,
    type: 'object',
    items: [],
  }

  const traverse = (
    pairs: Pair[],
    parentNode: OptionNode = tree,
    parentPath: string[] = [],
    parentPair?: Pair,
  ) => {
    pairs.forEach((item, index) => {
      // get key and value
      const key = item.key?.value
      const value = item.value?.toJSON ? item.value.toJSON() : undefined
      if (!key) return

      // get path
      const path = [...parentPath, key]

      // get comment
      const comment = extractItemComment(item, index, parentPair)

      // get type
      const probablyTypes = ['string', 'number', 'boolean', 'object']
      const type =
        (Array.isArray(value) ? 'array' : probablyTypes.find((t) => typeof value === t)) ||
        undefined

      if (!type) return

      // get default value
      const defaultValue = type !== 'object' ? value : undefined

      // create new node
      const node: OptionNode = {
        name: key,
        path: path.join('.'),
        level: parentNode ? parentNode.level + 1 : 0,
        ...extractComment(key, comment),
        default: defaultValue,
        type: type as any,
      }

      // traverse children
      if (type === 'object' && item.value?.items) {
        node.items = []
        traverse(item.value.items, node, path, item)
      }

      // add to parent
      if (!parentNode.items) parentNode.items = []
      parentNode.items.push(node)
    })
  }

  traverse((yamlObj.contents as YAML.YAMLMap<any, any>)?.items)

  return tree
}

/**
 * Get flatten meta data from yaml object
 *
 * @param yamlObj
 * @returns
 */
export function getFlattenNodes(tree: OptionNode): {
  [path: string]: OptionNode
} {
  const metas: { [path: string]: OptionNode } = {}

  const traverse = (node: OptionNode) => {
    metas[node.path] = node
    if (node.items) node.items.forEach(traverse)
  }

  traverse(tree)

  return metas
}

/**
 * Extract option info from comment
 *
 * @param name Option name
 * @param comment Option comment in YAML
 * @returns Option info
 */
function extractComment(name: string, comment: string) {
  comment = comment.trim()

  // ignore comments begin and end with `--`
  comment = comment.replace(/--(.*?)--/gm, '')

  let title = ''
  let subTitle = ''
  let selector: string[] | undefined

  const stReg = /\(.*?\)/gm
  title = comment.replace(stReg, '').trim()
  const stFind = stReg.exec(comment)
  subTitle = stFind ? stFind[0].substring(1, stFind[0].length - 1) : ''
  if (!title) {
    title = snakeToCamel(name)
  }

  const optReg = /\[.*?\]/gm
  const optFind = optReg.exec(title)
  if (optFind) {
    try {
      selector = JSON.parse(optFind[0])
    } catch (err) {
      console.error(err)
    }
    title = title.replace(optReg, '').trim()
  }

  return {
    title,
    subTitle,
    selector,
  }
}

function snakeToCamel(str: string) {
  return str.toLowerCase().replace(/([_][a-z]|^[a-z])/g, (group) => group.slice(-1).toUpperCase())
}

/**
 * Patch the option value by meta data
 *
 * @param value User custom value
 * @param meta Option meta data
 * @returns Patched value
 */
export function patchOptionValue(value: any, node: OptionNode) {
  console.log(value, node)
  switch (node.type) {
    case 'boolean':
      if (value === 'true') value = true
      else if (value === 'false') value = false
      break
    case 'string':
      if (!node.selector)
        // ignore option item
        value = String(value).trim()
      break
    case 'number':
      if (!isNaN(Number(value))) value = Number(value)
      break
    case 'array':
      // trim string array
      if (Array.isArray(value)) value = value.map((v) => (typeof v === 'string' ? v.trim() : v))
      break
  }

  return value
}

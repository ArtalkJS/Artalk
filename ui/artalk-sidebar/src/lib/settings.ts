import YAML from 'yaml'
import { getFlattenNodes, getTree, type OptionNode } from './settings-option'

export class Settings {
  private tree: OptionNode
  private flatten: { [path: string]: OptionNode }
  private customs = shallowRef<YAML.Document.Parsed<YAML.ParsedNode>>()

  constructor(yamlObj: YAML.Document.Parsed) {
    this.tree = getTree(yamlObj)
    this.flatten = getFlattenNodes(this.tree)
  }

  getTree() {
    return this.tree
  }

  getNode(path: string) {
    return this.flatten[path]
  }

  getCustoms() {
    return this.customs
  }

  setCustoms(yamlStr: string) {
    this.customs.value = YAML.parseDocument(yamlStr)
  }

  getCustom(path: string) {
    return this.customs.value?.getIn(path.split('.')) as any
  }

  setCustom(path: string, value: any) {
    const pathArr = path.split('.')

    this.makeSureObject(pathArr)

    this.customs.value?.setIn(pathArr, value)
  }

  // @see https://github.com/eemeli/yaml/issues/174#issuecomment-632281283
  private makeSureObject(pathArr: string[]) {
    for (let i = pathArr.length - 1; i >= 1; i--) {
      const parentPath = pathArr.slice(0, -i)

      const parentNode = this.customs.value?.getIn(parentPath)
      if (!parentNode) {
        this.customs.value?.setIn(parentPath, new YAML.YAMLMap())
      }
    }
  }
}

// -------------------------------------------------------

export * from './settings-option'

// Singleton instance
let instance: Settings

export default {
  init: (yamlObj: YAML.Document.Parsed) => (instance = new Settings(yamlObj)),
  get: () => instance,
}

import YAML from 'yaml'
import { getFlattenNodes, getTree, type OptionNode } from './settings-option'

export class Settings {
  private tree: OptionNode
  private flatten: { [path: string]: OptionNode }
  private customs = shallowRef<YAML.Document.Parsed<YAML.ParsedNode>>()
  private envs = shallowRef<{ [key: string]: string }>()

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

  setEnvs(envs: string[]) {
    const envsObj: { [key: string]: string } = {}
    envs.forEach((env) => {
      const [key, value] = env.split('=')
      envsObj[key] = value
    })
    this.envs.value = envsObj
  }

  getEnv(key: string) {
    return this.envs.value?.[key] || null
  }

  getEnvByPath(path: string) {
    // replace `.` to `_` and uppercase
    // replace `ATK_TRUSTED_DOMAINS_0` to `ATK_TRUSTED_DOMAINS`
    // replace `ATK_ADMIN_USERS_0_NAME` to `ATK_ADMIN_USERS`
    return this.getEnv(
      'ATK_' +
        path
          .replace(/\./g, '_')
          .toUpperCase()
          .replace(/(_\d+?_\w+|_\d+)$/, ''),
    )
  }

  getCustom(path: string) {
    const env = this.getEnvByPath(path)
    if (env) return env
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

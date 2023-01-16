import YAML from 'yaml'
import confTemplate from '../assets/artalk-go.example.yml?raw'

type YAMLPair = {
  key?: { value: string, commentBefore: string, comment: string },
  value?: { commentBefore: string, comment: string, value?: any, items?: YAMLPair[] }
}

export function createSettings() {
  const customs = shallowRef<YAML.Document.Parsed<YAML.ParsedNode>>()

  const yamlDocTpl = YAML.parseDocument(confTemplate)

  const comments: { [path: string]: string } = {}
  const defaultValues: { [path: string]: string } = {}

  loop((yamlDocTpl as any).contents.items)

  function loop(pairs: YAMLPair[], parent?: YAMLPair, path?: string[]) {
    pairs.forEach((item, index: number) => {
      const key = item?.key?.value
      if (!key) return

      const itemPath = (!path) ? [key] : [...path, key]

      let comment = ''
      if (index === 0 && parent) comment = parent?.value?.commentBefore || ''
      else comment = item?.key?.commentBefore || ''

      const pathStr = itemPath.join('.')
      comments[pathStr] = comment.trim()
      defaultValues[pathStr] = item?.value?.value

      // 继续迭代
      if (item?.value?.items) {
        loop(item.value.items, item, itemPath)
      }
    })
  }

  function extractItemDescFromComment(nodePath: string|(string|number)[]) {
    if (Array.isArray(nodePath)) nodePath = nodePath.join('.')
    const nodeName = nodePath.split('.').slice(-1)[0]
    const comment = (comments[nodePath] || '').trim()

    let title = ''
    let subTitle = ''
    let opts: string[]|null = null

    const stReg = /\(.*?\)$/gm
    title = comment.replace(stReg, '').trim()
    const stFind = stReg.exec(comment)
    subTitle = stFind ? stFind[0].substring(1, stFind[0].length-1) : ''
    if (!title) {
      const commonDict: any = { 'enabled': '启用' }
      title = commonDict[nodeName] || snakeToCamel(nodeName)
    }

    const optReg = /\[.*?\]/gm
    const optFind = optReg.exec(title)
    if (optFind) {
      try { opts = JSON.parse(optFind[0]) } catch (err) { console.error(err) }
      title = title.replace(optReg, '').trim()
    }

    return {
      title, subTitle, opts
    }
  }

  function snakeToCamel(str: string) {
    return str.toLowerCase()
      .replace(/([_][a-z]|^[a-z])/g, (group) =>
        group.slice(-1).toUpperCase()
      )
  }

  return { comments, customs, defaultValues, extractItemDescFromComment }
}

let instance: ReturnType<(typeof createSettings)>

export default {
  init: () => instance = createSettings(),
  get: () => instance
}

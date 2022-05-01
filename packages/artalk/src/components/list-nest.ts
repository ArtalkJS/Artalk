import { CommentData } from '~/types/artalk-data'

export type CommentNode = {
  id: number
  comment: CommentData
  children: CommentNode[]
  parent?: CommentNode
  level: number
}

export type SortByType = 'DATE_DESC'|'DATE_ASC'|'SRC_INDEX'|'VOTE_UP_DESC'

// 构建树状结构列表
export function makeNestCommentNodeList(srcData: CommentData[], sortBy: SortByType = 'DATE_DESC', nestMax = 2) {
  const nodeList: CommentNode[] = []

  const roots = srcData.filter((o) => o.rid === 0)
  roots.forEach((root: CommentData) => {
    const rootNode: CommentNode = {
      id: root.id,
      comment: root,
      children: [],
      level: 1,
    }

    rootNode.parent = rootNode
    nodeList.push(rootNode)

    ;(function loadChildren(parentNode: CommentNode) {
      const children = srcData.filter((o) => o.rid === parentNode.id)
      if (children.length === 0) return
      if (parentNode.level >= nestMax) parentNode = parentNode.parent!
      children.forEach((child: CommentData) => {
        const childNode: CommentNode = {
          id: child.id,
          comment: child,
          children: [],
          parent: parentNode,
          level: parentNode.level + 1,
        }

        parentNode.children.push(childNode)
        loadChildren(childNode)
      })
    })(rootNode)
  })

  // 排序
  const sortFunc = (a: CommentNode, b: CommentNode): number => {
    let v = a.id - b.id
    if (sortBy === 'DATE_ASC') v = +new Date(a.comment.date) - +new Date(b.comment.date)
    else if (sortBy === 'DATE_DESC') v = +new Date(b.comment.date) - +new Date(a.comment.date)
    else if (sortBy === 'SRC_INDEX') v = srcData.indexOf(a.comment) - srcData.indexOf(b.comment)
    else if (sortBy === 'VOTE_UP_DESC') v = b.comment.vote_up - a.comment.vote_up
    return v
  }

  (function sortLevels(nodes: CommentNode[]) {
    nodes.forEach((node) => {
      node.children = node.children.sort(sortFunc)
      sortLevels(node.children)
    })
  })(nodeList)

  return nodeList
}

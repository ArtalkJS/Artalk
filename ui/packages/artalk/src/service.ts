import Context from '~/types/context'
import CheckerLauncher from './lib/checker'
import Editor from './editor'
import List from './list'
import Layer from './layer'
import SidebarLayer from './layer/sidebar-layer'
import { initMarked } from './lib/marked'

type TService = (ctx: Context) => void

export default {
  // Markdown 组件
  markdown() {
    initMarked()
  },

  // CheckerLauncher
  checkerLauncher(ctx) {
    const checkerLauncher = new CheckerLauncher(ctx)
    ctx.setCheckerLauncher(checkerLauncher)
  },

  // 编辑器
  editor(ctx) {
    const editor = new Editor(ctx)
    ctx.setEditor(editor)
    ctx.$root.appendChild(editor.$el)
  },

  // 评论列表
  list(ctx) {
    // 评论列表
    const list = new List(ctx)
    ctx.setList(list)
    ctx.$root.appendChild(list.$el)

    // 评论获取
    list.fetchComments(0)
  },

  // 弹出层
  layer(ctx) {
    // 记录页面原始 CSS 属性
    Layer.BodyOrgOverflow = document.body.style.overflow
    Layer.BodyOrgPaddingRight = document.body.style.paddingRight
  },

  // 侧边栏 Layer
  sidebarLayer(ctx) {
    const sidebarLayer = new SidebarLayer(ctx)
    ctx.setSidebarLayer(sidebarLayer)
  },

  // 默认事件绑定
  eventsDefault(ctx) {
    // 锚点快速跳转评论
    window.addEventListener('hashchange', () => {
      ctx.listHashGotoCheck()
    })

    // 本地用户数据变更
    ctx.on('user-changed', () => {
      ctx.checkAdminShowEl()
      ctx.listRefreshUI()
    })
  }
} as { [name: string]: TService }

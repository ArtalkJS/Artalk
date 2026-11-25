import type { Context } from 'artalk'
import { AnnotationApi } from './api'
import {
  describeSelection,
  highlightThreads,
  selectionAllowed,
  sortByDocumentOrder,
} from './anchor'
import { resolveConfig, resolveElement } from './config'
import { messages, type Messages } from './i18n'
import {
  appendReply,
  buildThreads,
  filterThreads,
  initialState,
  removeReply,
  removeThread,
  searchThreads,
  sortForDisplay,
  upsertThread,
  type AnnotationState,
} from './state'
import type { AnnotationCommentData, AnnotationPluginOptions, AnnotationThread } from './types'
import { ConfirmDialog } from './ui/confirm'
import { Composer } from './ui/composer'
import { Drawer } from './ui/drawer'
import { Host } from './ui/host'
import { Toast } from './ui/toast'
import { Toolbar } from './ui/toolbar'

export class AnnotationController {
  private readonly api: AnnotationApi
  private readonly config: ReturnType<typeof resolveConfig>
  private msg: Messages
  private host?: Host
  private toolbar?: Toolbar
  private composer?: Composer
  private drawer?: Drawer
  private toast?: Toast
  private confirm?: ConfirmDialog
  private fab?: HTMLButtonElement
  private fabLabel?: HTMLSpanElement
  private state: AnnotationState = initialState()
  private scope?: HTMLElement
  private cleanupHighlights?: () => void
  private savedRange: Range | null = null
  private selectedQuote = ''
  private relocateId: number | null = null
  private selectionTimer?: number
  private mounted = false
  private destroyed = false

  constructor(
    private ctx: Context,
    private options: AnnotationPluginOptions = {},
  ) {
    this.config = resolveConfig(ctx, options)
    this.msg = messages(String(ctx.getConf().locale || 'en'))
    this.api = new AnnotationApi(ctx)
  }

  mount(): void {
    if (this.mounted || this.destroyed || !this.config.enabled) return
    this.mounted = true
    this.scope = resolveElement(this.config.root, document.body) || document.body
    this.host = new Host(
      resolveElement(this.options.panelContainer, document.body) || document.body,
    )
    this.host.setDarkMode(this.isDark())
    this.toolbar = new Toolbar(this.host.root, this.msg.create)
    this.composer = new Composer(this.host.root, this.msg)
    this.confirm = new ConfirmDialog(this.host.root, this.msg)
    this.toast = new Toast(this.host.root)
    this.drawer = new Drawer(
      this.host.root,
      {
        onClose: () => this.drawer?.close(),
        onFilter: (filter) => this.updateState({ filter }),
        onSearch: (search) => this.updateState({ search }),
        onSort: (sort) => this.updateState({ sort }),
        onFocus: (id) => this.focus(id),
        onReply: (id) => this.startReply(id),
        onResolve: (thread) => void this.toggleResolve(thread),
        onDelete: (thread) => this.confirmDelete(thread),
        onDeleteReply: (thread, reply) => this.confirmDeleteReply(thread, reply),
        onRelocate: (id) => this.startRelocate(id),
        onPrev: () => this.navigate(-1),
        onNext: () => this.navigate(1),
        onLoadMore: () => void this.loadMore(),
      },
      this.msg,
    )

    this.fab = document.createElement('button')
    this.fab.className = 'atk-annotation-fab'
    this.fab.type = 'button'
    this.fab.innerHTML = `<span aria-hidden="true">✎</span><span class="fab-label"></span>`
    this.fabLabel = this.fab.querySelector('.fab-label')!
    this.fabLabel.textContent = this.msg.panelTitle
    this.fab.addEventListener('click', () => this.drawer?.toggle())
    this.host.root.appendChild(this.fab)

    document.addEventListener('selectionchange', this.onSelectionChange)
    document.addEventListener('mouseup', this.onMouseUp)
    this.scope.addEventListener('click', this.onScopeClick)

    void this.load()
  }

  destroy(): void {
    if (this.destroyed) return
    this.destroyed = true
    window.clearTimeout(this.selectionTimer)
    document.removeEventListener('selectionchange', this.onSelectionChange)
    document.removeEventListener('mouseup', this.onMouseUp)
    this.scope?.removeEventListener('click', this.onScopeClick)
    this.cleanupHighlights?.()
    this.cleanupHighlights = undefined
    this.toolbar?.destroy()
    this.composer?.destroy()
    this.drawer?.destroy()
    this.toast?.destroy()
    this.confirm?.destroy()
    this.fab?.remove()
    this.host?.destroy()
  }

  onPageChange(): void {
    if (!this.mounted || this.destroyed) return
    this.cleanupHighlights?.()
    this.cleanupHighlights = undefined
    this.state = initialState()
    this.savedRange = null
    this.selectedQuote = ''
    this.relocateId = null
    this.toolbar?.hide()
    this.composer?.close()
    this.drawer?.close()
    void this.load()
  }

  onThemeChange(): void {
    this.host?.setDarkMode(this.isDark())
  }

  onLocaleChange(): void {
    this.msg = messages(String(this.ctx.getConf().locale || 'en'))
    this.drawer?.setMessages(this.msg)
    if (this.fabLabel) this.fabLabel.textContent = this.msg.panelTitle
    this.onPageChange()
  }

  onUserChange(): void {
    this.renderDrawer()
  }

  private isDark(): boolean {
    const dark = this.ctx.getConf().darkMode
    if (dark === 'auto') return window.matchMedia('(prefers-color-scheme: dark)').matches
    return dark === true
  }

  private async load(): Promise<void> {
    if (!this.mounted || this.destroyed) return
    this.updateState({ loading: true, error: '' })
    try {
      const page = await this.api.listThreads()
      if (this.destroyed) return
      const threads = buildThreads(page.comments)
      this.updateState({ threads, total: page.rootsCount, loading: false })
      await this.renderHighlights()
      this.openDeepLink()
    } catch (err: any) {
      this.updateState({ loading: false, error: err?.message || String(err) })
      this.toast?.show(this.msg.loadFail, 'error')
    }
  }

  private async loadMore(): Promise<void> {
    if (this.state.loadingMore || this.state.threads.length >= this.state.total) return
    this.updateState({ loadingMore: true })
    try {
      const page = await this.api.listThreads(this.state.threads.length)
      if (this.destroyed) return
      let threads = this.state.threads
      for (const thread of buildThreads(page.comments)) {
        threads = upsertThread(threads, thread)
      }
      this.updateState({ threads, total: page.rootsCount, loadingMore: false })
      await this.renderHighlights()
    } catch (err: any) {
      this.updateState({ loadingMore: false })
      this.toast?.show(err?.message || String(err), 'error')
    }
  }

  private async renderHighlights(): Promise<void> {
    if (!this.scope) return
    this.cleanupHighlights?.()
    this.cleanupHighlights = await highlightThreads(this.scope, this.state.threads, (id) => {
      this.state.orphans.add(id)
      this.state.threads = this.state.threads.map((thread) =>
        thread.root.id === id ? { ...thread, orphan: true } : thread,
      )
    })
    this.state.threads = sortByDocumentOrder(this.scope, this.state.threads)
    this.renderDrawer()
  }

  private renderDrawer(): void {
    const user = this.ctx.getUser().getData()
    const filtered = filterThreads(this.state.threads, this.state.filter, user.id)
    const searched = searchThreads(filtered, this.state.search)
    const sorted = sortForDisplay(searched, this.state.sort)
    this.drawer?.render(
      { ...this.state, threads: sorted },
      String(this.ctx.getConf().locale || 'en'),
      user.id,
      user.is_admin,
    )
  }

  private updateState(patch: Partial<AnnotationState>): void {
    this.state = { ...this.state, ...patch }
    this.renderDrawer()
  }

  private onSelectionChange = (): void => {
    if (!this.mounted) return
    window.clearTimeout(this.selectionTimer)
    this.selectionTimer = window.setTimeout(() => this.captureSelection(), 120)
  }

  private onMouseUp = (): void => {
    if (!this.mounted) return
    this.captureSelection()
  }

  private captureSelection(): void {
    const selection = window.getSelection()
    if (!selection || selection.rangeCount === 0 || selection.isCollapsed || !this.scope) {
      this.toolbar?.hide()
      return
    }
    const range = selection.getRangeAt(0)
    if (!selectionAllowed(range, this.scope, this.options.excludeSelector)) {
      this.toolbar?.hide()
      return
    }
    const text = selection.toString().trim()
    if (!text) {
      this.toolbar?.hide()
      return
    }

    if (this.relocateId !== null) {
      const id = this.relocateId
      this.relocateId = null
      void this.applyRelocate(id, range.cloneRange())
      return
    }

    this.savedRange = range.cloneRange()
    this.selectedQuote = text
    const rect = range.getBoundingClientRect()
    this.toolbar?.show(rect, () => void this.startCreate())
  }

  private async startCreate(): Promise<void> {
    if (!this.savedRange) return
    if (!this.config.allowGuest && !this.ctx.getUser().getData().token) {
      this.toast?.show(this.msg.loginRequired, 'error')
      return
    }
    const range = this.savedRange
    this.toolbar?.hide()
    this.composer?.open({
      title: this.msg.create,
      quote: this.selectedQuote,
      onSubmit: async (text) => {
        const anchor = await describeSelection(range, this.scope!)
        const comment = await this.api.createAnnotation(text, anchor)
        const thread: AnnotationThread = { root: comment, replies: [], orphan: false }
        this.updateState({ threads: upsertThread(this.state.threads, thread) })
        await this.renderHighlights()
        this.drawer?.open()
        this.focus(comment.id)
      },
      onCancel: () => undefined,
    })
  }

  private startReply(id: number): void {
    const thread = this.state.threads.find(
      (item) => item.root.id === id || item.replies.some((reply) => reply.id === id),
    )
    if (!thread) return
    this.composer?.open({
      title: this.msg.reply,
      quote: thread.root.anchor?.quote.exact || '',
      onSubmit: async (text) => {
        const reply = await this.api.reply(id, text)
        this.updateState({
          threads: appendReply(this.state.threads, thread.root.id, reply),
        })
      },
      onCancel: () => undefined,
    })
  }

  private async toggleResolve(thread: AnnotationThread): Promise<void> {
    const resolved = !thread.root.anchor_resolved
    try {
      const updated = await this.api.updateAnchor(thread.root.id, { resolved })
      this.updateState({
        threads: this.state.threads.map((item) =>
          item.root.id === thread.root.id ? { ...item, root: { ...item.root, ...updated } } : item,
        ),
      })
    } catch (err: any) {
      this.toast?.show(err?.message || String(err), 'error')
    }
  }

  private confirmDelete(thread: AnnotationThread): void {
    this.confirm?.open({
      title: this.msg.panelTitle,
      body: this.msg.deleteConfirm,
      onConfirm: async () => {
        try {
          await this.api.remove(thread.root.id)
          this.updateState({ threads: removeThread(this.state.threads, thread.root.id) })
          await this.renderHighlights()
        } catch (err: any) {
          this.toast?.show(err?.message || String(err), 'error')
        }
      },
    })
  }

  private confirmDeleteReply(thread: AnnotationThread, reply: AnnotationCommentData): void {
    this.confirm?.open({
      title: this.msg.panelTitle,
      body: this.msg.deleteConfirm,
      onConfirm: async () => {
        try {
          await this.api.remove(reply.id)
          this.updateState({
            threads: removeReply(this.state.threads, thread.root.id, reply.id),
          })
        } catch (err: any) {
          this.toast?.show(err?.message || String(err), 'error')
        }
      },
    })
  }

  private startRelocate(id: number): void {
    this.relocateId = id
    this.toolbar?.hide()
    this.toast?.show(this.msg.relocate, 'info')
  }

  private async applyRelocate(id: number, range: Range): Promise<void> {
    try {
      const anchor = await describeSelection(range, this.scope!)
      const updated = await this.api.updateAnchor(id, { anchor })
      this.updateState({
        threads: this.state.threads.map((thread) =>
          thread.root.id === id
            ? { ...thread, root: { ...thread.root, ...updated }, orphan: false }
            : thread,
        ),
      })
      await this.renderHighlights()
      this.toast?.show(this.msg.relocate, 'success')
    } catch (err: any) {
      this.toast?.show(err?.message || String(err), 'error')
    }
  }

  private onScopeClick = (evt: MouseEvent): void => {
    const target = (evt.target as HTMLElement).closest<HTMLElement>('mark.atk-annotation-mark')
    if (!target) return
    const id = Number(target.dataset.annotationId)
    if (!id) return
    this.drawer?.open()
    this.drawer?.setActive(id)
    this.scrollToThread(id)
  }

  private focus(id: number): void {
    this.drawer?.setActive(id)
    const mark = this.scope?.querySelector<HTMLElement>(
      `mark.atk-annotation-mark[data-annotation-id="${id}"]`,
    )
    if (mark) {
      mark.scrollIntoView({ behavior: 'smooth', block: 'center' })
      mark.classList.add('atk-annotation-flash')
      window.setTimeout(() => mark.classList.remove('atk-annotation-flash'), 1200)
    }
  }

  private scrollToThread(id: number): void {
    const el = this.drawer?.el.querySelector<HTMLElement>(`.thread[data-id="${id}"]`)
    el?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }

  private navigate(direction: -1 | 1): void {
    const list = this.state.threads
    if (list.length === 0) return
    const current = this.state.activeId
    const index = current === null ? -1 : list.findIndex((item) => item.root.id === current)
    const next = (index + direction + list.length) % list.length
    const id = list[next].root.id
    this.updateState({ activeId: id })
    this.focus(id)
  }

  private openDeepLink(): void {
    const match = window.location.hash.match(/^#atk-annotation-(\d+)$/)
    if (!match) return
    const id = Number(match[1])
    if (!this.state.threads.some((thread) => thread.root.id === id)) return
    this.drawer?.open()
    this.focus(id)
  }
}

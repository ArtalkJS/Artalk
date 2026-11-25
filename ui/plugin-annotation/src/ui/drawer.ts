import type { Messages } from '../i18n'
import type { AnnotationState } from '../state'
import type { AnnotationCommentData, AnnotationThread, Filter, SortMode } from '../types'
import { escapeHtml, relativeTime } from '../utils'

export interface DrawerHandlers {
  onClose: () => void
  onFilter: (filter: Filter) => void
  onSearch: (query: string) => void
  onSort: (sort: SortMode) => void
  onFocus: (id: number) => void
  onReply: (id: number) => void
  onResolve: (thread: AnnotationThread) => void
  onDelete: (thread: AnnotationThread) => void
  onDeleteReply: (thread: AnnotationThread, reply: AnnotationCommentData) => void
  onRelocate: (id: number) => void
  onPrev: () => void
  onNext: () => void
  onLoadMore: () => void
}

const FILTERS: Filter[] = ['all', 'mine', 'resolved', 'pending']

export class Drawer {
  readonly el: HTMLElement
  private readonly body: HTMLElement
  private readonly countEl: HTMLElement
  private readonly titleEl: HTMLElement
  private readonly tabsEl: HTMLElement
  private readonly searchInput: HTMLInputElement
  private readonly sortSelect: HTMLSelectElement
  private threads: AnnotationThread[] = []
  private msg: Messages
  private userId?: number
  private isAdmin = false

  constructor(
    root: ShadowRoot,
    private handlers: DrawerHandlers,
    msg: Messages,
  ) {
    this.msg = msg
    this.el = document.createElement('aside')
    this.el.className = 'atk-annotation-drawer'
    this.el.hidden = true
    this.el.innerHTML = `
      <header class="drawer-head">
        <strong class="drawer-title"></strong>
        <span class="drawer-count">0</span>
        <div class="drawer-nav">
          <button type="button" class="icon-btn prev" aria-label="previous">‹</button>
          <button type="button" class="icon-btn next" aria-label="next">›</button>
        </div>
        <button type="button" class="icon-btn close" aria-label="close">×</button>
      </header>
      <div class="drawer-tools">
        <input class="drawer-search" type="search" />
        <select class="drawer-sort"></select>
      </div>
      <div class="drawer-tabs"></div>
      <div class="drawer-body"></div>
    `
    root.appendChild(this.el)
    this.body = this.el.querySelector('.drawer-body')!
    this.countEl = this.el.querySelector('.drawer-count')!
    this.titleEl = this.el.querySelector('.drawer-title')!
    this.tabsEl = this.el.querySelector('.drawer-tabs')!
    this.searchInput = this.el.querySelector('.drawer-search')!
    this.sortSelect = this.el.querySelector('.drawer-sort')!

    this.el.querySelector('.close')!.addEventListener('click', () => this.handlers.onClose())
    this.el.querySelector('.prev')!.addEventListener('click', () => this.handlers.onPrev())
    this.el.querySelector('.next')!.addEventListener('click', () => this.handlers.onNext())
    this.tabsEl.addEventListener('click', (evt) => {
      const btn = (evt.target as HTMLElement).closest<HTMLElement>('[data-filter]')
      if (btn) this.handlers.onFilter(btn.dataset.filter as Filter)
    })
    this.searchInput.addEventListener('input', () => this.handlers.onSearch(this.searchInput.value))
    this.sortSelect.addEventListener('change', () =>
      this.handlers.onSort(this.sortSelect.value as SortMode),
    )
    this.body.addEventListener('click', (evt) => this.onBodyClick(evt))
  }

  render(state: AnnotationState, locale: string, userId?: number, isAdmin = false): void {
    this.threads = state.threads
    this.userId = userId
    this.isAdmin = isAdmin
    this.titleEl.textContent = this.msg.panelTitle
    this.countEl.textContent = String(state.total || state.threads.length)
    this.searchInput.placeholder = this.msg.searchPlaceholder
    this.searchInput.value = state.search
    this.sortSelect.innerHTML = `
      <option value="position">${escapeHtml(this.msg.sortPosition)}</option>
      <option value="newest">${escapeHtml(this.msg.sortNewest)}</option>
    `
    this.sortSelect.value = state.sort
    this.tabsEl.innerHTML = FILTERS.map((filter) => {
      const active = filter === state.filter ? ' active' : ''
      return `<button type="button" class="tab${active}" data-filter="${filter}">${escapeHtml(this.filterLabel(filter))}</button>`
    }).join('')

    if (state.threads.length === 0) {
      this.body.innerHTML = `<p class="empty">${escapeHtml(this.msg.empty)}</p>`
      return
    }

    this.body.innerHTML = state.threads.map((thread) => this.renderThread(thread, locale)).join('')
    if (state.threads.length < state.total) {
      const loadMore = document.createElement('button')
      loadMore.type = 'button'
      loadMore.className = 'load-more'
      loadMore.dataset.action = 'load-more'
      loadMore.textContent = state.loadingMore ? '...' : this.msg.loadMore
      loadMore.disabled = state.loadingMore
      this.body.appendChild(loadMore)
    }
  }

  open(): void {
    this.el.hidden = false
  }

  close(): void {
    this.el.hidden = true
  }

  toggle(): void {
    this.el.hidden = !this.el.hidden
  }

  isOpen(): boolean {
    return !this.el.hidden
  }

  setActive(id: number | null): void {
    this.body.querySelectorAll<HTMLElement>('.thread').forEach((el) => {
      el.classList.toggle('active', Number(el.dataset.id) === id)
    })
  }

  setMessages(msg: Messages): void {
    this.msg = msg
  }

  destroy(): void {
    this.el.remove()
  }

  private filterLabel(filter: Filter): string {
    switch (filter) {
      case 'mine':
        return this.msg.filterMine
      case 'resolved':
        return this.msg.filterResolved
      case 'pending':
        return this.msg.filterPending
      default:
        return this.msg.filterAll
    }
  }

  private canManage(userId?: number): boolean {
    if (this.isAdmin) return true
    // The viewer id arrives asynchronously with the user profile. Keep the
    // actions visible until it is known; the API still enforces ownership.
    if (this.userId === undefined) return true
    return userId !== undefined && userId === this.userId
  }

  private renderThread(thread: AnnotationThread, locale: string): string {
    const root = thread.root
    const quote = root.anchor?.quote.exact || ''
    const resolved = root.anchor_resolved === true
    const canManage = this.canManage(root.user_id)
    const replies = thread.replies
      .map((reply) => {
        const canDelete = this.canManage(reply.user_id)
        return `
        <div class="reply">
          <div class="reply-head"><strong>${escapeHtml(reply.nick)}</strong><span>${escapeHtml(relativeTime(reply.date, locale))}</span></div>
          <p>${escapeHtml(reply.content)}</p>
          <div class="reply-actions">
            <button type="button" data-action="reply" data-id="${reply.id}">${escapeHtml(this.msg.reply)}</button>
            ${canDelete ? `<button type="button" data-action="delete-reply" data-id="${reply.id}" data-root-id="${root.id}">×</button>` : ''}
          </div>
        </div>
      `
      })
      .join('')

    return `
      <article class="thread${resolved ? ' resolved' : ''}" data-id="${root.id}">
        <div class="thread-quote">${escapeHtml(quote)}</div>
        <div class="thread-head">
          <strong>${escapeHtml(root.nick)}</strong>
          <span>${escapeHtml(relativeTime(root.date, locale))}</span>
          ${resolved ? `<span class="badge">${escapeHtml(this.msg.resolved)}</span>` : ''}
          ${thread.orphan ? `<span class="badge warn">${escapeHtml(this.msg.orphan)}</span>` : ''}
        </div>
        <p class="thread-content">${escapeHtml(root.content)}</p>
        ${replies ? `<div class="replies">${replies}</div>` : ''}
        <div class="thread-actions">
          <button type="button" data-action="focus" data-id="${root.id}">${escapeHtml(this.msg.quote)}</button>
          <button type="button" data-action="reply" data-id="${root.id}">${escapeHtml(this.msg.reply)}</button>
          ${canManage ? `<button type="button" data-action="resolve" data-id="${root.id}">${escapeHtml(resolved ? this.msg.reopen : this.msg.resolve)}</button>` : ''}
          ${canManage && thread.orphan ? `<button type="button" data-action="relocate" data-id="${root.id}">${escapeHtml(this.msg.relocate)}</button>` : ''}
          ${canManage ? `<button type="button" data-action="delete" data-id="${root.id}">×</button>` : ''}
        </div>
      </article>
    `
  }

  private onBodyClick(evt: MouseEvent): void {
    const target = (evt.target as HTMLElement).closest<HTMLElement>('[data-action]')
    if (!target) return
    const action = target.dataset.action
    const id = Number(target.dataset.id)
    const rootId = Number(target.dataset.rootId)
    const thread = this.threads.find((item) => item.root.id === id || item.root.id === rootId)
    switch (action) {
      case 'focus':
        this.handlers.onFocus(id)
        break
      case 'reply':
        this.handlers.onReply(id)
        break
      case 'resolve':
        if (thread) this.handlers.onResolve(thread)
        break
      case 'delete':
        if (thread) this.handlers.onDelete(thread)
        break
      case 'delete-reply': {
        const reply = thread?.replies.find((item) => item.id === id)
        if (thread && reply) this.handlers.onDeleteReply(thread, reply)
        break
      }
      case 'relocate':
        this.handlers.onRelocate(id)
        break
      case 'load-more':
        this.handlers.onLoadMore()
        break
    }
  }
}

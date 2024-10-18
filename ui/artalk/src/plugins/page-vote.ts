import type { ArtalkPlugin, ListFetchedArgs } from '@/types'
import type { Api } from '@/api'
import $t from '@/i18n'

interface PageVoteOptions {
  /** Up Vote Button Selector */
  upBtnEl: string

  /** Down Vote Button Selector */
  downBtnEl: string

  /** Up Vote Count Selector */
  upCountEl: string

  /** Down Vote Count Selector */
  downCountEl: string

  /** Active class name if the vote is already cast */
  activeClass?: string
}

const defaults: PageVoteOptions = {
  upBtnEl: '.artalk-page-vote-up',
  downBtnEl: '.artalk-page-vote-down',
  upCountEl: '.artalk-page-vote-up-count',
  downCountEl: '.artalk-page-vote-down-count',
  activeClass: 'active',
}

type VoteBtnHandler = (evt: MouseEvent) => void

interface PageVoteState {
  upBtnEl?: HTMLElement | null
  downBtnEl?: HTMLElement | null
  upCountEl?: HTMLElement | null
  downCountEl?: HTMLElement | null
  voteUpHandler?: VoteBtnHandler
  voteDownHandler?: VoteBtnHandler
  activeClass?: string
}

export const PageVoteWidget: ArtalkPlugin = (ctx) => {
  const conf = ctx.inject('config')
  let state: PageVoteState = initState()
  let cleanup: (() => void) | undefined

  ctx.watchConf(['pageVote'], (conf) => {
    const options = { ...defaults, ...(typeof conf.pageVote === 'object' ? conf.pageVote : {}) }

    cleanup?.()
    if (conf.pageVote) {
      state = loadOptions(state, options, voteUpHandler, voteDownHandler)
      cleanup = () => (state = resetState(state))
    }
  })

  ctx.on('unmounted', () => {
    cleanup?.()
    ctx.off('list-fetched', listFetchedHandler)
  })

  ctx.on('list-fetched', listFetchedHandler)

  // List fetched handler
  let currPageId = 0
  function listFetchedHandler({ data }: ListFetchedArgs) {
    if (!conf.get().pageVote || !checkEls(state) || !data) return

    if (currPageId !== data.page.id) {
      // Initialize vote status in a new page
      currPageId = data.page.id
      updateVoteStatus(state, data.page.vote_up, data.page.vote_down, false, false) // reset initial status
      fetchVoteStatus(state, currPageId, ctx.getApi())
    } else {
      // Update vote status in the same page
      updateVoteStatus(state, data.page.vote_up, data.page.vote_down)
    }
  }

  // Vote button click handlers
  const handlerOptions = () => ({
    state,
    pageId: ctx.getData().getPage()?.id || 0,
    httpApi: ctx.getApi(),
  })
  const voteUpHandler = voteBtnHandler('up', handlerOptions)
  const voteDownHandler = voteBtnHandler('down', handlerOptions)
}

function loadOptions(
  _state: PageVoteState,
  opts: PageVoteOptions,
  voteUpHandler: VoteBtnHandler,
  voteDownHandler: VoteBtnHandler,
): PageVoteState {
  const state: PageVoteState = { ..._state } // clone to keep immutable
  state.upBtnEl = document.querySelector<HTMLElement>(opts.upBtnEl)
  state.downBtnEl = document.querySelector<HTMLElement>(opts.downBtnEl)
  state.upCountEl = document.querySelector<HTMLElement>(opts.upCountEl)
  state.downCountEl = document.querySelector<HTMLElement>(opts.downCountEl)
  state.activeClass = opts.activeClass
  state.voteUpHandler = voteUpHandler
  state.voteDownHandler = voteDownHandler
  state.upBtnEl?.addEventListener('click', voteUpHandler)
  state.downBtnEl?.addEventListener('click', voteDownHandler)
  return state
}

function checkEls(state: PageVoteState) {
  return state.upBtnEl || state.downBtnEl || state.upCountEl || state.downCountEl
}

function initState(): PageVoteState {
  return {}
}

function resetState(state: PageVoteState): PageVoteState {
  state.voteUpHandler && state.upBtnEl?.removeEventListener('click', state.voteUpHandler)
  state.voteDownHandler && state.downBtnEl?.removeEventListener('click', state.voteDownHandler)
  return initState()
}

function updateVoteStatus(
  state: PageVoteState,
  up: number | string,
  down: number | string,
  isUp?: boolean,
  isDown?: boolean,
) {
  // Up vote count
  if (state.upCountEl) {
    state.upCountEl.innerText = String(up)
  } else if (state.upBtnEl) {
    state.upBtnEl.innerText = `${$t('voteUp')} (${up})`
  }

  // Down vote count
  if (state.downCountEl) {
    state.downCountEl.innerText = String(down)
  } else if (state.downBtnEl) {
    state.downBtnEl.innerText = `${$t('voteDown')} (${down})`
  }

  if (typeof isUp === 'boolean')
    state.upBtnEl?.classList.toggle(state.activeClass || 'active', isUp)
  if (typeof isDown === 'boolean')
    state.downBtnEl?.classList.toggle(state.activeClass || 'active', isDown)
}

interface VoteBtnHandlerOptions {
  state: PageVoteState
  pageId: number
  httpApi: Api
}

const voteBtnHandler =
  (choice: 'up' | 'down', options: () => VoteBtnHandlerOptions) => (evt: MouseEvent) => {
    evt.preventDefault()
    const { state, pageId, httpApi } = options()
    if (!pageId) return
    httpApi.votes
      .createVote('page', pageId, choice, {
        ...httpApi.getUserFields(),
      })
      .then(({ data }) => {
        updateVoteStatus(state, data.up, data.down, data.is_up, data.is_down)
      })
      .catch((err) => {
        window.alert($t('voteFail'))
        console.error('[ArtalkPageVote]', err)
      })
  }

function fetchVoteStatus(state: PageVoteState, pageId: number, httpApi: Api) {
  httpApi.votes
    .getVote('page', pageId)
    .then(({ data }) => {
      updateVoteStatus(state, data.up, data.down, data.is_up, data.is_down)
    })
    .catch((err) => {
      console.error('[ArtalkPageVote]', err)
    })
}

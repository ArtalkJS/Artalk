import ActionBtn from '../components/action-btn'
import type { ArtalkPlugin, ContextApi } from '@/types'
import { Api } from '@/api'
import $t from '@/i18n'

export interface PageVoteOptions {
  getApi(): Api

  pageKey: string
  vote: boolean
  voteDown: boolean
  btnEl: string
  el: string
}

export const PageVoteWidget: ArtalkPlugin = (ctx) => {
  ctx.watchConf(['pageKey', 'pageVote'], (conf) => {
    const defaultOptions = {
      getApi: () => ctx.getApi(),
      pageKey: conf.pageKey,
      vote: true,
      voteDown: false,
      btnEl: '.artalk-page-vote',
      el: '.artalk-page-vote-count',
    }

    if (!conf.pageVote) return
    const pageVoteConfig = typeof conf.pageVote === 'object' ? conf.pageVote : {}
    initPageVoteWidget(ctx, { ...defaultOptions, ...pageVoteConfig })
  })
}

function initPageVoteWidget(ctx: ContextApi, options: PageVoteOptions) {
  if (!options.vote) return

  const btnContainer = document.querySelector(options.btnEl) as HTMLElement
  if (!btnContainer) throw Error(`page vote's config \`btnEl\` selector ${options.btnEl} not found`)

  const api = options.getApi()

  ctx.on('list-fetched', ({ data }) => {
    const voteUpBtn = new ActionBtn(() => `${$t('voteUp')}${data!.page.vote_up || 0}`).appendTo(btnContainer)
    voteUpBtn.setClick(() => {
      api.votes.vote('page_up', data!.page.id, {
        ...api.getUserFields()
      }).then(res => {
        // todo: update vote count
        console.log(res)
      })
    })

    if (options.voteDown) {
      const voteDownBtn = new ActionBtn(() => `${$t('voteDown')}`).appendTo(btnContainer)
      voteDownBtn.setClick(() => {
        api.votes.vote('page_down', data!.page.id, {
          ...api.getUserFields()
        })
      })
    }
  })
}

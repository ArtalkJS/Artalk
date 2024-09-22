import ActionBtn from '../components/action-btn'
import type { ArtalkPlugin, ContextApi } from '@/types'
import $t from '@/i18n'

interface PageVoteOptions {
  pageKey: string
  vote: boolean
  voteDown: boolean
  btnEl: string
  el: string
}

export const PageVoteWidget: ArtalkPlugin = (ctx) => {
  ctx.watchConf(['pageKey', 'pageVote'], (conf) => {
    if (!conf.pageVote) return

    const defaultOptions = {
      pageKey: conf.pageKey,
      vote: true,
      voteDown: false,
      btnEl: '.artalk-page-vote',
      el: '.artalk-page-vote-count',
    }

    const pageVoteConfig = typeof conf.pageVote === 'object' ? conf.pageVote : {}
    initPageVoteWidget(ctx, { ...defaultOptions, ...pageVoteConfig })
  })
}

function initPageVoteWidget(ctx: ContextApi, options: PageVoteOptions) {
  const btnContainer = document.querySelector(options.btnEl) as HTMLElement
  if (!btnContainer) return // throw Error(`page vote's config \`btnEl\` selector ${options.btnEl} not found`)
  btnContainer.classList.add('atk-layer-wrap')

  const api = ctx.getApi()

  ctx.on('list-fetched', ({ data }) => {
    const voteUpBtn = new ActionBtn(() => `${$t('voteUp')} (${data!.page.vote_up || 0})`).appendTo(btnContainer)
    voteUpBtn.setClick(() => {
      api.votes.vote('page_up', data!.page.id, {
        ...api.getUserFields()
      }).then(({ data: { up, down } }) => {
        data!.page.vote_up = up
        data!.page.vote_down = down
        voteUpBtn.updateText()
      }).catch((err) => {
        voteUpBtn.setError($t('voteFail'))
        console.error(err)
      })
    })

    if (options.voteDown) {
      const voteDownBtn = new ActionBtn(() => `${$t('voteDown')} (${data!.page.vote_down || 0})`).appendTo(btnContainer)
      voteDownBtn.setClick(() => {
        api.votes.vote('page_down', data!.page.id, {
          ...api.getUserFields()
        }).then(({ data: { up, down } }) => {
          data!.page.vote_up = up
          data!.page.vote_down = down
          voteUpBtn.updateText()
        }).catch((err) => {
          voteDownBtn.setError($t('voteFail'))
          console.error(err)
        })
      })
    }
  })
}

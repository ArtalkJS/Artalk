import type { Context as IContext, ArtalkPlugin, PageData } from '@/types'

export const WithEditor: ArtalkPlugin = (ctx) => {
  const list = ctx.inject('list')
  const editorPlugs = ctx.inject('editorPlugs')
  let $closeCommentBtn: HTMLElement | undefined

  // on Artalk mounted
  // (after all components had mounted)
  ctx.on('mounted', () => {
    $closeCommentBtn = list
      .getEl()
      .querySelector<HTMLElement>('[data-action="admin-close-comment"]')!

    // bind editor close button click event
    $closeCommentBtn.addEventListener('click', () => {
      const page = ctx.getData().getPage()
      if (!page) return

      adminPageEditSave(ctx, page)
    })
  })

  // on comment list loaded (it will include page data update)
  const syncPage = () => {
    const page = ctx.getData().getPage()
    // if page comment is closed
    if (page?.admin_only === true) {
      // then close editor
      editorPlugs?.getEvents().trigger('editor-close')
      $closeCommentBtn && ($closeCommentBtn.innerText = ctx.$t('openComment'))
    } else {
      // the open editor
      editorPlugs?.getEvents().trigger('editor-open')
      $closeCommentBtn && ($closeCommentBtn.innerText = ctx.$t('closeComment'))
    }
  }
  ctx.on('page-loaded', syncPage)
  ctx.on('updated', syncPage)
  ctx.on('user-changed', syncPage)

  ctx.on('list-loaded', (comments) => {
    // 防止评论框被吞
    ctx.editorResetState()
  })
}

/** 管理员设置页面信息 */
function adminPageEditSave(ctx: IContext, page: PageData) {
  const destination = () => {
    const { server, site, pageKey } = ctx.getConf()
    return JSON.stringify([server, site, pageKey])
  }
  const current = destination()
  const active = () => !ctx.isDestroyed() && destination() === current
  ctx.editorShowLoading()
  ctx
    .getApi()
    .pages.updatePage(page.id, { ...page, admin_only: !page.admin_only })
    .then(({ data }) => {
      if (!active()) return
      ctx.getData().updatePage(data)
    })
    .catch((err) => {
      if (!active()) return
      ctx.editorShowNotify(`${ctx.$t('editFail')}: ${err.message || String(err)}`, 'e')
    })
    .finally(() => {
      if (active()) ctx.editorHideLoading()
    })
}

import type { ContextApi, ArtalkPlugin, PageData } from '@/types'
import $t from '@/i18n'

export const WithEditor: ArtalkPlugin = (ctx) => {
  let $closeCommentBtn: HTMLElement | undefined

  // on Artalk mounted
  // (after all components had mounted)
  ctx.on('mounted', () => {
    const list = ctx.get('list')

    $closeCommentBtn = list.$el.querySelector<HTMLElement>('[data-action="admin-close-comment"]')!

    // bind editor close button click event
    $closeCommentBtn.addEventListener('click', () => {
      const page = ctx.getData().getPage()
      if (!page) throw new Error('Page data not found')

      page.admin_only = !page.admin_only
      adminPageEditSave(ctx, page)
    })
  })

  // on comment list loaded (it will include page data update)
  ctx.on('page-loaded', (page) => {
    const editor = ctx.get('editor')

    // if page comment is closed
    if (page?.admin_only === true) {
      // then close editor
      editor.getPlugs()?.getEvents().trigger('editor-close')
      $closeCommentBtn && ($closeCommentBtn.innerText = $t('openComment'))
    } else {
      // the open editor
      editor.getPlugs()?.getEvents().trigger('editor-open')
      $closeCommentBtn && ($closeCommentBtn.innerText = $t('closeComment'))
    }
  })

  ctx.on('list-loaded', (comments) => {
    // 防止评论框被吞
    ctx.editorResetState()
  })
}

/** 管理员设置页面信息 */
function adminPageEditSave(ctx: ContextApi, page: PageData) {
  ctx.editorShowLoading()
  ctx
    .getApi()
    .pages.updatePage(page.id, page)
    .then(({ data }) => {
      ctx.getData().updatePage(data)
    })
    .catch((err) => {
      ctx.editorShowNotify(`${$t('editFail')}: ${err.message || String(err)}`, 'e')
    })
    .finally(() => {
      ctx.editorHideLoading()
    })
}

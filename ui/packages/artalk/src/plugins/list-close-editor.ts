import type ContextApi from '~/types/context'
import type ArtalkPlugin from '~/types/plugin'
import { PageData } from '~/types/artalk-data'
import $t from '@/i18n'

function ensureListEditor(ctx: ContextApi) {
  const list = ctx.get('list')
  const editor = ctx.get('editor')

  if (!list) throw new Error('List instance not found')
  if (!editor) throw new Error('Editor instance not found')

  return { list, editor }
}

export const ListCloseEditor: ArtalkPlugin = (ctx) => {
  let $closeCommentBtn: HTMLElement|undefined

  // on Artalk inited
  // (after all components had mounted)
  ctx.on('inited', () => {
    const { list } = ensureListEditor(ctx)

    $closeCommentBtn = list.$el.querySelector<HTMLElement>('[data-action="admin-close-comment"]')!

    // bind editor close button click event
    $closeCommentBtn.addEventListener('click', () => {
      const page = ctx.getPage()
      if (!page) throw new Error('Page data not found')

      page.admin_only = !page.admin_only
      adminPageEditSave(ctx, page)
    })
  })

  // on comment list loaded (it will include page data update)
  ctx.on('page-loaded', (page) => {
    const { editor } = ensureListEditor(ctx)

    // if page comment is closed
    if (page?.admin_only === true) {
      // then close editor
      editor.getPlugs()?.getEvents().trigger('editor-close')
      $closeCommentBtn && ($closeCommentBtn.innerHTML = $t('openComment'))
    } else {
      // the open editor
      editor.getPlugs()?.getEvents().trigger('editor-open')
      $closeCommentBtn && ($closeCommentBtn.innerHTML = $t('closeComment'))
    }
  })
}

/** 管理员设置页面信息 */
function adminPageEditSave(ctx: ContextApi, page: PageData) {
  ctx.editorShowLoading()
  ctx.getApi().page.pageEdit(page)
    .then((respPage) => {
      ctx.updatePage(respPage)
    })
    .catch(err => {
      ctx.editorShowNotify(`${$t('editFail')}: ${err.msg || String(err)}`, 'e')
    })
    .finally(() => {
      ctx.editorHideLoading()
    })
}

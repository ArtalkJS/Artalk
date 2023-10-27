import type { ContextApi, ArtalkPlugin, PageData } from '~/types'
import $t from '@/i18n'

function ensureListEditor(ctx: ContextApi) {
  const list = ctx.get('list')
  const editor = ctx.get('editor')

  if (!list) throw new Error('List instance not found')
  if (!editor) throw new Error('Editor instance not found')

  return { list, editor }
}

export const WithEditor: ArtalkPlugin = (ctx) => {
  let $closeCommentBtn: HTMLElement|undefined

  // on Artalk inited
  // (after all components had mounted)
  ctx.on('inited', () => {
    const { list } = ensureListEditor(ctx)

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
    const { editor } = ensureListEditor(ctx)

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
  ctx.getApi().page.pageEdit(page)
    .then((respPage) => {
      ctx.getData().updatePage(respPage)
    })
    .catch(err => {
      ctx.editorShowNotify(`${$t('editFail')}: ${err.msg || String(err)}`, 'e')
    })
    .finally(() => {
      ctx.editorHideLoading()
    })
}

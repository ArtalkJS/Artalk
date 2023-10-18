import type ContextApi from '~/types/context'
import type ArtalkPlug from '~/types/plug'
import $t from '@/i18n'

function ensureListEditor(ctx: ContextApi) {
  const list = ctx.get('list')
  const editor = ctx.get('editor')

  if (!list) throw new Error('List instance not found')
  if (!editor) throw new Error('Editor instance not found')

  return { list, editor }
}

export const ListCloseEditor: ArtalkPlug = (ctx) => {
  let $closeCommentBtn: HTMLElement|undefined

  // on Artalk inited
  // (after all components had mounted)
  ctx.on('inited', () => {
    const { list } = ensureListEditor(ctx)

    $closeCommentBtn = list.$el.querySelector<HTMLElement>('[data-action="admin-close-comment"]')!

    // bind editor close button click event
    $closeCommentBtn.addEventListener('click', () => {
      const listData = list.getData()!

      listData.page.admin_only = !listData.page.admin_only
      list.adminPageEditSave(listData.page)
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

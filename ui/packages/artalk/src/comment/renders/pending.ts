import * as Utils from '../../lib/utils'
import RenderCtx from '../render-ctx'

/**
 * 待审核状态界面
 */
export default function renderPending(ctx: RenderCtx) {
  if (!ctx.data.is_pending) return

  const pendingEl = Utils.createElement(`<div class="atk-pending">${ctx.ctx.$t('pendingMsg')}</div>`)
  ctx.$body.prepend(pendingEl)
}

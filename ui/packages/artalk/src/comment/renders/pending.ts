import * as Utils from '../../lib/utils'
import type Render from '../render'

/**
 * 待审核状态界面
 */
export default function renderPending(r: Render) {
  if (!r.data.is_pending) return

  const pendingEl = Utils.createElement(
    `<div class="atk-pending">${r.ctx.$t('pendingMsg')}</div>`,
  )
  r.$body.prepend(pendingEl)
}

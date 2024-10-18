import { createSignal } from 'solid-js'
import { Dialog } from '../Dialog'
import type { AuthContext } from '../types'
import { DialogMergePageConfirm } from './DialogMergePageConfirm'

interface DialogMergeProps {
  ctx: AuthContext
  usernames: string[]
  onClose: () => void
}

export const DialogMerge = (props: DialogMergeProps) => {
  const { ctx, onClose, ...others } = props

  const pages = {
    confirm: () => (
      <DialogMergePageConfirm usernames={props.usernames} ctx={ctx} onClose={onClose} />
    ),
  }

  const homePage = 'confirm'
  const [page, setPage] = createSignal<keyof typeof pages>(homePage)

  return (
    <Dialog name="merge" onClose={onClose} title={() => 'Merge Tool'}>
      {() => pages[page()]()}
    </Dialog>
  )
}

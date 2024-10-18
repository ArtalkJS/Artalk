import { createSignal } from 'solid-js'
import { Dialog } from './Dialog'
import type { AuthContext } from './types'

interface DialogExampleProps {
  ctx: AuthContext
  onClose: () => void
}

export const DialogExample = (props: DialogExampleProps) => {
  const { ctx, onClose, ...others } = props

  const pages = {
    methods: () => <>Tool</>,
  }

  const homePage = 'methods'
  const [page, setPage] = createSignal<keyof typeof pages>(homePage)

  return (
    <Dialog onClose={onClose} title={() => ctx.$t('signIn')}>
      {() => pages[page()]()}
    </Dialog>
  )
}

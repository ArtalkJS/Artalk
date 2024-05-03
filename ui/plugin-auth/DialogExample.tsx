import type { ContextApi } from 'artalk'
import { createSignal } from 'solid-js'
import { Dialog } from './Dialog'

interface DialogExampleProps {
  ctx: ContextApi
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

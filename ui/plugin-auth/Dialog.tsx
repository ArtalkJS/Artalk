import { JSX } from 'solid-js'
import { DialogHeader } from './DialogHeader'

interface DialogProps {
  name?: string
  showBackBtn?: () => boolean
  onBack?: () => void
  onClose: () => void
  title: () => string
  children?: () => JSX.Element
  extras?: () => JSX.Element
}

export const Dialog = (props: DialogProps) => {
  const { showBackBtn, onBack, onClose, title, ...others } = props

  return (
    <div class="atk-auth-plugin-dialog-wrap">
      <div class="atk-auth-plugin-dialog" data-dialog-name={props.name}>
        <DialogHeader title={title} showBackBtn={showBackBtn} onBack={onBack} onClose={onClose} />
        <div class="atk-view-wrap atk-slim-scrollbar">{others.children?.()}</div>
      </div>
      {others.extras?.()}
    </div>
  )
}

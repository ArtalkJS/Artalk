export interface DialogHeaderProps {
  title: () => string
  showBackBtn?: () => boolean
  onBack?: () => void
  onClose: () => void
}

export const DialogHeader = (props: DialogHeaderProps) => {
  const { title, showBackBtn, onBack, onClose } = props

  return (
    <div class="atk-auth-dialog-title">
      {showBackBtn?.() && (
        <div class="atk-back-btn atk-icon atk-icon-arrow-left" onClick={onBack} />
      )}
      <div class="atk-text">{title()}</div>
      <div class="atk-close atk-icon-close atk-icon" onClick={() => onClose()}></div>
    </div>
  )
}

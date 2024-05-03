export interface DialogFooterProps {
  noText: string
  yesText: string
  onYes?: () => void
  onNo?: () => void
  yesDisabled?: () => boolean
}

export const DialogFooter = (props: DialogFooterProps) => (
  <div class="atk-dialog-footer">
    <div class="atk-btn atk-btn-no" onClick={props.onNo}>
      {props.noText}
    </div>
    <div
      class="atk-btn atk-btn-yes"
      classList={{ disabled: props.yesDisabled?.() }}
      onClick={() => !props.yesDisabled?.() && props.onYes?.()}
    >
      {props.yesText}
    </div>
  </div>
)

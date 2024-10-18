import type {
  Checker,
  CheckerCaptchaPayload,
  CheckerCtx,
  CheckerPayload,
} from '@/components/checker'

export interface CheckerManager {
  checkCaptcha: (payload: CheckerCaptchaPayload) => Promise<void>
  checkAdmin: (payload: CheckerPayload) => Promise<void>
  check: (checker: Checker, payload: CheckerPayload, beforeCheck?: (c: CheckerCtx) => void) => void
}

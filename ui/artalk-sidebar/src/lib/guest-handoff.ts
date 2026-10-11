import type { Context } from 'artalk'

// A fragment supplied by a direct link is not proof of a trusted parent.
export async function restoreGuestCredential(ctx: Context, token: string) {
  if (!token || window.parent === window) return
  try {
    const parent = new URL(document.referrer)
    if (!['http:', 'https:'].includes(parent.protocol)) return
    const { data } = await ctx.getApi().conf.getDomain({ url: parent.origin })
    if (data.is_trusted === true && data.origin === parent.origin) {
      ctx.getUser().setGuestToken?.(ctx.getConf().server, token)
    }
  } catch {
    // Missing referrers and failed trust checks cannot authorize a handoff.
  }
}

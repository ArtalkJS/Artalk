import type { LocalUser } from 'artalk'
import type { AuthContext } from '../types'

interface ResponseLoginData {
  user: Omit<LocalUser, 'token'>
  token: string
}

export const loginByApiRes = (ctx: AuthContext, { user, token }: ResponseLoginData) => {
  ctx.getUser().update({
    ...user,
    token,
  })
}

export const loginByToken = (ctx: AuthContext, token: string) => {
  ctx.getUser().update({ token })
  ctx
    .getApi()
    .user.getUser()
    .then((res) => {
      const { user } = res.data
      ctx.getUser().update({
        name: user.name,
        email: user.email,
        link: user.link,
        is_admin: user.is_admin,
      })
    })
}

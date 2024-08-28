import type { ContextApi, LocalUser } from 'artalk'

interface ResponseLoginData {
  user: LocalUser
  token: string
}

export const loginByApiRes = (ctx: ContextApi, { user, token }: ResponseLoginData) => {
  ctx.get('user').update({
    ...user,
    token,
  })
}

export const loginByToken = (ctx: ContextApi, token: string) => {
  ctx.get('user').update({ token })
  ctx
    .getApi()
    .user.getUser()
    .then((res) => {
      const { user } = res.data
      ctx.get('user').update({
        name: user.name,
        email: user.email,
        link: user.link,
        is_admin: user.is_admin,
      })
    })
}

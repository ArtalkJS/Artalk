import type { ContextApi } from 'artalk'

interface ResponseLoginData {
  user: {
    name: string
    email: string
    link: string
    is_admin: boolean
  }
  token: string
}

export const loginByApiRes = (ctx: ContextApi, data: ResponseLoginData) => {
  const { user, token } = data
  ctx.get('user').update({
    nick: user.name,
    email: user.email,
    link: user.link,
    isAdmin: user.is_admin,
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
        nick: user.name,
        email: user.email,
        link: user.link,
        isAdmin: user.is_admin,
      })
    })
}

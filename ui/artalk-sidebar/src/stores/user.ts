import type { LocalUser } from 'artalk'
import { defineStore } from 'pinia'
import sha256 from 'crypto-js/sha256'
import md5 from 'crypto-js/md5'
import { bootParams, getArtalk } from '../global'

interface UserState extends LocalUser {
  /**
   * Current site name
   */
  site: string
}

export const useUserStore = defineStore('user', {
  state: () =>
    <UserState>{
      site: bootParams.site || '',
      ...bootParams.user,
    },
  actions: {
    logout() {
      this.$reset()
      getArtalk()?.ctx.getUser().logout()
    },
    sync() {
      const user = getArtalk()?.ctx.getUser()
      if (!user) throw new Error('Artalk is not initialized')
      if (!user.checkHasBasicUserInfo()) throw new Error('User is not logged in')
      this.$patch({ ...user.getData(), site: '' })
    },
  },
  getters: {
    avatar: (state) => getGravatar(state.email),
  },
})

function getGravatar(email: string) {
  // TODO get avatar url from backend api
  const conf = getArtalk()?.ctx.conf?.gravatar
  if (!conf) return ''

  const emailHash =
    typeof conf.params == 'string' && conf.params.includes('sha256=1')
      ? sha256(email.toLowerCase()).toString()
      : md5(email.toLowerCase()).toString()

  return `${conf.mirror.replace(/\/$/, '')}/${emailHash}?${conf.params.replace(/^\?/, '')}`
}

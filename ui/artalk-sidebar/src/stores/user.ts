import { defineStore } from 'pinia'
import { bootParams, getArtalk } from '../global'
import sha256 from 'crypto-js/sha256'
import md5 from 'crypto-js/md5'

export const useUserStore = defineStore('user', {
  state: () => ({
    site: bootParams.site || '',
    name: bootParams.user.nick || '',
    email: bootParams.user.email || '',
    isAdmin: bootParams.user.isAdmin || false,
    token: bootParams.user.token || '',
  }),
  actions: {
    logout() {
      this.site = ''
      this.name = ''
      this.email = ''
      this.isAdmin = false
      this.token = ''

      getArtalk()?.ctx.get('user').logout()
    },
    sync() {
      const user = getArtalk()?.ctx.get('user')
      if (!user) throw new Error('Artalk is not initialized')
      if (!user.checkHasBasicUserInfo()) throw new Error('User is not logged in')
      const userData = user.getData()
      this.site = ''
      this.name = userData.nick
      this.email = userData.email
      this.isAdmin = userData.isAdmin
      this.token = userData.token
    },
  },
  getters: {
    avatar: (state) => {
      return getGravatar(state.email)
    },
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

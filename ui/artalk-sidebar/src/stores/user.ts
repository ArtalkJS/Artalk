import { defineStore } from 'pinia'
import { bootParams, getArtalk } from '../global'

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
    }
  }
})


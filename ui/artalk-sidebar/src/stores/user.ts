import { defineStore } from 'pinia'
import { bootParams } from '../global'

export const useUserStore = defineStore('user', () => {
  const site = ref(bootParams.site || '')
  const name = ref(bootParams.user.nick || '')
  const email = ref(bootParams.user.email || '')
  const isAdmin = ref(bootParams.user.isAdmin || false)
  const token = ref(bootParams.user.token || '')

  const logout = () => {
    site.value = ''
    name.value = ''
    email.value = ''
    isAdmin.value = false
    token.value = ''
  }

  return { site, name, email, isAdmin, token, logout }
})

import { defineStore } from 'pinia'
import { bootParams } from '../global'

export const useUserStore = defineStore('user', () => {
  const site = ref(bootParams.site || '')
  const name = ref(bootParams.user.nick || '')
  const token = ref(bootParams.user.token || '')

  return { site, name, token, }
})

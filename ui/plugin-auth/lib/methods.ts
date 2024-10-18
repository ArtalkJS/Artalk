import type { Api } from 'artalk'

export interface LoginMethod {
  name: string
  label: string
  icon: string
  link?: string
  onClick?: () => void
}

export const fetchMethods = async (api: Api) => {
  const { data } = await api.conf.getSocialLoginProviders()
  return data.providers
    .map<LoginMethod>(({ name, label, icon, path }) => {
      return { name, label, icon, link: path }
    })
    .sort((a, b) => {
      // email always on top
      if (a.name === 'email') return -1
      if (b.name === 'email') return 1
      // skip always on bottom
      if (a.name === 'skip') return 1
      if (b.name === 'skip') return -1
      // others by label
      return a.label.localeCompare(b.label)
    })
}

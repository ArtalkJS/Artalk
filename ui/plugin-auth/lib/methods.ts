import type { Api } from 'artalk'

export interface LoginMethod {
  name: string
  label: string
  icon: string
  link?: string
  callbackOrigin?: string
  onClick?: () => void
}

export const fetchMethods = async (api: Api) => {
  const { data } = await api.conf.getSocialLoginProviders()
  // Older backend/API types do not expose callback_origin.
  const callbackOrigin =
    'callback_origin' in data && typeof data.callback_origin === 'string'
      ? data.callback_origin
      : undefined
  return data.providers
    .map<LoginMethod>(({ name, label, icon, path }) => {
      return { name, label, icon, link: path, callbackOrigin }
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

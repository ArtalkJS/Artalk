import { defineConfig } from 'vitepress'
import { shared } from './shared'
import { zh } from './zh'
import { en } from './en'

export default defineConfig({
  ...shared,
  locales: {
    en: { label: 'English', ...en },
    zh: { label: '简体中文', ...zh },
  },
})

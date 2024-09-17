import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import 'normalize.css'
import './App.scss'
import './darkmode.scss'

import { useDarkMode } from './hooks/darkmode'
import { Header } from './components/Header'
import { Slogan } from './components/Slogan'
import { Footer } from './components/Footer'
import { SlightFeature } from './components/Features/SlightFeature'
import { QuickFeature } from './components/Features/QuickFeature'
import { FullFeature } from './components/Features/FullFeature'
import { FuncFeature } from './components/Features/FuncFeature'
import { SafeFeature } from './components/Features/SafeFeature'

function App() {
  const { t, i18n } = useTranslation()
  const { isDarkMode, toggle: toggleDarkMode } = useDarkMode()

  useEffect(() => {
    document.title = t('home_title')
    document.documentElement.lang = i18n.language
  }, [t, i18n.language])

  useEffect(() => {
    if (isDarkMode) document.documentElement.classList.add('app-dark-mode')
    else document.documentElement.classList.remove('app-dark-mode')
  }, [isDarkMode])

  return (
    <>
      <Header darkModeHandler={{ isDarkMode, toggle: toggleDarkMode }} />
      <Slogan />

      <SlightFeature />
      <QuickFeature />
      <FullFeature />
      <FuncFeature />
      <SafeFeature />

      <Footer />
    </>
  )
}

export default App

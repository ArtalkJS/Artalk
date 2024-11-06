import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import 'normalize.css'
import './App.scss'
import './darkmode.scss'

import { useDarkMode } from './hooks/darkmode'
import { Header } from './components/Header'
import { Slogan } from './components/Slogan'
import { Footer } from './components/Footer'
import { SlightFeature } from './components/Features/Slight'
import { QuickFeature } from './components/Features/Quick'
import { FullFeature } from './components/Features/Full'
import { FuncsFeature } from './components/Features/Funcs'
import { SafeFeature } from './components/Features/Safe'
import { Plugins } from './components/Features/Plugins'
import { Group } from './components/FeatureComps/Group'
import { StartFeature } from './components/Features/Start'
import { FrameworksFeature } from './components/Features/Frameworks'
import { GreyCardsFeature } from './components/Features/GreyCards'

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

      <StartFeature />
      <GreyCardsFeature />

      <FrameworksFeature />

      <Group>
        <SlightFeature />
        <QuickFeature />
      </Group>
      <Group>
        <SafeFeature />
      </Group>
      <Group>
        <Plugins />
      </Group>
      <FullFeature />
      <FuncsFeature />

      <Footer />
    </>
  )
}

export default App

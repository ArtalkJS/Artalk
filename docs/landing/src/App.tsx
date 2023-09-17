import 'normalize.css'
import './App.scss'

import { Header } from './components/Header'
import { Slogan } from './components/Slogan'
import { Footer } from './components/Footer'
import { SlightFeature } from './components/Features/SlightFeature'
import { QuickFeature } from './components/Features/QuickFeature'
import { FullFeature } from './components/Features/FullFeature'
import { FuncFeature } from './components/Features/FuncFeature'
import { SafeFeature } from './components/Features/SafeFeature'

function App() {
  return (
    <>
      <Header />
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

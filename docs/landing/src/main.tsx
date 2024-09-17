import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import { initI18n } from './i18n'
import './responsive.scss'

initI18n()

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)

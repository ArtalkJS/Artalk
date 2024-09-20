import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Popover, PopoverButton, PopoverPanel, useClose } from '@headlessui/react'
import { MdArrowDropDown, MdClose, MdDarkMode, MdMenu, MdSunny } from 'react-icons/md'
import { AiFillGithub } from 'react-icons/ai'
import { IoLanguage } from 'react-icons/io5'
import './Header.scss'

const LanguageList: React.FC = () => {
  const close = useClose()

  const { i18n } = useTranslation()

  const languages = [
    { name: 'English', code: 'en' },
    { name: '简体中文', code: 'zh' },
    { name: '繁體中文', code: 'zh-TW' },
    { name: 'Français', code: 'fr' },
    { name: '日本語', code: 'ja' },
    { name: '한국어', code: 'ko' },
    { name: 'Русский', code: 'ru' },
  ]

  const changeLanguage = (lang: string) => {
    i18n.changeLanguage(lang)
    close()
  }

  return (
    <>
      {languages.map((lang) => (
        <a
          key={lang.code}
          className="language-item"
          onClick={() => {
            changeLanguage(lang.code)
          }}
        >
          {lang.name}
        </a>
      ))}
    </>
  )
}

export interface HeaderProps {
  darkModeHandler: { isDarkMode: boolean; toggle: () => void }
}

export const Header: React.FC<HeaderProps> = ({ darkModeHandler }) => {
  const { t } = useTranslation()
  const [fixed, setFixed] = useState(false)
  const [mobileShow, setMobileShow] = useState(false)

  useEffect(() => {
    const scrollHandler = () => {
      if (window.scrollY > 70) setFixed(true)
      else setFixed(false)
    }

    scrollHandler()
    document.addEventListener('scroll', scrollHandler)

    return () => {
      document.removeEventListener('scroll', scrollHandler)
    }
  }, [])

  return (
    <div className={['app-header', fixed ? 'fixed' : ''].join(' ')}>
      <div className="container">
        {!mobileShow && <div className="brand">Artalk</div>}

        <div className={['links', mobileShow ? 'mobile-show' : ''].join(' ')}>
          <a className="link-item" href={t('nav_docs_link')}>
            {t('docs')}
          </a>
          <a className="link-item" href={t('nav_changelog_link')} target="__blank">
            {t('changelog')}
          </a>
          <a className="link-item" href={t('github_link')} target="__blank">
            <AiFillGithub size="1.7em" />
          </a>
          <a className="link-item" onClick={() => darkModeHandler.toggle()}>
            {darkModeHandler.isDarkMode ? (
              <MdSunny size="1.7em" style={{ color: '#f1c40f' }} />
            ) : (
              <MdDarkMode size="1.7em" style={{ color: '#f39c12' }} />
            )}
          </a>
          <Popover className="link-item">
            <PopoverButton className="language-toggle">
              <IoLanguage size="1.7em" style={{ color: '#3498db' }} />
              <MdArrowDropDown size="1.5em" style={{ color: '#3498db' }} />
            </PopoverButton>
            <PopoverPanel anchor="bottom" className="header-language-list">
              <LanguageList />
            </PopoverPanel>
          </Popover>
        </div>

        <a className="menu-toggle-mobile" onClick={() => setMobileShow((prev) => !prev)}>
          {!mobileShow ? <MdMenu size="1.7em" /> : <MdClose size="1.7em" />}
        </a>
      </div>
    </div>
  )
}

import React, { useEffect, useState } from 'react'
import './Header.scss'
import { AiFillGithub } from 'react-icons/ai'

export const Header: React.FC = () => {
  const [fixed, setFixed] = useState(false)

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
      <div className='container'>
        <div className='brand'>
          Artalk
        </div>

        <div className="links">
          <a className="link-item" href='https://artalk.js.org/guide/intro.html'>Docs</a>
          <a className="link-item" href='https://github.com/ArtalkJS/Artalk/blob/master/CHANGELOG.md' target="__blank">Changelog</a>
          <a className="link-item" href='https://github.com/ArtalkJS/Artalk' target="__blank">
            <AiFillGithub size='1.7em' />
          </a>
        </div>
      </div>
    </div>
  )
}

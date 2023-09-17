import React from 'react'
import './Footer.scss'

export const Footer: React.FC = () => {
  return (
    <div className='app-footer container'>
      <div className='brand'>Artalk</div>
      <div className='copyright'>The Artalk. Made with <span className='red'>♥️</span>.</div>
    </div>
  )
}

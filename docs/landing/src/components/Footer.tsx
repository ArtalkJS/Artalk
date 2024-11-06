import React from 'react'
import './Footer.scss'
import { RiArrowRightSLine } from 'react-icons/ri'
import { useTranslation } from 'react-i18next'

export const Footer: React.FC = () => {
  const { t } = useTranslation()

  return (
    <>
      <div className="footer-see-more-link container">
        <a
          className="link"
          href="https://github.com/ArtalkJS/Artalk"
          target="_blank"
          rel="noreferrer"
        >
          <span className="text">{t('see_more_on_github')}</span>
          <span className="icon">
            <RiArrowRightSLine />
          </span>
        </a>

        <span className="star-proposal">
          {t('star_proposal_line_1')}
          <br />
          {t('star_proposal_line_2')}
        </span>
      </div>

      <div className="app-grey-bg app-footer-wrap">
        <div className="app-footer container">
          <div className="brand">Artalk</div>
          <div className="copyright">
            The Artalk. Made with <span className="red">♥️</span>.
          </div>
        </div>
      </div>
    </>
  )
}

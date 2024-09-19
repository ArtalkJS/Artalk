import React from 'react'
import { FaArrowRight } from 'react-icons/fa'
import { useTranslation } from 'react-i18next'
import './LearnMoreLink.scss'

interface LearnMoreLinkProps {
  prompt: string
  link: string
  linkText?: string
}

export const LearnMoreLink: React.FC<LearnMoreLinkProps> = (props) => {
  const { t } = useTranslation()

  return (
    <div className="app-learn-more-link">
      <span className="prompt">{props.prompt}</span>&nbsp;&nbsp;
      <a className="link-btn" href={props.link} target="_blank" rel="noreferrer">
        {props.linkText || t('learn_more')}{' '}
        <span className="icon">
          <FaArrowRight size=".8em" />
        </span>
      </a>
    </div>
  )
}

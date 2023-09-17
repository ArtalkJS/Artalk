import React from 'react'
import { FaArrowRight } from 'react-icons/fa'
import './LearnMoreLink.scss'

interface LearnMoreLinkProps {
  prompt: string
  link: string
}

export const LearnMoreLink: React.FC<LearnMoreLinkProps> = (props) => {
  return (
    <div className="app-learn-more-link">
      <span className="prompt">{props.prompt}</span>
      <a className="link-btn" href={props.link} target='_blank'>了解更多 <span className='icon'><FaArrowRight size='.8em' /></span></a>
    </div>
  )
}

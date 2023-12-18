import './Slogan.scss'
import { FaArrowRight } from 'react-icons/fa'
import { Reveal } from './Reveal'

export const Slogan: React.FC = () => {
  return (
    <div className='slogan'>
      <div className="slogan-inner">
        <div className="text">
        <Reveal>
          A <span className='highlight'>
            Self-hosted
            <svg className='line-wrap' width="361" height="55" viewBox="0 0 361 55" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path className='line' d="M6 49C39.3707 29.6297 79.037 17.9573 117.534 28.5185C123.875 30.2581 129.86 32.8008 135.527 36.1667C139.986 38.8153 144.204 42.404 149.501 43.037C159.524 44.2349 168.187 33.3778 176.108 28.6481C191.347 19.5481 210.215 16.0215 227.536 21.1296C239.243 24.5822 249.671 29.5512 261.991 30.9815C295.794 34.9057 323.045 15.5757 354 7" stroke="#BDCFFF" stroke-width="12" stroke-linecap="round"/>
            </svg>
          </span>
          <br/>
        </Reveal>
        <Reveal delay={200}>Comment System</Reveal>
        </div>
        <div className='btns'>
          <a className='blue btn' href='https://artalk.js.org/guide/intro.html'>获取 Artalk</a>
          <a className='btn' href='https://github.com/ArtalkJS/Artalk' target="_blank">GitHub <FaArrowRight size='.8em' /></a>
        </div>
      </div>
      <div className='bg'>
      </div>
    </div>
  )
}

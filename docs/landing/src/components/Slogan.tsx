import './Slogan.scss'
import { useTranslation } from 'react-i18next'
import { useEffect, useRef, useState } from 'react'
import { FaArrowRight } from 'react-icons/fa'
import { Reveal } from './Reveal'

export const Slogan: React.FC = () => {
  const { t } = useTranslation()
  type VideoState = 'ready' | 'loading' | 'playing' | 'paused'
  const [videoState, setVideoState] = useState<VideoState>('ready')
  const videoLink =
    'https://github.com/user-attachments/assets/d8d1597a-7e1f-45f3-963e-371c0528bf91'
  const videoMimeType = 'video/webm' // Must set for Safari
  const videoElRef = useRef<HTMLVideoElement>(null)

  useEffect(() => {
    // load video in background
    if (videoState === 'ready') {
      setVideoState('loading')
      setTimeout(() => {
        const videoEl = videoElRef.current
        if (videoEl) {
          const source = document.createElement('source')
          source.src = videoLink
          source.type = videoMimeType
          videoEl.appendChild(source)
          videoEl.load()
          videoEl.onloadeddata = () => {
            setVideoState('playing')
          }
          videoEl.onended = () => {
            setTimeout(() => {
              setVideoState('paused')
              setTimeout(() => {
                setVideoState('playing')
                if (navigator.userAgent.includes('Safari')) {
                  videoEl.load()
                } else {
                  videoEl.currentTime = 0
                  videoEl.play()
                }
              }, 5000)
            }, 1000)
          }
        }
      }, 1000)
    }
  }, [videoElRef, videoState])

  return (
    <div className="slogan">
      <div className="slogan-inner">
        <div className="text">
          <Reveal>
            <span className="highlight">
              {t('slogan_line_1')}
              {/* <svg
                className="line-wrap"
                width="361"
                height="55"
                viewBox="0 0 361 55"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  className="line"
                  d="M6 49C39.3707 29.6297 79.037 17.9573 117.534 28.5185C123.875 30.2581 129.86 32.8008 135.527 36.1667C139.986 38.8153 144.204 42.404 149.501 43.037C159.524 44.2349 168.187 33.3778 176.108 28.6481C191.347 19.5481 210.215 16.0215 227.536 21.1296C239.243 24.5822 249.671 29.5512 261.991 30.9815C295.794 34.9057 323.045 15.5757 354 7"
                  stroke="#BDCFFF"
                  strokeWidth="12"
                  strokeLinecap="round"
                />
              </svg> */}
            </span>
            <br />
          </Reveal>
          <Reveal delay={200}>{t('slogan_line_2')}</Reveal>
          <Reveal delay={400}>{t('slogan_line_3')}</Reveal>

          <div className="btns">
            <a className="blue btn" href={t('get_artalk_link')}>
              {t('get_artalk')}
            </a>
            <a className="btn" href={t('github_link')} target="_blank" rel="noreferrer">
              GitHub <FaArrowRight size=".8em" />
            </a>
          </div>
        </div>
        <div className="demos">
          {videoState !== 'playing' && (
            <div className="img-wrap">
              <img className="demo-video anim-fade-in" src="/images/demo-video-1-thumbnail.png" />
            </div>
          )}
          <video
            style={{ display: videoState !== 'playing' ? 'none' : '' }}
            className="demo-video anim-fade-in"
            ref={videoElRef}
            autoPlay
            muted
            playsInline
          />
        </div>
      </div>
      <div className="bg"></div>
    </div>
  )
}

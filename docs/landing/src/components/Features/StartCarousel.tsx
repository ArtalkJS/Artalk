import { useTranslation } from 'react-i18next'
import { FiArrowUpRight } from 'react-icons/fi'
import { useCallback, useEffect, useState } from 'react'
import LazyImage from '../LazyImage'
import { getFeatures } from '../../features'
import './StartCarousel.scss'

interface IFeature {
  name: string
  desc: string
  icon: React.FC
  link: string
  active: boolean
  bgColor: string
  textColor: string
}

const allScreenshots = [
  '/images/screenshots/1.webp',
  '/images/screenshots/2.webp',
  '/images/screenshots/3.webp',
  '/images/screenshots/4.webp',
  '/images/screenshots/5.webp',
  '/images/screenshots/6.webp',
  '/images/screenshots/7.webp',
  '/images/screenshots/8.webp',
]

const colors = [
  // Teal
  {
    bgColor: '#B2DFDB',
    textColor: '#004D40',
  },
  // Blue
  {
    bgColor: '#BBDEFB',
    textColor: '#0D47A1',
  },
  // Green
  {
    bgColor: '#C8E6C9',
    textColor: '#388E3C',
  },
  // Cyan
  {
    bgColor: '#E0F7FA',
    textColor: '#006064',
  },
  // Purple
  {
    bgColor: '#D1C4E9',
    textColor: '#512DA8',
  },
  // Orange
  {
    bgColor: '#FFCCBC',
    textColor: '#E64A19',
  },
  // Red
  {
    bgColor: '#FFC1C1',
    textColor: '#C62828',
  },
  // Yellow
  {
    bgColor: '#FFF59D',
    textColor: '#F57F17',
  },
  // Pink
  {
    bgColor: '#F8BBD0',
    textColor: '#C2185B',
  },
  // Lime
  {
    bgColor: '#E6EE9C',
    textColor: '#827717',
  },
]

export const StartCarousel = () => {
  const { t } = useTranslation()

  const [allFeatures, setAllFeatures] = useState<
    {
      icon: React.FC
      name: string
      desc: string
      link: string
    }[]
  >([])

  useEffect(() => {
    setAllFeatures(getFeatures(t))
  }, [t])

  const [currentFeatures, setCurrentFeatures] = useState<IFeature[]>([])
  const [currentScreenshots, setCurrentScreenshots] = useState<string[]>([])

  const [dots, setDots] = useState<{ percent: number; active: boolean }[]>(
    new Array(12).fill(0).map(() => ({ percent: 0, active: false })),
  )
  const [currentIndex, setCurrentIndex] = useState(0)
  const [currentPercent, setCurrentPercent] = useState(0)

  useEffect(() => {
    const interval = setInterval(() => {
      setDots((prevDots) =>
        prevDots.map((_, i) => ({
          percent: i === currentIndex ? currentPercent : 0,
          active: i === currentIndex,
        })),
      )

      setCurrentPercent((p) => {
        if (p < 100) return p + 1
        setCurrentIndex((index) => (index + 1) % dots.length)
        return 0
      })
    }, 80)

    return () => clearInterval(interval)
  }, [dots.length, currentIndex, currentPercent])

  const perPage = 2

  const refreshFeatures = useCallback(() => {
    const newFeatures: IFeature[] = allFeatures
      .slice(
        (currentIndex * perPage) % allFeatures.length,
        (currentIndex * perPage + 2) % allFeatures.length,
      )
      .map((feature, index) => ({
        ...feature,
        active: false,
        bgColor: colors[(currentIndex * perPage + index) % colors.length].bgColor,
        textColor: colors[(currentIndex * perPage + index) % colors.length].textColor,
      }))

    if (currentPercent < 50) {
      if (newFeatures[0]) newFeatures[0].active = true
      if (newFeatures[1]) newFeatures[1].active = false
    } else {
      if (newFeatures[0]) newFeatures[0].active = false
      if (newFeatures[1]) newFeatures[1].active = true
    }

    setCurrentFeatures(newFeatures)
  }, [currentIndex, currentPercent, allFeatures])

  useEffect(() => {
    const start = (currentIndex * perPage) % allScreenshots.length
    let end = (currentIndex * perPage + 2) % allScreenshots.length
    if (start > end) end = allScreenshots.length
    const newScreenshots = allScreenshots.slice(start, end)
    setCurrentScreenshots(newScreenshots)
  }, [currentIndex])

  useEffect(() => {
    if (![0, 50].includes(currentPercent)) return
    refreshFeatures()
  }, [currentIndex, currentPercent, refreshFeatures])

  useEffect(() => {
    refreshFeatures()
  }, [allFeatures, refreshFeatures])

  const switchPage = useCallback(
    (index: number) => {
      setCurrentIndex(index)
      setCurrentPercent(0)
    },
    [setCurrentIndex, setCurrentPercent],
  )

  return (
    <div className="start-carousel">
      <div className="feature-cards">
        {currentFeatures.map((item) => (
          <a
            className={['feature-item', item.active ? 'active' : ''].join(' ')}
            key={item.name}
            href={item.link}
            target="_blank"
            rel="noreferrer"
          >
            <div className="icon" style={{ backgroundColor: item.bgColor, color: item.textColor }}>
              <item.icon />
            </div>
            <div className="main">
              <div className="title">{item.name}</div>
              <div className="desc">{item.desc}</div>
            </div>
            <div className="action">
              <div className="arrow-btn">
                <FiArrowUpRight size="2em" />
              </div>
            </div>
          </a>
        ))}
      </div>

      <div className="progress-indicator">
        <div className="dots">
          {dots.map((dot, index) => (
            <div
              key={index}
              className={['dot', dot.active ? 'active' : ''].join(' ')}
              onClick={() => switchPage(index)}
            >
              <div className="progress" style={{ width: `${dot.percent}%` }} />
            </div>
          ))}
        </div>
      </div>

      <div className="screenshot-group">
        {currentScreenshots.map((screenshot) => (
          <div className="screenshot" key={screenshot}>
            <LazyImage src={screenshot} alt="screenshot" referrerPolicy="no-referrer" />
          </div>
        ))}
      </div>
    </div>
  )
}

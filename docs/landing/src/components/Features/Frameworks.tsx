import './Frameworks.scss'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import VueIcon from '../../assets/techs/vue.svg'
import ReactIcon from '../../assets/techs/react.svg'
import SolidIcon from '../../assets/techs/solid.svg'
import SvelteIcon from '../../assets/techs/svelte.svg'
import HexoIcon from '../../assets/techs/hexo.svg'
import VitepressIcon from '../../assets/techs/vitepress.svg'
import WordpressIcon from '../../assets/techs/wordpress.svg'
import GoIcon from '../../assets/techs/go.svg'

const techs = [
  {
    name: 'Vue',
    icon: VueIcon,
    color: '#41b883',
    link: 'https://artalk.js.org/develop/import-framework.html#vue',
  },
  {
    name: 'React',
    icon: ReactIcon,
    color: '#116A92',
    link: 'https://artalk.js.org/develop/import-framework.html#react',
  },
  {
    name: 'SolidJS',
    icon: SolidIcon,
    color: '#27488B',
    link: 'https://artalk.js.org/develop/import-framework.html#solidjs',
  },
  {
    name: 'Svelte',
    icon: SvelteIcon,
    color: '#ff3e00',
    link: 'https://artalk.js.org/develop/import-framework.html#svelte',
  },
  {
    name: 'VitePress',
    icon: VitepressIcon,
    color: '#8F4CFF',
    link: 'https://artalk.js.org/develop/import-blog.html#vitepress',
  },
  {
    name: 'Hexo',
    icon: HexoIcon,
    color: '#136DC1',
    link: 'https://artalk.js.org/develop/import-blog.html#hexo',
  },
  {
    name: 'WordPress',
    icon: WordpressIcon,
    color: '#1C6189',
    link: 'https://artalk.js.org/develop/import-blog.html#wordpress',
  },
  {
    name: 'Go',
    icon: GoIcon,
    color: '#00acd7',
    link: 'https://artalk.js.org/guide/deploy.html',
  },
]

export const FrameworksFeature: React.FC = () => {
  const { t } = useTranslation()
  const [activeIndex, setActiveIndex] = useState(-1)
  const [stop, setStop] = useState(false)

  useEffect(() => {
    if (stop) return
    const timer = setInterval(() => {
      setActiveIndex((i) => {
        if (i + 1 == techs.length) return -1
        return (i + 1) % techs.length
      })
    }, 2000)
    return () => clearInterval(timer)
  }, [stop])

  const onMouseOver = (i: number) => {
    setStop(true)
    setActiveIndex(i)
  }

  const onMouseOut = () => {
    setStop(false)
  }

  return (
    <div className="frameworks-feature container">
      <div className="title">
        <div className="text-row">{t('use_artalk_with')}</div>
        <div
          className="text-row"
          style={{ color: activeIndex != -1 ? techs[activeIndex].color : undefined }}
        >
          {activeIndex == -1 ? t('any_website_or_blog') : techs[activeIndex].name}
        </div>
      </div>
      <div className="methods">
        {techs.map((tech, i) => (
          <a
            key={i}
            className={['icon', activeIndex == i ? 'active' : ''].join(' ')}
            style={{ backgroundImage: `url('${tech.icon.replace(/'/g, "\\'")}')` }}
            href={tech.link}
            target="_blank"
            rel="noopener noreferrer"
            onMouseOver={() => onMouseOver(i)}
            onMouseOut={() => onMouseOut()}
          ></a>
        ))}
      </div>
    </div>
  )
}

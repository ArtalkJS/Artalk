import { FiArrowUpRight } from 'react-icons/fi'
import { useTranslation } from 'react-i18next'
import './GreyCards.scss'
import LazyImage from '../LazyImage'

export const GreyCardsFeature: React.FC = () => {
  const { t } = useTranslation()

  return (
    <div className="grey-cards-feature container">
      <div className="section">
        <div className="section-title">{t('concept_section_locales_title')}</div>
        <div className="section-desc">{t('concept_section_locales_desc')}</div>

        <div className="window text-window">
          <div className="text-row">Artalk Self-hosted Comment System</div>
          <div className="text-row">Artalk 自托管评论系统</div>
          <div className="text-row">Artalk Auto-hébergé Commentaire Système</div>
          <div className="text-row">Artalk セルフホスティング コメント システム</div>
          <div className="text-row">Artalk 셀프 호스팅 댓글 시스템</div>
          <div className="text-row">Artalk Cамохостинг Комментарий Система</div>
        </div>
      </div>

      <div className="section">
        <div className="section-title">{t('concept_section_docker_title')}</div>
        <div className="section-desc">
          {t('concept_section_docker_desc_line_1')}
          <br />
          <br />
          {t('concept_section_docker_desc_line_2')}
        </div>
        <a
          className="section-learn-more"
          target="_blank"
          href="https://artalk.js.org/zh/guide/deploy.html"
          rel="noreferrer"
        >
          {t('learn_more')}
          <FiArrowUpRight />
        </a>
        <div className="window terminal-window">
          <div className="docker-icon">
            <LazyImage src="/images/docker-folder.png" alt="Docker" />
          </div>
          <div className="menubar">
            <div className="close"></div>
            <div className="minimize"></div>
            <div className="zoom"></div>
          </div>
          <div className="window-body">
            <div className="cmd-row">
              <div className="prompt">$</div>
              <div className="command">docker pull artalk/artalk</div>
            </div>
            <div className="cmd-row">
              <div className="prompt">$</div>
              <div className="command">docker run -d -p 8080:23366 artalk/artalk</div>
            </div>
          </div>
        </div>
      </div>

      <div className="section">
        <div className="section-title">{t('concept_section_open_source_title')}</div>
        <div className="section-desc">
          {t('concept_section_open_source_desc_line_1')}
          <br />
          <br />
          {t('concept_section_open_source_desc_line_2', { year: new Date().getFullYear() - 2018 })}
        </div>
        <LazyImage
          src="/images/github-repo-page.webp"
          alt="GitHub Repository"
          className="screenshot"
        />
      </div>
    </div>
  )
}

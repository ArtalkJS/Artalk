import React from 'react'
import { useTranslation } from 'react-i18next'
import './FullFeature.scss'
import { FeatureTitle } from '../FeatureTitle'
import { FeatureDesc } from '../FeatureDesc'
import { Reveal } from '../Reveal'
import { FeatureBase } from './FeatureBase'
import { FolderIcon } from './FullFeatureIcon'
import { FullFeatureList } from './FullFeatureList'

export const FullFeature: React.FC = () => {
  const { t } = useTranslation()

  return (
    <FeatureBase className="full">
      <FeatureTitle text={t('feature_full_title')} color="var(--color-font-secondary)" />

      <div className="content">
        <main>
          <div className="desc">
            <Reveal delay={200}>
              <FeatureDesc>{t('feature_full_desc_line_1')}</FeatureDesc>
            </Reveal>
          </div>
        </main>
        <aside>
          <FolderIcon />
        </aside>
      </div>

      <FullFeatureList />
    </FeatureBase>
  )
}

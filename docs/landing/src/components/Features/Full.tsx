import React from 'react'
import { useTranslation } from 'react-i18next'
import './Full.scss'
import { Title } from '../FeatureComps/Title'
import { Desc } from '../FeatureComps/Desc'
import { Reveal } from '../Reveal'
import { Base } from '../FeatureComps/Base'
import { FolderIcon } from '../FeatureComps/FolderIcon'
import { FullList } from './FullList'

export const FullFeature: React.FC = () => {
  const { t } = useTranslation()

  return (
    <div className="app-grey-bg">
      <Base className="full-feature">
        <Title text={t('feature_full_title')} color="var(--color-font-secondary)" />

        <div className="content">
          <main>
            <div className="desc">
              <Reveal delay={200}>
                <Desc>{t('feature_full_desc_line_1')}</Desc>
              </Reveal>
            </div>
          </main>
          <aside>
            <FolderIcon />
          </aside>
        </div>

        <FullList />
      </Base>
    </div>
  )
}

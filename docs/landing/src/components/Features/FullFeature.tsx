import React from 'react'
import './FullFeature.scss'
import { FeatureBase } from './FeatureBase'
import { FeatureTitle } from '../FeatureTitle'
import { FeatureDesc } from '../FeatureDesc'
import { Reveal } from '../Reveal'
import { FolderIcon } from './FullFeatureIcon'
import { FullFeatureList } from './FullFeatureList'

export const FullFeature: React.FC = () => {
  return (
    <FeatureBase className='full'>
      <FeatureTitle text='全面' color='#67757F' />

      <div className="content">
        <main>
          <div className="desc">
            <Reveal delay={200}>
              <FeatureDesc>
              Artalk 提供丰富的内置功能，我们尽力在简洁的同时保持功能的全面，为您带来开箱即用的体验。
              </FeatureDesc>
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

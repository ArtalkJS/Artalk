import React from 'react'
import './Quick.scss'
import { useTranslation } from 'react-i18next'
import { FiArrowUpRight } from 'react-icons/fi'
import { BlockBase } from '../FeatureComps/BlockBase'

export const QuickFeature: React.FC = () => {
  const { t } = useTranslation()

  const ArrowIcon = (
    <div className="arrow-icon">
      <div className="block" />
      <FiArrowUpRight size="2em" />
    </div>
  )

  return (
    <BlockBase className="quick-feature">
      <div className="feature-title" style={{ color: '#5DADEC' }}>
        {t('feature_swift_title')}

        <div className="airplane-wrap">
          <svg
            width="80"
            height="80"
            viewBox="0 0 361 343"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <g className="airplane">
              <path
                d="M327.5 133.04C325.317 135.228 321.788 135.228 319.605 133.04L311.71 125.145C309.522 122.962 309.522 119.427 311.71 117.25L327.5 101.46C329.683 99.2772 333.212 99.2772 335.395 101.46L343.29 109.355C345.478 111.538 345.478 115.072 343.29 117.25L327.5 133.04ZM243.75 49.2896C241.567 51.4783 238.038 51.4783 235.855 49.2896L227.96 41.3948C225.772 39.2117 225.772 35.6831 227.96 33.5L243.75 17.7103C245.933 15.5272 249.462 15.5272 251.645 17.7103L259.54 25.6051C261.728 27.7882 261.728 31.3169 259.54 33.5L243.75 49.2896Z"
                fill="#66757F"
              />
              <path
                d="M171.167 122.833C182.333 122.833 232.583 128.417 232.583 128.417C232.583 128.417 238.167 178.667 238.167 189.833C238.167 201 227 201 221.417 195.417C215.833 189.833 199.083 161.917 199.083 161.917C199.083 161.917 171.167 145.167 165.583 139.583C160 134 160 122.833 171.167 122.833ZM182.333 33.7178C199.083 33.5 321.917 39.0833 321.917 39.0833C321.917 39.0833 327.076 161.917 327.288 178.667C327.5 195.417 311.861 201.011 305.725 178.672C299.589 156.333 277.25 83.75 277.25 83.75C277.25 83.75 199.335 59.0884 182.305 55.2638C160 50.25 165.578 33.9299 182.333 33.7178Z"
                fill="#55ACEE"
              />
              <path
                d="M310.75 16.75C321.917 5.58331 349.833 -2.00272e-05 355.417 5.58331C361 11.1666 355.417 39.0833 344.25 50.25C333.083 61.4166 238.167 150.75 238.167 150.75C238.167 150.75 201.875 175.875 193.5 167.5C185.125 159.125 210.25 122.833 210.25 122.833C210.25 122.833 299.583 27.9166 310.75 16.75Z"
                fill="#CCD6DD"
              />
              <path
                d="M238.167 122.833C238.167 122.833 240.958 125.625 215.833 150.75C190.708 175.875 187.917 173.083 187.917 173.083C187.917 173.083 185.125 170.292 210.25 145.167C235.375 120.042 238.167 122.833 238.167 122.833ZM321.917 22.3333C331.168 22.3333 338.667 29.8317 338.667 39.0833H343.161C343.798 37.3246 344.25 35.4821 344.25 33.5C344.25 24.2484 336.752 16.75 327.5 16.75C325.518 16.75 323.675 17.2022 321.917 17.8387V22.3333Z"
                fill="#66757F"
              />
            </g>
          </svg>
        </div>
      </div>

      <div className="content">
        <main>
          <div className="deploy-links">
            <a className="item" href={t('deploy_bin_link')} target="_blank" rel="noreferrer">
              <span className="bold">{t('deploy_bin')}</span> <span>{t('deploy_bin_sub')}</span>
              {ArrowIcon}
            </a>
            <a className="item" href={t('deploy_docker_link')} target="_blank" rel="noreferrer">
              <span className="bold">{t('deploy_docker')}</span>{' '}
              <span>{t('deploy_docker_sub')}</span>
              {ArrowIcon}
            </a>
          </div>
        </main>
      </div>
    </BlockBase>
  )
}

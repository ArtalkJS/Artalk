import React from 'react'
import './QuickFeature.scss'
import { useTranslation } from 'react-i18next'
import { FiArrowUpRight } from 'react-icons/fi'
import { FeatureTitle } from '../FeatureTitle'
import { LearnMoreLink } from '../LearnMoreLink'
import { FeatureDesc } from '../FeatureDesc'
import { Reveal } from '../Reveal'
import { FeatureBase } from './FeatureBase'

export const QuickFeature: React.FC = () => {
  const { t } = useTranslation()

  const ArrowIcon = (
    <div className="arrow-icon">
      <div className="block" />
      <FiArrowUpRight size="2em" />
    </div>
  )

  const [isCopied, setIsCopied] = React.useState(false)

  let copiedResetTimer: any

  const copyDockerURLHandler = (e: React.MouseEvent) => {
    e.preventDefault()
    clearTimeout(copiedResetTimer)
    setIsCopied(false)

    const dockerURL = 'docker pull artalk/artalk-go'
    if (!navigator.clipboard.writeText) {
      alert('Browser does not support clipboard API')
      return
    }

    navigator.clipboard.writeText(dockerURL)
    setIsCopied(true)

    copiedResetTimer = setTimeout(() => {
      setIsCopied(false)
    }, 2000)
  }

  return (
    <FeatureBase className="quick">
      <FeatureTitle text={t('feature_swift_title')} color="#5DADEC" />

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

        <aside>
          <div className="airplane-wrap">
            <svg
              width="361"
              height="343"
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
        </aside>
      </div>

      {/* Golang Description */}
      <div className="extra-desc">
        <div className="left">
          <div className="text">
            <Reveal>
              <FeatureDesc>
                {t('feature_swift_desc_line_1')}
                <br />
                {t('feature_swift_desc_line_2')}
                <br />
              </FeatureDesc>
            </Reveal>
          </div>
          <a className="docker-url" href={t('docker_hub_link')} target="_blank" rel="noreferrer">
            <div className="content">docker pull artalk/artalk-go</div>
            <div
              className={['copy-btn', isCopied ? 'copied' : ''].join(' ')}
              onClick={copyDockerURLHandler}
            >
              {isCopied ? t('copied') : t('copy')}
            </div>
          </a>
          <div className="self-compile">
            <LearnMoreLink
              prompt={t('feature_swift_self_compile')}
              link={t('self_compile_guide_link')}
            />
          </div>
        </div>

        <div className="golang-icon">
          <svg
            width="292"
            height="109"
            viewBox="0 0 292 109"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <g clipPath="url(#clip0_406_120)">
              <path
                d="M22.035 32.9701C21.4664 32.9701 21.3242 32.6858 21.6086 32.2595L24.594 28.4225C24.8783 27.9961 25.5891 27.7119 26.1577 27.7119H76.9094C77.4781 27.7119 77.6203 28.1383 77.3359 28.5646L74.9192 32.2595C74.6349 32.6858 73.924 33.1122 73.4976 33.1122L22.035 32.9701Z"
                fill="#00ACD7"
              />
              <path
                d="M0.568647 46.0443C0 46.0443 -0.142162 45.7601 0.142162 45.3337L3.12756 41.4967C3.41188 41.0704 4.12269 40.7861 4.69133 40.7861H69.517C70.0857 40.7861 70.37 41.2125 70.2278 41.6388L69.0906 45.0495C68.9484 45.6179 68.3797 45.9022 67.8111 45.9022L0.568647 46.0443Z"
                fill="#00ACD7"
              />
              <path
                d="M34.9718 59.1185C34.4031 59.1185 34.2609 58.6922 34.5453 58.2658L36.5355 54.713C36.8199 54.2867 37.3885 53.8604 37.9571 53.8604H66.3895C66.9581 53.8604 67.2424 54.2867 67.2424 54.8551L66.9581 58.2658C66.9581 58.8343 66.3895 59.2606 65.963 59.2606L34.9718 59.1185Z"
                fill="#00ACD7"
              />
              <path
                d="M182.536 30.412C173.579 32.6857 167.466 34.3911 158.652 36.6649C156.52 37.2333 156.378 37.3754 154.53 35.2438C152.397 32.8279 150.834 31.2646 147.848 29.8435C138.892 25.438 130.22 26.717 122.117 31.9752C112.45 38.2281 107.474 47.4654 107.616 58.9765C107.759 70.3455 115.577 79.7249 126.808 81.2881C136.475 82.5671 144.578 79.1564 150.976 71.9087C152.255 70.3455 153.392 68.6401 154.814 66.6506C149.696 66.6506 143.299 66.6506 127.377 66.6506C124.391 66.6506 123.681 64.8031 124.676 62.3872C126.524 57.9817 129.936 50.5919 131.926 46.897C132.353 46.0443 133.348 44.6232 135.48 44.6232C142.73 44.6232 169.457 44.6232 187.227 44.6232C186.943 48.4602 186.943 52.2972 186.374 56.1343C184.81 66.3663 180.972 75.7457 174.717 83.9882C164.481 97.4889 151.118 105.873 134.201 108.147C120.269 109.995 107.332 107.295 95.9591 98.7679C85.4392 90.8096 79.4684 80.2933 77.9046 67.219C76.0565 51.7288 80.6057 37.8018 89.9883 25.5801C100.082 12.3637 113.445 3.97909 129.794 0.994739C143.157 -1.42117 155.951 0.142066 167.466 7.95823C175.001 12.9322 180.403 19.7535 183.957 27.996C184.81 29.2751 184.241 29.9856 182.536 30.412Z"
                fill="#00ACD7"
              />
              <path
                d="M229.591 109C216.654 108.716 204.855 105.021 194.904 96.4941C186.516 89.2464 181.256 80.0091 179.55 69.0664C176.991 53.0078 181.398 38.7965 191.065 26.1486C201.443 12.5058 213.953 5.4002 230.871 2.41584C245.371 -0.142175 259.019 1.27895 271.387 9.66356C282.617 17.3376 289.583 27.7118 291.431 41.3546C293.848 60.5397 288.304 76.172 275.083 89.5306C265.7 99.0521 254.185 105.021 240.964 107.721C237.126 108.431 233.287 108.574 229.591 109ZM263.426 51.5866C263.283 49.7392 263.283 48.3181 262.999 46.8969C260.44 32.8278 247.503 24.8696 233.998 27.996C220.777 30.9804 212.247 39.365 209.12 52.7235C206.561 63.8083 211.963 75.0351 222.199 79.5827C230.018 82.9934 237.836 82.5671 245.371 78.7301C256.602 72.9035 262.715 63.8083 263.426 51.5866Z"
                fill="#00ACD7"
              />
            </g>
            <defs>
              <clipPath id="clip0_406_120">
                <rect width="292" height="109" fill="white" />
              </clipPath>
            </defs>
          </svg>
        </div>
      </div>
    </FeatureBase>
  )
}

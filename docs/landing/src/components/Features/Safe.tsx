import React from 'react'
import { useTranslation } from 'react-i18next'
import './Safe.scss'
import { Title } from '../FeatureComps/Title'
import { Desc } from '../FeatureComps/Desc'
import { Reveal } from '../Reveal'
import { LearnMoreLink } from '../LearnMoreLink'
import { BlockBase } from '../FeatureComps/BlockBase'

export const SafeFeature: React.FC = () => {
  const { t } = useTranslation()

  return (
    <BlockBase className="safe-feature">
      <Title text={t('feature_safe_title')} color="#FFAC33" />

      <div className="content">
        <aside>
          <svg
            className="data-icon"
            width="236"
            height="257"
            viewBox="0 0 336 407"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              className="bg"
              fillRule="evenodd"
              clipRule="evenodd"
              d="M55 300V228.255L283.255 0H336V19.0001L55 300ZM55 345.255V371H101L336 136V64.2549L55 345.255ZM218 371H146.255L336 181.255V253L218 371ZM263.255 371H336V298.255L263.255 371ZM166.255 0H238L55 183V111.255L166.255 0ZM121 0H55V66L121 0Z"
              fill="#FFAC33"
            />
            <path
              className="text"
              d="M8.30273 195V85.084H53.0156C64.0352 85.084 73.4805 87.166 81.3516 91.3301C89.2227 95.4941 95.2402 101.613 99.4043 109.688C103.568 117.711 105.65 127.562 105.65 139.242V139.395C105.65 151.074 103.568 161.078 99.4043 169.406C95.2402 177.684 89.2227 184.031 81.3516 188.449C73.5312 192.816 64.0859 195 53.0156 195H8.30273ZM36.2578 172.453H49.2832C55.3262 172.453 60.4297 171.234 64.5938 168.797C68.7578 166.309 71.8809 162.627 73.9629 157.752C76.0957 152.877 77.1621 146.783 77.1621 139.471V139.318C77.1621 132.463 76.0703 126.674 73.8867 121.951C71.7031 117.229 68.5039 113.674 64.2891 111.287C60.125 108.85 55.123 107.631 49.2832 107.631H36.2578V172.453ZM107.936 195L144.955 85.084H170.701V109.002H161.789L137.338 195H107.936ZM129.568 170.93L135.814 150.744H188.373L194.619 170.93H129.568ZM186.85 195L162.322 109.002V85.084H179.232L216.252 195H186.85ZM33.8203 345V257.555H3.50391V235.084H92.0156V257.555H61.7754V345H33.8203ZM83.1035 345L120.123 235.084H145.869V259.002H136.957L112.506 345H83.1035ZM104.736 320.93L110.982 300.744H163.541L169.787 320.93H104.736ZM162.018 345L137.49 259.002V235.084H154.4L191.42 345H162.018Z"
              fill="var(--color-font)"
            />
            <path
              className="frame"
              fillRule="evenodd"
              clipRule="evenodd"
              d="M34 47H258V383H34V374H10V383V407H34H258H282V383V47V23H258H34H10V47V56H34V47Z"
              fill="#AAB8C2"
            />
          </svg>
        </aside>
        <main>
          <Desc>
            <Reveal>
              {t('feature_safe_desc_line_1')}
              <br />
              {t('feature_safe_desc_line_2')}
              <br />
              <br />
            </Reveal>
            <Reveal delay={500}>{t('feature_safe_desc_line_3')}</Reveal>
            <div className="report-audit">
              <LearnMoreLink
                prompt={t('feature_safe_report_audit_hint')}
                link={t('feature_safe_report_audit_link')}
                linkText={t('feature_safe_report_audit_btn')}
              />
            </div>
          </Desc>
        </main>
      </div>
    </BlockBase>
  )
}

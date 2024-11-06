import { PiMoonFill, PiSmileyWinkBold } from 'react-icons/pi'
import { useTranslation } from 'react-i18next'
import './Plugins.scss'
import {
  TbCardsFilled,
  TbLayoutSidebarRightExpandFilled,
  TbMailFilled,
  TbPlug,
  TbSocial,
  TbTerminal,
  TbTransformFilled,
} from 'react-icons/tb'
import { BiSolidBadgeCheck } from 'react-icons/bi'
import { LuNewspaper } from 'react-icons/lu'
import { FaArrowTrendUp } from 'react-icons/fa6'
import { FaMarkdown } from 'react-icons/fa'
import { RiRobot2Fill } from 'react-icons/ri'
import { BlockBase } from '../FeatureComps/BlockBase'

export const Plugins: React.FC = () => {
  const { t } = useTranslation()

  return (
    <BlockBase className="plugins-feature">
      <div className="feature-title" style={{ color: '#fff' }}>
        {t('feature_community_title')}
      </div>
      <div className="left-right">
        <main>
          <div className="feature-desc">{t('feature_community_desc')}</div>
          <div className="stamps">
            <div className="row">
              <PiSmileyWinkBold />
              <TbCardsFilled />
              <BiSolidBadgeCheck />
              <TbPlug />
            </div>
            <div className="row">
              <LuNewspaper />
              <TbLayoutSidebarRightExpandFilled />
              <TbSocial />
              <PiMoonFill />
            </div>
            <div className="row">
              <TbMailFilled />
              <FaArrowTrendUp />
              <TbTransformFilled />
              <FaMarkdown />
            </div>
            <div className="row">
              <TbTerminal />
              <RiRobot2Fill />
            </div>
          </div>
        </main>

        <aside></aside>
      </div>
    </BlockBase>
  )
}
